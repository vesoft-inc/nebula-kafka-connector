
package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.NODE_VALUES_TEMPLATE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.PROPERTY_TEMPLATE;

import com.vesoft.nebula.connector.exceptions.DataFormatException;
import com.vesoft.nebula.connector.util.NebulaUtils;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class NebulaNode {
    private Map<String, String> properties;

    public NebulaNode(Map<String, String> properties) {
        this.properties = properties;
    }



    public Map<String, String> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, String> properties) {
        this.properties = properties;
    }

    @Override
    public String toString() {
        return "NebulaNode{"
                + "properties=" + properties
                + '}';
    }
}
