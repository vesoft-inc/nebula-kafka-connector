import logging
from typing import Any, Dict, List, Literal, Optional, Union

from nebulagraph_python.client.base_executor import NebulaBaseExecutor
from nebulagraph_python.client.connection import (
    ConnectionConfig,
    SessionConfig,
    _Connection,
    _Session,
)
from nebulagraph_python.data import HostAddress, SSLParam
from nebulagraph_python.result_set import ResultSet

logger = logging.getLogger(__name__)


class NebulaClient(NebulaBaseExecutor):
    """A high level client for connecting to NebulaGraph. It maintains a connection and a session.

    Required to explicitly call `close()` to release all resources.
    """

    # Owned Resources
    _conn: _Connection
    _session: _Session

    def __init__(
        self,
        hosts: Union[str, List[str], List[HostAddress]],
        username: str,
        password: str,
        *,
        ssl_param: Union[SSLParam, Literal[True], None] = None,
        auth_options: Optional[Dict[str, Any]] = None,
        conn_config: Optional[ConnectionConfig] = None,
        session_config: Optional[SessionConfig] = None,
    ):
        """Initialize NebulaGraph client

        Args:
        ----
            hosts: Single host string ("hostname:port"), list of host strings,
                  or list of HostAddress objects
            username: Username for authentication
            password: Password for authentication
            ssl_param: SSL configuration
            auth_options: dict of authentication options
            conn_config: Connection configuration. If provided, it overrides hosts and ssl_param.
            session_config: Session configuration.
        """
        self._conn = _Connection(
            conn_config or ConnectionConfig.from_defults(hosts, ssl_param)
        )
        self._session = self._conn.authenticate(
            username=username,
            password=password,
            session_config=session_config or SessionConfig(),
            auth_options=auth_options or {},
        )

    def execute(self, statement: str, *, timeout: Optional[float] = None) -> ResultSet:
        return self._session.execute(statement, timeout=timeout)

    def ping(self, timeout: Optional[float] = None) -> bool:
        try:
            res = (
                self.execute(statement="RETURN 1", timeout=timeout).one().as_primitive()
            )
            if not res == {"1": 1}:
                raise ValueError(f"Unexpected result from ping: {res}")
            return True
        except Exception:
            logger.exception("Failed to ping NebulaGraph")
            return False

    def close(self):
        """Close the client connection and session. No Exception will be raised but an error will be logged."""
        self._session.close()
        self._conn.close()
