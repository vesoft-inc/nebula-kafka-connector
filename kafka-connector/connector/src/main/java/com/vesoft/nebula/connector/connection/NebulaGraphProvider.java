/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.connection;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import com.vesoft.nebula.client.graph.net.NebulaClient;
import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import com.vesoft.nebula.connector.sink.NebulaEdgeSchema;
import com.vesoft.nebula.connector.sink.NebulaNodeSchema;
import java.io.Serializable;
import java.net.UnknownHostException;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * NebulaGraph graph client provider, responsible for sending request with graphd server.
 */
public class NebulaGraphProvider implements Serializable {
    private final String host;
    private final String user;
    private final String passwd;
    private NebulaClient client = null;

    public NebulaGraphProvider(NebulaSinkConnectConfig config) {
        this.host = config.graphServers;
        this.user = config.user;
        this.passwd = config.passwd;
        try {
            client = NebulaClient.builder(host, user, passwd)
                    .setConnectTimeoutMills(config.connectTimeout)
                    .setRequestTimeoutMills(config.requestTimeout)
                    .setMaxSessionSize(config.sinkPartition)
                    .setMinSessionSize(1)
                    .setRetryTimes(config.retryTimes)
                    .setIntervalTimeMills(config.intervalTimeMill)
                    .setReconnect(config.reconnect)
                    .setBlockWhenExhausted(true)
                    .setMaxWaitMills(Integer.MAX_VALUE)
                    .setStrictlyServerHealthy(true)
                    .build();
        } catch (IOErrorException e) {
            throw new RuntimeException("connect to NebulaGraph server failed, please check " +
                    "the connectivity between client and server.", e);
        } catch (UnknownHostException e) {
            throw new IllegalArgumentException(String.format("wrong graph server config %s",
                    config.graphServers), e);
        }
    }

    public ResultSet execute(String statement) throws IOErrorException, NoValidSessionException {
        return client.execute(statement);
    }

    public NebulaNodeSchema getNodeSchema(String graphName, String nodeType) throws IOErrorException, NoValidSessionException {
        NebulaNodeSchema nodeSchema = new NebulaNodeSchema();
        Map<String, String> schema = new HashMap<>();
        ResultSet result = getGraphDesc(graphName);
        List<ResultSet.Record> records = result.getRows();

        for (ResultSet.Record record : records) {
            if (record.get("Field").asString().equalsIgnoreCase(nodeType)) {
                String propertyString = record.get("Properties").asString();
                String[] proeprties =
                        propertyString.substring(1, propertyString.length() - 1).split(",");
                for (String prop : proeprties) {
                    String[] nameAndType = prop.trim().split(" ");
                    schema.put(nameAndType[0], nameAndType[1]);
                }
                nodeSchema.setNodeTypeName(nodeType);
                nodeSchema.setNodeIdType(schema.get("id"));
                nodeSchema.setNodeProperties(schema);
                return nodeSchema;
            }
        }
        throw new IllegalArgumentException("nodeType " + nodeType + " does not exist.");

    }

    public NebulaEdgeSchema getEdgeSchema(String graphName, String edgeType) throws IOErrorException, NoValidSessionException {
        NebulaEdgeSchema edgeSchema = new NebulaEdgeSchema();
        Map<String, String> edgeInfo = getEdgeInfo(graphName, edgeType);
        String sourceType = edgeInfo.get("sourceType");
        String targetType = edgeInfo.get("targetType");
        String sourceIdType = getNodeSchema(graphName, sourceType).getNodeIdType();
        String targetIdType = getNodeSchema(graphName, targetType).getNodeIdType();

        edgeSchema.setEdgeTypeName(edgeType);
        edgeSchema.setSourceNodeTypeName(sourceType);
        edgeSchema.setSourceNodeIdType(sourceIdType);
        edgeSchema.setTargetNodeTypeName(targetType);
        edgeSchema.setTargetNodeIdType(targetIdType);

        edgeInfo.remove("sourceType");
        edgeInfo.remove("targetType");

        edgeSchema.setProperties(edgeInfo);
        return edgeSchema;
    }


    private Map<String, String> getEdgeInfo(String graphName, String edgeType) throws IOErrorException, NoValidSessionException {
        Map<String, String> schema = new HashMap<>();
        ResultSet result = getGraphDesc(graphName);
        List<ResultSet.Record> records = result.getRows();
        String sourceNodeType;
        String targetNodeType;

        for (ResultSet.Record record : records) {
            if (record.get("Kind").asString().equals("Edge")) {
                String fullEdgeType = record.get("Field").asString();
                String regex = "\\((.*)\\)-\\[(.*)\\]->\\((.*)\\)";
                Pattern r = Pattern.compile(regex);
                Matcher m = r.matcher(fullEdgeType);
                if (m.find()) {
                    if (edgeType.equalsIgnoreCase(m.group(2))) {
                        sourceNodeType = m.group(1);
                        targetNodeType = m.group(3);
                        schema.put("sourceType", sourceNodeType);
                        schema.put("targetType", targetNodeType);
                        String propertiesString = record.get("Properties").asString();
                        String[] properties = propertiesString.substring(1,
                                propertiesString.length() - 1).split(",");
                        for (String prop : properties) {
                            String[] nameAndType = prop.trim().split(" ");
                            schema.put(nameAndType[0], nameAndType[1]);
                        }
                        return schema;
                    }
                }
            }
        }
        throw new IllegalArgumentException("edgeType " + edgeType + " does not exist.");
    }

    private ResultSet getGraphDesc(String graphName) throws IOErrorException,
            NoValidSessionException {
        ResultSet resultSet = client.execute("DESCRIBE GRAPH " + graphName);
        String graphType;
        if (resultSet.isSucceeded() && !resultSet.isEmpty()) {
            graphType = resultSet.getRows().get(0).values().get(1).asString();
        } else {
            throw new IllegalArgumentException("graphName " + graphName + " does not exist.");
        }

        String queryStatement = "DESCRIBE GRAPH TYPE " + graphType;
        resultSet = client.execute(queryStatement);
        if (!resultSet.isSucceeded()) {
            throw new RuntimeException("query error with " + queryStatement + " for " + resultSet.getGqlStatus());
        }
        return resultSet;
    }


    public void close() {
        client.close();
    }
}
