/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.google.common.base.Charsets;
import com.google.protobuf.ByteString;
import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.proto.common.Date;
import com.vesoft.nebula.proto.common.Duration;
import com.vesoft.nebula.proto.common.Edge;
import com.vesoft.nebula.proto.common.LocalDatetime;
import com.vesoft.nebula.proto.common.LocalTime;
import com.vesoft.nebula.proto.common.Node;
import com.vesoft.nebula.proto.common.Path;
import com.vesoft.nebula.proto.common.Record;
import com.vesoft.nebula.proto.common.Status;
import com.vesoft.nebula.proto.common.Value;
import com.vesoft.nebula.proto.graph.ExecuteResponse;
import com.vesoft.nebula.proto.graph.ResultTable;
import com.vesoft.nebula.proto.graph.Row;
import com.vesoft.nebula.proto.graph.Summary;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.stream.Collectors;
import org.junit.Assert;
import org.junit.Test;

public class DataTest {

    @Test
    public void testNode() {
        try {
            Vertex vertex = new Vertex(getNode(1L));
            assert Objects.equals(vertex.getId(), 1L);
            assert (vertex.toString().startsWith("(1@person"));


            List<String> names = Arrays.asList("prop0", "prop1", "prop2", "prop3", "prop4");
            assert Objects.equals(
                    vertex.getColumnNames().stream().sorted().collect(Collectors.toList()),
                    names.stream().sorted().collect(Collectors.toList()));
            List<Value> propValues = Arrays.asList(
                    Value.newBuilder().setInt64Value(1L).build(),
                    Value.newBuilder().setInt16Value(2).build(),
                    Value.newBuilder().setInt16Value(3).build(),
                    Value.newBuilder().setInt8Value(4).build());
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
            List<ValueWrapper> values = relationShip.getPropertyValues();
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
            ResultTable resultTable = getDateset();
            ExecuteResponse response = ExecuteResponse
                    .newBuilder()
                    .setResult(resultTable)
                    .setSummary(Summary.newBuilder()
                            .setTotalServerTimeUs(1000).build())
                    .setStatus(Status.newBuilder()
                            .setCode(ByteString.copyFrom("00000", Charsets.UTF_8)).build())
                    .build();

            ResultSet resultSet = new ResultSet(response);
            assert resultSet.isSucceeded();
            assert resultSet.getErrorCode() == ErrorCode.SUCCESSFUL_COMPLETION;
            assert !resultSet.isEmpty();
            Assert.assertEquals(1000, resultSet.getLatency());
            List<String> expectColNames = Arrays.asList("col0_empty", "col1_bool", "col2_int64",
                    "col3_int32", "col4_int16", "col5_int8", "col6_float", "col7_double",
                    "col8_string", "col9_list", "col10_vertex", "col11_edge");
            assert Objects.equals(resultSet.getColumnNames(), expectColNames);

            assert resultSet.rowSize() == 1;
            assert resultSet.hasNext();
            ResultSet.Record record = resultSet.next();
            assert record.size() == 12;
            assert record.get(0).isEmpty();

            assert record.get(1).isBoolean();
            assert !record.get(1).asBoolean();

            assert record.get(2).isLong();
            assert record.get(2).asLong() == 1;

            assert record.get(3).isInt();
            assert record.get(3).asInt() == 2;

            assert record.get(4).isInt();
            assert record.get(4).asInt() == 3;

            assert record.get(5).isInt();
            assert record.get(5).asInt() == 4;

            assert record.get(6).isFloat();
            assert Math.abs(record.get(6).asFloat() - 10.01) < 0.001;

            assert record.get(7).isDouble();
            assert Math.abs(record.get(7).asDouble() - 20.01) < 0.001;

            assert record.get(8).isString();
            assert Objects.equals("value1", record.get(8).asString());

            Assert.assertArrayEquals(
                    record.get(9).asList().stream().map(ValueWrapper::asLong).toArray(),
                    Arrays.asList(10L, 11L, 12L, 13L).toArray());


            assert record.get(10).isNode();
            assert Objects.equals(record.get(10).asNode(),
                    new Vertex(getNode(1)));

            assert record.get(11).isEdge();
            assert Objects.equals(record.get(11).asEdge(),
                    new Relationship(getEdge(1, 2, 10)));
        } catch (Exception e) {
            e.printStackTrace();
            assert (false);
        }
    }

