/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.proto.graph.Path;
import com.vesoft.nebula.proto.graph.Value;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class NPath {
    private Path path;
    private String decodeType = "utf-8";

    private List<Vertex> nodes = new ArrayList<>();
    private List<Relationship> relationships = new ArrayList<>();
    private List<ValueWrapper> values = new ArrayList<>();

    public NPath(Path path) {
        if (path == null) {
            return;
        }
        this.path = path;
        for (Value value : path.getValuesList()) {
            values.add(new ValueWrapper(value));
            if (new ValueWrapper(value).isNode()) {
                nodes.add(new Vertex(value.getNodeValue()));
            }
            if (new ValueWrapper(value).isEdge()) {
                relationships.add(new Relationship(value.getEdgeValue()));
            }
        }
    }

    /**
     * Create a list over the nodes in this path, nodes will appear in the same order as they appear
     * in the path.
     *
     * @return a List of all nodes in this path
     */
    public List<Vertex> nodes() {
        return nodes;
    }


    /**
     * Create a list over the relationships in this path. The relationships will appear
     * in the same order as they appear in the path.
     *
     * @return a List of all relationships in this path
     */
    public List<Relationship> relationships() {
        return relationships;
    }

    /**
     * Create a list over the nodes and relationships in this path. The value will appear
     * in the same order as they appear in the path. The first value will be Node type, then the
     * next one will be RelationShip type, and next one will be Node type.
     *
     * @return a List of all values in this path
     */
    public List<ValueWrapper> values() {
        return values;
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj) {
            return true;
        }
        if (obj == null || getClass() != obj.getClass()) {
            return false;
        }
        NPath that = (NPath) obj;
        return values.equals(that.values);
    }

    @Override
    public int hashCode() {
        return values.hashCode();
    }

    public String toString() {
        List<String> edgeStrs = new ArrayList<>();
        for (int i = 0; i < relationships.size(); i++) {
            Relationship relationship = relationships.get(i);

            List<String> edgePropStrs = new ArrayList<>();
            Map<String, ValueWrapper> props = relationship.getProperties();
            for (String key : props.keySet()) {
                edgePropStrs.add(key + ":" + props.get(key).toString());
            }

            Vertex prefixNode = nodes.get(i);
            List<String> prefixNodePropStrs = new ArrayList<>();
            Map<String, ValueWrapper> prefixNodeProps = prefixNode.getProperties();
            for (String key : prefixNodeProps.keySet()) {
                prefixNodePropStrs.add(key + ":" + prefixNodeProps.get(key).toString());
            }

            Vertex suffixNode = nodes.get(i + 1);
            List<String> suffixNodePropStrs = new ArrayList<>();
            Map<String, ValueWrapper> suffixNodeProps = suffixNode.getProperties();
            for (String key : suffixNodeProps.keySet()) {
                suffixNodePropStrs.add(key + ":" + suffixNodeProps.get(key).toString());
            }

            String template;
            if (i == 0) {
                template = "(%d@%s:%s{%s})~[%d@%s:%s{%s}]~(%d@%s:%s{%s})";
                if (relationship.isDirected() && relationship.getSrcId() == prefixNode.getId()) {
                    template = "(%d@%s:%s{%s})-[%d@%s:%s{%s}]->(%d@%s:%s{%s})";
                }
                if (relationship.isDirected() && relationship.getSrcId() != prefixNode.getId()) {
                    template = "(%d@%s:%s{%s})<-[%d@%s:%s{%s}]-(%d@%s:%s{%s})";
                }

                edgeStrs.add(String.format(template,
                        prefixNode.getId(),
                        prefixNode.getType(),
                        String.join("&", prefixNode.getLabels()),
                        String.join(",", prefixNodePropStrs),
                        relationship.getRank(),
                        relationship.getType(),
                        String.join("&", relationship.getLabels()),
                        String.join(",", edgePropStrs),
                        suffixNode.getId(),
                        suffixNode.getType(),
                        String.join("&", suffixNode.getLabels()),
                        String.join(",", suffixNodePropStrs)));
            } else {
                template = "~[%d@%s:%s{%s}]~(%d@%s:%s{%s})";
                if (relationship.isDirected() && relationship.getSrcId() == prefixNode.getId()) {
                    template = "-[%d@%s:%s{%s}]->(%d@%s:%s{%s})";
                }
                if (relationship.isDirected() && relationship.getSrcId() != prefixNode.getId()) {
                    template = "<-[%d@%s:%s{%s}]-(%d@%s:%s{%s})";
                }
                edgeStrs.add(String.format(template,
                        relationship.getRank(),
                        relationship.getType(),
                        String.join("&", relationship.getLabels()),
                        String.join(",", edgePropStrs),
                        suffixNode.getId(),
                        suffixNode.getType(),
                        String.join("&", suffixNode.getLabels()),
                        String.join(",", suffixNodePropStrs)));
            }
        }
        return String.join("", edgeStrs);
    }
}
