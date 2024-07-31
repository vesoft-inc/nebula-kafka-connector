package com.vesoft.nebula.client.graph.scan;

import com.vesoft.nebula.client.graph.data.ValueWrapper;
import java.util.List;

public class TableRow {
    private final List<ValueWrapper> values;

    public TableRow(List<ValueWrapper> values) {
        this.values = values;
    }

    /**
     * number of elements in TableRow
     */
    public int size() {
        return values.size();
    }

    /**
     * check whether the value at position i is null
     */
    public boolean isNullAt(int i) {
        return values.get(i) == null;
    }

    public List<ValueWrapper> getValues() {
        return values;
    }
}
