
package com.vesoft.nebula.connector.connection;

import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import com.vesoft.nebula.connector.config.NebulaSourceConnectConfig;
import com.vesoft.nebula.connector.sink.NebulaEdgeSchema;
import com.vesoft.nebula.connector.sink.NebulaNodeSchema;
import com.vesoft.nebula.driver.graph.data.ResultSet;
import com.vesoft.nebula.driver.graph.exception.AuthFailedException;
import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import com.vesoft.nebula.driver.graph.exception.NoValidSessionException;
import com.vesoft.nebula.driver.graph.net.NebulaClient;
import com.vesoft.nebula.driver.graph.scan.ScanEdgeResultIterator;
import com.vesoft.nebula.driver.graph.scan.ScanNodeResultIterator;
import java.io.Serializable;
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
    private final Map<String, Object> authOptions;
    private NebulaClient client = null;

    public NebulaGraphProvider(NebulaSinkConnectConfig config) {
        this.host = config.graphServers;
        this.user = config.user;
        this.authOptions = config.authOptions;
        try {
            client = NebulaClient.builder(host, user)
                    .withAuthOptions(authOptions)
                    .withRequestTimeoutMills(config.requestTimeout)
                    .build();
        } catch (AuthFailedException e) {
            throw new RuntimeException("auth failed, please check your user and passwd");
        } catch (IOErrorException e) {
            throw new RuntimeException("connect to NebulaGraph server failed, please check "
                    + "the connectivity between client and server.", e);
        }
    }

    public NebulaGraphProvider(NebulaSourceConnectConfig config) {
        this.host = config.graphServers;
        this.user = config.user;
        this.authOptions = config.authOptions;
        try {
            client = NebulaClient.builder(host, user)
                    .withAuthOptions(authOptions)
                    .withRequestTimeoutMills(config.requestTimeout)
                    .build();
        } catch (AuthFailedException e) {
            throw new RuntimeException("auth failed, please check your user and passwd");
        } catch (IOErrorException e) {
            throw new RuntimeException("connect to NebulaGraph server failed, please check "
                    + "the connectivity between client and server.", e);
        }
    }

    public ResultSet execute(String statement) throws IOErrorException, NoValidSessionException {
        return client.execute(statement);
    }

    public ScanNodeResultIterator scanNode(String graphName,
                                           String nodeType,
                                           List<String> returnColumns,
                                           List<Integer> partsId,
                                           int limit) {
        return client.scanNode(graphName, nodeType, returnColumns, partsId, limit);
    }


    public ScanEdgeResultIterator scanEdge(String graphName,
                                           String edgeType,
                                           List<String> returnColumns,
                                           List<Integer> partsId,
                                           int limit) {
        return client.scanEdge(graphName, edgeType, returnColumns, partsId, limit);
    }


    /**
     * get node schema
     *
     * @param graphName graph name
     * @param nodeType  node type name
     */
    public NebulaNodeSchema getNodeSchema(String graphName, String nodeType)
            throws IOErrorException, NoValidSessionException {
        NebulaNodeSchema nodeSchema = new NebulaNodeSchema();
        Map<String, String> schema = new HashMap<>();
        String graphType = getGraphType(graphName);

        ResultSet result = client.execute(
                String.format("DESCRIBE NODE TYPE %s OF %s", nodeType, graphType));
        if (!result.isSucceeded() || result.isEmpty()) {
            throw new IllegalArgumentException(
                    "node type " + nodeType + " does not exist in " + graphName);
        }

        String pk = null;
        String pkDataType = null;
        while (result.hasNext()) {
            ResultSet.Record record = result.next();
            schema.put(record.get("property_name").asString(), record.get("data_type").asString());
            if ("Y".equals(record.get("primary_key").asString())) {
                pk = record.get("property_name").asString();
                pkDataType = record.get("data_type").asString();
            }
        }
        nodeSchema.setNodeTypeName(nodeType);
        nodeSchema.setNodePkName(pk);
        nodeSchema.setNodePkType(pkDataType);
        nodeSchema.setNodeProperties(schema);
        return nodeSchema;
    }

    /**
     * get edge type schema
     *
     * @param graphName graph name
     * @param edgeType  edge type name
     */
    public NebulaEdgeSchema getEdgeSchema(String graphName, String edgeType)
            throws IOErrorException, NoValidSessionException {
        Map<String, String> schema = new HashMap<>();

        String graphType = getGraphType(graphName);

        String descEdgeType = String.format(
                "CALL describe_graph_type('%s') filter type_name='%s' return type_pattern "
                        + "next CALL describe_edge_type('%s','%s') return *",
                graphType, edgeType, graphType, edgeType);

        ResultSet result = client.execute(descEdgeType);
        if (!result.isSucceeded() || result.isEmpty()) {
            throw new IllegalArgumentException(
                    "edge type " + edgeType + " does not exist in " + graphName);
        }

        String edgeTypePattern = null;
        while (result.hasNext()) {
            ResultSet.Record record = result.next();
            if (edgeTypePattern == null) {
                edgeTypePattern = record.get("type_pattern").asString();
            }
            schema.put(record.get("property_name").asString(), record.get("data_type").asString());
        }

        // get the src node type and dst node type according to edge pattern
        String srcNodeType = null;
        String dstNodeType = null;
        String edgeDirectionPattern = "\\((.*?)\\)-\\[.*?\\]->\\((.*?)\\)";
        String edgeUnDirectionPattern = "\\((.*?)\\)~\\[.*?\\]~\\((.*?)\\)";
        Pattern patternWithEdgeDirection = Pattern.compile(edgeDirectionPattern);
        Pattern patternWithoutEdgeDirection = Pattern.compile(edgeUnDirectionPattern);
        Matcher matcherWithEdgeDirection = patternWithEdgeDirection.matcher(edgeTypePattern);
        Matcher matcherWithoutEdgeDirection = patternWithoutEdgeDirection.matcher(edgeTypePattern);
        if (matcherWithEdgeDirection.matches()) {
            srcNodeType = matcherWithEdgeDirection.group(1);
            dstNodeType = matcherWithEdgeDirection.group(2);
        } else if (matcherWithoutEdgeDirection.matches()) {
            srcNodeType = matcherWithoutEdgeDirection.group(1);
            dstNodeType = matcherWithoutEdgeDirection.group(2);
        } else {
            throw new RuntimeException("Cannot parse the edge type pattern.");
        }

        NebulaNodeSchema srcNodeSchema = getNodeSchema(graphName, srcNodeType);
        NebulaNodeSchema dstNodeSchema = getNodeSchema(graphName, dstNodeType);

        NebulaEdgeSchema edgeSchema = new NebulaEdgeSchema();
        edgeSchema.setEdgeTypeName(edgeType);
        edgeSchema.setSourceNodeTypeName(srcNodeType);
        edgeSchema.setSourceNodePkName(srcNodeSchema.getNodePkName());
        edgeSchema.setSourceNodePkType(srcNodeSchema.getNodePkType());
        edgeSchema.setTargetNodeTypeName(dstNodeType);
        edgeSchema.setTargetNodePkName(dstNodeSchema.getNodePkName());
        edgeSchema.setTargetNodePkType(dstNodeSchema.getNodePkType());
        edgeSchema.setProperties(schema);
        return edgeSchema;
    }


    /**
     * get the graph type of graph
     *
     * @param graphName graph name
     * @return graph type name
     */
    private String getGraphType(String graphName) throws IOErrorException {
        ResultSet resultSet = client.execute("DESCRIBE GRAPH " + graphName);
        String graphType;
        if (resultSet.isSucceeded() && !resultSet.isEmpty()) {
            graphType = resultSet.next().values().get(1).asString();
        } else {
            throw new IllegalArgumentException("graphName " + graphName + " does not exist.");
        }
        return graphType;
    }


    public void close() {
        client.close();
    }
}
