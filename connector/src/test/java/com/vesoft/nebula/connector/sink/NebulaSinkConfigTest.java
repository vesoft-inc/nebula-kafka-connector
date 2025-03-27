
package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_BATCH_SIZE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_DST_PKS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_ADDRESS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_DATA_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_EDGE_TYPE_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NODE_TYPE_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_PASSWD;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_USER;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_KAFKA_EDGE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_KAFKA_NODE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_NEBULA_EDGE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_NEBULA_NODE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_PRIMARY_KEYS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_REQUEST_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_MODE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_PARTITION;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_RETRY_TIMES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SRC_PKS;
import static junit.framework.TestCase.assertEquals;

import com.vesoft.nebula.connector.config.NebulaConnectConfigName;
import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import org.apache.kafka.common.config.ConfigException;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

public class NebulaSinkConfigTest {
    private Map<String, Object>     props = new HashMap<>();
    private NebulaSinkConnectConfig config;

    @Before
    public void beforeEach() {
        // add the minimum settings only
        props.put(CONNECT_GRAPH_ADDRESS, "127.0.0.1:9669");
        props.put(CONNECT_GRAPH_PASSWD, "nebula");
        props.put(CONNECT_GRAPH_NAME, "nba");
        props.put(CONNECT_GRAPH_DATA_TYPE, "EDGE");
        props.put(CONNECT_GRAPH_NODE_TYPE_NAME, "person");
        props.put(CONNECT_PRIMARY_KEYS, "id");
    }

    @After
    public void afterEach() {
        props.clear();
        config = null;
    }

    @Test
    public void testSucceedToCreateConfig() {
        createConfig();
    }

    @Test(expected = ConfigException.class)
    public void testFailedToCreateConfigWithoutGraphServers() {
        props.remove(CONNECT_GRAPH_ADDRESS);
        createConfig();
    }

    @Test(expected = ConfigException.class)
    public void testFailedToCreateConfigWithEmptyGraphName() {
        props.remove(CONNECT_GRAPH_NAME);
        createConfig();
    }

    @Test
    public void testCreateConfigWithMinimalConfigs() {
        createConfig();
        assertEquals("nba", config.graphName);
        assertEquals(NebulaSinkConnectConfig.DataType.EDGE.name(), config.dataType.name());
        assertEquals("person", config.graphNodeType);
        assert (config.primaryKeys.contains("id"));
        assertEquals(2000, config.sinkBatchSize);
        assertEquals("root", config.user);
        assertEquals("nebula", config.authOptions.get("password"));
    }

    @Test
    public void testCreateConfigWithAllConfigs() {
        props.put(CONNECT_GRAPH_USER, "test");
        props.put(CONNECT_GRAPH_PASSWD, "12345");
        props.put(CONNECT_SINK_MODE, "UPDATE");
        props.put(CONNECT_GRAPH_EDGE_TYPE_NAME, "friend");
        props.put(CONNECT_SRC_PKS, "src");
        props.put(CONNECT_DST_PKS, "dst");
        props.put(CONNECT_KAFKA_NODE_PROPERTIES, Arrays.asList("name", "age"));
        props.put(CONNECT_KAFKA_EDGE_PROPERTIES, Arrays.asList("type", "duration"));
        props.put(CONNECT_NEBULA_NODE_PROPERTIES, Arrays.asList("Name", "Age"));
        props.put(CONNECT_NEBULA_EDGE_PROPERTIES, Arrays.asList("Type", "Duration"));
        props.put(CONNECT_REQUEST_TIMEOUT, 5000);
        props.put(CONNECT_SINK_PARTITION, 2000);
        props.put(CONNECT_SINK_RETRY_TIMES, 5);
        props.put(CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY, 2000);
        props.put(CONNECT_BATCH_SIZE, 1000);
        createConfig();
        assertEquals("test", config.user);
        assertEquals("12345", config.authOptions.get("password"));
        assertEquals("UPDATE", config.sinkMode.name());
        assert (config.primaryKeys.contains("id"));
        assert (config.srcKeys.contains("src"));
        assert (config.dstKeys.contains("dst"));
        assert (config.kafkaNodePropertyNames.containsAll(Arrays.asList("name", "age")));
        assert (config.kafkaEdgePropertyNames.containsAll(Arrays.asList("type", "duration")));
        assert (config.nebulaNodePropertyNames.containsAll(Arrays.asList("Name", "Age")));
        assert (config.nebulaEdgePropertyNames.containsAll(Arrays.asList("Type", "Duration")));
        assertEquals(5000, config.requestTimeout);
        assertEquals(2000, config.sinkPartition);
        assertEquals(5, config.retryTimes);
        assertEquals(2000, config.intervalTimeMill);
        assertEquals(1000, config.sinkBatchSize);
    }

    private void createConfig() {
        config = new NebulaSinkConnectConfig(props);
    }
}
