/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.driver.graph.decode.datatype;

import com.vesoft.nebula.driver.graph.decode.ColumnType;
import java.util.Map;

public class EdgeType extends DataType {
    private final Map<Integer, Map<String, DataType>> edgeTypes;

    public EdgeType(Map<Integer, Map<String, DataType>> edgeTypes) {
        super(ColumnType.COLUMN_TYPE_EDGE);
        this.edgeTypes = edgeTypes;
    }

    public Map<Integer, Map<String, DataType>> getEdgeTypes() {
        return edgeTypes;
    }
}
