/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.driver.graph.decode.datatype;

import static com.vesoft.nebula.driver.graph.decode.ColumnType.COLUMN_TYPE_NODE;

import java.util.Map;

public class NodeType extends DataType {
    // nodeTypeId -> (prop name -> prop data type)
    private final Map<Integer, Map<String, DataType>> nodeTypes;

    public NodeType(Map<Integer, Map<String, DataType>> nodeTypes) {
        super(COLUMN_TYPE_NODE);
        this.nodeTypes = nodeTypes;
    }

    public Map<Integer, Map<String, DataType>> getNodeTypes() {
        return nodeTypes;
    }
}
