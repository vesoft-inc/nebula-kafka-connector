
package com.vesoft.nebula.connector.config;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_BATCH_SIZE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_CONNECT_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_ADDRESS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_DATA_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_EDGE_TYPE_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NODE_TYPE_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_PASSWD;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_USER;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_MAX_TASK;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_NEBULA_EDGE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_NEBULA_NODE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_NEBULA_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_REQUEST_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_RETRY_TIMES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SOURCE_EDGE_TYPES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SOURCE_NODE_TYPES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SOURCE_TOPIC_PREFIX;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_BATCH_SIZE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECTOR_CONNECT_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECTOR_REQUEST_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECT_GRAPH_USER;

import com.vesoft.nebula.connector.util.ConfigUtils;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.ConfigException;

public class NebulaSourceConnectConfig extends AbstractConfig implements Serializable {

    public final String connectorName;
    public final String graphServers;
    public final String user;
    public final Map<String, Object> authOptions;
    public final String graphName;
    public final NebulaConnectDataTypeEnum dataType;
    public final List<String> graphNodeTypes = new ArrayList<>();
    public final List<String> graphEdgeTypes = new ArrayList<>();
    public List<String> nebulaNodePropertyNames = new ArrayList<>();
    public List<String> nebulaEdgePropertyNames = new ArrayList<>();
    public final int requestTimeout;
    public final int retryTimes;
    public final int intervalTimeMill;

    public final int batchSize;

    public final int maxTask;
    public final String topicPrefix;

    public NebulaSourceConnectConfig(Map<?, ?> props) {
        super(configDef(), props);
        connectorName = ConfigUtils.connectorName(props);
        graphServers = getString(CONNECT_GRAPH_ADDRESS);
        user = getString(CONNECT_GRAPH_USER);
        String passwd = getString(CONNECT_GRAPH_PASSWD);
        authOptions = ConfigUtils.authOptions(props);
        if (passwd != null && !passwd.isEmpty()) {
            authOptions.put("password", passwd);
        }
        graphName = getString(CONNECT_GRAPH_NAME);
        dataType = NebulaConnectDataTypeEnum.getDataType(getString(CONNECT_GRAPH_DATA_TYPE));
        Collections.addAll(graphNodeTypes, getString(CONNECT_SOURCE_NODE_TYPES).split(","));
        Collections.addAll(graphEdgeTypes, getString(CONNECT_SOURCE_EDGE_TYPES).split(","));
        nebulaNodePropertyNames = getList(CONNECT_NEBULA_NODE_PROPERTIES);
        nebulaEdgePropertyNames = getList(CONNECT_NEBULA_EDGE_PROPERTIES);
        requestTimeout = getInt(CONNECT_REQUEST_TIMEOUT);
        retryTimes = getInt(CONNECT_SINK_RETRY_TIMES);
        intervalTimeMill = getInt(CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY);
        batchSize = getInt(CONNECT_BATCH_SIZE);
        maxTask = getInt(CONNECT_MAX_TASK);
        topicPrefix = getString(CONNECT_SOURCE_TOPIC_PREFIX);
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
                        new ConfigDef.NonEmptyString(),
                        ConfigDef.Importance.HIGH,
                        "NebulaGraph passwd")
                .define(CONNECT_GRAPH_DATA_TYPE,
                        ConfigDef.Type.STRING,
                        NebulaSourceConnectConfig.EnumValidator
                                .in(NebulaSinkConnectConfig.DataType.values()),
                        ConfigDef.Importance.HIGH,
                        "the data type, node or edge")
                .define(CONNECT_GRAPH_NODE_TYPE_NAME,
                        ConfigDef.Type.STRING,
                        new ConfigDef.NonEmptyString(),
                        ConfigDef.Importance.HIGH,
                        "Name of the node type")
                .define(CONNECT_GRAPH_EDGE_TYPE_NAME,
                        ConfigDef.Type.STRING,
                        new ConfigDef.NonEmptyString(),
                        ConfigDef.Importance.HIGH,
                        "Name of the edge type")
                .define(CONNECT_NEBULA_PROPERTIES,
                        ConfigDef.Type.LIST,
                        null,
                        ConfigDef.Importance.HIGH,
                        "property name list in Nebula")
                .define(CONNECT_CONNECT_TIMEOUT,
                        ConfigDef.Type.INT,
                        DEFAULT_CONNECTOR_CONNECT_TIMEOUT,
                        ConfigDef.Importance.LOW,
                        "connect timeout for connect to NebulaGraph")
                .define(CONNECT_REQUEST_TIMEOUT,
                        ConfigDef.Type.INT,
                        DEFAULT_CONNECTOR_REQUEST_TIMEOUT,
                        ConfigDef.Importance.LOW,
                        "request timeout for scan's response time")
                .define(CONNECT_BATCH_SIZE,
                        ConfigDef.Type.INT,
                        DEFAULT_BATCH_SIZE,
                        ConfigDef.Importance.HIGH,
                        "read batch size from Nebula, default is " + DEFAULT_BATCH_SIZE)
                .define(CONNECT_SOURCE_TOPIC_PREFIX,
                        ConfigDef.Type.STRING,
                        "nebula",
                        ConfigDef.Importance.HIGH,
                        "kafka topic name prefix.");
    }


    private static class EnumValidator implements ConfigDef.Validator {
        private final List<String> canonicalValues;
        private final Set<String> validValues;

        private EnumValidator(List<String> canonicalValues, Set<String> validValues) {
            this.canonicalValues = canonicalValues;
            this.validValues = validValues;
        }

        public static <E> NebulaSourceConnectConfig.EnumValidator in(E[] enumerators) {
            final List<String> canonicalValues = new ArrayList<>(enumerators.length);
            final Set<String> validValues = new HashSet<>(enumerators.length * 2);
            for (E e : enumerators) {
                canonicalValues.add(e.toString().toLowerCase());
                validValues.add(e.toString().toUpperCase());
                validValues.add(e.toString().toLowerCase());
            }
            return new NebulaSourceConnectConfig.EnumValidator(canonicalValues, validValues);
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
