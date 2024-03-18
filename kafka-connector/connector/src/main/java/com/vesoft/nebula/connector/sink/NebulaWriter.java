/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.BATCH_INSERT_EDGE_TEMPLATE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.BATCH_INSERT_NODE_TEMPLATE;
import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import com.vesoft.nebula.connector.config.NebulaConnectDataTypeEnum;
import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import com.vesoft.nebula.connector.connection.NebulaGraphProvider;
import com.vesoft.nebula.connector.converter.NebulaRecordConverter;
import com.vesoft.nebula.connector.exceptions.DataFormatException;
import com.vesoft.nebula.connector.util.NebulaUtils;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.apache.kafka.connect.sink.SinkRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * a writer to batch write the records into NebulaGraph
 * refer： https://github.dev/ClickHouse/clickhouse-kafka-connect/blob/main/src/main/java/com
 * /clickhouse/kafka/connect/sink/db/ClickHouseWriter.java
 */
public class NebulaWriter {
    private static final Logger log = LoggerFactory.getLogger(NebulaWriter.class);

    private final NebulaSinkConnectConfig config;
    private final NebulaGraphProvider graphProvider;
    private List<NebulaNode> nodes = new ArrayList<>();
    private List<NebulaEdge> edges = new ArrayList<>();

    private final String nodeType = "node";
    private final String edgeType = "edge";

    private NebulaNodeSchema nebulaNodeSchema = null;
    private NebulaEdgeSchema nebulaEdgeSchema = null;

    NebulaWriter(final NebulaSinkConnectConfig config) {
        this.config = config;
        graphProvider = new NebulaGraphProvider(config);
        try {
            if (config.dataType == NebulaConnectDataTypeEnum.NODE) {
                nebulaNodeSchema = graphProvider.getNodeSchema(config.graphName,
                        config.graphNodeType);
            } else if (config.dataType == NebulaConnectDataTypeEnum.EDGE) {
                nebulaEdgeSchema = graphProvider.getEdgeSchema(config.graphName,
                        config.graphEdgeType);
            } else {
                nebulaNodeSchema = graphProvider.getNodeSchema(config.graphName,
                        config.graphNodeType);
                nebulaEdgeSchema = graphProvider.getEdgeSchema(config.graphName,
                        config.graphEdgeType);
            }
        } catch (Exception e) {
            throw new RuntimeException("failed to get NebulaGraph schema", e);
        }
    }

    public void write(final Collection<SinkRecord> records) throws IOErrorException,
            NoValidSessionException, InterruptedException {
        for (SinkRecord record : records) {
            if (record.value() == null) {
                log.warn(String.format("Getting empty record skip the insert topic[%s] offset[%s]",
                        record.topic(), record.kafkaOffset()));
            } else {
                // the mapping for kafka property name and property value. get values
                // according to propertyNames config
                Map<String, Object> properties = NebulaRecordConverter.convertRecord(record);
                switch (config.dataType) {
                    case NODE:
                        nodes.add(getNode(properties));
                        break;
                    case EDGE:
                        edges.add(getEdge(properties));
                        break;
                    case BOTH:
                        nodes.add(getNode(properties));
                        edges.add(getEdge(properties));
                        break;
                }
            }
        }
        commit();
    }

    /**
     * commit to execute the nodes and edges when their size beyond the sinkBatchSize
     */
    public void commit() throws IOErrorException, NoValidSessionException, InterruptedException {
        int batchSize = config.sinkBatchSize;
        // write nodes
        if (nodes.size() >= batchSize) {
            ResultSet result = execute(getNodeStatement(nodes), nodeType, nodes.size());
            if (result.isSucceeded()) {
                log.info(">> write {} batch({}), latency({})", nodeType, nodes.size(),
                        result.getLatency());
            } else {
                log.warn(">> write failed for batch nodes, now retry the write one by one.");
                for (NebulaNode node : nodes) {
                    execute(getNodeStatement(Arrays.asList(node)), nodeType, 1);
                }
            }
            nodes.clear();
        }
        // write edges
        if (edges.size() >= batchSize) {
            ResultSet result = execute(getEdgeStatement(edges), edgeType, edges.size());
            if (result.isSucceeded()) {
                log.info(">> write {} batch({}), latency({})", edgeType, edges.size(),
                        result.getLatency());
            } else {
                log.warn(">> write failed for batch edges, now retry the write one by one.");
                for (NebulaEdge edge : edges) {
                    execute(getEdgeStatement(Arrays.asList(edge)), edgeType, 1);
                }
            }
            edges.clear();
        }
    }

    /**
     * execute the write operation for nodes and edges
     *
     * @param statement gql statement
     * @param type      the statement data type, node or edge
     * @param batch     batch size of nodes or edges for param statement
     * @return {@link ResultSet}
     * TODO change the match for result gqlStatus
     */
    private ResultSet execute(String statement, String type, int batch) throws IOErrorException,
            NoValidSessionException, InterruptedException {
        ResultSet result = graphProvider.execute(statement);
        if (result.isSucceeded()) {
            log.info(">> write {} batch({}), latency({})", type, batch, result.getLatency());
            return result;
        }

        int retry = 0;
        // retry the write execution for RPC_ERROR, LEADER_CHANGED, RAFT_ERROR
        while (retry++ < config.retryTimes && (result.getErrorCode().isRpcError()
                || result.getErrorCode() == ErrorCode.LEADER_CHANGED
                || result.getErrorCode().isRaftError())) {
            Thread.sleep(config.intervalTimeMill);
            result = graphProvider.execute(statement);
            if (result.isSucceeded()) {
                log.info(">> write {} batch({}), latency({})", type, batch, result.getLatency());
                return result;
            }
        }

        log.error(">> write {} failed for {}", type, result.getErrorMessage());
        return result;
    }

