# Copyright (c) 2024 vesoft inc. All rights reserved.
#
# This source code is licensed under Apache 2.0 License.

from typing import Any, List

from nebulagraph_python.data_types import ByteOrder, ResultGraphSchemas
from nebulagraph_python.decode import Batch, BytesReader
from nebulagraph_python.error import InternalError
from nebulagraph_python.proto.vector_pb2 import VectorResultTable
from nebulagraph_python.value_parser import DataType, ValueParser, ValueTypeParser
from nebulagraph_python.value_wrapper import Row


class ResultTable:
    result_table: VectorResultTable
    byte_order: ByteOrder
    parser: ValueParser
    num_batches: int
    column_names: List[str]
    column_data_types: List[DataType]
    total_num_records: int

    def __init__(self, table: Any):
        if not isinstance(table, VectorResultTable):
            raise InternalError("table must be a VectorResultTable")

        self.result_table = table

        graph_schemas = ResultGraphSchemas(table.meta.graph_schema)
        time_zone_offset = table.meta.time_zone_offset

        if table.meta.is_little_endian:
            self.byte_order = ByteOrder.LITTLE_ENDIAN
        else:
            self.byte_order = ByteOrder.BIG_ENDIAN

        self.parser = ValueParser(graph_schemas, time_zone_offset, self.byte_order)
        self.num_batches = table.meta.num_batches
        if self.num_batches != len(table.batch):
            raise RuntimeError("the number of batch is not equal to numBatches")

        self.column_names = list(table.meta.row_type.column_names)
        self.column_data_types = []
        value_type_parser = ValueTypeParser(self.byte_order)
        for col_type in table.meta.row_type.column_types:
            if not col_type.value_type:
                raise ValueError("Invalid column type: empty value_type")
            data_type = value_type_parser.get_data_type(
                BytesReader(col_type.value_type),
            )
            self.column_data_types.append(data_type)

        self.total_num_records = table.meta.num_records

    def _get_row_by_index(self, batch: Batch, index: int) -> Row:
        """Parse row record from batch

        Args:
        ----
            batch: the batch to parse from
            index: the position of each vector in current batch

        Returns:
        -------
            Row: row record

        """

        row = Row()
        for i in range(batch.get_vectors_count()):
            value = self.parser.decode_value_wrapper(
                batch.get_vectors(i),
                self.column_data_types[i],
                index,
            )
            row.add_value(value)
        return row

    def __iter__(self):
        return self.rows()

    def rows(self):
        """Generator that yields rows from the result table.

        Returns
        -------
            Iterator[Row]: iterator of row records

        Raises
        ------
            InternalError: if no result table data
        """
        for batch_index in range(self.num_batches):
            current_batch = Batch(self.result_table.batch[batch_index], self.byte_order)

            # each VectorMetaData has the same numRecords value,
            # just use the first one to get the numRecord for this batch
            current_batch_row_size = 0
            if current_batch.get_vectors_count() != 0:
                current_batch_row_size = current_batch.get_batch_row_size()

            # Skip empty batches
            if current_batch.get_vectors_count() == 0:
                continue

            # Process rows in current batch
            for row_index in range(current_batch_row_size):
                row = self._get_row_by_index(current_batch, row_index)
                yield row
