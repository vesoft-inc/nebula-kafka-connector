/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.RawRecord;
import com.vesoft.nebula.Value;
import com.vesoft.nebula.graph.ExecutionResponse;
import com.vesoft.nebula.graph.PlanDescription;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.Spliterator;
import java.util.function.Consumer;

public class ResultSet {
    private final ExecutionResponse response;
    private final List<String> columnNames = new ArrayList<>();
    private final String decodeType = "utf-8";

    public static class Record implements Iterable<ValueWrapper> {
        private final List<ValueWrapper> colValues = new ArrayList<>();
        private List<String> columnNames = new ArrayList<>();

        public Record(List<String> columnNames, RawRecord row, String decodeType) {
            if (columnNames == null) {
                return;
            }

            if (row == null || row.values == null) {
                return;
            }

            for (Value value : row.values) {
                this.colValues.add(new ValueWrapper(value, decodeType));
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
         * get the size of record
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

    public ResultSet(ExecutionResponse resp) {
        if (resp == null) {
            throw new RuntimeException("Input an null `ExecutionResponse' object");
        }
        this.response = resp;
        if (resp.executionOutcome != null
                && resp.executionOutcome.result != null
                && resp.executionOutcome.result.columnNames != null) {
            // space name's charset is 'utf-8'
            for (byte[] column : resp.executionOutcome.result.columnNames) {
                this.columnNames.add(new String(column));
            }
        }
    }

    /**
     * the execute result is succeeded
     *
     * @return boolean
     */
    public boolean isSucceeded() {
        return "SUCCESS".equals(new String(response.executionOutcome.gqlStatus.status));
    }

    /**
     * the result data is empty
     *
     * @return boolean
     */
    public boolean isEmpty() {
        return response.executionOutcome.result == null
                || response.executionOutcome.result.records.isEmpty();
    }

    /**
     * get errorCode of execute result
     *
     * @return int
     */
    public String getGqlStatus() {
        return new String(response.executionOutcome.gqlStatus.status);
    }


    /**
     * get latency of the query execute time
     *
     * @return int
     */
    public long getLatency() {
        return response.latencyInUs;
    }

    /**
     * get the PalnDesc
     *
     * @return PlanDescription
     */
    public PlanDescription getPlanDesc() {
        return response.executionOutcome.plan_desc;
    }

    /**
     * get keys of the dataset
     *
     * @return the list of String
     */
    public List<String> keys() {
        return columnNames;
    }

    /**
     * get column names of the dataset
     *
     * @return the
     */
    public List<String> getColumnNames() {
        return columnNames;
    }

    /**
     * get the size of rows
     *
     * @return int
     */
    public int rowsSize() {
        if (response.executionOutcome.result.records == null) {
            return 0;
        }
        return response.executionOutcome.result.records.size();
    }

    /**
     * get row values with the row index
     *
     * @param index the index of the rows
     * @return Record
     */
    public Record rowValues(int index) {
        if (response.executionOutcome.result == null) {
            throw new RuntimeException("Empty data");
        }
        List<RawRecord> records = response.executionOutcome.result.records;
        if (index >= records.size()) {
            throw new ArrayIndexOutOfBoundsException();
        }
        return new Record(columnNames, records.get(index), decodeType);
    }

    /**
     * get col values on the column key
     *
     * @param columnName the column name
     * @return the list of ValueWrapper
     */
    public List<ValueWrapper> colValues(String columnName) {
        if (response.executionOutcome.result == null) {
            throw new RuntimeException("Empty data");
        }
        int index = columnNames.indexOf(columnName);
        if (index < 0) {
            throw new ArrayIndexOutOfBoundsException();
        }
        List<ValueWrapper> values = new ArrayList<>();
        List<RawRecord> records = response.executionOutcome.result.records;
        for (int i = 0; i < records.size(); i++) {
            values.add(new ValueWrapper(records.get(i).values.get(index), decodeType));
        }
        return values;
    }

    /**
     * get all rows, the value is the origin value, the string is binary
     *
     * @return the list of Row
     */
    public List<RawRecord> getRows() {
        if (response.executionOutcome.result == null) {
            throw new RuntimeException("Empty data");
        }
        return response.executionOutcome.result.records;
    }

    @Override
    public String toString() {
        // When error, print the raw data directly
        if (!isSucceeded()) {
            return response.toString();
        }
        int i = 0;
        List<String> rowStrs = new ArrayList<>();
        while (i < rowsSize()) {
            List<String> valueStrs = new ArrayList<>();
            for (ValueWrapper value : rowValues(i)) {
                valueStrs.add(value.toString());
            }
            rowStrs.add(String.join(",", valueStrs));
            i++;
        }
        return String.format("ColumnName: %s, Rows: %s",
                columnNames.toString(), rowStrs.toString());
    }
}
