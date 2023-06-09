/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.schema;

public class Node extends Schema{
    private String nodeType;
    private String vidType;
    private String vidField;

    public String getNodeType() {
        return nodeType;
    }

    public Node setNodeType(String nodeType) {
        this.nodeType = nodeType;
        return this;
    }

    public String getVidType() {
        return vidType;
    }

    public Node setVidType(String vidType) {
        this.vidType = vidType;
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
        return String.format("(%s(%s) LABEL %s(%s))", nodeType, vidField, nodeType, props);
    }

    @Override
    public String toString() {
        return "Node{" +
                "nodeType='" + nodeType + '\'' +
                ", vidType='" + vidType + '\'' +
                ", vidField='" + vidField + '\'' +
                ", properties=" + super.properties +
                '}';
    }
}
