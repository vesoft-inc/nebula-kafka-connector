/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.google.common.base.Charsets;
import com.google.protobuf.ByteString;
import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.proto.graph.ExecuteResponse;
import com.vesoft.nebula.proto.graph.ExtraInfoElement;
import com.vesoft.nebula.proto.graph.ExtraInfoKind;
import com.vesoft.nebula.proto.graph.Row;
import com.vesoft.nebula.proto.graph.Value;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.Spliterator;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class ResultSet {
    private final ExecuteResponse response;
    private final List<String> columnNames = new ArrayList<>();
    private final Charset charset = Charsets.UTF_8;

    private boolean isEmpty = false;

    private volatile AtomicInteger index = new AtomicInteger(0);

    private final int size;

    public static class Record implements Iterable<ValueWrapper> {
        private final List<ValueWrapper> colValues = new ArrayList<>();
        private List<String> columnNames = new ArrayList<>();

        public Record(List<String> columnNames, Row row) {
            if (columnNames == null) {
                return;
            }

            if (row == null || row.getValuesList().isEmpty()) {
                return;
            }

            for (Value value : row.getValuesList()) {
                this.colValues.add(new ValueWrapper(value));
            }

            this.columnNames = columnNames;
        }

        @Override
        public Iterator<ValueWrapper> iterator() {
            return this.colValues.iterator();
        }

        @Override
        public void forEach(Consumer<? super ValueWrapper> action) {
            this.colValues.forEach(action);
        }

        @Override
        public Spliterator<ValueWrapper> spliterator() {
            return this.colValues.spliterator();
        }


        @Override
        public String toString() {
            List<String> valueStr = new ArrayList<>();
            for (ValueWrapper v : colValues) {
                valueStr.add(v.toString());
            }
            return String.format("ColumnName: %s, Values: %s",
                    columnNames.toString(), valueStr.toString());
        }

        /**
         * get column value by index
         *
         * @param index the index of the rows
         * @return ValueWrapper
         */
        public ValueWrapper get(int index) {
            if (index >= columnNames.size()) {
                throw new IllegalArgumentException(
                        String.format("Cannot get field because the key '%d' out of range", index));
            }
            return this.colValues.get(index);
        }

        /**
         * get column value by column name
         *
         * @param columnName the columna name
         * @return ValueWrapper
         */
        public ValueWrapper get(String columnName) {
            int index = columnNames.indexOf(columnName);
            if (index == -1) {
                throw new IllegalArgumentException(
                        "Cannot get field because the columnName '"
                                + columnName + "' is not exists");
            }
            return this.colValues.get(index);
        }

        /**
         * get all values
         *
         * @return the list of ValueWrapper
         */
        public List<ValueWrapper> values() {
            return colValues;
        }

        /**
         * get the size of record columns
         *
         * @return int the size of columns
         */
        public int size() {
            return this.columnNames.size();
        }

        /**
         * if the column name exists
         *
         * @param columnName the column name
         * @return boolean
         */
        public boolean contains(String columnName) {
            return this.columnNames.contains(columnName);
        }

    }

    public ResultSet(ExecuteResponse resp) {
        if (resp == null || !resp.hasExecutionOutcome()) {
            throw new RuntimeException("got null object for server's response");
        }
        this.response = resp;
        if (!resp.getExecutionOutcome().hasResult()
                || resp.getExecutionOutcome().getResult().getRecordsList().isEmpty()) {
            isEmpty = true;
            size = 0;
        } else {
            size = resp.getExecutionOutcome().getResult().getRecordsCount();
        }
        for (ByteString column : resp.getExecutionOutcome().getResult().getColumnNamesList()) {
            this.columnNames.add(column.toString(charset));
        }
    }

    /**
     * the execute result is succeeded
     *
     * @return boolean
     */
    public boolean isSucceeded() {
        return ErrorCode.SUCCESSFUL_COMPLETION.code.equals(
                response.getExecutionOutcome().getGqlStatus().getCode().toString(charset));
    }

    /**
     * the result data is empty
     *
     * @return boolean
     */
    public boolean isEmpty() {
        return isEmpty;
    }

    /**
     * get error code of execute result
     * TODO return {@link ErrorCode}
     *
     * @return String
     */
    public ErrorCode getErrorCode() {
        return ErrorCode.find(
                response.getExecutionOutcome().getGqlStatus().getCode().toString(charset));
    }

    /**
     * get error message of execute result
     *
     * @return String
     */
    public String getErrorMessage() {
        return response.getExecutionOutcome().getGqlStatus().getMessage().toString(charset);
    }


    /**
     * get latency of the query execute time
     *
     * @return int
     */
    public long getLatency() {
        return response.getLatencyInUs();
    }

    /**
     * get the PlanDesc
     *
     * @return string
     */
    public String getPlanDesc() {
        return response.getExecutionOutcome().getPlanDesc().toString(charset);
    }

    /**
     * get column names of the dataset
     *
     * @return the list of result columns
     */
    public List<String> getColumnNames() {
        return columnNames;
    }

    /**
     * get the size of rows
     *
     * @return int
     */
    public int rowSize() {
        if (isEmpty) {
            return 0;
        }
        return size;
    }

    /**
     * if the ResultSet has row Record
     */
    public boolean hasNext() {
        if (isEmpty) {
            return false;
        }
        return index.get() < size;
    }

    /**
     * get the next row Record
     */
    public Record next() {
        if (!hasNext()) {
            throw new NoSuchElementException("no more row record data");
        }
        Row row = response.getExecutionOutcome().getResult().getRecords(index.getAndIncrement());
        return new Record(columnNames, row);
    }

    /**
     * get col values on the column key
     *
     * @param columnName the column name
     * @return the list of ValueWrapper
     */
    public List<ValueWrapper> colValues(String columnName) {
        if (isEmpty) {
            throw new RuntimeException("Empty data");
        }
        int index = columnNames.indexOf(columnName);
        if (index < 0) {
            throw new ArrayIndexOutOfBoundsException();
        }
        List<ValueWrapper> values = new ArrayList<>();
        List<Row> records = response.getExecutionOutcome().getResult().getRecordsList();
        for (int i = 0; i < records.size(); i++) {
            values.add(new ValueWrapper(records.get(i).getValues(index)));
        }
        return values;
    }


    /**
     * get extra info for result
     *
     * @return map for extra info
     */
    public ExtraInfo getExtraInfo() {
        if (!response.hasExecutionOutcome()) {
            return new ExtraInfo();
        }
        ExtraInfo extraInfo = new ExtraInfo();
        List<ExtraInfoElement> extraInfos = response.getExecutionOutcome().getExtraInfoList();
        for (ExtraInfoElement info : extraInfos) {
            switch (info.getKind()) {
                case kCursor: {
                    extraInfo.setCursor((new ValueWrapper(info.getValue())).asString());
                    continue;
                }
                case kAffectedNodes: {
                    extraInfo.setAffectedNodes(new ValueWrapper(info.getValue()).asLong());
                    continue;
                }
                case kAffectedForwardEdges: {
                    extraInfo.setAffectedForwardEdges(new ValueWrapper(info.getValue()).asLong());
                    continue;
                }
                case kAffectedReverseEdges: {
                    extraInfo.setAffectedReverseEdges(new ValueWrapper(info.getValue()).asLong());
                    continue;
                }
                default:
            }
        }
        return extraInfo;
    }

    @Override
    public String toString() {
        if (!isSucceeded()) {
            return response.getExecutionOutcome().getGqlStatus().getMessage().toString(charset);
        }
        List<String> rowStrs = new ArrayList<>();
        while (hasNext()) {
            List<String> valueStrs = new ArrayList<>();
            for (ValueWrapper value : next().values()) {
                valueStrs.add(value.toString());
            }
            String values = "[" + String.join(",", valueStrs) + "]";
            rowStrs.add(values);
        }
        return String.format("ColumnName: %s,\n Rows: %s",
                columnNames.toString(), rowStrs.toString());
    }
}
