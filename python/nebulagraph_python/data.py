import enum
from dataclasses import dataclass
from typing import List, Optional

from .proto.graph_pb2 import PlanInfo


@dataclass
class HostAddress:
    """Represents a NebulaGraph service address"""

    host: str
    port: int

    def __str__(self):
        return f"{self.host}:{self.port}"


class SSLParam:
    """Base class for SSL parameters"""

    class SignMode(enum.Enum):
        NONE = 0
        SELF_SIGNED = 1
        CA_SIGNED = 2

    def __init__(self, sign_mode: SignMode):
        self.sign_mode = sign_mode


@dataclass
class SelfSignedSSLParam(SSLParam):
    """SSL parameters for self-signed certificates"""

    crt_file_path: str
    key_file_path: str
    password: str

    def __init__(self, crt_file_path: str, key_file_path: str, password: str):
        super().__init__(SSLParam.SignMode.SELF_SIGNED)
        self.crt_file_path = crt_file_path
        self.key_file_path = key_file_path
        self.password = password


@dataclass
class CASignedSSLParam(SSLParam):
    """SSL parameters for CA-signed certificates"""

    ca_crt_file_path: str
    crt_file_path: str
    key_file_path: str

    def __init__(self, ca_crt_file_path: str, crt_file_path: str, key_file_path: str):
        super().__init__(SSLParam.SignMode.CA_SIGNED)
        self.ca_crt_file_path = ca_crt_file_path
        self.crt_file_path = crt_file_path
        self.key_file_path = key_file_path


class PlanInfoNode:
    def __init__(self, plan_info: PlanInfo):
        self.plan_info = plan_info
        self.id = plan_info.id.decode()
        self.name = plan_info.name.decode()
        self.details = plan_info.details.decode()
        self.time_ms = plan_info.time_ms
        self.rows = plan_info.rows
        self.memory_kib = plan_info.memory_kib
        self.blocked_ms = plan_info.blocked_ms
        self.queued_ms = plan_info.queued_ms
        self.consume_ms = plan_info.consume_ms
        self.produce_ms = plan_info.produce_ms
        self.finish_ms = plan_info.finish_ms
        self.batches = plan_info.batches
        self.concurrency = plan_info.concurrency
        self.other_stats_json = plan_info.other_stats_json.decode()
        self.children = [PlanInfoNode(plan) for plan in plan_info.children]

    def get_plan_id(self) -> str:
        return self.id

    def get_id(self) -> str:
        return self.id

    def get_name(self) -> str:
        return self.name

    def get_details(self) -> str:
        return self.details

    def get_time_ms(self) -> float:
        return self.time_ms

    def get_rows(self) -> int:
        return self.rows

    def get_memory_kib(self) -> float:
        return self.memory_kib

    def get_blocked_ms(self) -> float:
        return self.blocked_ms

    def get_children(self) -> List["PlanInfoNode"]:
        return self.children


class ExtraInfo:
    """Class that maintains additional information for execution result."""

    def __init__(self):
        self.cursor = None
        self.affected_nodes = 0
        self.affected_edges = 0
        self.total_server_time_us = 0
        self.build_time_us = 0
        self.optimize_time_us = 0
        self.serialize_time_us = 0

    def set_cursor(self, cursor: str) -> None:
        self.cursor = cursor

    def set_affected_nodes(self, affected_nodes: int) -> None:
        self.affected_nodes = affected_nodes

    def set_affected_edges(self, affected_edges: int) -> None:
        self.affected_edges = affected_edges

    def set_total_server_time_us(self, total_server_time_us: int) -> None:
        self.total_server_time_us = total_server_time_us

    def set_build_time_us(self, build_time_us: int) -> None:
        self.build_time_us = build_time_us

    def set_optimize_time_us(self, optimize_time_us: int) -> None:
        self.optimize_time_us = optimize_time_us

    def set_serialize_time_us(self, serialize_time_us: int) -> None:
        self.serialize_time_us = serialize_time_us

    def get_cursor(self) -> Optional[str]:
        return self.cursor

    def get_affected_nodes(self) -> int:
        return self.affected_nodes

    def get_affected_edges(self) -> int:
        return self.affected_edges

    def get_total_server_time_us(self) -> int:
        return self.total_server_time_us

    def get_build_time_us(self) -> int:
        return self.build_time_us

    def get_optimize_time_us(self) -> int:
        return self.optimize_time_us

    def get_serialize_time_us(self) -> int:
        return self.serialize_time_us

    def __str__(self) -> str:
        return (
            f"ExtraInfo{{cursor='{self.cursor}', "
            f"affectedNodes={self.affected_nodes}, "
            f"affectedEdges={self.affected_edges}, "
            f"totalServerTimeUs={self.total_server_time_us}, "
            f"buildTimeUs={self.build_time_us}, "
            f"optimizeTimeUs={self.optimize_time_us}, "
            f"serializeTimeUs={self.serialize_time_us}}}"
        )
