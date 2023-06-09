/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.schema;

public class Edge extends Schema{
    private String edgeType;
    private String srcNodeType;
    private String srcField;
    private String dstNodeType;
    private String dstField;

    public String getEdgeType() {
        return edgeType;
    }

    public Edge setEdgeType(String edgeType) {
        this.edgeType = edgeType;
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

    public String getSrcNodeType() {
        return srcNodeType;
    }

    public void setSrcNodeType(String srcNodeType) {
        this.srcNodeType = srcNodeType;
    }

    public String getDstNodeType() {
        return dstNodeType;
    }

    public void setDstNodeType(String dstNodeType) {
        this.dstNodeType = dstNodeType;
    }

    public String getSchemaString() {
        String props = super.getSchemaString();

        return String.format("(%s)-[%s LABEL %s {%s}]->(%s)",
                srcNodeType, edgeType, edgeType, props, dstNodeType);
    }

    @Override
    public String toString() {
        return "Edge{" +
                "edgeType='" + edgeType + '\'' +
                ", srcNodeType='" + srcNodeType + '\'' +
                ", srcField='" + srcField + '\'' +
                ", dstNodeType='" + dstNodeType + '\'' +
                ", dstField='" + dstField + '\'' +
                ", properties=" + super.properties +
                '}';
    }
}