    @Test
    public void testToString() {
        try {
            // test node
            ValueWrapper valueWrapper = new ValueWrapper(Value
                    .newBuilder()
                    .setNodeValue(getSimpleNode(1))
                    .build());
            String expectString =
                    "(1@person:teacher{prop:Bob})";

            Assert.assertEquals(expectString, valueWrapper.asNode().toString());
            // test relationship
            valueWrapper = new ValueWrapper(Value
                    .newBuilder()
                    .setEdgeValue(getSimpleEdge(false, 1, 2))
                    .build());
            expectString = "(1)-[10@knows:knows1&knows2{edge_prop:100}]->(2)";
            Assert.assertEquals(expectString, valueWrapper.asEdge().toString());

            valueWrapper = new ValueWrapper(Value
                    .newBuilder()
                    .setEdgeValue(getSimpleEdge(true, 1, 2))
                    .build());
            expectString = "(1)-[10@knows:knows1&knows2{edge_prop:100}]->(2)";
            Assert.assertEquals(expectString, valueWrapper.asEdge().toString());

            // test local time
            valueWrapper = new ValueWrapper(Value
                    .newBuilder()
                    .setLocalTimeValue(getSimpleLocalTime())
                    .build());
            expectString = "12:20:15.000030";
            Assert.assertEquals(expectString, valueWrapper.asLocalTime().toString());

            // test local datetime
            valueWrapper = new ValueWrapper(Value
                    .newBuilder()
                    .setLocalDatetimeValue(getSimpleLocalDateTime())
                    .build());
            expectString = "2024-01-01T12:20:15.000030";
            Assert.assertEquals(expectString, valueWrapper.asLocalDateTime().toString());

            // test date
            valueWrapper = new ValueWrapper(Value
                    .newBuilder()
                    .setDateValue(getSimpleDate())
                    .build());
            expectString = "2024-01-01";
            Assert.assertEquals(expectString, valueWrapper.asDate().toString());
        } catch (Exception e) {
            e.printStackTrace();
            assert (false);
        }
    }

    @Test
    public void testDuration() {
        // test duration of time-based duration
        ValueWrapper valueWrapper = new ValueWrapper(Value
                .newBuilder()
                .setDurationValue(Duration.newBuilder()
                        .setIsMonthBased(false)
                        .setDay(1)
                        .setHour(2)
                        .setMinute(3)
                        .setSec(4)
                        .setMicrosec(5)
                        .build())
                .build());
        String expectString = "P1DT2H3M4.000005S";
        Assert.assertEquals(expectString, valueWrapper.asDuration().toString());

        // test duration with 5000 ms, test the number of digits after the decimal point
        valueWrapper = new ValueWrapper(Value
                .newBuilder()
                .setDurationValue(Duration.newBuilder()
                        .setIsMonthBased(false)
                        .setSec(4)
                        .setMicrosec(5000)
                        .build())
                .build());
        expectString = "PT4.005S";
        Assert.assertEquals(expectString, valueWrapper.asDuration().toString());

        // tet duration of time-based duration only with day
        valueWrapper = new ValueWrapper(Value
                .newBuilder()
                .setDurationValue(Duration.newBuilder()
                        .setIsMonthBased(false)
                        .setDay(1)
                        .build())
                .build());
        expectString = "P1D";
        Assert.assertEquals(expectString, valueWrapper.asDuration().toString());

        // test duration of month-based duration
        valueWrapper = new ValueWrapper(Value
                .newBuilder()
                .setDurationValue(Duration.newBuilder()
                        .setIsMonthBased(true)
                        .setYear(-1)
                        .setMonth(-1)
                        .build())
                .build());
        expectString = "P-1Y-1M";
        Assert.assertEquals(expectString, valueWrapper.asDuration().toString());
    }

