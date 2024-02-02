/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.converter;

import com.vesoft.nebula.connector.exceptions.RecordConversionException;
import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;
import org.apache.kafka.connect.data.Field;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.sink.SinkRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class NebulaRecordConverter {
    private static final Logger log = LoggerFactory.getLogger(NebulaRecordConverter.class);

    private static final Set<Class<?>> BASIC_TYPES = new HashSet<>(
            Arrays.asList(
                    Boolean.class, Character.class, Byte.class, Short.class,
                    Integer.class, Long.class, Float.class, Double.class, String.class)
    );

    public static Map<String, Object> convertRecord(SinkRecord record) throws RecordConversionException {
        Map<String, Object> properties = new HashMap<>();
        if (record.valueSchema() != null && record.value() instanceof Struct) {
            Struct struct = (Struct) record.value();
            log.debug("kafka record: {}", struct.toString());
            Schema valueSchema = record.valueSchema();

            // the mapping for kafka property name and property value. get values
            // according to propertyNames config
            for (Field field : valueSchema.fields()) {
                boolean isEmptyStruct =
                        field.schema().type() == Schema.Type.STRUCT && field.schema().fields().isEmpty();
                if (isEmptyStruct) continue;
                Object structValue = convertObject(struct.get(field.name()), field.schema());
                properties.put(field.name(), structValue);
            }
        } else if (record.value() instanceof Map) {
            properties = (Map<String, Object>) convertSchemalessRecord(record.value());
            log.debug("kafka record: {}", properties);
        } else {
            throw new RecordConversionException("only Map object supported in absence of schema for " +
                    "converting kafka record into Nebula format.");
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
            return ((List<?>) value).stream().map(NebulaRecordConverter::convertSchemalessRecord).collect(Collectors.toList());
        }
        if (value instanceof Map) {
            return ((Map<Object, Object>) value)
                    .entrySet()
                    .stream()
                    .collect(HashMap::new,
                            (m, e) -> {
                                if (!(e.getKey() instanceof String)) {
                                    throw new RecordConversionException("Failed to convert record to Nebula format, " +
                                            "Map object needs to have String type key.");
                                }
                                m.put(e.getKey(), convertSchemalessRecord(e.getValue()));
                            },
                            HashMap::putAll);
        }
        throw new RecordConversionException("Unsupported class " + value.getClass() +
                " found in schemaless record data.");
    }


    private static Object convertObject(Object value, Schema kafkaValueSchema) {
        if (value == null) {
            if (kafkaValueSchema.isOptional()) {
                return null;
            } else {
                throw new RecordConversionException(kafkaValueSchema.name() + " is not optional, but converting " +
                        "object had null value.");
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
            boolean isEmptyStruct = field.schema().type() == Schema.Type.STRUCT && field.schema().fields().isEmpty();
            if (isEmptyStruct) continue;
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
     *
     * */
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
}
