/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import java.util.ArrayList;
import java.util.List;

public class NebulaEdges {
    private NebulaEdgeSchema edgeSchema;
    private String           kafkaSrcField;
    private String           kafkaDstField;
    private List<String>     kafkaEdgePropertyNames;
    private List<String>     nebulaEdgePropertyNames;
    private List<NebulaEdge> edges;
    private List<String>     propertyNames = new ArrayList<>();

    public NebulaEdges(NebulaEdgeSchema edgeSchema,
                       String kafkaSrcField,
                       String kafkaDstField,
                       List<String> kafkaEdgePropertyNames,
                       List<String> nebulaEdgePropertyNames,
                       List<NebulaEdge> edges) {
        this.edgeSchema = edgeSchema;
        this.kafkaSrcField = kafkaSrcField;
        this.kafkaDstField = kafkaDstField;
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
                + "MATCH (src@%s) WHERE src.%s=r.%s "
                + "MATCH (dst@%s) WHERE dst.%s=r.%s \n"
                + "%s (src)-[@%s{%s}]->(dst)";

        return String.format(format,
                             getTableHeaders(),
                             getTableValues(),
                             graphName,
                             edgeSchema.getSourceNodeTypeName(),
                             edgeSchema.getSourceNodePkName(),
                             kafkaSrcField,
                             edgeSchema.getTargetNodeTypeName(),
                             edgeSchema.getTargetNodePkName(),
                             kafkaDstField,
                             insertModeString,
                             edgeSchema.getEdgeTypeName(),
                             getProperties());
    }


    public String getUpdateStatement(String graphName) {
        String format = "TABLE t{%s} = \n"
                + "%s \n"
                + "USE %s \n"
                + "FOR r IN t \n"
                + "MATCH (src@%s)-[e@%s]->(dst@%s) WHERE src.%s=r.%s AND dst.%s=r.%s \n"
                + "SET %s";
        return String.format(format,
                             getTableHeaders(),
                             getTableValues(),
                             graphName,
                             edgeSchema.getSourceNodeTypeName(),
                             edgeSchema.getEdgeTypeName(),
                             edgeSchema.getTargetNodeTypeName(),
                             edgeSchema.getSourceNodePkName(),
                             kafkaSrcField,
                             edgeSchema.getTargetNodePkName(),
                             kafkaDstField,
                             getUpdateProperties());
    }

    public String getDeleteStatement(String graphName) {
        String format = "TABLE t{%s} = \n"
                + "%s \n"
                + "USE %s \n"
                + "FOR r IN t \n"
                + "MATCH (src@%s)-[e@%s]->(dst@%s) WHERE src.%s=r.%s AND dst.%s=r.%s \n"
                + "DELETE e";
        return String.format(format,
                             getTableHeaders(),
                             getTableValues(),
                             graphName,
                             edgeSchema.getSourceNodeTypeName(),
                             edgeSchema.getEdgeTypeName(),
                             edgeSchema.getTargetNodeTypeName(),
                             edgeSchema.getSourceNodePkName(),
                             kafkaSrcField,
                             edgeSchema.getTargetNodePkName(),
                             kafkaDstField);
    }

    private String getTableHeaders() {
        List<String> headerNames = new ArrayList<>();
        headerNames.add(kafkaSrcField);
        headerNames.add(kafkaDstField);
        for (String kafkaFieldName : kafkaEdgePropertyNames) {
            if (kafkaFieldName.equals(kafkaSrcField) || kafkaFieldName.equals(kafkaDstField)) {
                continue;
            }
            headerNames.add(kafkaFieldName);
        }
        return String.join(",", headerNames);
    }

    private String getTableValues() {
        List<String> tableRows = new ArrayList<>();
        for (NebulaEdge edge : edges) {
            List<String>  rowValues      = new ArrayList<>();
            rowValues.add(edge.getSrcPk());
            rowValues.add(edge.getDstPk());
            for (String propName : nebulaEdgePropertyNames) {
                int index = nebulaEdgePropertyNames.indexOf(propName);
                if (kafkaEdgePropertyNames.get(index).equals(kafkaSrcField)
                        || kafkaEdgePropertyNames.get(index).equals(kafkaDstField)) {
                    continue;
                }
                rowValues.add(edge.getProperties().get(propName).toString());
            }
            StringBuilder rowValueString = new StringBuilder();
            rowValueString.append("(").append(String.join(",", rowValues)).append(")");
            tableRows.add(rowValueString.toString());
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
            propertyString
                    .append("e.")
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


}