    /**
     * close the NebulaGraph client
     */
    public void close() {
        graphProvider.close();
    }


    /**
     * @param properties mapping of kafka property name and kafka property value
     */
    private NebulaNode getNode(Map<String, Object> properties) {
        List<String> nodePropertyNames = config.kafkaNodePropertyNames;
        List<String> nebulaNodePropertyNames = config.nebulaNodePropertyNames;

        Map<String, Object> nodeProperties = new HashMap<>();
        String pk = config.primaryKey;
        for (int index = 0; index < nodePropertyNames.size(); index++) {
            nodeProperties.put(nebulaNodePropertyNames.get(index),
                    properties.get(nodePropertyNames.get(index)));
        }

        if (properties.get(pk) == null) {
            log.error(">>>>> record {} has null value for node primary key", pk);
            return null;
        }
        if (nebulaNodeSchema.getNodePkType().equals("INT64")
                && !NebulaUtils.isNumeric(properties.get(pk).toString())) {
            log.error(">>>>> record {} value {} is not INT64 for node primary key", pk, properties.get(pk));
            return null;
        }
        // 将kafka properties 转换成nebula properties
        NebulaNode node = new NebulaNode(nodeProperties);

        log.debug("nebula node: {}", node);
        return node;
    }

    private NebulaEdge getEdge(Map<String, Object> properties) {
        List<String> edgePropertyNames = config.kafkaEdgePropertyNames;
        List<String> nebulaEdgePropertyNames = config.nebulaEdgePropertyNames;
        String srcPkName = nebulaEdgeSchema.getSourceNodePkName();
        String srcPk = config.srcKey;
        String dstPkName = nebulaEdgeSchema.getTargetNodePkName();
        String dstPk = config.dstKey;

        if (properties.get(srcPk) == null) {
            log.error(">>>>> record {} has null value for source node primary key", srcPk);
            return null;
        }
        if (nebulaEdgeSchema.getSourceNodePkType().equals("INT64")
                && !NebulaUtils.isNumeric(properties.get(srcPk).toString())) {
            log.error(">>>>> record {} value {} is not INT64 for source node primary key", srcPk, properties.get(srcPk));
            return null;
        }

        if (properties.get(dstPk) == null) {
            log.error(">>>>> record {} has null value for target node primary key", dstPk);
            return null;
        }
        if (nebulaEdgeSchema.getTargetNodePkType().equals("INT64")
                && !NebulaUtils.isNumeric(properties.get(dstPk).toString())) {
            log.error(">>>>> record {} value {} is not INT64 for target node primary key", dstPk, properties.get(dstPk));
            return null;
        }

        Map<String, Object> edgeProperties = new HashMap<>();

        for (int index = 0; index < edgePropertyNames.size(); index++) {
            edgeProperties.put(nebulaEdgePropertyNames.get(index),
                    properties.get(edgePropertyNames.get(index)));
        }
        NebulaEdge edge = new NebulaEdge(srcPkName, String.valueOf(properties.get(srcPk)),
                dstPkName, String.valueOf(properties.get(dstPk)), edgeProperties);
        return edge;
    }

    private String getNodeStatement(List<NebulaNode> nodes) {
        List<String> nodeStringList = new ArrayList<>();
        for (NebulaNode node : nodes) {
            String nodeStatement = null;
            try {
                nodeStatement = node.getNodeStatement(nebulaNodeSchema);
            } catch (DataFormatException e) {
                log.error(">>>> dirty data {exception: {}, record: {}}", e.getMessage(),
                        node.getNodeString());
                continue;
            }
            nodeStringList.add(nodeStatement);
        }
        String graphName = config.graphName;
        String nodeName = config.graphNodeType;

        switch (config.sinkMode) {
            case INSERT:
                return String.format(BATCH_INSERT_NODE_TEMPLATE, graphName, nodeName,
                        NebulaUtils.join(nodeStringList, ","));
            case UPDATE:
                return null; // TODO implement the update statement
            case DELETE:
                return null; // TODO implement the delete statement
        }
        return null; // placeholder
    }

    private String getEdgeStatement(List<NebulaEdge> edges) {
        List<String> edgeStringList = new ArrayList<>();
        for (NebulaEdge edge : edges) {
            String edgeStatement = null;
            try {
                edgeStatement = edge.getEdgeStatement(nebulaEdgeSchema);
            } catch (DataFormatException e) {
                log.error(">>>> dirty data {exception: {}, record: {}}", e.getMessage(),
                        edge.getEdgeString());
                continue;
            }
            edgeStringList.add(edgeStatement);
        }
        String graphName = config.graphName;
        String edgeName = config.graphEdgeType;

        return String.format(BATCH_INSERT_EDGE_TEMPLATE, graphName, edgeName,
                NebulaUtils.join(edgeStringList, ","));
    }
}
