/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.config;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_BATCH_SIZE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_BLOCK_WHEN_EXHAUSTED;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_CONNECT_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_DST_KEY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_ADDRESS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_EDGE_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NODE_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_PASSWD;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_USER;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_HEALTH_CHECK_TIME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_KAFKA_EDGE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_KAFKA_NODE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_MAX_WAIT_TIME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_NEBULA_EDGE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_NEBULA_NODE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_PRIMARY_KEY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_REQUEST_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_MODE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_PARTITION;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_RECONNECT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_RETRY_TIMES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SRC_KEY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_STRICTLY_SERVER_HEALTHY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_BATCH_SIZE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECT_GRAPH_PASSWD;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECT_GRAPH_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECT_GRAPH_USER;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECT_SINK_PARTITION;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_INSERT_MODE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_SINK_CONNECT_RECONNECT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_SINK_CONNECT_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_SINK_REQUEST_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_SINK_RETRY_TIMES;
import com.vesoft.nebula.connector.util.ConfigUtils;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.ConfigException;

/**
 * configurations for Nebula kafka connect
 */
public class NebulaSinkConnectConfig extends AbstractConfig implements Serializable {

    public enum InsertMode {
        INSERT,
        UPDATE,
        DELETE;
    }

    public enum DataType {
        NODE,
        EDGE,
        BOTH;
    }

    public final String connectorName;
    public final String graphServers;
    public final String user;
    public final String passwd;
    public final String graphName;
    public final NebulaConnectDataTypeEnum dataType;
    public final String graphNodeType;
    public final String graphEdgeType;
    public final String primaryKey;
    public final String srcKey;
    public final String dstKey;
    public List<String> kafkaNodePropertyNames = new ArrayList<>();
    public List<String> nebulaNodePropertyNames = new ArrayList<>();
    public List<String> kafkaEdgePropertyNames = new ArrayList<>();
    public List<String> nebulaEdgePropertyNames = new ArrayList<>();
    public final int connectTimeout;
    public final int requestTimeout;
    public final int retryTimes;
    public final int intervalTimeMill;
    public final boolean reconnect;
    public final int healthCheckTime;
    public final boolean blockWithExhausted;
    public final int maxWaitTimeWhenSessionExhausted;
    public final boolean strictlyServerHealth;
    public final int sinkPartition;

    public final InsertMode sinkMode;

    public final int sinkBatchSize;

    public static final ConfigDef CONFIG_DEF = configDef();


    public NebulaSinkConnectConfig(Map<?, ?> props) {
        super(CONFIG_DEF, props);
        connectorName = ConfigUtils.connectorName(props);
        graphServers = getString(CONNECT_GRAPH_ADDRESS);
        user = getString(CONNECT_GRAPH_USER);
        passwd = getString(CONNECT_GRAPH_PASSWD);
        graphName = getString(CONNECT_GRAPH_NAME);
        dataType = NebulaConnectDataTypeEnum.getDataType(getString(CONNECT_GRAPH_TYPE));
        graphNodeType = getString(CONNECT_GRAPH_NODE_TYPE);
        graphEdgeType = getString(CONNECT_GRAPH_EDGE_TYPE);
        primaryKey = getString(CONNECT_PRIMARY_KEY);
        srcKey = getString(CONNECT_SRC_KEY);
        dstKey = getString(CONNECT_DST_KEY);
        kafkaNodePropertyNames = getList(CONNECT_KAFKA_NODE_PROPERTIES);
        kafkaEdgePropertyNames = getList(CONNECT_KAFKA_EDGE_PROPERTIES);
        nebulaNodePropertyNames = getList(CONNECT_NEBULA_NODE_PROPERTIES);
        nebulaEdgePropertyNames = getList(CONNECT_NEBULA_EDGE_PROPERTIES);
        connectTimeout = getInt(CONNECT_CONNECT_TIMEOUT);
        requestTimeout = getInt(CONNECT_REQUEST_TIMEOUT);
        retryTimes = getInt(CONNECT_SINK_RETRY_TIMES);
        intervalTimeMill = getInt(CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY);
        reconnect = getBoolean(CONNECT_SINK_RECONNECT);
        healthCheckTime = getInt(CONNECT_HEALTH_CHECK_TIME);
        blockWithExhausted = getBoolean(CONNECT_BLOCK_WHEN_EXHAUSTED);
        maxWaitTimeWhenSessionExhausted = getInt(CONNECT_MAX_WAIT_TIME);
        strictlyServerHealth = getBoolean(CONNECT_STRICTLY_SERVER_HEALTHY);
        sinkPartition = getInt(CONNECT_SINK_PARTITION);
        sinkMode = ConfigUtils.sinkMode(props);
        sinkBatchSize = getInt(CONNECT_BATCH_SIZE);
    }

