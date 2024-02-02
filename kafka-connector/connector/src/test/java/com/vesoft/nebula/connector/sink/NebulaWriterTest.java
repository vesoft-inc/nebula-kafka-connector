package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_BATCH_SIZE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_DST_KEY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_ADDRESS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_EDGE_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NODE_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_KAFKA_EDGE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_KAFKA_NODE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_NEBULA_EDGE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_NEBULA_NODE_PROPERTIES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_PRIMARY_KEY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SRC_KEY;
import static org.junit.Assert.assertEquals;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.mock;
import static org.powermock.api.mockito.PowerMockito.when;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import com.vesoft.nebula.connector.connection.NebulaGraphProvider;
import com.vesoft.nebula.graph.ExecutionOutcome;
import com.vesoft.nebula.graph.ExecutionResponse;
import com.vesoft.nebula.graph.GQLStatus;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.SchemaBuilder;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.sink.SinkRecord;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

@RunWith(PowerMockRunner.class)
@PrepareForTest(NebulaWriter.class)
public class NebulaWriterTest {
    private Map<String, Object> props = new HashMap<>();
    private NebulaSinkConnectConfig config;
    @Mock
    NebulaGraphProvider mockProvider;

    NebulaWriter writer;

    @Before
    public void setUp() throws Exception {
        props.put(CONNECT_GRAPH_ADDRESS, "127.0.0.1:9669");
        props.put(CONNECT_GRAPH_NAME, "nba");
        props.put(CONNECT_GRAPH_TYPE, "BOTH");
        // node config
        props.put(CONNECT_GRAPH_NODE_TYPE, "person");
        props.put(CONNECT_GRAPH_EDGE_TYPE, "friend");
        props.put(CONNECT_PRIMARY_KEY, "id");
        props.put(CONNECT_NEBULA_NODE_PROPERTIES, Arrays.asList("Name", "Age"));
        props.put(CONNECT_KAFKA_NODE_PROPERTIES, Arrays.asList("name", "age"));
        // edge config
        props.put(CONNECT_SRC_KEY, "src");
        props.put(CONNECT_DST_KEY, "dst");
        props.put(CONNECT_NEBULA_EDGE_PROPERTIES, Arrays.asList("Degree", "Type"));
        props.put(CONNECT_KAFKA_EDGE_PROPERTIES, Arrays.asList("degree", "type"));

        props.put(CONNECT_BATCH_SIZE, 1);
        config = new NebulaSinkConnectConfig(props);
        MockitoAnnotations.initMocks(this);
        PowerMockito.whenNew(NebulaGraphProvider.class).withAnyArguments().thenReturn(mockProvider);
        writer = new NebulaWriter(config);
    }

    @After
    public void tearDown() throws Exception {
        writer.close();
    }

    @Test
    public void testWrite() {
        ExecutionOutcome outcome = new ExecutionOutcome(new GQLStatus("SUCCESS".getBytes()));
        ExecutionResponse response = new ExecutionResponse(outcome, 1000);
        ResultSet result = new ResultSet(response);
        try {
            PowerMockito.when(mockProvider.execute(anyString())).thenReturn(result);
            writer.write(getSinkRecords());
        } catch (IOErrorException | NoValidSessionException | InterruptedException e) {
            e.printStackTrace();
            assert (false);
        }
    }


    @Test
    public void testGetNode() {
        Map<String, Object> properties = new HashMap<>();
        properties.put("id", "1");
        properties.put("name", "Bob");
        properties.put("age", 10);

        try {
            Method getNodeMethod = writer.getClass().getDeclaredMethod("getNode", Map.class);
            getNodeMethod.setAccessible(true);
            NebulaNode node = (NebulaNode) getNodeMethod.invoke(writer, properties);
            assertEquals("1", node.getVid());
            assert (node.getProperties().containsKey("Name"));
            assert (node.getProperties().containsKey("Age"));
            assert (node.getProperties().containsValue("Bob"));
            assert (node.getProperties().containsValue(10));
        } catch (NoSuchMethodException | InvocationTargetException | IllegalAccessException e) {
            e.printStackTrace();
            assert (false);
        }

    }

    @Test
    public void testGetEdge() {
        Map<String, Object> properties = new HashMap<>();
        properties.put("src", 1);
        properties.put("dst", 2);
        properties.put("degree", 10);
        properties.put("type", "friend");
        try {
            Method getEdgeMethod = writer.getClass().getDeclaredMethod("getEdge", Map.class);
            getEdgeMethod.setAccessible(true);
            NebulaEdge edge = (NebulaEdge) getEdgeMethod.invoke(writer, properties);
            assertEquals("1", edge.getSrcPk());
            assertEquals("2", edge.getDstPk());
            assert (edge.getProperties().containsKey("Degree"));
            assert (edge.getProperties().containsKey("Type"));
            assert (edge.getProperties().containsValue(10));
            assert (edge.getProperties().containsValue("friend"));
        } catch (NoSuchMethodException | InvocationTargetException | IllegalAccessException e) {
            e.printStackTrace();
            assert (false);
        }
    }

