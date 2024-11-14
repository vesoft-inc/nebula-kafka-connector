# Copyright (c) 2024 vesoft inc. All rights reserved.
#
# This source code is licensed under Apache 2.0 License.

from typing import List

from nebulagraph_python.data_types import ByteOrder, ResultGraphSchemas
from nebulagraph_python.decode import Batch, BytesReader
from nebulagraph_python.proto.vector_pb2 import VectorResultTable
from nebulagraph_python.value_parser import ValueParser, ValueTypeParser
from nebulagraph_python.value_wrapper import Row


class ResultTable:
    def __init__(self, table: VectorResultTable):
        if table is None or not table.HasField("meta"):
            self.result_table = None
            return

        self.result_table = table
        graph_schemas = ResultGraphSchemas(table.meta.graph_schema)
        time_zone_offset = table.meta.time_zone_offset

        if table.meta.is_little_endian:
            self.byte_order = ByteOrder.LITTLE_ENDIAN
        else:
            self.byte_order = ByteOrder.BIG_ENDIAN

        self.parser = ValueParser(graph_schemas, time_zone_offset, self.byte_order)
        self.num_batches = table.meta.num_batches
        batches = list(table.batch)
        if self.num_batches != len(batches):
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

        self.batch_index = 0
        self.current_batch = None
        self.current_batch_row_index = 0

        if table.batch:
            self.current_batch = Batch(table.batch[0], self.byte_order)

    def get_column_names(self) -> List[str]:
        """Get the column names of the response

        Returns
        -------
            List[str]: list of column names

        """
        return self.column_names

    def get_total_num_records(self) -> int:
        """Get the total data records size of the response

        Returns
        -------
            int: total number of records

        """
        if self.result_table is None:
            raise RuntimeError("no result table data")
        return self.result_table.meta.num_records

    def _get_row_by_index(self, index: int) -> Row:
        """Parse row record from batch

        Args:
        ----
            index: the position of each vector in current batch

        Returns:
        -------
            Row: row record

        """
        if self.current_batch is None:
            raise RuntimeError("no more batch data")
        row = Row()
        for i in range(self.current_batch.get_vectors_count()):
            value = self.parser.decode_value_wrapper(
                self.current_batch.get_vectors(i),
                self.column_data_types[i],
                index,
            )
            row.add_value(value)
        return row

    def next(self) -> Row:
        """Get the next row data

        Returns
        -------
            Row: next row record

        Raises
        ------
            RuntimeError: if no more batch data

        """
        if self.result_table is None:
            raise RuntimeError("no result table data")
        if self.current_batch is None:
            raise RuntimeError("no more batch data")
        # each VectorMetaData has the same numRecords value,
        # just use the first one to get the numRecord for this batch
        current_batch_row_size = 0
        if self.current_batch.get_vectors_count() != 0:
            current_batch_row_size = self.current_batch.get_batch_row_size()

        # the current batch is empty or already finished the batch, jump the batch
        if (
            self.current_batch.get_vectors_count() == 0
            or self.current_batch_row_index >= current_batch_row_size
        ):
            self.batch_index += 1
            if self.batch_index >= self.num_batches:
                raise RuntimeError("no more batch data")

            # reset currentBatchRowIndex to 0
            self.current_batch_row_index = 0
            self.current_batch = Batch(
                self.result_table.batch[self.batch_index],
                self.byte_order,
            )

        # resolve the current batch
        row = self._get_row_by_index(self.current_batch_row_index)
        self.current_batch_row_index += 1
        return row