    public static ConfigDef configDef() {
        return new ConfigDef()
                .define(CONNECT_GRAPH_ADDRESS,
                        ConfigDef.Type.STRING,
                        ConfigDef.Importance.HIGH,
                        "graph server address, comma-split for multiple servers. ")
                .define(CONNECT_GRAPH_USER,
                        ConfigDef.Type.STRING,
                        DEFAULT_CONNECT_GRAPH_USER,
                        new ConfigDef.NonEmptyString(),
                        ConfigDef.Importance.HIGH,
                        "NebulaGraph user, default user: " + DEFAULT_CONNECT_GRAPH_USER)
                .define(CONNECT_GRAPH_PASSWD,
                        ConfigDef.Type.STRING,
                        DEFAULT_CONNECT_GRAPH_PASSWD,
                        new ConfigDef.NonEmptyString(),
                        ConfigDef.Importance.HIGH,
                        "NebulaGraph passwd, default passwd: " + DEFAULT_CONNECT_GRAPH_PASSWD)
                .define(CONNECT_GRAPH_NAME,
                        ConfigDef.Type.STRING,
                        "nba",
                        new ConfigDef.NonEmptyString(),
                        ConfigDef.Importance.HIGH,
                        "Name of the graph to streaming objects to")
                .define(CONNECT_GRAPH_TYPE,
                        ConfigDef.Type.STRING,
                        DEFAULT_CONNECT_GRAPH_TYPE,
                        EnumValidator.in(DataType.values()),
                        ConfigDef.Importance.HIGH,
                        "the data type, node or edge, default is " + DEFAULT_CONNECT_GRAPH_TYPE)
                .define(CONNECT_GRAPH_NODE_TYPE,
                        ConfigDef.Type.STRING,
                        null,
                        ConfigDef.Importance.HIGH,
                        "Name of the node type")
                .define(CONNECT_GRAPH_EDGE_TYPE,
                        ConfigDef.Type.STRING,
                        null,
                        ConfigDef.Importance.HIGH,
                        "Name of the edge type")
                .define(CONNECT_SINK_MODE,
                        ConfigDef.Type.STRING,
                        DEFAULT_INSERT_MODE,
                        EnumValidator.in(InsertMode.values()),
                        ConfigDef.Importance.HIGH,
                        "sink mode, available mode is: INSERT, UPDATE, DELETE, default is " +
                                DEFAULT_INSERT_MODE)
                .define(CONNECT_PRIMARY_KEY,
                        ConfigDef.Type.STRING,
                        null,
                        ConfigDef.Importance.HIGH,
                        "primaryKey field for kafka data")
                .define(CONNECT_SRC_KEY,
                        ConfigDef.Type.STRING,
                        null,
                        ConfigDef.Importance.HIGH,
                        "src key field for kafka data")
                .define(CONNECT_DST_KEY,
                        ConfigDef.Type.STRING,
                        null,
                        ConfigDef.Importance.HIGH,
                        "dst key field for kafka data")
                .define(CONNECT_KAFKA_NODE_PROPERTIES,
                        ConfigDef.Type.LIST,
                        null,
                        ConfigDef.Importance.HIGH,
                        "property name list for node type in kafka")
                .define(CONNECT_KAFKA_EDGE_PROPERTIES,
                        ConfigDef.Type.LIST,
                        null,
                        ConfigDef.Importance.HIGH,
                        "property name list for edge type in kafka")
                .define(CONNECT_NEBULA_NODE_PROPERTIES,
                        ConfigDef.Type.LIST,
                        null,
                        ConfigDef.Importance.HIGH,
                        "property name list for node type in Nebula")
                .define(CONNECT_NEBULA_EDGE_PROPERTIES,
                        ConfigDef.Type.LIST,
                        null,
                        ConfigDef.Importance.HIGH,
                        "property name list for node type in Nebula")
                .define(CONNECT_CONNECT_TIMEOUT,
                        ConfigDef.Type.INT,
                        DEFAULT_SINK_CONNECT_TIMEOUT,
                        ConfigDef.Importance.LOW,
                        "connect timeout for connection between client and server")
                .define(CONNECT_REQUEST_TIMEOUT,
                        ConfigDef.Type.INT,
                        DEFAULT_SINK_REQUEST_TIMEOUT,
                        ConfigDef.Importance.LOW,
                        "request timeout for insert query's response time"
                )
                .define(CONNECT_SINK_PARTITION,
                        ConfigDef.Type.INT,
                        DEFAULT_CONNECT_SINK_PARTITION,
                        ConfigDef.Importance.LOW,
                        "max session number for nebula client")
                .define(CONNECT_SINK_RETRY_TIMES,
                        ConfigDef.Type.INT,
                        DEFAULT_SINK_RETRY_TIMES,
                        ConfigDef.Importance.LOW,
                        "retry times for failed insert query, default is " + DEFAULT_SINK_RETRY_TIMES)
                .define(CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY,
                        ConfigDef.Type.INT,
                        DEFAULT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY,
                        ConfigDef.Importance.LOW,
                        "interval time between each retry execution, default is " +
                                DEFAULT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY)
                .define(CONNECT_SINK_RECONNECT,
                        ConfigDef.Type.BOOLEAN,
                        DEFAULT_SINK_CONNECT_RECONNECT,
                        ConfigDef.Importance.MEDIUM,
                        "whether reconnect when session or connection is broken, default is " +
                                DEFAULT_SINK_CONNECT_RECONNECT)
                .define(CONNECT_HEALTH_CHECK_TIME,
                        ConfigDef.Type.INT,
                        3000,
                        ConfigDef.Importance.LOW,
                        "schedule time for checking the health of session, default is 3000")
                .define(CONNECT_BLOCK_WHEN_EXHAUSTED,
                        ConfigDef.Type.BOOLEAN,
                        true,
                        ConfigDef.Importance.LOW,
                        "if block when session in pool is exhausted, false to throw exception, " +
                                "true to wait. default is false")
                .define(CONNECT_MAX_WAIT_TIME,
                        ConfigDef.Type.INT,
                        -1,
                        ConfigDef.Importance.LOW,
                        "the max wait time if `blockWhenExhausted` is true, if value less than 0," +
                                " always wait. default is -1")
                .define(CONNECT_STRICTLY_SERVER_HEALTHY,
                        ConfigDef.Type.BOOLEAN,
                        false,
                        ConfigDef.Importance.LOW,
                        "if need all servers are strictly healthy, if true, all addresses must be" +
                                " " +
                                "available, if false, at least one address is available")
                .define(CONNECT_BATCH_SIZE,
                        ConfigDef.Type.INT,
                        DEFAULT_BATCH_SIZE,
                        ConfigDef.Importance.HIGH,
                        "write batch size for Nebula, default is " + DEFAULT_BATCH_SIZE);
    }


