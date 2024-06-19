/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.converter;

import com.alibaba.fastjson.JSON;
import com.vesoft.nebula.client.graph.data.ValueWrapper;
import com.vesoft.nebula.client.graph.scan.TableRow;
import com.vesoft.nebula.connector.exceptions.RecordConversionException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;
import org.apache.kafka.connect.data.Field;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.SchemaBuilder;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.sink.SinkRecord;
import org.apache.kafka.connect.source.SourceRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class NebulaRecordConverter {
    private static final Logger log = LoggerFactory.getLogger(NebulaRecordConverter.class);

    private static final Set<Class<?>> BASIC_TYPES = new HashSet<>(
        Arrays.asList(
            Boolean.class, Character.class, Byte.class, Short.class, Integer.class,
                Long.class, BigInteger.class, Float.class, Double.class, BigDecimal.class,
                String.class)
    );

    public static Map<String, Object> convertRecord(SinkRecord record)
        throws RecordConversionException {
        Map<String, Object> properties = new HashMap<>();
        if (record.valueSchema() != null && record.value() instanceof Struct) {
            Struct struct = (Struct) record.value();
            log.debug("kafka record: {}", struct.toString());
            Schema valueSchema = record.valueSchema();

            // the mapping for kafka property name and property value. get values
            // according to propertyNames config
            for (Field field : valueSchema.fields()) {
                boolean isEmptyStruct =
                    field.schema().type() == Schema.Type.STRUCT
                        && field.schema().fields().isEmpty();
                if (isEmptyStruct) {
                    continue;
                }
                Object structValue = convertObject(struct.get(field.name()), field.schema());
                properties.put(field.name(), structValue);
            }
        } else if (record.value() instanceof String) {
            Map map = JSON.parseObject((String) record.value());
            properties = (Map<String, Object>) convertSchemalessRecord(map);
        } else if (record.value() instanceof Map) {
            properties = (Map<String, Object>) convertSchemalessRecord(record.value());
            log.debug("kafka record: {}", properties);
        } else {
            throw new RecordConversionException("only Map object supported in absence of schema "
                + "for converting kafka record into Nebula format.");
        }

        return properties;
    }

    @SuppressWarnings("unchecked")
    private static Object convertSchemalessRecord(Object value) {
        if (value == null) {
            return null;
        }
        if (value instanceof byte[] || value instanceof ByteBuffer) {
            return convertBytes(value);
        }
        if (BASIC_TYPES.contains(value.getClass())) {
            return value;
        }
        if (value instanceof List) {
            return ((List<?>) value)
                .stream()
                .map(NebulaRecordConverter::convertSchemalessRecord)
                .collect(Collectors.toList());
        }
        if (value instanceof Map) {
            return ((Map<Object, Object>) value)
                .entrySet()
                .stream()
                .collect(HashMap::new,
                    (m, e) -> {
                        if (!(e.getKey() instanceof String)) {
                            throw new RecordConversionException(
                                "Failed to convert record to Nebula format, "
                                    + "Map object needs to have String type key.");
                        }
                        m.put(e.getKey(), convertSchemalessRecord(e.getValue()));
                    },
                    HashMap::putAll);
        }
        throw new RecordConversionException("Unsupported class " + value.getClass()
            + " found in schemaless record data.");
    }


    private static Object convertObject(Object value, Schema kafkaValueSchema) {
        if (value == null) {
            if (kafkaValueSchema.isOptional()) {
                return null;
            } else {
                throw new RecordConversionException(kafkaValueSchema.name()
                    + " is not optional, but converting object had null value.");
            }
        }

        Schema.Type schemaType = kafkaValueSchema.type();
        switch (schemaType) {
            case ARRAY:
                return convertArray(value, kafkaValueSchema);
            case MAP:
                return convertMap(value, kafkaValueSchema);
            case STRUCT:
                return convertObject(value, kafkaValueSchema);
            case BYTES:
                return convertBytes(value);
            case FLOAT32:
            case FLOAT64:
            case BOOLEAN:
            case INT8:
            case INT16:
            case INT32:
            case INT64:
            case STRING:
                return value;
            default:
                throw new RecordConversionException("Unrecognized schema type: " + schemaType);
        }
    }

    private static Map<String, Object> convertStruct(Object value, Schema schema) {
        Map<String, Object> propertyRecord = new HashMap<>();
        List<Field> kafkaSchemaFields = schema.fields();
        Struct struct = (Struct) value;
        for (Field field : schema.fields()) {
            boolean isEmptyStruct =
                field.schema().type() == Schema.Type.STRUCT
                    && field.schema().fields().isEmpty();
            if (isEmptyStruct) {
                continue;
            }
            Object structValue = convertObject(struct.get(field.name()), field.schema());
            propertyRecord.put(field.name(), structValue);
        }
        return propertyRecord;
    }

    /**
     * convert kafka Array data into String type for NebulaGraph
     * because NebulaGraph property does not support List type.
     */
    private static String convertArray(Object value, Schema schema) {

        List<Object> kafkaList = (List<Object>) value;
        StringBuilder sb = new StringBuilder();
        for (Object elementValue : kafkaList) {
            Object nebulaValue = convertObject(elementValue, schema);
            sb.append(nebulaValue.toString());
            sb.append(",");
        }
        sb.deleteCharAt(sb.length() - 1);
        return sb.toString();
    }

    /**
     * convert kafka Map data into String type for NebulaGraph
     * because NebulaGraph property does not support Map type.
     */
    private static String convertMap(Object value, Schema schema) {
        Map<String, Object> nebulaMap = new HashMap<>();
        StringBuilder sb = new StringBuilder();
        for (Map.Entry<String, Object> entry : ((Map<String, Object>) value).entrySet()) {
            sb.append(entry.getKey());
            sb.append("->");
            sb.append(convertObject(entry.getValue(), schema));
            sb.append(",");
        }
        sb.deleteCharAt(sb.length() - 1);
        return sb.toString();
    }

    private static Object convertBytes(Object value) {
        byte[] bytes;
        if (value instanceof ByteBuffer) {
            ByteBuffer byteBuffer = (ByteBuffer) value;
            bytes = byteBuffer.array();
        } else {
            bytes = (byte[]) value;
        }
        return new String(bytes);
    }


    /**
     * converter for kafka source connect.
     * convert NebulaGraph {@link TableRow} to Kafka {@link SourceRecord}
     *
     * @param rowNames        row value names
     * @param rows            rows
     * @param topicName       topic name
     * @param nebulaRowSchema NebulaGraph row schema
     * @return list of SourceRecord
     */
    public static List<SourceRecord> convertTableRows(List<String> rowNames,
                                                      List<TableRow> rows,
                                                      String topicName,
                                                      Map<String, String> nebulaRowSchema) {
        List<SourceRecord> records = new ArrayList<>();
        if (rows == null || rows.isEmpty()) {
            return records;
        }
        Map<String, String> sourcePartitions = Collections.singletonMap("type", topicName);
        // 定义row 的schema
        SchemaBuilder schemaBuilder = SchemaBuilder.struct().name(topicName);
        for (String name : rowNames) {
            schemaBuilder.field(name, convertNebulaType2KafkaType(name, nebulaRowSchema));
        }
        Schema schema = schemaBuilder.build();

        for (TableRow row : rows) {
            Struct struct = new Struct(schema);
            for (int i = 0; i < row.size(); i++) {
                struct.put(rowNames.get(i), convertValueWrapper2KafkaValue(row.getValues().get(i)));
            }
            SourceRecord record = new SourceRecord(
                sourcePartitions,
                null,
                topicName,
                schema,
                struct);
            records.add(record);
        }

        return records;
    }


    /**
     * data type converter for kafka source connect.
     * convert NebulaGraph data type to Kafka data type.
     *
     * @param name         property name
     * @param nebulaSchema NebulaGraph property schema map. key is the property name,
     *                     value is the data type for property name.
     * @return kafka data {@link Schema}
     */
    private static Schema convertNebulaType2KafkaType(String name,
                                                      Map<String, String> nebulaSchema) {
        String nebulaDataType = nebulaSchema.get(name);
        switch (nebulaDataType.toUpperCase()) {
            case "INT8":
            case "INT16":
            case "INT32":
                return Schema.INT32_SCHEMA;
            case "INT64":
                return Schema.INT64_SCHEMA;
            case "FLOAT":
                return Schema.FLOAT32_SCHEMA;
            case "DOUBLE":
                return Schema.FLOAT64_SCHEMA;
            case "BOOL":
                return Schema.BOOLEAN_SCHEMA;
            default:
                return Schema.STRING_SCHEMA;
        }
    }

    /**
     * converter for kafka source connect.
     * convert the NebulaGraph value to kafka value.
     *
     * @param value NebulaGraph property value
     * @return Object
     */
    private static Object convertValueWrapper2KafkaValue(ValueWrapper value) {
        if (value.isEmpty()) {
            return null;
        }
        if (value.isInt()) {
            return value.asInt();
        }
        if (value.isLong()) {
            return value.asLong();
        }
        if (value.isBoolean()) {
            return value.asBoolean();
        }
        if (value.isFloat()) {
            return value.asFloat();
        }
        if (value.isDouble()) {
            return value.asDouble();
        }
        return value.toString();
    }
}
