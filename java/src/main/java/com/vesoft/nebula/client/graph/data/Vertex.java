/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.Node;
import com.vesoft.nebula.Value;
import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

public class Vertex extends BaseDataObject {
    private final Node node;
    private final long vid;

    private String label;

    /**
     * Node is a wrapper around the Vertex type returned by nebula-graph
     *
     * @param node the vertex returned by nebula-graph
     */
    public Vertex(Node node) {
        if (node == null) {
            throw new RuntimeException("Input an null node object");
        }
        vid = node.nodeID;
        this.node = node;
    }

    /**
     * get vid from the node
     *
     * @return long id
     */
    public long getId() {
        return vid;
    }

    /**
     * Used to be compatible with older versions of interfaces
     *
     * @return the list of tag name
     */
    public String label() {
        return label;
    }


    public Map<String, ValueWrapper> properties() {
        Map<String, ValueWrapper> props = new HashMap<>();
        for (Map.Entry<byte[], Value> p : node.getProperties().entrySet()) {
            props.put(new String(p.getKey()), new ValueWrapper(p.getValue(), getDecodeType()));
        }
        return props;
    }

    /**
     * get property names from the node
     *
     * @return the list of property names
     * @throws UnsupportedEncodingException decode error exception
     */
    public List<String> keys() throws UnsupportedEncodingException {
        List<String> keys = new ArrayList<>();
        for (byte[] name : node.properties.keySet()) {
            keys.add(new String(name, getDecodeType()));
        }
        return keys;
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
        List<String> propStrs = new ArrayList<>();

        return String.format("(%d:%d {%s})", vid, node.nodeTypeID, String.join(" ", propStrs));
    }
}
