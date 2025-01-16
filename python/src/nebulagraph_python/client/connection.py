import json
import time
from dataclasses import dataclass, field
from typing import Any, Dict, List, Literal, Optional, Union

import grpc

from nebulagraph_python.client import constants
from nebulagraph_python.client.base_executor import NebulaBaseExecutor
from nebulagraph_python.client.logger import logger
from nebulagraph_python.data import HostAddress, SSLParam
from nebulagraph_python.error import (
    AuthenticatingError,
    ConnectingError,
    ErrorCode,
    ExecutingError,
    InternalError,
    NebulaGraphRemoteError,
)
from nebulagraph_python.proto import common_pb2, graph_pb2, graph_pb2_grpc
from nebulagraph_python.result_set import ResultSet


def _parse_hosts(hosts: Union[str, List[str], List[HostAddress]]) -> List[HostAddress]:
    """Convert various host formats to list of HostAddress objects"""
    if isinstance(hosts, str):
        hosts = hosts.split(",")

    addresses = []
    for host in hosts:
        if isinstance(host, HostAddress):
            addresses.append(host)
        else:
            addr, port = host.split(":")
            addresses.append(HostAddress(addr, int(port)))
    return addresses


@dataclass
class ConnectionConfig:
    hosts: List[HostAddress]
    ssl_param: Optional[SSLParam] = None
    connect_timeout: Optional[float] = constants.DEFAULT_CONNECT_TIMEOUT
    request_timeout: Optional[float] = constants.DEFAULT_REQUEST_TIMEOUT

    @classmethod
    def from_defults(
        cls,
        hosts: Union[str, List[str], List[HostAddress]],
        ssl_param: Union[SSLParam, Literal[True], None] = None,
        connect_timeout: Optional[float] = None,
        request_timeout: Optional[float] = None,
    ):
        if ssl_param is True:
            ssl_param = SSLParam()
        return cls(
            hosts=_parse_hosts(hosts),
            ssl_param=ssl_param,
            connect_timeout=connect_timeout,
            request_timeout=request_timeout,
        )


@dataclass
class SessionConfig:
    schema: Optional[str] = None
    graph: Optional[str] = None
    timezone: Optional[str] = None
    parameters: Dict[str, Any] = field(default_factory=dict)


