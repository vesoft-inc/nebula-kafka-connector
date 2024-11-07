from collections.abc import Iterable, Iterator
from threading import Lock
from typing import Any, List, Optional, Union

from .data import ExtraInfo, PlanInfoNode
from .error_code import ErrorCode
from .proto.graph_pb2 import ExecuteResponse, QueryStats
from .result_table import ResultTable
from .value_wrapper import Row, ValueWrapper


class Record(Iterable[ValueWrapper]):
    def __init__(self, column_names: Optional[List[str]], row: Row):
        self.column_names: List[str] = []
        self.col_values: List[ValueWrapper] = []

        if column_names is None:
            return

        if row is None or not row.values:
            return

        for value in row.values:
            self.col_values.append(value)

        self.column_names = column_names

    def __iter__(self) -> Iterator[ValueWrapper]:
        return iter(self.col_values)

    def __str__(self) -> str:
        value_strs = [str(v.cast()) for v in self.col_values]
        return f"ColumnName: {self.column_names}, Values: {value_strs}"

    def get(self, key: Union[int, str]) -> ValueWrapper:
        if isinstance(key, int):
            if key >= len(self.column_names):
                raise ValueError(
                    f"Cannot get field because the key '{key}' out of range",
                )
            return self.col_values[key]
        index = self.column_names.index(key)
        if index == -1:
            raise ValueError(
                f"Cannot get field because the columnName '{key}' is not exists",
            )
        return self.col_values[index]

    def values(self) -> List[ValueWrapper]:
        return self.col_values

    def size(self) -> int:
        return len(self.column_names)

    def contains(self, column_name: str) -> bool:
        return column_name in self.column_names

    def for_each(self, action):
        for value in self.col_values:
            action(value)

    def spliterator(self):
        return self.col_values.__iter__()


