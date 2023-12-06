/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

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
