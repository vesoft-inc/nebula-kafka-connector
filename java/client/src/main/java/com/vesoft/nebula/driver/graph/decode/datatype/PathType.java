/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.driver.graph.decode.datatype;

import com.vesoft.nebula.driver.graph.decode.ColumnType;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class PathType extends DataType {
    private final List<DataType> dataTypes;

    private final List<NodeType> nodeTypes = new ArrayList<>();
    private final List<EdgeType> edgeTypes = new ArrayList<>();


    public PathType(List<DataType> dataTypes) {
        super(ColumnType.COLUMN_TYPE_PATH);
        this.dataTypes = dataTypes;
        for (DataType dataType : dataTypes) {
            if (dataType.getType() == ColumnType.COLUMN_TYPE_NODE) {
                nodeTypes.add((NodeType) dataType);
            }
            if (dataType.getType() == ColumnType.COLUMN_TYPE_EDGE) {
                edgeTypes.add((EdgeType) dataType);
            }
        }
    }

    public List<DataType> getDataTypes() {
        return dataTypes;
    }

    public Map<Integer, Map<String, DataType>> getNodeTypes() {
        Map<Integer, Map<String, DataType>> nodeTypesMap = new HashMap<>();
        for (NodeType nodeType : nodeTypes) {
            nodeTypesMap.putAll(nodeType.getNodeTypes());
        }
        return nodeTypesMap;
    }

    public Map<Integer, Map<String, DataType>> getEdgeTypes() {
        Map<Integer, Map<String, DataType>> edgeTypesMap = new HashMap<>();
        for (EdgeType edgeType : edgeTypes) {
            edgeTypesMap.putAll(edgeType.getEdgeTypes());
        }
        return edgeTypesMap;
    }
}
