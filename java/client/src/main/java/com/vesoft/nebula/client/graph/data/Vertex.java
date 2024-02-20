/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.proto.Node;
import com.vesoft.nebula.proto.Value;
import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

public class Vertex extends BaseDataObject {
    private final Node node;
    private final long vid;

    private String label = "";
    private int typeId;

    /**
     * Node is a wrapper around the Vertex type returned by nebula-graph
     *
     * @param node the vertex returned by nebula-graph
     */
    public Vertex(Node node) {
        if (node == null) {
            throw new RuntimeException("Input an null node object");
        }
        vid = node.getNodeId();
        this.node = node;
        this.typeId = node.getNodeTypeId();
    }

    /**
     * get node type name, TODO core server should return the node type name
     *
     * @return String
     */
    public String getNodeType() {
        return label;
    }

    /**
     * get node type id
     *
     * @return int
     */
    public int getNodeTypeId() {
        return typeId;
    }

    /**
     * get vid
     *
     * @return long id
     */
    public long getId() {
        return vid;
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
        for (Map.Entry<String, Value> p : node.getProperties().entrySet()) {
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
        return vid == node.vid;
    }

    @Override
    public int hashCode() {
        return Objects.hash(node, vid, getDecodeType(), label);
    }

    @Override
    public String toString() {
        Map<String, ValueWrapper> props = getProperties();
        List<String> propStrs = new ArrayList<>();
        for (String propName : props.keySet()) {
            propStrs.add(propName + ": " + props.get(propName).toString());
        }
        return String.format("(%d:%d {%s})",
                vid, node.getNodeTypeId(), String.join(", ", propStrs));
    }
}
