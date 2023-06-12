/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.schema;

public class Node extends Schema{
    private String nodeTypeName;
    private String vidDataType;
    private String vidField;

    public String getNodeTypeName() {
        return nodeTypeName;
    }

    public Node setNodeTypeName(String nodeTypeName) {
        this.nodeTypeName = nodeTypeName;
        return this;
    }

    public String getVidDataType() {
        return vidDataType;
    }

    public Node setVidDataType(String vidDataType) {
        this.vidDataType = vidDataType;
        return this;
    }

    public String getVidField() {
        return vidField;
    }

    public Node setVidField(String vidField) {
        this.vidField = vidField;
        return this;
    }

    public String getSchemaString() {
        String props = super.getSchemaString();
        return String.format("(%s(%s) LABEL %s(%s))", nodeTypeName, vidField, nodeTypeName, props);
    }

    @Override
    public String toString() {
        return "Node{" +
                "nodeTypeName='" + nodeTypeName + '\'' +
                ", vidDataType='" + vidDataType + '\'' +
                ", vidField='" + vidField + '\'' +
                ", properties=" + super.properties +
                '}';
    }
}
