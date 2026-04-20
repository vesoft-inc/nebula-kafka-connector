/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import java.util.ArrayList;
import java.util.List;

public class NebulaEdges {
    private static final String SRC_NODE_ALIAS = "n_src";
    private static final String DST_NODE_ALIAS = "n_dst";
    private static final String EDGE_ALIAS     = "n_e";

    private NebulaEdgeSchema edgeSchema;
    private List<String>     kafkaSrcFields;
    private List<String>     kafkaDstFields;
    private List<String>     kafkaEdgePropertyNames;
    private List<String>     nebulaEdgePropertyNames;
    private List<NebulaEdge> edges;
    private List<String>     propertyNames = new ArrayList<>();

    public NebulaEdges(NebulaEdgeSchema edgeSchema,
                       List<String> kafkaSrcFields,
                       List<String> kafkaDstFields,
                       List<String> kafkaEdgePropertyNames,
                       List<String> nebulaEdgePropertyNames,
                       List<NebulaEdge> edges) {
        this.edgeSchema = edgeSchema;
        this.kafkaSrcFields = kafkaSrcFields;
        this.kafkaDstFields = kafkaDstFields;
        this.kafkaEdgePropertyNames = kafkaEdgePropertyNames;
        this.nebulaEdgePropertyNames = nebulaEdgePropertyNames;
        this.edges = edges;
        if (!edges.isEmpty()) {
            propertyNames.addAll(edges.get(0).getProperties().keySet());
        }
    }

    public String getInsertStatement(String graphName,
                                     NebulaSinkConnectConfig.InsertMode insertMode) {
        if (edges.isEmpty()) {
            return null;
        }
        String insertModeString;
        switch (insertMode) {
            case INSERT:
                insertModeString = "INSERT";
                break;
            case INSERTIGNORE:
                insertModeString = "INSERT OR IGNORE";
                break;
            case INSERTREPLACE:
                insertModeString = "INSERT OR REPLACE";
                break;
            default:
                throw new IllegalArgumentException("insert mode is illegal for insert:"
                                                           + insertMode);
        }
        String format = "TABLE t{%s} = \n"
                + "%s \n"
                + "USE %s \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (%s@%s) WHERE %s "
                + "OPTIONAL MATCH (%s@%s) WHERE %s \n"
                + "%s (%s)-[@%s{%s}]->(%s)";

        return String.format(format,
                             getTableHeaders(),
                             getTableValues(),
                             graphName,
                             SRC_NODE_ALIAS,
                             edgeSchema.getSourceNodeTypeName(),
                             getSrcPkFilter(),
                             DST_NODE_ALIAS,
                             edgeSchema.getTargetNodeTypeName(),
                             getDstPkFilter(),
                             insertModeString,
                             SRC_NODE_ALIAS,
                             edgeSchema.getEdgeTypeName(),
                             getProperties(),
                             DST_NODE_ALIAS);
    }


    public String getUpdateStatement(String graphName) {
        String format = "TABLE t{%s} = \n"
                + "%s \n"
                + "USE %s \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (%s@%s)-[%s@%s]->(%s@%s) "
                + "WHERE %s AND %s%s \n"
                + "SET %s";
        return String.format(format,
                             getTableHeaders(),
                             getTableValues(),
                             graphName,
                             SRC_NODE_ALIAS,
                             edgeSchema.getSourceNodeTypeName(),
                             EDGE_ALIAS,
                             edgeSchema.getEdgeTypeName(),
                             DST_NODE_ALIAS,
                             edgeSchema.getTargetNodeTypeName(),
                             getSrcPkFilter(),
                             getDstPkFilter(),
                             getMultiEdgeKeysFilter(),
                             getUpdateProperties());
    }

    public String getDeleteStatement(String graphName) {
        String format = "TABLE t{%s} = \n"
                + "%s \n"
                + "USE %s \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (%s@%s)-[%s@%s]->(%s@%s) "
                + "WHERE %s AND %s%s \n"
                + "DELETE %s";
        return String.format(format,
                             getDeleteTableHeaders(),
                             getDeleteTableValues(),
                             graphName,
                             SRC_NODE_ALIAS,
                             edgeSchema.getSourceNodeTypeName(),
                             EDGE_ALIAS,
                             edgeSchema.getEdgeTypeName(),
                             DST_NODE_ALIAS,
                             edgeSchema.getTargetNodeTypeName(),
                             getSrcPkFilter(),
                             getDstPkFilter(),
                             getMultiEdgeKeysFilter(),
                             EDGE_ALIAS);
    }

    private String getTableHeaders() {
        List<String> headerNames = new ArrayList<>();
        for (String pk : edgeSchema.getSourcePkNameAndType().keySet()) {
            headerNames.add("src_" + pk);
        }
        for (String pk : edgeSchema.getTargetPkNameAndType().keySet()) {
            headerNames.add("dst_" + pk);
        }
        headerNames.addAll(kafkaEdgePropertyNames);
        return String.join(",", headerNames);
    }

