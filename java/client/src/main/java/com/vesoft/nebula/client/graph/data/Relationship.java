/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.proto.graph.Direction;
import com.vesoft.nebula.proto.graph.Edge;
import com.vesoft.nebula.proto.graph.Value;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

public class Relationship extends BaseDataObject {
    private final Edge edge;

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
     * get edge type name
     *
     * @return String
     */
    public String getEdgeType() {
        return edge.getType();
    }

    /**
     * get edge labels
     *
     * @return list of edge labels
     */
    public List<String> getEdgeLabels() {
        return edge.getLabelsList();
    }

    /**
     * get the src id
     *
     * @return long
     */
    public long getSrcId() {
        return edge.getSrcId();
    }

    /**
     * get the dst id from the relationship
     *
     * @return long
     */
    public long getDstId() {
        return edge.getDstId();
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
        return new ArrayList<>(edge.getPropertiesMap().keySet());
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
        if (edge.getDirection() == Direction.DIRECTED) {
            return getRank() == that.getRank()
                    && getSrcId() == that.getSrcId()
                    && getDstId() == that.getDstId()
                    && Objects.equals(getEdgeType(), that.getEdgeType());
        } else {
            return getRank() == that.getRank()
                    && ((getSrcId() == that.getSrcId() && getDstId() == that.getDstId())
                    || (getSrcId() == that.getDstId() && getDstId() == that.getSrcId()))
                    && Objects.equals(getEdgeType(), that.getEdgeType());
        }
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
        if (edge.getDirection() == Direction.DIRECTED) {
            return String.format("(%d)-[:%s{%s}]->(%d)",
                    getSrcId(), getEdgeType(), String.join(", ", propStrs), getDstId());
        } else {
            return String.format("(%d)-[:%s{%s}]-(%d)",
                    getSrcId(), getEdgeType(), String.join(", ", propStrs), getDstId());
        }

    }
}
