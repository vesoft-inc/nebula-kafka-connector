/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.Edge;
import com.vesoft.nebula.Value;
import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

public class Relationship extends BaseDataObject {
    private final Edge edge;
    private String label;

    /**
     * Relationship is a wrapper around the Edge type returned by nebula-graph
     *
     * @param edge the Edge type returned by nebula-graph
     */
    public Relationship(Edge edge) {
        if (edge == null) {
            throw new RuntimeException("Input an null edge object");
        }
        this.edge = edge;
    }

    /**
     * get the src id from the relationship
     *
     * @return long
     */
    public long srcId() {
        return edge.edgeTypeID > 0 ? edge.srcID : edge.dstID;
    }

    /**
     * get the dst id from the relationship
     *
     * @return long
     */
    public long dstId() {
        return edge.edgeTypeID > 0 ? edge.dstID : edge.srcID;
    }

    /**
     * get edge name from the relationship
     *
     * @return String
     */
    // todo 转换成edge name
    public String edgeName() {
        return label;
    }

    /**
     * get ranking from the relationship
     *
     * @return long
     */
    public long ranking() {
        return edge.rank;
    }

    /**
     * get all property name from the relationship
     *
     * @return the List of String
     * @throws UnsupportedEncodingException if decode binary failed
     */
    public List<String> keys() throws UnsupportedEncodingException {
        List<String> propNames = new ArrayList<>();
        for (byte[] name : edge.properties.keySet()) {
            propNames.add(new String(name, getDecodeType()));
        }
        return propNames;
    }

    /**
     * get property values from the relationship
     *
     * @return the List of ValueWrapper
     */
    public List<ValueWrapper> values() {
        List<ValueWrapper> propVals = new ArrayList<>();
        for (Value val : edge.properties.values()) {
            propVals.add(new ValueWrapper(val, getDecodeType()));
        }
        return propVals;
    }

    /**
     * get property names and values from the relationship
     *
     * @return the HashMap, key is String, value is ValueWrapper>
     */
    public HashMap<String, ValueWrapper> properties() throws UnsupportedEncodingException {
        HashMap<String, ValueWrapper> properties = new HashMap<>();
        for (byte[] key : edge.properties.keySet()) {
            properties.put(new String(key, getDecodeType()),
                    new ValueWrapper(edge.properties.get(key), getDecodeType()));
        }
        return properties;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        Relationship that = (Relationship) o;
        return ranking() == that.ranking()
                && srcId() == that.srcId()
                && dstId() == that.dstId()
                && Objects.equals(edgeName(), that.edgeName());
    }

    @Override
    public int hashCode() {
        return Objects.hash(edge, getDecodeType());
    }

    @Override
    public String toString() {
        try {
            List<String> propStrs = new ArrayList<>();
            Map<String, ValueWrapper> props = properties();
            for (String key : props.keySet()) {
                propStrs.add(key + ": " + props.get(key).toString());
            }
            // TODO change the edgeTypeID to edge name
            return String.format("(%d)-[:%d@%d{%s}]->(%d)",
                    srcId(), edge.edgeTypeID, ranking(), String.join(", ", propStrs), dstId());
        } catch (UnsupportedEncodingException e) {
            return e.getMessage();
        }
    }
}
