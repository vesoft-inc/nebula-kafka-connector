/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.schema;

public class Edge extends Schema{
    private String edgeTypeName;
    private String srcNodeTypeName;
    private String srcField;
    private String dstNodeTypeName;
    private String dstField;

    public String getEdgeTypeName() {
        return edgeTypeName;
    }

    public Edge setEdgeTypeName(String edgeTypeName) {
        this.edgeTypeName = edgeTypeName;
        return this;
    }

    public String getSrcField() {
        return srcField;
    }

    public Edge setSrcField(String srcField) {
        this.srcField = srcField;
        return this;
    }

    public String getDstField() {
        return dstField;
    }

    public Edge setDstField(String dstField) {
        this.dstField = dstField;
        return this;
    }

    public String getSrcNodeTypeName() {
        return srcNodeTypeName;
    }

    public void setSrcNodeTypeName(String srcNodeTypeName) {
        this.srcNodeTypeName = srcNodeTypeName;
    }

    public String getDstNodeTypeName() {
        return dstNodeTypeName;
    }

    public void setDstNodeTypeName(String dstNodeTypeName) {
        this.dstNodeTypeName = dstNodeTypeName;
    }

    public String getSchemaString() {
        String props = super.getSchemaString();

        return String.format("(%s)-[%s LABEL %s {%s}]->(%s)",
                srcNodeTypeName, edgeTypeName, edgeTypeName, props, dstNodeTypeName);
    }

    @Override
    public String toString() {
        return "Edge{" +
                "edgeTypeName='" + edgeTypeName + '\'' +
                ", srcNodeTypeName='" + srcNodeTypeName + '\'' +
                ", srcField='" + srcField + '\'' +
                ", dstNodeTypeName='" + dstNodeTypeName + '\'' +
                ", dstField='" + dstField + '\'' +
                ", properties=" + super.properties +
                '}';
    }
}