    /**
     * TODO: now the nebula schema is mock as all data type is String,
     * so the property values in expect is enclosed in double quotes.
     * When schema is real ,change the expect too.
     */
    @Test
    public void testGetNodeStatement() {
        List<NebulaNode> nodes = new ArrayList<>();
        Map<String, Object> properties1 = new HashMap<>();
        properties1.put("name", "A");
        properties1.put("age", 10);
        properties1.put("weight", 12);
        NebulaNode node1 = new NebulaNode("1", properties1);
        Map<String, Object> properties2 = new HashMap<>();
        properties2.put("name", "B");
        properties2.put("age", 20);
        properties2.put("weight", 22);
        NebulaNode node2 = new NebulaNode("2", properties2);
        nodes.add(node1);
        nodes.add(node2);

        try {
            Method getNodeStatementMethod = writer.getClass().getDeclaredMethod("getNodeStatement"
                    , List.class);
            getNodeStatementMethod.setAccessible(true);
            String statement = (String) getNodeStatementMethod.invoke(writer, nodes);
            String expect = ("USE nba INSERT NODE `person` ({id:1,name:\"A\",age:\"10\",weight:\"12\"})," +
                    "({id:2,name:\"B\",age:\"20\",weight:\"22\"})");
            String statementStr = statement.chars().sorted().collect(StringBuilder::new,
                    StringBuilder::appendCodePoint, StringBuilder::append).toString();
            String expectStr = expect
                    .chars()
                    .sorted()
                    .collect(StringBuilder::new,
                            StringBuilder::appendCodePoint,
                            StringBuilder::append)
                    .toString();
            assertEquals(statementStr, expectStr);
        } catch (NoSuchMethodException | InvocationTargetException | IllegalAccessException e) {
            e.printStackTrace();
            assert (false);
        }
    }

    @Test
    public void testGetEdgeStatement() {
        List<NebulaEdge> edges = new ArrayList<>();
        Map<String, Object> properties1 = new HashMap<>();
        properties1.put("degree", 1);
        properties1.put("type", "friend");
        NebulaEdge edge1 = new NebulaEdge("1", "2", properties1);
        Map<String, Object> properties2 = new HashMap<>();
        properties2.put("degree", 2);
        properties2.put("type", "friend");
        NebulaEdge edge2 = new NebulaEdge("2", "3", properties2);
        edges.add(edge1);
        edges.add(edge2);

        try {
            Method getEdgeStatementMethod = writer.getClass().getDeclaredMethod("getEdgeStatement", List.class);
            getEdgeStatementMethod.setAccessible(true);
            String statement = (String) getEdgeStatementMethod.invoke(writer, edges);
            String expect = ("USE nba INSERT EDGE `friend` ({id:1})-[{degree:\"1\",type:\"friend\"}]->({id:2})," +
                    "({id:2})-[{degree:\"2\",type:\"friend\"}]->({id:3})");
            String statementStr = statement.chars().sorted().collect(StringBuilder::new,
                    StringBuilder::appendCodePoint, StringBuilder::append).toString();
            String expectStr = expect
                    .chars()
                    .sorted()
                    .collect(StringBuilder::new,
                            StringBuilder::appendCodePoint,
                            StringBuilder::append)
                    .toString();
            assertEquals(statementStr, expectStr);
        } catch (NoSuchMethodException | InvocationTargetException | IllegalAccessException e) {
            e.printStackTrace();
            assert (false);
        }
    }

    private List<SinkRecord> getSinkRecords() {
        final Schema SCHEMA = SchemaBuilder.struct()
                .field("id", Schema.STRING_SCHEMA)
                .field("src", Schema.STRING_SCHEMA)
                .field("dst", Schema.STRING_SCHEMA)
                .field("name", Schema.STRING_SCHEMA)
                .field("age", Schema.INT32_SCHEMA)
                .field("gender", Schema.STRING_SCHEMA)
                .field("degree", Schema.INT32_SCHEMA)
                .field("type", Schema.STRING_SCHEMA)
                .build();
        Struct struct1 = new Struct(SCHEMA)
                .put("id", "1")
                .put("src", "1")
                .put("dst", "2")
                .put("name", "Tom")
                .put("age", 10)
                .put("gender", "male")
                .put("degree", 20)
                .put("type", "friend");

        Struct struct2 = new Struct(SCHEMA)
                .put("id", "2")
                .put("src", "2")
                .put("dst", "3")
                .put("name", "Bob")
                .put("age", 20)
                .put("gender", "female")
                .put("degree", 50)
                .put("type", "friend");

        SinkRecord sinkRecord1 = new SinkRecord("test", 0, null, null, SCHEMA, struct1, 2);
        SinkRecord sinkRecord2 = new SinkRecord("test", 0, null, null, SCHEMA, struct2, 2);
        SinkRecord nullSinkRecord = new SinkRecord("test", 0, null, null, SCHEMA, null, 3);
        List<SinkRecord> records = new ArrayList<>();
        records.add(sinkRecord1);
        records.add(sinkRecord2);
        records.add(nullSinkRecord);
        return records;
    }
}