    @Test
    public void testPath() {
        ValueWrapper valueWrapper = new ValueWrapper(
                Value.newBuilder().setPathValue(getPath()).build());
        String expectString = "(1@person:teacher{prop:Bob})-[10@knows:knows1&knows2{edge_prop:100}]"
                + "->(2@person:teacher{prop:Bob})"
                + "-[10@knows:knows1&knows2{edge_prop:100}]->(3@person:teacher{prop:Bob})";
        Assert.assertEquals(expectString, valueWrapper.asPath().toString());

        NPath path = valueWrapper.asPath();
        Assert.assertEquals(3, path.nodes().size());
        Assert.assertEquals(2, path.relationships().size());
        Assert.assertEquals(5, path.values().size());
    }

    @Test
    public void testReversePath() {
        ValueWrapper valueWrapper = new ValueWrapper(
                Value.newBuilder().setPathValue(getReversePath()).build());
        String expectString = "(1@person:teacher{prop:Bob})"
                + "<-[10@knows:knows1&knows2{edge_prop:100}]"
                + "-(2@person:teacher{prop:Bob})"
                + "<-[10@knows:knows1&knows2{edge_prop:100}]-(3@person:teacher{prop:Bob})";
        Assert.assertEquals(expectString, valueWrapper.asPath().toString());

        NPath path = valueWrapper.asPath();
        Assert.assertEquals(3, path.nodes().size());
        Assert.assertEquals(2, path.relationships().size());
        Assert.assertEquals(5, path.values().size());
    }

    @Test
    public void testRecord() {
        ValueWrapper valueWrapper = new ValueWrapper(Value
                .newBuilder()
                .setRecordValue(getRowRecord())
                .build());

        NRecord record = valueWrapper.asRecord();
        assert (!record.isEmpty());
        Assert.assertEquals(6, record.size());
        Assert.assertEquals(1, record.getValue("prop1").asInt());
        Assert.assertEquals(2L, record.getValue("prop2").asLong());
        Assert.assertEquals("Tom", record.getValue("prop3").asString());
        Assert.assertTrue(record.getValue("prop4").asBoolean());
        Assert.assertEquals(1, record.getValue("prop5").asNode().getId());
        Assert.assertEquals("person",
                record.getValue("prop5").asNode().getType());
        Assert.assertEquals("teacher",
                record.getValue("prop5").asNode().getLabels().get(0));
        Assert.assertEquals("prop",
                record.getValue("prop5").asNode().getProperties().keySet().toArray()[0]);

        Map<String, ValueWrapper> values = record.getValuesMap();
        List<String> expectMapKeys = Arrays.asList("prop1", "prop2", "prop3", "prop4", "prop5",
                "prop6");
        Assert.assertEquals(expectMapKeys.stream().sorted().collect(Collectors.toList()),
                values.keySet().stream().sorted().collect(Collectors.toList()));
    }

    @Test
    public void testEmptyRecord() {
        ValueWrapper valueWrapper = new ValueWrapper(Value
                .newBuilder()
                .setRecordValue(Record.newBuilder().build())
                .build());
        NRecord record = valueWrapper.asRecord();
        assert (record.isEmpty());
    }

    @Test
    public void testEqualRecord() {
        Map<String, Value> map = new HashMap<>();
        map.put("prop", Value.newBuilder().setInt32Value(1).build());

        ValueWrapper valueWrapper1 = new ValueWrapper(Value
                .newBuilder()
                .setRecordValue(Record.newBuilder().putAllValues(map).build())
                .build());
        NRecord record1 = valueWrapper1.asRecord();

        ValueWrapper valueWrapper2 = new ValueWrapper(Value
                .newBuilder()
                .setRecordValue(Record.newBuilder().putAllValues(map).build())
                .build());
        NRecord record2 = valueWrapper2.asRecord();
        Assert.assertEquals(record1, record2);
        Assert.assertEquals(record1.hashCode(), record2.hashCode());
    }