    /**
     * check the validation for configurations
     */
    public void check() {
        if (user == null || user.isEmpty()) {
            throw new IllegalArgumentException("user cannot be blank.");
        }
        if (passwd == null || passwd.isEmpty()) {
            throw new IllegalArgumentException("passwd cannot be blank.");
        }

    }


    private static class EnumValidator implements ConfigDef.Validator {
        private final List<String> canonicalValues;
        private final Set<String> validValues;

        private EnumValidator(List<String> canonicalValues, Set<String> validValues) {
            this.canonicalValues = canonicalValues;
            this.validValues = validValues;
        }

        public static <E> EnumValidator in(E[] enumerators) {
            final List<String> canonicalValues = new ArrayList<>(enumerators.length);
            final Set<String> validValues = new HashSet<>(enumerators.length * 2);
            for (E e : enumerators) {
                canonicalValues.add(e.toString().toLowerCase());
                validValues.add(e.toString().toUpperCase());
                validValues.add(e.toString().toLowerCase());
            }
            return new EnumValidator(canonicalValues, validValues);
        }

        @Override
        public void ensureValid(String key, Object value) {
            if (!validValues.contains(String.valueOf(value))) {
                throw new ConfigException(key, value, "Invalid enumerator " + value);
            }
        }

        @Override
        public String toString() {
            return canonicalValues.toString();
        }

    }
}
