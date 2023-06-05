/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.Value;
import com.vesoft.nebula.client.graph.exception.InvalidValueException;
import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

public class ValueWrapper {
    public static class NullType {
        public static final int __NULL__ = 0;
        public static final int NaN = 1;
        public static final int BAD_DATA = 2;
        public static final int BAD_TYPE = 3;
        public static final int ERR_OVERFLOW = 4;
        public static final int UNKNOWN_PROP = 5;
        public static final int DIV_BY_ZERO = 6;
        public static final int OUT_OF_RANGE = 7;
        int nullType;

        public NullType(int nullType) {
            this.nullType = nullType;
        }

        public int getNullType() {
            return nullType;
        }

        @Override
        public String toString() {
            switch (nullType) {
                case __NULL__:
                    return "NULL";
                case NaN:
                    return "NaN";
                case BAD_DATA:
                    return "BAD_DATA";
                case BAD_TYPE:
                    return "BAD_TYPE";
                case ERR_OVERFLOW:
                    return "ERR_OVERFLOW";
                case UNKNOWN_PROP:
                    return "UNKNOWN_PROP";
                case DIV_BY_ZERO:
                    return "DIV_BY_ZERO";
                case OUT_OF_RANGE:
                    return "OUT_OF_RANGE";
                default:
                    return "Unknown type: " + nullType;
            }
        }
    }

    private final Value value;
    private String decodeType = "utf-8";

    public String getDataType() {
        switch (value.getSetField()) {
            case Value.BOOLVAL:
                return "BOOLEAN";
            case Value.INT8VAL:
                return "BYTE";
            case Value.INT16VAL:
                return "SHORT";
            case Value.INT32VAL:
                return "INT";
            case Value.INT64VAL:
                return "LONG";
            case Value.FLOATVAL:
                return "FLOAT";
            case Value.DOUBLEVAL:
                return "DOUBLE";
            case Value.STRINGVAL:
                return "STRING";
            case Value.NODEVAL:
                return "NODE";
            case Value.EDGEVAL:
                return "EDGE";
            case Value.LISTVAL:
                return "LIST";
            case Value.DURATIONVAL:
                return "DURATION";
            case Value.LOCALTIMEVAL:
                return "LOCALTIME";
            case Value.LOCALDATETIMEVAL:
                return "LOCALDATETIME";
            case Value.DATEVAL:
                return "DATE";
            case Value.RECORDVAL:
                return "Map";
            default:
                throw new IllegalArgumentException("Unknown field id " + value.getSetField());
        }
    }

    /**
     * @param value      the Value get from service
     * @param decodeType the decodeType get from the service to decode the byte array,
     *                   but now the service no return the decodeType, so use the utf-8
     */
    public ValueWrapper(Value value, String decodeType) {
        this.value = value;
        this.decodeType = decodeType;
    }

    /**
     * get the original data structure, the Value is the return from nebula-graph
     *
     * @return Value
     */
    public Value getValue() {
        return value;
    }

    /**
     * judge the Value is Empty type, the Empty type is the nebula's type
     *
     * @return boolean
     */
    public boolean isEmpty() {
        return value.getSetField() == 0;
    }

    public boolean isNull() {
        return value == null;
    }

    /**
     * judge the Value is Boolean type
     *
     * @return boolean
     */
    public boolean isBoolean() {
        return value.getSetField() == Value.BOOLVAL;
    }

    /**
     * judge the Value is Long type
     *
     * @return boolean
     */
    public boolean isLong() {
        return value.getSetField() == Value.INT64VAL;
    }

    /**
     * judge the Value is Int type
     *
     * @return boolean
     */
    public boolean isInt() {
        return value.getSetField() == Value.INT32VAL;
    }

    /**
     * judge the Value is Short type
     *
     * @return boolean
     */
    public boolean isShort() {
        return value.getSetField() == Value.INT16VAL;
    }

    /**
     * judge the Value is Byte type
     *
     * @return boolean
     */
    public boolean isByte() {
        return value.getSetField() == Value.INT8VAL;
    }

    /**
     * judge the Value is Double type
     *
     * @return boolean
     */
    public boolean isDouble() {
        return value.getSetField() == Value.FLOATVAL || value.getSetField() == Value.DOUBLEVAL;
    }

