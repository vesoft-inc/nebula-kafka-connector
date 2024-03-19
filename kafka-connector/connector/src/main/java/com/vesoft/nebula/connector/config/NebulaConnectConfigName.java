/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.config;

public class NebulaConnectConfigName {
    public static final String CONNECT_PREFIX = "nebula.";
    public static final String CONNECT_GRAPH_ADDRESS = CONNECT_PREFIX + "graph.servers";
    public static final String CONNECT_GRAPH_USER = CONNECT_PREFIX + "user";
    public static final String CONNECT_GRAPH_PASSWD = CONNECT_PREFIX + "passwd";
    public static final String CONNECT_GRAPH_NAME = CONNECT_PREFIX + "graph.name";
    public static final String CONNECT_GRAPH_TYPE = CONNECT_PREFIX + "data.type";
    public static final String CONNECT_GRAPH_NODE_TYPE = CONNECT_PREFIX + "node.type";
    public static final String CONNECT_GRAPH_EDGE_TYPE = CONNECT_PREFIX + "edge.type";
    public static final String CONNECT_PRIMARY_KEY = CONNECT_PREFIX + "node.primaryKey";
    public static final String CONNECT_SRC_KEY = CONNECT_PREFIX + "edge.srcKey";
    public static final String CONNECT_DST_KEY = CONNECT_PREFIX + "edge.dstKey";
    public static final String CONNECT_KAFKA_NODE_PROPERTIES = "kafka.node.property.names";
    public static final String CONNECT_NEBULA_NODE_PROPERTIES = CONNECT_PREFIX + "node.property.names";
    public static final String CONNECT_KAFKA_EDGE_PROPERTIES = "kafka.edge.property.names";
    public static final String CONNECT_NEBULA_EDGE_PROPERTIES = CONNECT_PREFIX + "edge.property.names";
    public static final String CONNECT_REQUEST_TIMEOUT = CONNECT_PREFIX + "request.timeout";
    public static final String CONNECT_SINK_RETRY_TIMES = CONNECT_PREFIX + "sink.retry.times";
    public static final String CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY = CONNECT_PREFIX +
            "sink.interval.time.mills";
    public static final String CONNECT_SINK_PARTITION = CONNECT_PREFIX + "sink.partitions";
    public static final String CONNECT_SINK_MODE = CONNECT_PREFIX + "sink.mode";
    public static final String CONNECT_HEALTH_CHECK_TIME = CONNECT_PREFIX + "healthCheckTime";
    public static final String CONNECT_BLOCK_WHEN_EXHAUSTED = CONNECT_PREFIX + "blockWhenExhausted";
    public static final String CONNECT_MAX_WAIT_TIME = CONNECT_PREFIX + "maxWaitMills";
    public static final String CONNECT_STRICTLY_SERVER_HEALTHY = CONNECT_PREFIX + "strictlyServerHealthy";
    public static final String CONNECT_BATCH_SIZE = CONNECT_PREFIX + "sink.batchSize";

    // #########  default config value
    public static final String DEFAULT_CONNECT_GRAPH_USER = "root";
    public static final String DEFAULT_CONNECT_GRAPH_PASSWD = "nebula";
    public static final String DEFAULT_CONNECT_GRAPH_TYPE = "NODE";

    public static final int DEFAULT_SINK_RETRY_TIMES = 0;
    public static final int DEFAULT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY = 0;
    public static final boolean DEFAULT_SINK_CONNECT_RECONNECT = false;
    public static final String DEFAULT_INSERT_MODE = "INSERT";
    public static final int DEFAULT_SINK_CONNECT_TIMEOUT = 3000;
    public static final int DEFAULT_SINK_REQUEST_TIMEOUT = 3000;
    public static final int DEFAULT_CONNECT_SINK_PARTITION = 10;
    public static final int DEFAULT_BATCH_SIZE = 2000;


    // ########## NebulaGraph Sink statement template
    public static String BATCH_INSERT_NODE_TEMPLATE = "USE `%s` INSERT NODE `%s` %s";
    public static String BATCH_INSERT_EDGE_TEMPLATE = "USE `%s` INSERT EDGE `%s` %s";
    public static String NODE_VALUES_TEMPLATE = "({%s})";
    public static String EDGE_VALUES_TEMPLATE = "({`%s`:%s})-[{%s}]->({`%s`:%s})";
    public static String PROPERTY_TEMPLATE = "`%s`:%s";


}