@dataclass
class _Connection:
    """Represents a connection to a NebulaGraph server.

    Required to explicitly call `close()` to release all resources.
    """

    # Config
    config: ConnectionConfig

    # Owned Resources
    _stub: Optional[graph_pb2_grpc.GraphServiceStub] = field(default=None, init=False)
    _channel: Optional[grpc.Channel] = field(default=None, init=False)

    def __post_init__(self):
        self.connect()

    def connect(self):
        """Establish connection to NebulaGraph"""
        last_error = None

        # Try each address until one succeeds
        for host_addr in self.config.hosts:
            try:
                channel_options = [
                    ("grpc.max_send_message_length", -1),
                    ("grpc.max_receive_message_length", -1),
                    ("grpc.enable_deadline_checking", 1),
                ]

                if self.config.ssl_param:
                    self._channel = grpc.secure_channel(
                        f"{host_addr.host}:{host_addr.port}",
                        credentials=grpc.ssl_channel_credentials(
                            root_certificates=self.config.ssl_param.ca_crt,
                            private_key=self.config.ssl_param.private_key,
                            certificate_chain=self.config.ssl_param.cert,
                        ),
                        options=channel_options,
                    )
                else:
                    self._channel = grpc.insecure_channel(
                        f"{host_addr.host}:{host_addr.port}",
                        options=channel_options,
                    )

                # Wait for channel to be ready with timeout
                start_time = time.time()
                while True:
                    # Get current state and try to connect
                    state = self._channel._channel.check_connectivity_state(True)
                    if state == grpc.ChannelConnectivity.READY.value[0]:
                        break
                    if (
                        self.config.connect_timeout is not None
                        and time.time() - start_time > self.config.connect_timeout
                    ):
                        raise ConnectingError(
                            f"Connection timeout after {self.config.connect_timeout} seconds"
                        )
                    self._channel._channel.watch_connectivity_state(
                        state,
                        (
                            self.config.connect_timeout + start_time
                            if self.config.connect_timeout is not None
                            else None
                        ),
                    )

            except Exception as e:
                logger.warning(
                    f"Failed to connect to {(host_addr.host, host_addr.port) if host_addr else 'No Available Addr'}: {e}",
                )
                last_error = e
                if self._channel:
                    self._channel.close()
                    self._channel = None
                self._stub = None

            self._stub = graph_pb2_grpc.GraphServiceStub(self._channel)
            return

        # If we get here, all connection attempts failed
        raise ConnectingError(
            f"Failed to connect to any of the provided hosts. Last error: {last_error}",
        )

    def close(self):
        """Close the connection. No Exception will be raised but an error will be logged."""
        try:
            if self._channel:
                self._channel.close()
                self._channel = None
            self._stub = None
        except Exception:
            logger.exception("Failed to close connection")

    def execute(
        self, session_id: int, statement: str, *, timeout: Optional[float] = None
    ) -> ResultSet:
        if not self._stub:
            raise InternalError("Connection not established")

        logger.debug(f"Executing statement: {statement}")

        try:
            request = graph_pb2.ExecuteRequest(
                session_id=session_id,
                stmt=statement.encode("utf-8"),
            )
            logger.debug(f"Request: {request}")
            # Use request_timeout as default if timeout is not specified
            effective_timeout = (
                timeout if timeout is not None else self.config.request_timeout
            )
            response = self._stub.Execute(request, timeout=effective_timeout)
            logger.debug(f"Response: {response}")
        except grpc.RpcError as e:
            raise ExecutingError() from e

        return ResultSet(response)

    def authenticate(
        self,
        username: str,
        password: str,
        *,
        auth_options: Dict[str, Any],
        session_config: SessionConfig,
    ) -> "_Session":
        """Authenticate with NebulaGraph and return session ID"""
        if not self._stub:
            raise InternalError("Connection not established")

        client_info = common_pb2.ClientInfo(
            lang=common_pb2.ClientInfo.PYTHON,
            protocol_version=b"5.0.0",
        )

        auth_info = json.dumps(
            {"password": password, **auth_options},
        ).encode("utf-8")

        request = graph_pb2.AuthRequest(
            username=username.encode("utf-8"),
            auth_info=auth_info,
            client_info=client_info,
        )

        try:
            response = self._stub.Authenticate(
                request, timeout=self.config.request_timeout
            )
        except grpc.RpcError as e:
            raise AuthenticatingError() from e

        if response.status.code != b"00000":
            raise NebulaGraphRemoteError(
                code=ErrorCode.from_str(response.status.code.decode("utf-8")),
                message=response.status.message.decode("utf-8"),
            )
        return _Session(self, int(response.session_id), session_config)


@dataclass
class _Session(NebulaBaseExecutor):
    """Represents a session with the NebulaGraph server.

    Required to explicitly call `close()` to release all resources.
    """

    # Borrowed Resources
    _conn: _Connection

    # Config
    session_id: int
    config: SessionConfig = field(default_factory=SessionConfig)

    def __post_init__(self):
        if self.config.schema is not None:
            self._conn.execute(
                self.session_id, f"SESSION SET SCHEMA `{self.config.schema}`"
            ).raise_on_error()
        if self.config.graph is not None:
            self._conn.execute(
                self.session_id, f"SESSION SET GRAPH `{self.config.graph}`"
            ).raise_on_error()
        if self.config.timezone is not None:
            self._conn.execute(
                self.session_id, f"SESSION SET TIME ZONE `{self.config.timezone}`"
            ).raise_on_error()
        if self.config.parameters:
            self._conn.execute(
                self.session_id,
                f"SESSION SET VALUE {','.join(f'${key}={value}' for key, value in self.config.parameters.items())}",
            ).raise_on_error()

    def execute(self, statement: str, *, timeout: Optional[float] = None) -> ResultSet:
        return self._conn.execute(
            self.session_id, statement, timeout=timeout
        ).raise_on_error()

    def close(self):
        """Close the session. No Exception will be raised but an error will be logged."""
        try:
            self.execute("SESSION CLOSE")
        except Exception:
            logger.exception("Failed to close session")
