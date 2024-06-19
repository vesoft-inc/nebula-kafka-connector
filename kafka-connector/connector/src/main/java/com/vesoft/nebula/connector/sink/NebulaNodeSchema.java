/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import java.util.Map;

public class NebulaNodeSchema {
    private String nodeTypeName;
    private String nodePkName;
    private String nodePkType;

    // map of property name and property data type
    private Map<String, String> nodeProperties;

    public String getNodeTypeName() {
        return nodeTypeName;
    }

    public void setNodeTypeName(String nodeTypeName) {
        this.nodeTypeName = nodeTypeName;
    }

    public String getNodePkName() {
        return nodePkName;
    }

    public void setNodePkName(String nodePkName) {
        this.nodePkName = nodePkName;
    }

    public String getNodePkType() {
        return nodePkType;
    }

    public void setNodePkType(String nodePkType) {
        this.nodePkType = nodePkType;
    }

    public Map<String, String> getNodeProperties() {
        return nodeProperties;
    }

    public void setNodeProperties(Map<String, String> nodeProperties) {
        this.nodeProperties = nodeProperties;
    }

    @Override
    public String toString() {
        return "NebulaNodeSchema{"
                + "nodeTypeName='" + nodeTypeName
                + ", nodePkName='" + nodePkName
                + ", nodePkType='" + nodePkType
                + ", nodeProperties=" + nodeProperties
                + '}';
    }
}
