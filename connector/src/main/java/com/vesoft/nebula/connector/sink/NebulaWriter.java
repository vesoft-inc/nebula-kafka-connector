
package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_DST_PKS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_PRIMARY_KEYS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SRC_PKS;

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
import java.util.HashSet;
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
        // check the if the fields for nebula in config file is valid
        if (nebulaNodeSchema != null) {
            if (config.primaryKeys.isEmpty()) {
                if (nebulaNodeSchema.getPkNames().size() == 1) {
                    config.primaryKeys.addAll(nebulaNodeSchema.getPkNames());
                } else if (config.sinkMode == NebulaSinkConnectConfig.InsertMode.UPDATE
                        || config.sinkMode == NebulaSinkConnectConfig.InsertMode.DELETE
                        || config.sinkMode == NebulaSinkConnectConfig.InsertMode.DETACHDELETE) {
                    throw new IllegalArgumentException(
                            "node type " + nebulaNodeSchema.getNodeTypeName()
                                    + " has multiple primary keys, please config "
                                    + CONNECT_PRIMARY_KEYS);
                }
            } else {
                assert (new HashSet<>(nebulaNodeSchema.getPkNames())
                        .containsAll(config.primaryKeys));
            }
        }
        if (nebulaEdgeSchema != null) {
            if (config.srcKeys.isEmpty()) {
                if (nebulaEdgeSchema.getSourcePkNameAndType().size() == 1) {
                    config.srcKeys.addAll(nebulaEdgeSchema.getSourcePkNameAndType().keySet());
                } else {
                    throw new IllegalArgumentException(
                            "source node type " + nebulaEdgeSchema.getSourceNodeTypeName()
                                    + " of " + nebulaEdgeSchema.getEdgeTypeName()
                                    + " has multiple primary keys, please config "
                                    + CONNECT_SRC_PKS);
                }
            } else {
                if (!nebulaEdgeSchema.getSourcePkNameAndType()
                        .keySet().containsAll(config.srcKeys)) {
                    throw new IllegalArgumentException(
                            "source node type " + nebulaEdgeSchema.getSourceNodeTypeName()
                                    + " of " + nebulaEdgeSchema.getEdgeTypeName()
                                    + " does not contain pk:" + config.srcKeys);
                }
            }
            if (config.dstKeys.isEmpty()) {
                if (nebulaEdgeSchema.getTargetPkNameAndType().size() == 1) {
                    config.dstKeys.addAll(nebulaEdgeSchema.getTargetPkNameAndType().keySet());
                } else {
                    throw new IllegalArgumentException(
                            "target node type " + nebulaEdgeSchema.getTargetNodeTypeName()
                                    + " of " + nebulaEdgeSchema.getEdgeTypeName()
                                    + " has multiple primary keys, please config "
                                    + CONNECT_DST_PKS);
                }
            } else {
                if (!nebulaEdgeSchema.getTargetPkNameAndType()
                        .keySet()
                        .containsAll(config.dstKeys)) {
                    throw new IllegalArgumentException(
                            "target node type " + nebulaEdgeSchema.getTargetNodeTypeName()
                                    + " of " + nebulaEdgeSchema.getEdgeTypeName()
                                    + " does not contain pk:" + config.dstKeys);
                }
            }
        }
    }

    public void write(final Collection<SinkRecord> records) throws IOErrorException,
                                                                   NoValidSessionException,
                                                                   InterruptedException {
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
        if (nodes.size() >= batchSize) {
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
        String    gql    = getNodeStatement(nodes);
        ResultSet result = execute(gql, config.graphNodeType, nodes.size());
        if (!result.isSucceeded()) {
            if (nodes.size() == 1) {
                log.error(">> write nodes failed: {}", result.getErrorMessage());
                log.error(">> failed gql: \n{}", gql);
                return;
            }
            log.warn(">> write failed for batch nodes, now retry the write one by one.");
            for (NebulaNode node : nodes) {
                execute(getNodeStatement(Arrays.asList(node)), config.graphNodeType, 1);
            }
        }
    }


    private void batchWriteEdge(List<NebulaEdge> edges)
            throws IOErrorException, NoValidSessionException, InterruptedException {
        if (edges.isEmpty()) {
            return;
        }
        ResultSet result = execute(getEdgeStatement(edges), config.graphEdgeType, edges.size());
        if (!result.isSucceeded()) {
            log.error(">> write edges failed: {}", result.getErrorMessage());
            if (edges.size() == 1) {
                return;
            }
            log.warn(">> write failed for batch edges, now retry the write one by one.");
            for (NebulaEdge edge : edges) {
                execute(getEdgeStatement(Arrays.asList(edge)), config.graphEdgeType, 1);
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
    private ResultSet execute(String statement, String type, int batch)
            throws IOErrorException, NoValidSessionException, InterruptedException {
        ResultSet result = graphProvider.execute(statement);
        if (result.isSucceeded()) {
            log.info(">> write ({}), batchSize({}), latency({}ms)",
                     type,
                     batch,
                     result.getLatency() / 1000.0);
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
                log.info(">> write ({}), batchSize({}), latency({}ms)",
                         type,
                         batch,
                         result.getLatency() / 1000.0);
                return result;
            }
        }
        log.error(">> write {} failed for {}", type, result.getErrorMessage());
        log.error(">> failed gql: {}", statement);
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
        List<String>        pks            = config.kafkaPrimaryKeys;
        for (int index = 0; index < nodePropertyNames.size(); index++) {
            Object kafkaValue = kafkaRecordProperties.get(nodePropertyNames.get(index));
            String propName   = nebulaNodePropertyNames.get(index);
            String value = NebulaUtils.extractPropertyValue(nebulaNodeSchema.getNodeProperties(),
                                                            nebulaNodePropertyNames.get(index),
                                                            String.valueOf(kafkaValue),
                                                            config.nullValue);
            nodeProperties.put(propName, value);
        }

        for (String pk : pks) {
            if (kafkaRecordProperties.get(pk) == null) {
                log.error(">>>>> record {} has null value for node primary key", pk);
                return null;
            }
        }

        // convert kafka properties to nebula properties
        NebulaNode node = new NebulaNode(nodeProperties);

        log.debug("nebula node: {}", node);
        return node;
    }

    private NebulaEdge getEdge(Map<String, Object> kafkaRecordProperties) {
        List<String> srcPkFields = config.kafkaSrcKeys;
        List<String> dstPkFields = config.kafkaDstKeys;

        for (String srcPk : srcPkFields) {
            if (kafkaRecordProperties.get(srcPk) == null) {
                log.error(">>>>> record {} has null value for source node primary key", srcPk);
                return null;
            }
        }

        for (String dstPk : dstPkFields) {
            if (kafkaRecordProperties.get(dstPk) == null) {
                log.error(">>>>> record {} has null value for target node primary key", dstPk);
                return null;
            }
        }

        Map<String, String> edgeProperties = new HashMap<>();

        List<String> kafkaEdgePropertyNames  = config.kafkaEdgePropertyNames;
        List<String> nebulaEdgePropertyNames = config.nebulaEdgePropertyNames;
        for (int index = 0; index < kafkaEdgePropertyNames.size(); index++) {
            Object kafkaValue = kafkaRecordProperties.get(kafkaEdgePropertyNames.get(index));
            String propName   = nebulaEdgePropertyNames.get(index);
            String value = NebulaUtils.extractPropertyValue(nebulaEdgeSchema.getProperties(),
                                                            nebulaEdgePropertyNames.get(index),
                                                            String.valueOf(kafkaValue),
                                                            config.nullValue);
            edgeProperties.put(propName, value);
        }

        List<String> nebulaSrcPkNames = config.srcKeys;
        List<String> nebulaDstPkNames = config.dstKeys;
        Map<String, String> srcPkAndValue = new HashMap<>();
        for (int i = 0; i < nebulaSrcPkNames.size(); i++) {
            String srcValue = NebulaUtils
                    .extractValue(
                            nebulaEdgeSchema.getSourcePkNameAndType().get(nebulaSrcPkNames.get(i)),
                            String.valueOf(kafkaRecordProperties.get(srcPkFields.get(i))),
                            config.nullValue);
            srcPkAndValue.put(nebulaSrcPkNames.get(i), srcValue);
        }

        Map<String, String> dstPkAndValue = new HashMap<>();
        for (int i = 0; i < nebulaDstPkNames.size(); i++) {
            String dstValue = NebulaUtils
                    .extractValue(
                            nebulaEdgeSchema.getTargetPkNameAndType().get(nebulaDstPkNames.get(i)),
                            String.valueOf(kafkaRecordProperties.get(dstPkFields.get(i))),
                            config.nullValue);
            dstPkAndValue.put(nebulaSrcPkNames.get(i), dstValue);
        }

        NebulaEdge edge = new NebulaEdge(
                srcPkAndValue,
                dstPkAndValue,
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
                                                  config.kafkaSrcKeys,
                                                  config.kafkaDstKeys,
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
