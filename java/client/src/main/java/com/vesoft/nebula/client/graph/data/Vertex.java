/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.proto.graph.Node;
import com.vesoft.nebula.proto.graph.Value;
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
     * get node type name
     *
     * @return String
     */
    public String getNodeType() {
        return node.getType();
    }


    /**
     *
     */
    public List<String> getNodeLabels() {
        return node.getLabelsList();
    }

    /**
     * get vid
     *
     * @return long id
     */
    public long getNodeId() {
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
        return getNodeId() == node.getNodeId();
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
            propStrs.add(propName + ": " + props.get(propName).toString());
        }
        return String.format("(%d:%s {%s})",
                getNodeId(), getNodeType(), String.join(", ", propStrs));
    }
}
