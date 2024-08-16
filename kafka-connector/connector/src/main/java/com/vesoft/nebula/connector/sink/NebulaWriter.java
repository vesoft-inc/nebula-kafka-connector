
package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.BATCH_INSERT_EDGE_TEMPLATE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.BATCH_INSERT_NODE_TEMPLATE;

import com.vesoft.nebula.connector.config.NebulaConnectDataTypeEnum;
import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import com.vesoft.nebula.connector.connection.NebulaGraphProvider;
import com.vesoft.nebula.connector.converter.NebulaRecordConverter;
import com.vesoft.nebula.connector.util.GroupUtils;
import com.vesoft.nebula.connector.util.NebulaUtils;
import com.vesoft.nebula.driver.graph.ErrorCode;
import com.vesoft.nebula.driver.graph.data.ResultSet;
import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import com.vesoft.nebula.driver.graph.exception.NoValidSessionException;
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
    private final NebulaGraphProvider     graphProvider;
    private       List<NebulaNode>        nodes = new ArrayList<>();
    private       List<NebulaEdge>        edges = new ArrayList<>();

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
                continue;
            }
            // the mapping for kafka property name and property value. get values
            // according to propertyNames config
            Map<String, Object> kafkaRecordProperties = NebulaRecordConverter
                    .convertRecord(record);
            switch (config.dataType) {
                case NODE:
                    nodes.add(getNode(kafkaRecordProperties));
                    break;
                case EDGE:
                    edges.add(getEdge(kafkaRecordProperties));
                    break;
                case BOTH:
                    nodes.add(getNode(kafkaRecordProperties));
                    edges.add(getEdge(kafkaRecordProperties));
                    break;
                default: // nothing
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
        if (nodes.size() > batchSize) {
            // write nodes batch by batchSize
            List<List<NebulaNode>> groupNodes = GroupUtils.getGroups(nodes, batchSize);
            for (List<NebulaNode> batchNodes : groupNodes) {
                batchWriteNode(batchNodes);
            }
        } else {
            batchWriteNode(nodes);
        }
        nodes.clear();
        // write edges
        if (edges.size() > batchSize) {
            // write edges batch by batchSize
            List<List<NebulaEdge>> groupEdges = GroupUtils.getGroups(edges, batchSize);
            for (List<NebulaEdge> batchEdges : groupEdges) {
                batchWriteEdge(batchEdges);
            }
        } else {
            batchWriteEdge(edges);
        }
        edges.clear();
    }


    private void batchWriteNode(List<NebulaNode> nodes)
            throws IOErrorException, NoValidSessionException, InterruptedException {
        if (nodes.isEmpty()) {
            return;
        }
        ResultSet result = execute(getNodeStatement(nodes), nodeType, nodes.size());
        if (!result.isSucceeded()) {
            if (nodes.size() == 1) {
                return;
            }
            log.warn(">> write failed for batch nodes, now retry the write one by one.");
            for (NebulaNode node : nodes) {
                execute(getNodeStatement(Arrays.asList(node)), nodeType, 1);
            }
        }
    }


    private void batchWriteEdge(List<NebulaEdge> edges)
            throws IOErrorException, NoValidSessionException, InterruptedException {
        if (edges.isEmpty()) {
            return;
        }
        ResultSet result = execute(getEdgeStatement(edges), edgeType, edges.size());
        if (!result.isSucceeded()) {
            if (edges.size() == 1) {
                return;
            }
            log.warn(">> write failed for batch edges, now retry the write one by one.");
            for (NebulaEdge edge : edges) {
                execute(getEdgeStatement(Arrays.asList(edge)), edgeType, 1);
            }
        }
    }

    /**
     * execute the write operation for nodes and edges
     *
     * @param statement gql statement
     * @param type      the statement data type, node or edge
     * @param batch     batch size of nodes or edges for param statement
     * @return {@link ResultSet}
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
     * @param kafkaRecordProperties mapping of kafka property name and kafka property value
     */
    private NebulaNode getNode(Map<String, Object> kafkaRecordProperties) {
        List<String> nodePropertyNames       = config.kafkaNodePropertyNames;
        List<String> nebulaNodePropertyNames = config.nebulaNodePropertyNames;

        Map<String, String> nodeProperties = new HashMap<>();
        String              pk             = config.primaryKey;
        for (int index = 0; index < nodePropertyNames.size(); index++) {
            Object kafkaValue = kafkaRecordProperties.get(nodePropertyNames.get(index));
            String propName   = nebulaNodePropertyNames.get(index);
            String value = NebulaUtils.extractPropertyValue(nebulaNodeSchema.getNodeProperties(),
                                                            nebulaNodePropertyNames.get(index),
                                                            String.valueOf(kafkaValue));
            nodeProperties.put(propName, value);
        }

        if (kafkaRecordProperties.get(pk) == null) {
            log.error(">>>>> record {} has null value for node primary key", pk);
            return null;
        }
        if (nebulaNodeSchema.getNodePkType().equals("INT64")
                && !NebulaUtils.isNumeric(kafkaRecordProperties.get(pk).toString())) {
            log.error(">>>>> record {} value {} is not INT64 for node primary key",
                      pk, kafkaRecordProperties.get(pk));
            return null;
        }
        // convert kafka properties to nebula properties
        NebulaNode node = new NebulaNode(nodeProperties);

        log.debug("nebula node: {}", node);
        return node;
    }

    private NebulaEdge getEdge(Map<String, Object> kafkaRecordProperties) {
        String srcPk = config.srcKey;
        String dstPk = config.dstKey;

        if (kafkaRecordProperties.get(srcPk) == null) {
            log.error(">>>>> record {} has null value for source node primary key", srcPk);
            return null;
        }
        if (nebulaEdgeSchema.getSourceNodePkType().equals("INT64")
                && !NebulaUtils.isNumeric(kafkaRecordProperties.get(srcPk).toString())) {
            log.error(">>>>> record {} value {} is not INT64 for source node primary key",
                      srcPk, kafkaRecordProperties.get(srcPk));
            return null;
        }

        if (kafkaRecordProperties.get(dstPk) == null) {
            log.error(">>>>> record {} has null value for target node primary key", dstPk);
            return null;
        }
        if (nebulaEdgeSchema.getTargetNodePkType().equals("INT64")
                && !NebulaUtils.isNumeric(kafkaRecordProperties.get(dstPk).toString())) {
            log.error(">>>>> record {} value {} is not INT64 for target node primary key",
                      dstPk, kafkaRecordProperties.get(dstPk));
            return null;
        }

        Map<String, String> edgeProperties = new HashMap<>();

        List<String> kafkaEdgePropertyNames  = config.kafkaEdgePropertyNames;
        List<String> nebulaEdgePropertyNames = config.nebulaEdgePropertyNames;
        for (int index = 0; index < kafkaEdgePropertyNames.size(); index++) {
            Object kafkaValue = kafkaRecordProperties.get(kafkaEdgePropertyNames.get(index));
            String propName   = nebulaEdgePropertyNames.get(index);
            String value = NebulaUtils.extractPropertyValue(nebulaEdgeSchema.getProperties(),
                                                            nebulaEdgePropertyNames.get(index),
                                                            String.valueOf(kafkaValue));
            edgeProperties.put(propName, value);
        }


        String srcValue = NebulaUtils
                .extractIdValue(nebulaEdgeSchema.getSourceNodePkType(),
                                String.valueOf(kafkaRecordProperties.get(srcPk)));
        String dstValue = NebulaUtils
                .extractIdValue(nebulaEdgeSchema.getTargetNodePkType(),
                                String.valueOf(kafkaRecordProperties.get(dstPk)));
        NebulaEdge edge = new NebulaEdge(
                srcValue,
                dstValue,
                edgeProperties);
        return edge;
    }

    private String getNodeStatement(List<NebulaNode> nodes) {
        String      graphName   = config.graphName;
        NebulaNodes nebulaNodes = new NebulaNodes(nebulaNodeSchema, nodes);
        switch (config.sinkMode) {
            case INSERT:
            case INSERTIGNORE:
            case INSERTREPLACE:
                return nebulaNodes.getInsertStatement(graphName, config.sinkMode);
            case UPDATE:
                return nebulaNodes.getUpdateStatement(graphName);
            case DELETE:
            case DETACHDELETE:
                return nebulaNodes.getDeleteStatement(graphName, config.sinkMode);
            default:
                throw new IllegalArgumentException("unsupported sink mode.");
        }
    }

    private String getEdgeStatement(List<NebulaEdge> edges) {
        List<String> kafkaEdgePropertyNames  = config.kafkaEdgePropertyNames;
        List<String> nebulaEdgePropertyNames = config.nebulaEdgePropertyNames;
        String       graphName               = config.graphName;
        NebulaEdges nebulaEdges = new NebulaEdges(nebulaEdgeSchema,
                                                  config.srcKey,
                                                  config.dstKey,
                                                  kafkaEdgePropertyNames,
                                                  nebulaEdgePropertyNames,
                                                  edges);
        switch (config.sinkMode) {
            case INSERT:
            case INSERTIGNORE:
            case INSERTREPLACE:
                return nebulaEdges.getInsertStatement(graphName, config.sinkMode);
            case UPDATE:
                return nebulaEdges.getUpdateStatement(graphName);
            case DELETE:
            case DETACHDELETE:
                return nebulaEdges.getDeleteStatement(graphName);
            default:
                throw new IllegalArgumentException("unsupported sink mode.");
        }
    }
}
