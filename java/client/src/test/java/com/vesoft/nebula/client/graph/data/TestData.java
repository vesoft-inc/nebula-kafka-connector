/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;


import com.vesoft.nebula.Date;
import com.vesoft.nebula.Edge;
import com.vesoft.nebula.ExecutionOutcome;
import com.vesoft.nebula.ExecutionResponse;
import com.vesoft.nebula.GQLStatus;
import com.vesoft.nebula.LocalDatetime;
import com.vesoft.nebula.LocalTime;
import com.vesoft.nebula.NList;
import com.vesoft.nebula.NRecord;
import com.vesoft.nebula.Node;
import com.vesoft.nebula.Path;
import com.vesoft.nebula.ResultTable;
import com.vesoft.nebula.Row;
import com.vesoft.nebula.Value;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.stream.Collectors;
import org.junit.Assert;
import org.junit.Test;


public class TestData {

    @Test
    public void testNode() {
        try {
            Vertex vertex = new Vertex(getNode(1L));
            assert Objects.equals(vertex.getId(), 1L);
            assert (vertex.toString().startsWith("(1:0"));


            List<String> names = Arrays.asList("prop0", "prop1", "prop2", "prop3", "prop4");
            assert Objects.equals(
                    vertex.getColumnNames().stream().sorted().collect(Collectors.toList()),
                    names.stream().sorted().collect(Collectors.toList()));
            List<Value> propValues = Arrays.asList(new Value(Value.INT64VAL, 0L),
                    new Value(Value.INT64VAL, 1L),
                    new Value(Value.INT64VAL, 2L),
                    new Value(Value.INT64VAL, 3L),
                    new Value(Value.INT64VAL, 4L));
            for (String name : vertex.getProperties().keySet()) {
                if (name.equals("prop0")) {
                    assert vertex.getProperties().get(name).asLong() == 0;
                } else if (name.equals("prop1")) {
                    assert vertex.getProperties().get(name).asLong() == 1;
                } else if (name.equals("prop2")) {
                    assert vertex.getProperties().get(name).asLong() == 2;
                } else if (name.equals("prop3")) {
                    assert vertex.getProperties().get(name).asLong() == 3;
                } else if (name.equals("prop4")) {
                    assert vertex.getProperties().get(name).asLong() == 4;
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
            assert (false);
        }
    }

    @Test
    public void testRelationShip() {
        try {
            Edge edge = getEdge(101L, 102L, 10L);
            Relationship relationShip = new Relationship(edge);
            assert relationShip.getSrcId() == 101L;
            assert relationShip.getDstId() == 102L;
            // TODO assert Objects.equals(relationShip.edgeName(), "classmate");
            assert relationShip.getRank() == 10;

            // check keys
            List<String> names = Arrays.asList("prop0", "prop1", "prop2", "prop3", "prop4");
            assert Objects.equals(
                    relationShip.getColumnNames().stream().sorted().collect(Collectors.toList()),
                    names.stream().sorted().collect(Collectors.toList()));

            // check get values
            List<ValueWrapper> values = relationShip.getValues();
            assert values.get(0).isInt();
            ArrayList<Integer> longVals = new ArrayList<>();
            for (ValueWrapper val : values) {
                assert val.isInt();
                longVals.add(val.asInt());
            }
            List<Integer> expectVals = Arrays.asList(0, 1, 2, 3, 4);
            assert Objects.equals(expectVals.stream().sorted().collect(Collectors.toList()),
                    longVals.stream().sorted().collect(Collectors.toList()));

            // check properties
            HashMap<String, ValueWrapper> properties = relationShip.getProperties();
            assert properties.containsKey("prop0");
            assert properties.get("prop0").isInt();
            Assert.assertEquals(properties.get("prop0").asInt(), 0);
            assert properties.containsKey("prop1");
            assert properties.get("prop1").isInt();
            Assert.assertEquals(properties.get("prop1").asInt(), 1);
            assert properties.containsKey("prop2");
            assert properties.get("prop2").isInt();
            Assert.assertEquals(properties.get("prop2").asInt(), 2);
            assert properties.containsKey("prop3");
            assert properties.get("prop3").isInt();
            Assert.assertEquals(properties.get("prop3").asInt(), 3);
            assert properties.containsKey("prop4");
            assert properties.get("prop4").isInt();
            Assert.assertEquals(properties.get("prop4").asInt(), 4);
            assert properties.containsKey("prop4");
        } catch (Exception e) {
            e.printStackTrace();
            assert (false);
        }
    }


    @Test
    public void testResult() {
        try {
            ExecutionResponse resp = new ExecutionResponse();
            resp.setExecutionOutcome(new ExecutionOutcome(new GQLStatus("SUCCESS".getBytes()),
                    getDateset(), null)).setLatencyInUs(1000);

            ResultSet resultSet = new ResultSet(resp);
            assert resultSet.isSucceeded();
            assert resultSet.getGqlStatus().equals("SUCCESS");
            assert !resultSet.isEmpty();
            Assert.assertEquals(1000, resultSet.getLatency());
            assert resultSet.getPlanDesc() == null;
            List<String> expectColNames = Arrays.asList("col0_empty", "col1_bool", "col2_int64",
                    "col3_int32", "col4_int16", "col5_int8", "col6_float", "col7_double",
                    "col8_string", "col9_list", "col10_vertex", "col11_edge");
            assert Objects.equals(resultSet.getColumnNames(), expectColNames);

            assert resultSet.getRows().size() == 1;
            ResultSet.Record record = resultSet.rowValues(0);
            assert record.size() == 12;
            assert record.get(0).isEmpty();

            assert record.get(1).isBoolean();
            assert !record.get(1).asBoolean();

            assert record.get(2).isLong();
            assert record.get(2).asLong() == 1;

            assert record.get(3).isInt();
            assert record.get(3).asInt() == 2;

            assert record.get(4).isShort();
            assert record.get(4).asShort() == 3;

            assert record.get(5).isByte();
            assert record.get(5).asByte() == 4;

            assert record.get(6).isDouble();
            assert Double.compare(record.get(6).asDouble(), 10.01) == 0;

            assert record.get(7).isDouble();
            assert Double.compare(record.get(7).asDouble(), 20.01) == 0;

            assert record.get(8).isString();
            assert Objects.equals("value1", record.get(8).asString());

            Assert.assertArrayEquals(
                    record.get(9).asList().stream().map(ValueWrapper::asLong).toArray(),
                    Arrays.asList((long) 1, (long) 2).toArray());


            assert record.get(10).isNode();
            assert Objects.equals(record.get(10).asNode(),
                    new Vertex(getNode(100)));

            assert record.get(11).isEdge();
            assert Objects.equals(record.get(11).asEdge(),
                    new Relationship(getEdge(101L, 102L, 0)));
        } catch (Exception e) {
            e.printStackTrace();
            assert (false);
        }
    }

    @Test
    public void testToString() {
        try {
            // test node
            ValueWrapper valueWrapper = new ValueWrapper(
                    new Value(Value.NODEVAL, getSimpleNode()), "utf-8");
            String expectString =
                    "(1:0 {prop: 100})";
            Assert.assertEquals(expectString, valueWrapper.asNode().toString());

            // test relationship
            valueWrapper = new ValueWrapper(
                    new Value(Value.EDGEVAL, getSimpleEdge(false)), "utf-8");
            expectString = "(1)-[:1@10{edge_prop: 100}]->(2)";
            Assert.assertEquals(expectString, valueWrapper.asEdge().toString());

            valueWrapper = new ValueWrapper(
                    new Value(Value.EDGEVAL, getSimpleEdge(true)), "utf-8");
            expectString = "(2)-[:-1@10{edge_prop: 100}]->(1)";
            Assert.assertEquals(expectString, valueWrapper.asEdge().toString());

            // test local time
            valueWrapper = new ValueWrapper(new Value(Value.LOCALTIMEVAL, getSimpleLocalTime()),
                    "utf-8");
            expectString = "12:20:15.000030";
            Assert.assertEquals(expectString, valueWrapper.asLocalTime().toString());

            // test local datetime
            valueWrapper = new ValueWrapper(new Value(Value.LOCALDATETIMEVAL,
                    getSimpleLocalDateTime()), "utf-8");
            expectString = "2023-01-01T12:20:15.000030";
            Assert.assertEquals(expectString, valueWrapper.asLocalDateTime().toString());

            // test date
            valueWrapper = new ValueWrapper(new Value(Value.DATEVAL, getSimpleDate()), "utf-8");
            expectString = "2023-01-01";
            Assert.assertEquals(expectString, valueWrapper.asDate().toString());
        } catch (Exception e) {
            e.printStackTrace();
            assert (false);
        }
    }

    @Test
    public void testPath() {
        ValueWrapper valueWrapper = new ValueWrapper(new Value(Value.PATHVAL, getPath()), "utf-8");
        String expectString = "81-[:1{prop2:2,prop1:1,prop0:0,prop4:4,prop3:3}]->82"
                + "-[:1{prop2:2,prop1:1,prop0:0,prop4:4,prop3:3}]->83";
        Assert.assertEquals(expectString, valueWrapper.asPath().toString());

        NPath path = valueWrapper.asPath();
        Assert.assertEquals(3, path.nodes().size());
        Assert.assertEquals(2, path.relationships().size());
        Assert.assertEquals(5, path.values().size());
    }

    @Test
    public void testMap() {
        ValueWrapper valueWrapper = new ValueWrapper(
                new Value(Value.RECORDVAL, getRowRecord()), "utf-8");
        List<String> expectMapKeys = Arrays.asList("prop1", "prop2", "prop3", "prop4", "prop5",
                "prop6");
        Map<String, ValueWrapper> values = valueWrapper.asMap();
        Assert.assertEquals(expectMapKeys.stream().sorted().collect(Collectors.toList()),
                values.keySet().stream().sorted().collect(Collectors.toList()));
        Assert.assertEquals((short) 1, values.get("prop1").asShort());
        Assert.assertEquals(2L, values.get("prop2").asLong());
        Assert.assertEquals("Tom", values.get("prop3").asString());
        Assert.assertTrue(values.get("prop4").asBoolean());
        Assert.assertEquals(1L, values.get("prop5").asNode().getId());
        Assert.assertEquals(1, values.get("prop5").asNode().getNodeTypeId());
        Assert.assertEquals("Bob",
                values.get("prop5").asNode().getProperties().keySet().toArray()[0]);
    }


    @Test
    public void testList() {
        ValueWrapper valueWrapper = new ValueWrapper(new Value(Value.LISTVAL, getList()), "utf-8");
        List<ValueWrapper> values = valueWrapper.asList();
        Assert.assertEquals(4, values.size());
    }


    // mock data
    private Node getNode(long vid) {

        Map<byte[], Value> props = new HashMap<>();
        for (int j = 0; j < 5; j++) {
            Value value = new Value();
            value.setInt64Val(j);
            props.put(String.format("prop%d", j).getBytes(), value);
        }
        short nodeType = 0;
        return new Node(vid, nodeType, props);
    }

    private Edge getEdge(long srcId, long dstId, long rank) {
        Map<byte[], Value> props = new HashMap<>();
        for (int i = 0; i < 5; i++) {
            Value value = new Value();
            value.setInt32Val(i);
            props.put(String.format("prop%d", i).getBytes(), value);
        }
        int edgeType = 1;
        return new Edge(srcId, dstId, edgeType, rank, props);
    }

    private Node getSimpleNode() {
        Map<byte[], Value> props1 = new HashMap<>();
        props1.put("prop".getBytes(), new Value(Value.INT64VAL, (long) 100));
        long nodeId = 1;
        short nodeType = 0;
        return new Node(nodeId, nodeType, props1);
    }

    private Edge getSimpleEdge(boolean isReverse) {
        Map<byte[], Value> props = new HashMap<>();
        props.put("edge_prop".getBytes(), new Value(Value.INT64VAL, (long) 100));
        int type = 1;
        if (isReverse) {
            type = -1;
        }
        long srcId = 1;
        long dstId = 2;
        int edgeType = type;
        long rank = 10;
        return new Edge(srcId, dstId, edgeType, rank, props);
    }

    private NRecord getRowRecord() {
        Map<byte[], Value> values = new HashMap<>();
        values.put("prop1".getBytes(), new Value(Value.INT8VAL, (byte) 1));
        values.put("prop2".getBytes(), new Value(Value.INT64VAL, 2L));
        values.put("prop3".getBytes(), new Value(Value.STRINGVAL, "Tom".getBytes()));
        values.put("prop4".getBytes(), new Value(Value.BOOLVAL, true));

        Map<byte[], Value> properties = new HashMap<>();
        properties.put("name".getBytes(), new Value(Value.STRINGVAL, "Bob".getBytes()));
        Node node = new Node(1L, (short) 1, properties);
        values.put("prop5".getBytes(), new Value(Value.NODEVAL, node));
        values.put("prop6".getBytes(), new Value(Value.DOUBLEVAL, 12.1));
        return new NRecord(values);
    }

    private NList getList() {
        List<Value> values = new ArrayList<>();
        values.add(new Value(Value.STRINGVAL, "player".getBytes()));
        values.add(new Value(Value.INT64VAL, 10L));
        values.add(new Value(Value.DOUBLEVAL, 10.1));
        values.add(new Value(Value.BOOLVAL, true));
        return new NList(values);
    }

    private LocalTime getSimpleLocalTime() {
        LocalTime time = new LocalTime((byte) 12, (byte) 20, (byte) 15, 30);
        return time;
    }

    private LocalDatetime getSimpleLocalDateTime() {
        LocalDatetime datetime = new LocalDatetime((short) 2023, (byte) 1, (byte) 1, (byte) 12,
                (byte) 20, (byte) 15, 30);
        return datetime;
    }

    private Date getSimpleDate() {
        Date date = new Date((short) 2023, (byte) 1, (byte) 1);
        return date;
    }

    private Path getPath() {
        List<Value> values = new ArrayList<>();
        values.add(new Value(Value.NODEVAL, getNode(81)));
        values.add(new Value(Value.EDGEVAL, getEdge(81, 82, 1)));
        values.add(new Value(Value.NODEVAL, getNode(82)));
        values.add(new Value(Value.EDGEVAL, getEdge(82, 83, 1)));
        values.add(new Value(Value.NODEVAL, getNode(83)));
        Path path = new Path(values);
        return path;
    }

    private ResultTable getDateset() {
        final ArrayList<Value> list = new ArrayList<>();
        list.add(new Value(Value.INT64VAL, 1L));
        list.add(new Value(Value.INT64VAL, 2L));

        final HashMap<byte[], Value> map = new HashMap();
        map.put("key1".getBytes(), new Value(Value.INT64VAL, 1L));
        map.put("key2".getBytes(), new Value(Value.INT64VAL, 2L));

        final Row row = new Row(Arrays.asList(
                new Value(),
                new Value(Value.BOOLVAL, false),
                new Value(Value.INT64VAL, 1L),
                new Value(Value.INT32VAL, 2),
                new Value(Value.INT16VAL, (short) 3),
                new Value(Value.INT8VAL, (byte) 4),
                new Value(Value.FLOATVAL, 10.01),
                new Value(Value.DOUBLEVAL, 20.01),
                new Value(Value.STRINGVAL, "value1".getBytes()),
                new Value(Value.LISTVAL, new NList(list)),
                new Value(Value.NODEVAL, getNode(100L)),
                new Value(Value.EDGEVAL, getEdge(101L, 102L, 0L))));
        // TODO add NMap datatype, depends on the server finish

        final List<byte[]> columnNames = Arrays.asList(
                "col0_empty".getBytes(),
                "col1_bool".getBytes(),
                "col2_int64".getBytes(),
                "col3_int32".getBytes(),
                "col4_int16".getBytes(),
                "col5_int8".getBytes(),
                "col6_float".getBytes(),
                "col7_double".getBytes(),
                "col8_string".getBytes(),
                "col9_list".getBytes(),
                "col10_vertex".getBytes(),
                "col11_edge".getBytes());
        return new ResultTable(columnNames, Collections.singletonList(row));
    }
}
