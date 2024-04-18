/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_BATCH_SIZE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_BLOCK_WHEN_EXHAUSTED;
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
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_RETRY_TIMES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SRC_KEY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_STRICTLY_SERVER_HEALTHY;
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
    private Map<String, Object> props = new HashMap<>();
    private NebulaSinkConnectConfig config;

    @Before
    public void beforeEach() {
        // add the minimum settings only
        props.put(CONNECT_GRAPH_ADDRESS, "127.0.0.1:9669");
        props.put(CONNECT_GRAPH_NAME, "nba");
        props.put(CONNECT_GRAPH_TYPE, "EDGE");
        props.put(CONNECT_GRAPH_NODE_TYPE, "person");
        props.put(CONNECT_PRIMARY_KEY, "id");
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
        props.put(NebulaConnectConfigName.CONNECT_GRAPH_NAME, "");
        createConfig();
    }

    @Test
    public void testCreateConfigWithMinimalConfigs() {
        createConfig();
        assertEquals("nba", config.graphName);
        assertEquals(NebulaSinkConnectConfig.DataType.EDGE.name(), config.dataType.name());
        assertEquals("person", config.graphNodeType);
        assertEquals("id", config.primaryKey);
        assertEquals(2000, config.sinkBatchSize);
        assertEquals("root", config.user);
        assertEquals("nebula", config.passwd);
    }

    @Test
    public void testCreateConfigWithAllConfigs() {
        props.put(CONNECT_GRAPH_USER, "test");
        props.put(CONNECT_GRAPH_PASSWD, "12345");
        props.put(CONNECT_SINK_MODE, "UPDATE");
        props.put(CONNECT_GRAPH_EDGE_TYPE, "friend");
        props.put(CONNECT_SRC_KEY, "src");
        props.put(CONNECT_DST_KEY, "dst");
        props.put(CONNECT_KAFKA_NODE_PROPERTIES, Arrays.asList("name", "age"));
        props.put(CONNECT_KAFKA_EDGE_PROPERTIES, Arrays.asList("type", "duration"));
        props.put(CONNECT_NEBULA_NODE_PROPERTIES, Arrays.asList("Name", "Age"));
        props.put(CONNECT_NEBULA_EDGE_PROPERTIES, Arrays.asList("Type", "Duration"));
        props.put(CONNECT_REQUEST_TIMEOUT, 5000);
        props.put(CONNECT_SINK_PARTITION, 2000);
        props.put(CONNECT_SINK_RETRY_TIMES, 5);
        props.put(CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY, 2000);
        props.put(CONNECT_HEALTH_CHECK_TIME, 5000);
        props.put(CONNECT_BLOCK_WHEN_EXHAUSTED, false);
        props.put(CONNECT_MAX_WAIT_TIME, 10000);
        props.put(CONNECT_STRICTLY_SERVER_HEALTHY, true);
        props.put(CONNECT_BATCH_SIZE, 1000);
        createConfig();
        assertEquals("test", config.user);
        assertEquals("12345", config.passwd);
        assertEquals("UPDATE", config.sinkMode.name());
        assertEquals("id", config.primaryKey);
        assertEquals("src", config.srcKey);
        assertEquals("dst", config.dstKey);
        assert (config.kafkaNodePropertyNames.containsAll(Arrays.asList("name", "age")));
        assert (config.kafkaEdgePropertyNames.containsAll(Arrays.asList("type", "duration")));
        assert (config.nebulaNodePropertyNames.containsAll(Arrays.asList("Name", "Age")));
        assert (config.nebulaEdgePropertyNames.containsAll(Arrays.asList("Type", "Duration")));
        assertEquals(5000, config.requestTimeout);
        assertEquals(2000, config.sinkPartition);
        assertEquals(5, config.retryTimes);
        assertEquals(2000, config.intervalTimeMill);
        assertEquals(5000, config.healthCheckTime);
        assertEquals(false, config.blockWithExhausted);
        assertEquals(10000, config.maxWaitTimeWhenSessionExhausted);
        assertEquals(true, config.strictlyServerHealth);
        assertEquals(1000, config.sinkBatchSize);
    }

    private void createConfig() {
        config = new NebulaSinkConnectConfig(props);
    }
}
