/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.proto.Edge;
import com.vesoft.nebula.proto.Value;
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
     * get edge name, TODO core server should return the edge name
     *
     * @return String
     */
    public String getEdgeName() {
        return String.valueOf(edge.getEdgeTypeId());
    }

    /**
     * get edge type id
     *
     * @return int
     */
    public int getEdgeId() {
        return edge.getEdgeTypeId();
    }

    /**
     * get the src id
     *
     * @return long
     */
    public long getSrcId() {
        return edge.getEdgeTypeId() > 0 ? edge.getSrcId() : edge.getDstId();
    }

    /**
     * get the dst id from the relationship
     *
     * @return long
     */
    public long getDstId() {
        return edge.getEdgeTypeId() > 0 ? edge.getDstId() : edge.getSrcId();
    }

    /**
     * get rank from the relationship
     *
     * @return long
     */
    public long getRank() {
        return edge.getRank();
    }

    /**
     * get all property name
     *
     * @return the List of property names
     */
    public List<String> getColumnNames() {
        List<String> propNames = new ArrayList<>();
        for (String name : edge.getPropertiesMap().keySet()) {
            propNames.add(name);
        }
        return propNames;
    }

    /**
     * get all property values
     *
     * @return the List of property values
     */
    public List<ValueWrapper> getPropertyValues() {
        List<ValueWrapper> values = new ArrayList<>();
        for (Map.Entry<String, Value> kv : edge.getPropertiesMap().entrySet()) {
            values.add(new ValueWrapper(kv.getValue()));
        }
        return values;
    }

    /**
     * get property names and values from the relationship
     *
     * @return HashMap, property name -> property value
     */
    public HashMap<String, ValueWrapper> getProperties() {
        HashMap<String, ValueWrapper> properties = new HashMap<>();
        for (String key : edge.getPropertiesMap().keySet()) {
            properties.put(key, new ValueWrapper(edge.getPropertiesMap().get(key)));
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
        return getRank() == that.getRank()
                && getSrcId() == that.getSrcId()
                && getDstId() == that.getDstId()
                && Objects.equals(getEdgeName(), that.getEdgeName());
    }

    @Override
    public int hashCode() {
        return Objects.hash(edge, getDecodeType());
    }

    @Override
    public String toString() {
        List<String> propStrs = new ArrayList<>();
        Map<String, ValueWrapper> props = getProperties();
        for (String key : props.keySet()) {
            propStrs.add(key + ": " + props.get(key).toString());
        }
        return String.format("(%d)-[:%s@%d{%s}]->(%d)",
                getSrcId(), getEdgeName(), getRank(), String.join(", ", propStrs),
                getDstId());
    }
}
