package com.vesoft.nebula.driver.graph.data;

import com.vesoft.nebula.proto.common.Node;
import com.vesoft.nebula.proto.common.Value;
import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

public class Vertex extends BaseDataObject {
    private final Node node;


    /**
     * Node is a wrapper around the Vertex type returned by nebula-graph
     *
     * @param node the vertex returned by nebula-graph
     */
    public Vertex(Node node) {
        if (node == null) {
            throw new RuntimeException("Input an null node object");
        }
        this.node = node;
    }

    /**
     * get graph
     *
     * @return String
     */
    public String getGraph() {
        return node.getGraph();
    }


    /**
     * get node type name
     *
     * @return String
     */
    public String getType() {
        return node.getType();
    }


    /**
     * get node label list
     *
     * @return list of label
     */
    public List<String> getLabels() {
        return node.getLabelsList();
    }

    /**
     * get vid
     *
     * @return long id
     */
    public long getId() {
        return node.getNodeId();
    }


    /**
     * get property names from the node
     *
     * @return the list of property names
     * @throws UnsupportedEncodingException decode error exception
     */
    public List<String> getColumnNames() throws UnsupportedEncodingException {
        List<String> keys = new ArrayList<>();
        for (String name : node.getPropertiesMap().keySet()) {
            keys.add(new String(name));
        }
        return keys;
    }

    /**
     * get all property values
     *
     * @return the List of property values
     */
    public List<ValueWrapper> getValues() {
        List<ValueWrapper> values = new ArrayList<>();
        for (Map.Entry<String, Value> kv : node.getPropertiesMap().entrySet()) {
            values.add(new ValueWrapper(kv.getValue()));
        }
        return values;
    }

    /**
     * get all properties for vertex
     *
     * @return HashMap, property name -> property value
     */
    public Map<String, ValueWrapper> getProperties() {
        Map<String, ValueWrapper> props = new HashMap<>();
        for (Map.Entry<String, Value> p : node.getPropertiesMap().entrySet()) {
            props.put(p.getKey(), new ValueWrapper(p.getValue()));
        }
        return props;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        Vertex node = (Vertex) o;
        return getId() == node.getId();
    }

    @Override
    public int hashCode() {
        return Objects.hash(node, getDecodeType());
    }

    @Override
    public String toString() {
        Map<String, ValueWrapper> props = getProperties();
        List<String> propStrs = new ArrayList<>();
        for (String propName : props.keySet()) {
            propStrs.add(propName + ":" + props.get(propName).toString());
        }
        return String.format("(%d@%s:%s{%s})",
                getId(),
                getType(),
                String.join("&", getLabels()),
                String.join(",", propStrs));
    }
}
