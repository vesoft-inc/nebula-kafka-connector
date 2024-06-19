/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.converter;

import com.vesoft.nebula.connector.exceptions.RecordConversionException;
import java.util.Arrays;
import java.util.Map;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.SchemaBuilder;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.sink.SinkRecord;
import org.junit.Test;

public class NebulaRecordConverterTest {


    @Test
    public void testConvertRecord() {
        SinkRecord sinkRecord = mockSinkRecord();
        try {
            Map<String, Object> properties = NebulaRecordConverter.convertRecord(sinkRecord);
            assert (properties.keySet().containsAll(
                    Arrays.asList("a", "b", "c", "d", "e", "f", "g")));
            assert (properties.get("a").equals("Tom"));
            assert (Integer.valueOf((byte) properties.get("b")) == 1);
            assert (Integer.valueOf((short) properties.get("c")) == 2);
            assert ((int) properties.get("d") == 3);
            assert ((long) properties.get("e") == 10);
            assert ((boolean) properties.get("f"));
            assert ((float) properties.get("h") == 2.0);
            assert ((double) properties.get("i") == 1.0);
        } catch (RecordConversionException e) {
            e.printStackTrace();
            assert (false);
        }
    }

    @Test
    public void testConvertStringRecord() {
        SinkRecord sinkRecord = new SinkRecord("test", 0, null, null, null,
                "{\"score\":12.0,\"gender\":\"男\","
                        + "\"rate\":1.0,\"name\":\"Tom\",\"id\":\"200\"}", 0);
        try {
            Map<String, Object> properties = NebulaRecordConverter.convertRecord(sinkRecord);
            assert (properties.get("name").equals("Tom"));
        } catch (RecordConversionException e) {
            e.printStackTrace();
            assert false;
        }
    }

    private SinkRecord mockSinkRecord() {
        final Schema SCHEMA = SchemaBuilder.struct()
                .field("a", Schema.STRING_SCHEMA)
                .field("b", Schema.INT8_SCHEMA)
                .field("c", Schema.INT16_SCHEMA)
                .field("d", Schema.INT32_SCHEMA)
                .field("e", Schema.INT64_SCHEMA)
                .field("f", Schema.BOOLEAN_SCHEMA)
                .field("g", Schema.BYTES_SCHEMA)
                .field("h", Schema.FLOAT32_SCHEMA)
                .field("i", Schema.FLOAT64_SCHEMA)
                .build();
        Struct struct = new Struct(SCHEMA)
                .put("a", "Tom")
                .put("b", (byte) 1)
                .put("c", (short) 2)
                .put("d", 3)
                .put("e", 10L)
                .put("f", true)
                .put("g", "abc".getBytes())
                .put("h", 2.0f)
                .put("i", 1.0);
        return new SinkRecord("test", 0, null, null, SCHEMA, struct, 0);
    }
}