    @Test
    public void testList() {
        ValueWrapper valueWrapper = new ValueWrapper(Value
                .newBuilder()
                .setListValue(getList())
                .build());
        List<ValueWrapper> values = valueWrapper.asList();
        Assert.assertEquals(4, values.size());
    }


    // mock data
    private Node getNode(long vid) {

        Map<String, Value> props = new HashMap<>();
        for (int j = 0; j < 5; j++) {
            Value value = Value.newBuilder().setInt64Value(j).build();
            props.put(String.format("prop%d", j), value);
        }
        String nodeType = "person";
        Node.Builder builder = Node.newBuilder().setNodeId(vid).setType(nodeType);
        builder.putAllProperties(props);
        return builder.build();
    }

    private Edge getEdge(long srcId, long dstId, long rank) {
        String edgeType = "knows";
        String graphName = "test";
        Edge.Builder builder = Edge
                .newBuilder()
                .setType(edgeType)
                .setGraph(graphName)
                .setDirectionValue(1)
                .setSrcId(srcId)
                .setDstId(dstId)
                .setRank(rank);
        for (int i = 0; i < 5; i++) {
            Value value = Value.newBuilder().setInt32Value(i).build();
            builder.putProperties(String.format("prop%d", i), value);
        }
        return builder.build();
    }

    private Node getSimpleNode(long nodeId) {
        Map<String, Value> props = new HashMap<>();
        props.put("prop", Value.newBuilder()
                .setStringValue(ByteString.copyFrom("Bob", Charsets.UTF_8)).build());
        String nodeType = "person";

        Node node = Node.newBuilder()
                .setNodeId(nodeId)
                .setType(nodeType)
                .addLabels("teacher")
                .setGraph("test")
                .putAllProperties(props)
                .build();
        return node;
    }

    private Edge getSimpleEdge(boolean isReverse, long srcId, long dstId) {
        Map<String, Value> props = new HashMap<>();
        props.put("edge_prop", Value.newBuilder().setInt64Value(100L).build());
        String edgeType = "knows";
        long rank = 10;

        return Edge.newBuilder()
                .setType(edgeType)
                .addLabels("knows1")
                .addLabels("knows2")
                .setSrcId(srcId)
                .setDstId(dstId)
                .setRank(rank)
                .putAllProperties(props)
                .build();
    }

    private Path getPath() {
        Path.Builder builder = Path.newBuilder();
        builder.addValues(Value.newBuilder().setNodeValue(getSimpleNode(1)));
        builder.addValues(Value.newBuilder().setEdgeValue(getSimpleEdge(false, 1, 2)).build());
        builder.addValues(Value.newBuilder().setNodeValue(getSimpleNode(2)).build());
        builder.addValues(Value.newBuilder().setEdgeValue(getSimpleEdge(false, 2, 3)).build());
        builder.addValues(Value.newBuilder().setNodeValue(getSimpleNode(3)).build());
        return builder.build();
    }

    private Path getReversePath() {
        Path.Builder builder = Path.newBuilder();
        builder.addValues(Value.newBuilder().setNodeValue(getSimpleNode(1)));
        builder.addValues(Value.newBuilder().setEdgeValue(getSimpleEdge(false, 2, 1)).build());
        builder.addValues(Value.newBuilder().setNodeValue(getSimpleNode(2)).build());
        builder.addValues(Value.newBuilder().setEdgeValue(getSimpleEdge(false, 3, 2)).build());
        builder.addValues(Value.newBuilder().setNodeValue(getSimpleNode(3)).build());
        return builder.build();
    }

