
package com.vesoft.nebula.connector.sink;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class NebulaNodeSchema {
    private String       nodeTypeName;
    private List<String> pkNames = new ArrayList<>();

    // map of property name and property data type
    private Map<String, String> nodeProperties = new HashMap<>();

    public String getNodeTypeName() {
        return nodeTypeName;
    }

    public void setNodeTypeName(String nodeTypeName) {
        this.nodeTypeName = nodeTypeName;
    }


    public List<String> getPkNames() {
        return pkNames;
    }

    public void setPkNames(List<String> pkNames) {
        this.pkNames = pkNames;
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
                + "nodeTypeName='" + nodeTypeName + '\''
                + ", pkNames=" + pkNames
                + ", nodeProperties=" + nodeProperties
                + '}';
    }
}
