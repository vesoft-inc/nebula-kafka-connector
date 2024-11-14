import json
import logging
import os
import sys
from dataclasses import dataclass
from typing import Any, Dict, List, Optional, Union

import grpc

from nebulagraph_python.data import HostAddress, SSLParam
from nebulagraph_python.error import (
    ErrorCode,
    NebulaGraphRemoteError,
    ConnectionError,
)
from nebulagraph_python.proto import common_pb2, graph_pb2, graph_pb2_grpc
from nebulagraph_python.result_set import ResultSet

logger = logging.getLogger(__name__)

# Configure logging based on environment variables
log_level = os.getenv("NG_PYTHON_LOG_LEVEL", "INFO")
logger_sink = os.getenv("NG_PYTHON_LOG_SINK", "stdout")
debug_flag = os.getenv("NG_PYTHON_DEBUG", "false").lower() == "true"

# Set base log level
logger.setLevel(log_level)

# Add debug handler if debug logging enabled
if log_level == "DEBUG":
    # Create handler based on sink config
    handler = (
        logging.StreamHandler(sys.stdout)
        if logger_sink == "stdout"
        else logging.FileHandler(logger_sink)
    )

    # Add formatter for debug logs
    formatter = logging.Formatter(
        "%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )
    handler.setFormatter(formatter)

    logger.addHandler(handler)


@dataclass
class AuthResult:
    session_id: int


class NebulaClientConfig:
    """Configuration for NebulaGraph client"""

    username: str
    password: str
    auth_options: Dict[str, Any]
    connect_timeout_ms: int
    request_timeout_ms: int
    ssl_param: Optional[SSLParam]

    def __init__(
        self,
        hosts: Union[str, List[str], List[HostAddress]],
        username: str,
        password: str,
        *,  # Force keyword arguments after these
        auth_options: Optional[Dict[str, Any]] = None,
        connect_timeout_ms: int = 1000,
        request_timeout_ms: int = 3000,
        ssl_param: Optional[SSLParam] = None,
    ):
        self.username = username
        self.password = password
        self.auth_options = auth_options or {}
        self.connect_timeout_ms = connect_timeout_ms
        self.request_timeout_ms = request_timeout_ms
        self.ssl_param = ssl_param

        # Convert hosts to HostAddress objects
        self.addresses = self._parse_hosts(hosts)

    def _parse_hosts(self, hosts) -> List[HostAddress]:
        """Convert various host formats to list of HostAddress objects"""
        if isinstance(hosts, str):
            hosts = [hosts]

        addresses = []
        for host in hosts:
            if isinstance(host, HostAddress):
                addresses.append(host)
            else:
                host, port = host.split(":")
                addresses.append(HostAddress(host, int(port)))
        return addresses


class Session:
    """Represents a session with the NebulaGraph server"""

    session_id: int
    username: str
    graph_addr: HostAddress
    timeout: int
    _stub: Optional[graph_pb2_grpc.GraphServiceStub]
    _channel: Optional[grpc.Channel]

    def __init__(
        self,
        session_id: int,
        username: str,
        graph_addr: HostAddress,
        timeout: int = 0,
    ):
        self.session_id = session_id
        self.username = username
        self.graph_addr = graph_addr
        self.timeout = timeout
        self._stub = None
        self._channel = None

    def execute(self, statement: str) -> ResultSet:
        """Execute a nGQL statement

        Args:
        ----
            statement: The nGQL statement to execute

        Returns:
        -------
            ResultSet containing the execution results

        Raises:
        ------
            Exception if execution fails

        """
        if not self._stub:
            raise Exception("Session not connected")

        logger.debug(f"Executing statement: {statement}")

        request = graph_pb2.ExecuteRequest(
            session_id=self.session_id,
            stmt=statement.encode("utf-8"),
        )

        logger.debug(f"Request: {request}")

        try:
            response = self._stub.Execute(request)
            logger.debug(f"Response: {response}")
            return ResultSet(response)
        except grpc.RpcError as e:
            logger.error(f"Failed to execute statement: {e}")
            raise

    def print_query_result(
        self,
        query: str,
        style: str = "table",
        width: Optional[int] = None,
        min_width: int = 8,
        max_width: Optional[int] = None,
        padding: int = 1,
        collapse_padding: bool = False,
    ) -> None:
        """Execute a query and print the results in a formatted way using rich

        Args:
        ----
            query: The nGQL query to execute
            style: Output style - either "table" (default) or "rows"
            width: Fixed width for all columns. If None, width will be auto-calculated
            min_width: Minimum width of columns when using table style
            max_width: Maximum width of columns. If None, no maximum is enforced
            padding: Number of spaces around cell contents in table style
            collapse_padding: Reduce padding when cell contents are too wide

        Raises:
        ------
            Exception if execution fails

        """
        try:
            result = self.execute(query)
            result.print(
                style=style,
                width=width,
                min_width=min_width,
                max_width=max_width,
                padding=padding,
                collapse_padding=collapse_padding,
            )
        except Exception as e:
            from rich.console import Console
            from rich.traceback import Traceback

            console = Console()
            console.print(f"[bold red]Error executing query:[/bold red] {e!s}")
            if debug_flag:
                console.print(Traceback.from_exception(type(e), e, e.__traceback__))

    def pq(self, query: str, **kwargs):
        """Print query result using rich"""
        self.print_query_result(query, **kwargs)

    def close(self):
        """Close the session and cleanup resources"""
        if self._channel:
            self._channel.close()
            self._channel = None
            self._stub = None


class NebulaClient:
    """Client for connecting to NebulaGraph"""

    config: NebulaClientConfig
    _session: Optional[Session]
    _round_robin_index: int

    def __init__(
        self,
        hosts: Union[str, List[str], List[HostAddress]],
        username: str,
        password: str,
        **kwargs,
    ):
        """Initialize NebulaGraph client

        Args:
        ----
            hosts: Single host string ("hostname:port"), list of host strings,
                  or list of HostAddress objects
            username: Username for authentication
            password: Password for authentication
            **kwargs: Additional configuration options:
                - auth_options: dict of authentication options
                - connect_timeout_ms: Connection timeout in milliseconds
                - request_timeout_ms: Request timeout in milliseconds
                - ssl_param: SSL configuration

        """
        self.config = NebulaClientConfig(hosts, username, password, **kwargs)
        self._session = None
        self._round_robin_index = 0
        self._connect()

    def _get_next_address(self) -> HostAddress:
        """Get next address using round-robin"""
        addr = self.config.addresses[self._round_robin_index]
        self._round_robin_index = (self._round_robin_index + 1) % len(
            self.config.addresses,
        )
        return addr

    def _connect(self):
        """Establish connection to NebulaGraph"""
        last_error = None

        # Try each address until one succeeds
        host_addr = None
        channel = None
        for _ in range(len(self.config.addresses)):
            try:
                host_addr = self._get_next_address()
                channel_options = [
                    ("grpc.max_send_message_length", -1),
                    ("grpc.max_receive_message_length", -1),
                ]

                if self.config.ssl_param:
                    # TODO: Add SSL credentials
                    pass

                channel = grpc.insecure_channel(
                    f"{host_addr.host}:{host_addr.port}",
                    options=channel_options,
                )

                # Create stub
                stub = graph_pb2_grpc.GraphServiceStub(channel)

                # Create auth request
                client_info = common_pb2.ClientInfo(
                    lang=common_pb2.ClientInfo.PYTHON,
                    protocol_version=b"5.0.0",
                )

                auth_info = json.dumps(
                    {"password": self.config.password, **self.config.auth_options},
                ).encode("utf-8")

                request = graph_pb2.AuthRequest(
                    username=self.config.username.encode("utf-8"),
                    auth_info=auth_info,
                    client_info=client_info,
                )

                response = stub.Authenticate(request)

                if response.status.code != b"00000":
                    raise NebulaGraphRemoteError(
                        code=ErrorCode.from_str(response.status.code.decode("utf-8")),
                        message=response.status.message.decode("utf-8"),
                    )

                # Create internal session
                self._session = Session(
                    session_id=response.session_id,
                    username=self.config.username,
                    graph_addr=host_addr,
                    timeout=self.config.request_timeout_ms,
                )
                self._session._channel = channel
                self._session._stub = stub

                # Successfully connected
                return

            except Exception as e:
                logger.warning(
                    f"Failed to connect to {(host_addr.host, host_addr.port) if host_addr else 'No Available Addr'}: {e}",
                )
                last_error = e
                if channel:
                    channel.close()

        # If we get here, all connection attempts failed
        raise ConnectionError(
            f"Failed to connect to any of the provided hosts. Last error: {last_error}",
        )

    def execute(self, statement: str) -> ResultSet:
        """Execute a nGQL statement

        Args:
        ----
            statement: The nGQL statement to execute

        Returns:
        -------
            ResultSet containing the execution results

        Raises:
        ------
            Exception: If client is not connected or execution fails

        """
        if not self._session:
            raise Exception("Client not connected")
        return self._session.execute(statement)

    def close(self):
        """Close the client and cleanup resources"""
        if self._session:
            self._session.close()
            self._session = None

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()