    /**
     * judge the Value is String type
     *
     * @return boolean
     */
    public boolean isString() {
        return value.getSetField() == Value.STRINGVAL;
    }

    /**
     * judge the Value is List type, the List type is the nebula's type
     *
     * @return boolean
     */
    public boolean isList() {
        return value.getSetField() == Value.LISTVAL;
    }


    public boolean isNode() {
        return value.getSetField() == Value.NODEVAL;
    }

    public boolean isEdge() {
        return value.getSetField() == Value.EDGEVAL;
    }

    public boolean isLocalTime() {
        return value.getSetField() == Value.LOCALTIMEVAL;
    }

    public boolean isLocalDateTime() {
        return value.getSetField() == Value.LOCALDATETIMEVAL;
    }

    public boolean isDate() {
        return value.getSetField() == Value.DATEVAL;
    }

    public boolean isMap() {
        return value.getSetField() == Value.RECORDVAL;
    }

    public boolean isDuration() {
        return value.getSetField() == Value.DURATIONVAL;
    }

    /**
     * Convert the original data type Value to Object
     *
     * @return Object
     */
    public Object asObject() {
        return value.getFieldValue();
    }

    /**
     * Convert the original data type Value to boolean
     *
     * @return boolean
     * @throws InvalidValueException if the value type is not boolean
     */
    public boolean asBoolean() throws InvalidValueException {
        if (value.getSetField() == Value.BOOLVAL) {
            return (boolean) (value.getFieldValue());
        }
        throw new InvalidValueException(
                "Cannot get field `boolean` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to byte
     *
     * @return byte
     * @throws InvalidValueException if the value type is not byte
     */
    public byte asByte() throws InvalidValueException {
        if (value.getSetField() == Value.INT8VAL) {
            return (value.getInt8Val());
        } else {
            throw new InvalidValueException(
                    "Cannot get field `byte` because value's type is " + getDataType());
        }
    }

    /**
     * Convert the original data type Value to short
     *
     * @return short
     * @throws InvalidValueException if the value type is not short
     */
    public short asShort() throws InvalidValueException {
        if (value.getSetField() == Value.INT16VAL) {
            return value.getInt16Val();
        } else {
            throw new InvalidValueException(
                    "Cannot get field `short` because value's type is " + getDataType());
        }
    }

    /**
     * Convert the original data type Value to int
     *
     * @return int
     * @throws InvalidValueException if the value type is not int32
     */
    public int asInt() throws InvalidValueException {
        if (value.getSetField() == Value.INT32VAL) {
            return value.getInt32Val();
        } else {
            throw new InvalidValueException(
                    "Cannot get field `int` because value's type is " + getDataType());
        }
    }

    /**
     * Convert the original data type Value to long
     *
     * @return long
     * @throws InvalidValueException if the value type is not long
     */
    public long asLong() throws InvalidValueException {
        if (value.getSetField() == Value.INT64VAL) {
            return value.getInt64Val();
        } else {
            throw new InvalidValueException(
                    "Cannot get field `long` because value's type is " + getDataType());
        }
    }


    /**
     * Convert the original data type Value to String
     *
     * @return String
     * @throws InvalidValueException        if the value type is not string
     * @throws UnsupportedEncodingException if decode failed
     */
    public String asString() throws InvalidValueException {
        if (value.getSetField() == Value.STRINGVAL) {
            return new String((byte[]) value.getFieldValue());
        }
        throw new InvalidValueException(
                "Cannot get field `string` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to double
     *
     * @return double
     * @throws InvalidValueException if the value type is not double
     */
    public double asDouble() throws InvalidValueException {
        if (value.getSetField() == Value.DOUBLEVAL || value.getSetField() == Value.FLOATVAL) {
            return (double) value.getFieldValue();
        }
        throw new InvalidValueException(
                "Cannot get field `double` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link List}
     *
     * @return list
     * @throws InvalidValueException if the value type is not list
     */
    public List<ValueWrapper> asList() throws InvalidValueException {
        if (value.getSetField() == Value.LISTVAL) {
            List<ValueWrapper> values = new ArrayList<>();
            for (Value value : value.getListVal().getValues()) {
                values.add(new ValueWrapper(value, decodeType));
            }
            return values;
        }
        throw new InvalidValueException(
                "Cannot get field `list` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link Vertex}
     *
     * @return Vertex
     * @throws InvalidValueException if the value type is not node
     */
    public Vertex asNode() throws InvalidValueException {
        if (value.getSetField() == Value.NODEVAL) {
            return new Vertex(value.getNodeVal());
        }
        throw new InvalidValueException(
                "cannot get field `node` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link Relationship}
     *
     * @return Relationship
     * @throws InvalidValueException if the value type is not edge
     */
    public Relationship asEdge() throws InvalidValueException {
        if (value.getSetField() == Value.EDGEVAL) {
            return new Relationship(value.getEdgeVal());
        }
        throw new InvalidValueException(
                "cannot get field `edge` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link NTime}
     *
     * @return NTime
     * @throws InvalidValueException if the value type is not localtime
     */
    public NTime asLocalTime() throws InvalidValueException {
        if (value.getSetField() == Value.LOCALTIMEVAL) {
            return new NTime(value.getLocalTimeVal());
        }
        throw new InvalidValueException(
                "cannot get field `LocalTime` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link NDate}
     *
     * @return NDate
     * @throws InvalidValueException if the value type is not localtime
     */
    public NDate asDate() throws InvalidValueException {
        if (value.getSetField() == Value.DATEVAL) {
            return new NDate(value.getDateVal());
        }
        throw new InvalidValueException(
                "cannot get field `Date` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link NDateTime}
     *
     * @return NDateTime
     * @throws InvalidValueException if the value type is not localDatetime
     */
    public NDateTime asLocalDateTime() throws InvalidValueException {
        if (value.getSetField() == Value.LOCALDATETIMEVAL) {
            return new NDateTime(value.getLocalDatetimeVal());
        }
        throw new InvalidValueException(
                "cannot get field `LocalDatetime` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link NDuration}
     *
     * @return NDuration
     * @throws InvalidValueException if the value type is not duration
     */
    public NDuration asDuration() throws InvalidValueException {
        if (value.getSetField() == Value.DURATIONVAL) {
            return new NDuration(value.getDurationVal());
        }
        throw new InvalidValueException(
                "cannot get field `LocalDatetime` because value's type is " + getDataType());
    }

    public Map<String, ValueWrapper> asMap() throws InvalidValueException {
        if (value.getSetField() == Value.RECORDVAL) {
            Map<byte[], Value> values = value.getRecordVal().getValues();
            Map<String, ValueWrapper> recordValues = new HashMap<>(values.size());
            for (Map.Entry<byte[], Value> kv : values.entrySet()) {
                recordValues.put(new String(kv.getKey()), new ValueWrapper(kv.getValue(), "utf-8"));
            }
            return recordValues;
        }
        throw new InvalidValueException(
                "cannot get field `map` because value's type is " + getDataType());
    }


    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        ValueWrapper that = (ValueWrapper) o;
        return Objects.equals(value, that.value)
                && Objects.equals(decodeType, that.decodeType);
    }

    @Override
    public int hashCode() {
        return Objects.hash(value, decodeType);
    }

    /**
     * Convert Value to String format
     *
     * @return String
     */
    @Override
    public String toString() {
        if (isEmpty()) {
            return "";
        } else if (isNull()) {
            return null;
        } else if (isBoolean()) {
            return String.valueOf(asBoolean());
        } else if (isByte()) {
            return String.valueOf(asByte());
        } else if (isShort()) {
            return String.valueOf(asShort());
        } else if (isInt()) {
            return String.valueOf(asInt());
        } else if (isLong()) {
            return String.valueOf(asLong());
        } else if (isDouble()) {
            return String.valueOf(asDouble());
        } else if (isString()) {
            return asString();
        } else if (isList()) {
            return asList().toString();
        } else if (isMap()) {
            return asMap().toString();
        } else if (isNode()) {
            return asNode().toString();
        } else if (isEdge()) {
            return asEdge().toString();
        } else if (isLocalTime()) {
            return asLocalTime().toString();
        } else if (isLocalDateTime()) {
            return asLocalDateTime().toString();
        } else if (isDate()) {
            return asDate().toString();
        } else if (isDuration()) {
            return asDuration().toString();
        }
        return "Unknown type: " + getDataType();

    }
}
