/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import java.util.Map;

public class NebulaEdgeSchema {
    private String edgeTypeName;
    private String sourceNodeTypeName;
    private String sourceNodeIdType;
    private String targetNodeTypeName;
    private String targetNodeIdType;
    private Map<String,String> properties;

    public String getEdgeTypeName() {
        return edgeTypeName;
    }

    public void setEdgeTypeName(String edgeTypeName) {
        this.edgeTypeName = edgeTypeName;
    }

    public String getSourceNodeTypeName() {
        return sourceNodeTypeName;
    }

    public void setSourceNodeTypeName(String sourceNodeTypeName) {
        this.sourceNodeTypeName = sourceNodeTypeName;
    }

    public String getSourceNodeIdType() {
        return sourceNodeIdType;
    }

    public void setSourceNodeIdType(String sourceNodeIdType) {
        this.sourceNodeIdType = sourceNodeIdType;
    }

    public String getTargetNodeTypeName() {
        return targetNodeTypeName;
    }

    public void setTargetNodeTypeName(String targetNodeTypeName) {
        this.targetNodeTypeName = targetNodeTypeName;
    }

    public String getTargetNodeIdType() {
        return targetNodeIdType;
    }

    public void setTargetNodeIdType(String targetNodeIdType) {
        this.targetNodeIdType = targetNodeIdType;
    }

    public Map<String, String> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, String> properties) {
        this.properties = properties;
    }

    @Override
    public String toString() {
        return "NebulaEdgeSchema{" +
                "edgeTypeName='" + edgeTypeName + '\'' +
                ", sourceNodeTypeName='" + sourceNodeTypeName + '\'' +
                ", sourceNodeIdType='" + sourceNodeIdType + '\'' +
                ", targetNodeTypeName='" + targetNodeTypeName + '\'' +
                ", targetNodeIdType='" + targetNodeIdType + '\'' +
                ", properties=" + properties +
                '}';
    }
}