    private String getTableValues() {
        List<String> tableRows = new ArrayList<>();
        for (NebulaEdge edge : edges) {
            List<String> rowValues = new ArrayList<>();
            for (String pk : edgeSchema.getSourcePkNameAndType().keySet()) {
                rowValues.add(edge.getSrcPks().get(pk));
            }
            for (String pk : edgeSchema.getTargetPkNameAndType().keySet()) {
                rowValues.add(edge.getDstPks().get(pk));
            }
            for (String propName : nebulaEdgePropertyNames) {
                rowValues.add(edge.getProperties().get(propName));
            }
            tableRows.add("(" + String.join(",", rowValues) + ")");
        }
        return String.join(",", tableRows);
    }

    private String getDeleteTableHeaders() {
        List<String> headerNames = new ArrayList<>();
        for (String pk : edgeSchema.getSourcePkNameAndType().keySet()) {
            headerNames.add("src_" + pk);
        }
        for (String pk : edgeSchema.getTargetPkNameAndType().keySet()) {
            headerNames.add("dst_" + pk);
        }
        for (String key : edgeSchema.getMultipleEdgeKeys()) {
            headerNames.add(getKafkaFieldByNebulaField(key));
        }
        return String.join(",", headerNames);
    }

    private String getDeleteTableValues() {
        List<String> tableRows = new ArrayList<>();
        for (NebulaEdge edge : edges) {
            List<String> rowValues = new ArrayList<>();
            for (String pk : edgeSchema.getSourcePkNameAndType().keySet()) {
                rowValues.add(edge.getSrcPks().get(pk));
            }
            for (String pk : edgeSchema.getTargetPkNameAndType().keySet()) {
                rowValues.add(edge.getDstPks().get(pk));
            }
            for (String key : edgeSchema.getMultipleEdgeKeys()) {
                rowValues.add(edge.getProperties().get(key));
            }
            tableRows.add("(" + String.join(",", rowValues) + ")");
        }
        return String.join(",", tableRows);
    }

    private String getProperties() {
        StringBuilder propertyString = new StringBuilder();
        for (int index = 0; index < nebulaEdgePropertyNames.size(); index++) {
            propertyString
                    .append(nebulaEdgePropertyNames.get(index))
                    .append(":r.")
                    .append(kafkaEdgePropertyNames.get(index))
                    .append(",");
        }
        if (propertyString.length() > 0) {
            propertyString.deleteCharAt(propertyString.length() - 1);
        }
        return propertyString.toString();
    }

    private String getUpdateProperties() {
        StringBuilder propertyString = new StringBuilder();
        for (int index = 0; index < nebulaEdgePropertyNames.size(); index++) {
            if (edgeSchema.getMultipleEdgeKeys().contains(nebulaEdgePropertyNames.get(index))) {
                continue;
            }
            propertyString
                    .append(EDGE_ALIAS)
                    .append(".")
                    .append(nebulaEdgePropertyNames.get(index))
                    .append("=r.")
                    .append(kafkaEdgePropertyNames.get(index))
                    .append(",");
        }
        if (propertyString.length() > 0) {
            propertyString.deleteCharAt(propertyString.length() - 1);
        }
        return propertyString.toString();
    }

    private String getMultiEdgeKeysFilter() {
        if (edgeSchema.getMultipleEdgeKeys() == null
                || edgeSchema.getMultipleEdgeKeys().isEmpty()) {
            return "";
        }
        StringBuilder filter = new StringBuilder();
        for (String key : edgeSchema.getMultipleEdgeKeys()) {
            String kafkaField = getKafkaFieldByNebulaField(key);
            filter.append(" AND ")
                    .append(EDGE_ALIAS)
                    .append(".")
                    .append(key)
                    .append("=r.")
                    .append(kafkaField);
        }
        return filter.toString();
    }

    private String getKafkaFieldByNebulaField(String nebulaField) {
        int index = nebulaEdgePropertyNames.indexOf(nebulaField);
        if (index < 0) {
            throw new IllegalArgumentException(
                    "edge property " + nebulaField + " is not configured");
        }
        return kafkaEdgePropertyNames.get(index);
    }

    private String getSrcPkFilter() {
        StringBuilder pkFilterString = new StringBuilder();
        for (String srcPk : edgeSchema.getSourcePkNameAndType().keySet()) {
            if (pkFilterString.length() > 0) {
                pkFilterString.append(" AND ");
            }
            pkFilterString.append(SRC_NODE_ALIAS)
                    .append(".")
                    .append(srcPk)
                    .append("=r.src_")
                    .append(srcPk);
        }
        return pkFilterString.toString();
    }


    private String getDstPkFilter() {
        StringBuilder pkFilterString = new StringBuilder();
        for (String dstPk : edgeSchema.getTargetPkNameAndType().keySet()) {
            if (pkFilterString.length() > 0) {
                pkFilterString.append(" AND ");
            }
            pkFilterString.append(DST_NODE_ALIAS)
                    .append(".")
                    .append(dstPk)
                    .append("=r.dst_")
                    .append(dstPk);
        }
        return pkFilterString.toString();
    }
}
