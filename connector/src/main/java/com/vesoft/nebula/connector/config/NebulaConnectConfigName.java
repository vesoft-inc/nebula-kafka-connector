
package com.vesoft.nebula.connector.config;

public class NebulaConnectConfigName {
    public static final String CONNECT_PREFIX                 = "nebula.";
    public static final String CONNECT_GRAPH_ADDRESS          = CONNECT_PREFIX + "graph.servers";
    public static final String CONNECT_GRAPH_USER             = CONNECT_PREFIX + "user";
    public static final String CONNECT_GRAPH_PASSWD           = CONNECT_PREFIX + "passwd";
    public static final String CONNECT_AUTH_OPTIONS           = CONNECT_PREFIX + "authOptions";

    public static final String CONNECT_ENABLE_TLS = CONNECT_PREFIX + "enable.tls";
    public static final String CONNECT_DISABLE_VERIFY_CERT = CONNECT_PREFIX + "disable.verify.cert";
    public static final String CONNECT_TLS_CA_PATH = CONNECT_PREFIX + "tls.ca.path";
    public static final String CONNECT_TLS_CERT_PATH = CONNECT_PREFIX + "tls.cert.path";
    public static final String CONNECT_TLS_KEY_PATH = CONNECT_PREFIX + "tls.key.path";

    public static final String CONNECT_GRAPH_NAME             = CONNECT_PREFIX + "graph.name";
    public static final String CONNECT_GRAPH_DATA_TYPE        = CONNECT_PREFIX + "data.type";
    public static final String CONNECT_GRAPH_NODE_TYPE_NAME   = CONNECT_PREFIX + "node.typeName";
    public static final String CONNECT_GRAPH_EDGE_TYPE_NAME   = CONNECT_PREFIX + "edge.typeName";
    public static final String CONNECT_PRIMARY_KEYS           = CONNECT_PREFIX + "node.primaryKeys";
    public static final String CONNECT_KAFKA_PRIMARY_KEYS     = "kafka.node.primaryKeys";
    public static final String CONNECT_SRC_PKS                = CONNECT_PREFIX + "edge.srcPks";
    public static final String CONNECT_KAFKA_SRC_PKS          = "kafka.edge.srcPks";
    public static final String CONNECT_DST_PKS                = CONNECT_PREFIX + "edge.dstPks";
    public static final String CONNECT_KAFKA_DST_PKS          = "kafka.edge.dstPks";
    public static final String CONNECT_KAFKA_NODE_PROPERTIES  = "kafka.node.property.names";
    public static final String CONNECT_NEBULA_NODE_PROPERTIES = CONNECT_PREFIX
            + "node.property.names";
    public static final String CONNECT_KAFKA_EDGE_PROPERTIES  = "kafka.edge.property.names";
    public static final String CONNECT_NEBULA_EDGE_PROPERTIES = CONNECT_PREFIX
            + "edge.property.names";
    public static final String CONNECT_KAFKA_NULL_VALUE = "kafka.null.value";
    public static final String CONNECT_CONNECT_TIMEOUT        = CONNECT_PREFIX + "connect.timeout";
    public static final String CONNECT_REQUEST_TIMEOUT        = CONNECT_PREFIX + "request.timeout";
    public static final String CONNECT_SINK_RETRY_TIMES       = CONNECT_PREFIX + "sink.retry.times";

    public static final String CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY = CONNECT_PREFIX
            + "sink.interval.time.mills";

    public static final String CONNECT_NEBULA_PROPERTIES = CONNECT_PREFIX + "return.columns";
    public static final String CONNECT_SINK_PARTITION    = CONNECT_PREFIX + "sink.partitions";
    public static final String CONNECT_SINK_MODE         = CONNECT_PREFIX + "sink.mode";
    public static final String CONNECT_BATCH_SIZE        = CONNECT_PREFIX + "batchSize";

    // ######## source connect config
    public static final String CONNECT_SOURCE_NODE_TYPES   = "connect.source.nodeTypes";
    public static final String CONNECT_SOURCE_EDGE_TYPES   = "connect.source.edgeTypes";
    public static final String CONNECT_MAX_TASK            = "maxTask";
    public static final String CONNECT_SOURCE_TOPIC_PREFIX = "topic.prefix";

    public static final String CONNECT_POLLING_SLEEP_MS = "connect.source.pollingMs";


    // #########  default config value
    public static final NebulaSinkConnectConfig.InsertMode DEFAULT_INSERT_MODE =
            NebulaSinkConnectConfig.InsertMode.INSERTIGNORE;

    public static final String  DEFAULT_CONNECT_GRAPH_USER                    = "root";
    public static final int     DEFAULT_SINK_RETRY_TIMES                      = 0;
    public static final int     DEFAULT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY = 0;
    public static final boolean DEFAULT_SINK_CONNECT_RECONNECT                = false;
    public static final int     DEFAULT_CONNECTOR_CONNECT_TIMEOUT             = 3000;
    public static final int     DEFAULT_CONNECTOR_REQUEST_TIMEOUT             = 3000;
    public static final int     DEFAULT_CONNECT_SINK_PARTITION                = 10;
    public static final int     DEFAULT_BATCH_SIZE                            = 2000;
    public static final int     DEFAULT_CONNECT_RETRY_TIMES                   = 1;
    public static final long    DEFAULT_INTERVAL_BETWEEN_RETRY                = 0;


    // ########## NebulaGraph Sink statement template
    public static String BATCH_INSERT_NODE_TEMPLATE = "USE `%s` INSERT NODE `%s` %s";
    public static String BATCH_INSERT_EDGE_TEMPLATE = "USE `%s` INSERT EDGE `%s` %s";
    public static String NODE_VALUES_TEMPLATE       = "({%s})";
    public static String EDGE_VALUES_TEMPLATE       = "({`%s`:%s})-[{%s}]->({`%s`:%s})";
    public static String PROPERTY_TEMPLATE          = "`%s`:%s";

    // constant for internal
    public static final String NEBULA_PARTS_FOR_EACH_TASK = "nebula.parts";


}
