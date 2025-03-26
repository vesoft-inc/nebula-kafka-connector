
package com.vesoft.nebula.connector.sink;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class NebulaEdgeSchema {
    private String       edgeTypeName;
    private String       sourceNodeTypeName;
    private Map<String, String> sourcePkNameAndType = new HashMap<>();
    private String       targetNodeTypeName;
    private Map<String, String> targetPkNameAndType = new HashMap<>();

    // map of property name and property data type
    private Map<String, String> properties = new HashMap<>();

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


    public String getTargetNodeTypeName() {
        return targetNodeTypeName;
    }

    public void setTargetNodeTypeName(String targetNodeTypeName) {
        this.targetNodeTypeName = targetNodeTypeName;
    }


    public Map<String, String> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, String> properties) {
        this.properties = properties;
    }

    public Map<String, String> getSourcePkNameAndType() {
        return sourcePkNameAndType;
    }

    public void setSourcePkNameAndType(Map<String, String> sourcePkNameAndType) {
        this.sourcePkNameAndType = sourcePkNameAndType;
    }

    public Map<String, String> getTargetPkNameAndType() {
        return targetPkNameAndType;
    }

    public void setTargetPkNameAndType(Map<String, String> targetPkNameAndType) {
        this.targetPkNameAndType = targetPkNameAndType;
    }
}
