/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import java.util.ArrayList;
import java.util.List;

public class NebulaNodes {
    private static final String NODE_ALIAS = "n_v";
    private NebulaNodeSchema    nodeSchema;
    private List<NebulaNode> nodes;

    private List<String> propertyNames = new ArrayList<>();

    public NebulaNodes(NebulaNodeSchema nodeSchema, List<NebulaNode> nodes) {
        this.nodeSchema = nodeSchema;
        this.nodes = nodes;
        if (!nodes.isEmpty()) {
            propertyNames.addAll(nodes.get(0).getProperties().keySet());
        }
    }

    public String getInsertStatement(String graphName,
                                     NebulaSinkConnectConfig.InsertMode insertMode) {
        if (nodes.isEmpty()) {
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
                + "%s (@%s{%s})";

        return String.format(format,
                             getTableHeaders(),
                             getTableValues(),
                             graphName,
                             insertModeString,
                             nodeSchema.getNodeTypeName(),
                             getProperties());
    }


    public String getUpdateStatement(String graphName) {
        if (nodes.isEmpty()) {
            return null;
        }
        String format = "TABLE t{%s} = \n"
                + "%s \n"
                + "USE %s \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (%s@%s) WHERE %s \n"
                + "SET %s";
        return String.format(format,
                             getTableHeaders(),
                             getTableValues(),
                             graphName,
                             NODE_ALIAS,
                             nodeSchema.getNodeTypeName(),
                             getPkFilters(),
                             getUpdateProperties());
    }

    public String getDeleteStatement(String graphName,
                                     NebulaSinkConnectConfig.InsertMode deleteMode) {
        if (nodes.isEmpty()) {
            return null;
        }
        String deleteModeString;
        switch (deleteMode) {
            case DELETE:
                deleteModeString = "DELETE";
                break;
            case DETACHDELETE:
                deleteModeString = "DETACH DELETE";
                break;
            default:
                throw new IllegalArgumentException("insert mode is illegal for delete:"
                                                           + deleteMode);
        }
        String format = "TABLE t{%s} = \n"
                + "%s \n"
                + "USE %s \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (%s@%s) WHERE %s \n"
                + "%s %s";
        return String.format(format,
                             getTableHeaders(),
                             getTableValues(),
                             graphName,
                             NODE_ALIAS,
                             nodeSchema.getNodeTypeName(),
                             getPkFilters(),
                             deleteModeString,
                             NODE_ALIAS);
    }

    private String getTableHeaders() {
        return String.join(",", propertyNames);
    }

    private String getTableValues() {
        List<String> tableRows = new ArrayList<>();
        for (NebulaNode node : nodes) {
            List<String>  rowValues      = new ArrayList<>();
            StringBuilder rowValueString = new StringBuilder();
            for (String propName : propertyNames) {
                rowValues.add(node.getProperties().get(propName));
            }
            rowValueString.append("(").append(String.join(",", rowValues)).append(")");
            tableRows.add(rowValueString.toString());
        }
        return String.join(",", tableRows);
    }

    private String getProperties() {
        StringBuilder propertyString = new StringBuilder();
        for (String propertyName : propertyNames) {
            propertyString.append(propertyName).append(":r.").append(propertyName).append(",");
        }
        if (propertyString.length() > 0) {
            propertyString.deleteCharAt(propertyString.length() - 1);
        }
        return propertyString.toString();
    }

    private String getUpdateProperties() {
        StringBuilder propertyString = new StringBuilder();
        for (String propertyName : propertyNames) {
            if (nodeSchema.getPkNames().contains(propertyName)) {
                continue;
            }
            propertyString
                    .append(NODE_ALIAS)
                    .append(".")
                    .append(propertyName)
                    .append("=r.")
                    .append(propertyName)
                    .append(",");
        }
        if (propertyString.length() > 0) {
            propertyString.deleteCharAt(propertyString.length() - 1);
        }
        return propertyString.toString();
    }

    private String getPkFilters() {
        StringBuilder pkFilterString = new StringBuilder();
        for (String pk : nodeSchema.getPkNames()) {
            pkFilterString
                    .append(NODE_ALIAS)
                    .append(".")
                    .append(pk)
                    .append("=r.")
                    .append(pk)
                    .append(",");
        }
        if (pkFilterString.length() > 0) {
            pkFilterString.deleteCharAt(pkFilterString.length() - 1);
        }
        return pkFilterString.toString();
    }
}
