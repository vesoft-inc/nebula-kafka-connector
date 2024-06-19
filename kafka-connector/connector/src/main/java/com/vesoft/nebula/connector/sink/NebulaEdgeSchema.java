/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import java.util.Map;

public class NebulaEdgeSchema {
    private String edgeTypeName;
    private String sourceNodeTypeName;
    private String sourceNodePkName;
    private String sourceNodePkType;
    private String targetNodePkName;
    private String targetNodeTypeName;
    private String targetNodePkType;

    // map of property name and property data type
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

    public String getSourceNodePkName() {
        return sourceNodePkName;
    }

    public void setSourceNodePkName(String sourceNodePkName) {
        this.sourceNodePkName = sourceNodePkName;
    }

    public String getSourceNodePkType() {
        return sourceNodePkType;
    }

    public void setSourceNodePkType(String sourceNodePkType) {
        this.sourceNodePkType = sourceNodePkType;
    }

    public String getTargetNodeTypeName() {
        return targetNodeTypeName;
    }

    public void setTargetNodeTypeName(String targetNodeTypeName) {
        this.targetNodeTypeName = targetNodeTypeName;
    }

    public String getTargetNodePkName() {
        return targetNodePkName;
    }

    public void setTargetNodePkName(String targetNodePkName) {
        this.targetNodePkName = targetNodePkName;
    }

    public String getTargetNodePkType() {
        return targetNodePkType;
    }

    public void setTargetNodePkType(String targetNodePkType) {
        this.targetNodePkType = targetNodePkType;
    }

    public Map<String, String> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, String> properties) {
        this.properties = properties;
    }

    @Override
    public String toString() {
        return "NebulaEdgeSchema{"
                + "edgeTypeName='" + edgeTypeName + '\''
                + ", sourceNodeTypeName='" + sourceNodeTypeName + '\''
                + ", sourceNodePkName='" + sourceNodePkName + '\''
                + ", sourceNodePkType='" + sourceNodePkType + '\''
                + ", targetNodeTypeName='" + targetNodeTypeName + '\''
                + ", targetNodePkType='" + targetNodePkName + '\''
                + ", targetNodePkType='" + targetNodePkType + '\''
                + ", properties=" + properties
                + '}';
    }
}