class ResultSet:
    def __init__(self, response: ExecuteResponse):
        if response is None:
            raise RuntimeError("got null object for server's response")

        self.response = response
        self.column_names: List[str] = []
        self.empty = False
        self._index = 0
        self._index_lock = Lock()

        if not response.HasField("result") or not response.result.HasField("meta"):
            self.size = 0
            self.result_table = None
            self.empty = True
            return

        self.result_table = ResultTable(response.result)
        self.column_names.extend(self.result_table.get_column_names())
        self.size = self.result_table.get_total_num_records()

        self.empty = self.size == 0

    def is_succeeded(self) -> bool:
        return self.response.status.code == b"00000"

    def is_empty(self) -> bool:
        return self.empty

    def get_error_code(self) -> ErrorCode:
        return ErrorCode.find(self.response.status.code.decode("utf-8"))

    def get_error_message(self) -> str:
        return self.response.status.message.decode("utf-8")

    def get_latency(self) -> int:
        return self.response.summary.elapsed_time.total_server_time_us

    def get_plan_desc(self) -> PlanInfoNode:
        return PlanInfoNode(self.response.summary.plan_info)

    def get_column_names(self) -> List[str]:
        return self.column_names

    def row_size(self) -> int:
        if self.empty:
            return 0
        return self.size

    def has_next(self) -> bool:
        if self.empty:
            return False
        return self._index < self.size

    def next(self) -> Record:
        if not self.has_next():
            raise StopIteration("no more row record data")
        row = self.result_table.next()
        with self._index_lock:
            self._index += 1
        return Record(self.column_names, row)

    def get_extra_info(self) -> ExtraInfo:
        if not self.response.HasField("summary"):
            return ExtraInfo()

        extra_info = ExtraInfo()
        query_stats: QueryStats = self.response.summary.query_stats
        extra_info.affected_nodes = query_stats.num_affected_nodes
        extra_info.affected_edges = query_stats.num_affected_edges
        extra_info.cursor = self.response.cursor.decode("utf-8")
        extra_info.build_time_us = self.response.summary.elapsed_time.build_time_us
        extra_info.optimize_time_us = (
            self.response.summary.elapsed_time.optimize_time_us
        )
        extra_info.serialize_time_us = (
            self.response.summary.elapsed_time.serialize_time_us
        )
        extra_info.total_server_time_us = (
            self.response.summary.elapsed_time.total_server_time_us
        )
        return extra_info

    def __str__(self) -> str:
        if not self.is_succeeded():
            return self.response.status.message.decode("utf-8")
        return f"ColumnName: {self.column_names}, RowSize: {self.row_size()}, Latency: {self.get_latency()}"

    def as_pandas_df(self) -> Any:
        """Convert result set to pandas DataFrame.

        Returns:
            pandas.DataFrame: DataFrame containing the query results

        Raises:
            ImportError: If pandas is not installed

        """
        try:
            import pandas as pd
        except ImportError:
            raise ImportError(
                "pandas is required to use this method. Please install it using 'pip install pandas'",
            )

        if not self.is_succeeded():
            raise RuntimeError(f"Query failed: {self.get_error_message()}")

        if self.is_empty():
            return pd.DataFrame(columns=self.column_names)

        # Reset index to start
        self._index = 0

        # Build list of rows
        rows = []
        while self.has_next():
            record = self.next()
            row = []
            for val in record.values():
                row.append(val.cast_primitive() if val is not None else None)
            rows.append(row)

        return pd.DataFrame(rows, columns=self.column_names)

    def as_ascii_table(
        self,
        style: str = "table",
        width: int = None,
        min_width: int = 8,
        max_width: int = None,
        padding: int = 1,
        collapse_padding: bool = False,
    ) -> str:
        """Print query results in a formatted table or row-by-row format.

        Args:
            style: Output style - either "table" (default) or "rows"
            width: Fixed width for all columns. If None, width will be auto-calculated
            min_width: Minimum width of columns when using table style
            max_width: Maximum width of columns. If None, no maximum is enforced
            padding: Number of spaces around cell contents in table style
            collapse_padding: Reduce padding when cell contents are too wide

        Examples:
            # Print as table (default)
            result.as_ascii_table()

            # Print as rows
            result.as_ascii_table(style="rows")

            # Customize table formatting
            result.as_ascii_table(width=20, max_width=30, padding=2)

        Returns:
            str: Formatted representation of the results, or error message if query failed

        """
        try:
            from io import StringIO

            from rich import box
            from rich.console import Console
            from rich.table import Table
        except ImportError:
            raise ImportError(
                "The 'rich' library is required to use this method. Please install it using 'pip install rich'.",
            )

        console = Console(file=StringIO(), force_terminal=False)

        if not self.is_succeeded():
            console.print(f"[bold red]Error:[/bold red] {self.get_error_message()}")
            return console.file.getvalue()

        if self.is_empty():
            console.print("[yellow]Empty result set[/yellow]")
            return console.file.getvalue()

        # Reset index to start
        self._index = 0
        column_names = self.get_column_names()

        if style == "rows":
            # Row-by-row format
            row_num = 1
            try:
                while self.has_next():
                    record = self.next()
                    console.print(f"\n[bold blue]Row {row_num}[/bold blue]")
                    for col, val in zip(column_names, record.values()):
                        console.print(f"  [cyan]{col}:[/cyan] {val}")
                    row_num += 1
            except StopIteration:
                pass

        else:
            # Table format
            table = Table(
                box=box.DOUBLE_EDGE,
                show_header=True,
                header_style="bold cyan",
                width=width,
                min_width=min_width,
                padding=padding,
                collapse_padding=collapse_padding,
            )

            for header in column_names:
                table.add_column(header, max_width=max_width)

            try:
                while self.has_next():
                    record = self.next()
                    table.add_row(*[str(v.cast()) for v in record.values()])
            except StopIteration:
                pass

            console.print(table)

        # Print summary
        console.print("\n[bold green]Summary[/bold green]")
        console.print(f"├── [green]Rows:[/green] {self.row_size()}")
        console.print(f"└── [blue]Latency:[/blue] {self.get_latency()}μs")

        return console.file.getvalue()

    def print(
        self,
        style: str = "table",
        width: int = None,
        min_width: int = 8,
        max_width: int = None,
        padding: int = 1,
        collapse_padding: bool = False,
    ):
        """Print the results directly to console with rich formatting.

        Args are the same as as_ascii_table().
        """
        from rich.console import Console

        console = Console()
        console.print(
            self.as_ascii_table(
                style,
                width,
                min_width,
                max_width,
                padding,
                collapse_padding,
            ),
        )