    private Record getRowRecord() {
        Map<String, Value> values = new HashMap<>();
        values.put("prop1", Value.newBuilder().setInt8Value(1).build());
        values.put("prop2", Value.newBuilder().setInt64Value(2).build());
        values.put("prop3", Value.newBuilder()
                .setStringValue(ByteString.copyFrom("Tom", Charsets.UTF_8)).build());
        values.put("prop4", Value.newBuilder().setBoolValue(true).build());
        values.put("prop5", Value.newBuilder().setNodeValue(getSimpleNode(1)).build());
        values.put("prop6", Value.newBuilder().setDoubleValue(20.1).build());

        return Record.newBuilder().putAllValues(values).build();
    }

    private com.vesoft.nebula.proto.common.List getList() {
        List<Value> values = new ArrayList<>();
        values.add(Value.newBuilder().setInt64Value(10).build());
        values.add(Value.newBuilder().setInt64Value(11).build());
        values.add(Value.newBuilder().setInt64Value(12).build());
        values.add(Value.newBuilder().setInt64Value(13).build());
        com.vesoft.nebula.proto.common.List.Builder builder =
                com.vesoft.nebula.proto.common.List.newBuilder();
        builder.addAllValues(values);
        return builder.build();
    }

    private LocalTime getSimpleLocalTime() {
        return LocalTime.newBuilder().setHour(12).setMinute(20).setSec(15).setMicrosec(30).build();
    }

    private LocalDatetime getSimpleLocalDateTime() {
        return LocalDatetime.newBuilder()
                .setYear(2024)
                .setMonth(1)
                .setDay(1)
                .setHour(12)
                .setMinute(20)
                .setSec(15)
                .setMicrosec(30)
                .build();
    }

    private Date getSimpleDate() {
        return Date.newBuilder().setYear(2024).setMonth(1).setDay(1).build();
    }

    private ResultTable getDateset() {
        List<Value> values = Arrays.asList(
                Value.newBuilder().build(),
                Value.newBuilder().setBoolValue(false).build(),
                Value.newBuilder().setInt64Value(1).build(),
                Value.newBuilder().setInt32Value(2).build(),
                Value.newBuilder().setInt16Value(3).build(),
                Value.newBuilder().setInt8Value(4).build(),
                Value.newBuilder().setFloatValue(10.01f).build(),
                Value.newBuilder().setDoubleValue(20.01).build(),
                Value.newBuilder()
                        .setStringValue(ByteString.copyFrom("value1", Charsets.UTF_8))
                        .build(),
                Value.newBuilder().setListValue(getList()).build(),
                Value.newBuilder().setNodeValue(getSimpleNode(1)).build(),
                Value.newBuilder().setEdgeValue(getSimpleEdge(false, 1, 2)).build());

        Row row = Row.newBuilder().addAllValues(values).build();

        final List<ByteString> columnNames = Arrays.asList(
                ByteString.copyFrom("col0_empty", Charsets.UTF_8),
                ByteString.copyFrom("col1_bool", Charsets.UTF_8),
                ByteString.copyFrom("col2_int64", Charsets.UTF_8),
                ByteString.copyFrom("col3_int32", Charsets.UTF_8),
                ByteString.copyFrom("col4_int16", Charsets.UTF_8),
                ByteString.copyFrom("col5_int8", Charsets.UTF_8),
                ByteString.copyFrom("col6_float", Charsets.UTF_8),
                ByteString.copyFrom("col7_double", Charsets.UTF_8),
                ByteString.copyFrom("col8_string", Charsets.UTF_8),
                ByteString.copyFrom("col9_list", Charsets.UTF_8),
                ByteString.copyFrom("col10_vertex", Charsets.UTF_8),
                ByteString.copyFrom("col11_edge", Charsets.UTF_8));
        ResultTable.Builder tableBuilder = ResultTable.newBuilder().addAllColumnNames(columnNames);

        tableBuilder.addRecords(row);
        return tableBuilder.build();
    }
}
