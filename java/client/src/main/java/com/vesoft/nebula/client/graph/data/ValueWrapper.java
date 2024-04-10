/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.google.common.base.Charsets;
import com.vesoft.nebula.client.graph.exception.InvalidValueException;
import com.vesoft.nebula.proto.common.Value;
import java.io.UnsupportedEncodingException;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

public class ValueWrapper {

    private final Value value;
    private final Charset charset = Charsets.UTF_8;

    public String getDataType() {
        if (value.getDataCase() == Value.DataCase.DATA_NOT_SET) {
            return "DATA_NOT_SET";
        }
        if (value.getDataCase() == Value.DataCase.BOOL_VALUE) {
            return "BOOLEAN";
        }
        if (value.getDataCase() == Value.DataCase.INT8_VALUE) {
            return "BYTE";
        }
        if (value.getDataCase() == Value.DataCase.INT16_VALUE) {
            return "SHORT";
        }
        if (value.getDataCase() == Value.DataCase.INT32_VALUE) {
            return "INT";
        }
        if (value.getDataCase() == Value.DataCase.INT64_VALUE) {
            return "LONG";
        }
        if (value.getDataCase() == Value.DataCase.FLOAT_VALUE) {
            return "FLOAT";
        }
        if (value.getDataCase() == Value.DataCase.DOUBLE_VALUE) {
            return "DOUBLE";
        }
        if (value.getDataCase() == Value.DataCase.STRING_VALUE) {
            return "STRING";
        }
        if (value.getDataCase() == Value.DataCase.NODE_VALUE) {
            return "NODE";
        }
        if (value.getDataCase() == Value.DataCase.EDGE_VALUE) {
            return "EDGE";
        }
        if (value.getDataCase() == Value.DataCase.LIST_VALUE) {
            return "LIST";
        }
        if (value.getDataCase() == Value.DataCase.DURATION_VALUE) {
            return "DURATION";
        }
        if (value.getDataCase() == Value.DataCase.LOCAL_TIME_VALUE) {
            return "LOCAL_TIME";
        }
        if (value.getDataCase() == Value.DataCase.LOCAL_DATETIME_VALUE) {
            return "LOCAL_DATETIME";
        }
        if (value.getDataCase() == Value.DataCase.ZONED_TIME_VALUE) {
            return "ZONED_TIME";
        }
        if (value.getDataCase() == Value.DataCase.ZONED_DATETIME_VALUE) {
            return "ZONED_DATETIME";
        }
        if (value.getDataCase() == Value.DataCase.DATE_VALUE) {
            return "DATE";
        }
        if (value.getDataCase() == Value.DataCase.RECORD_VALUE) {
            return "RECORD";
        }
        if (value.getDataCase() == Value.DataCase.PATH_VALUE) {
            return "PATH";
        }
        throw new IllegalArgumentException("Unknown field id " + value.getDataCase());
    }

    /**
     * @param value the Value get from service
     */
    public ValueWrapper(Value value) {
        this.value = value;
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
        return value.getDataCase() == Value.DataCase.DATA_NOT_SET;
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
        return value.getDataCase() == Value.DataCase.BOOL_VALUE;
    }

    /**
     * judge the Value is Long type
     *
     * @return boolean
     */
    public boolean isLong() {
        return value.getDataCase() == Value.DataCase.INT64_VALUE
                || value.getDataCase() == Value.DataCase.UINT64_VALUE;
    }

    /**
     * judge the Value is Int type
     *
     * @return boolean
     */
    public boolean isInt() {
        return value.getDataCase() == Value.DataCase.INT32_VALUE
                || value.getDataCase() == Value.DataCase.UINT32_VALUE
                || value.getDataCase() == Value.DataCase.INT16_VALUE
                || value.getDataCase() == Value.DataCase.UINT16_VALUE
                || value.getDataCase() == Value.DataCase.INT8_VALUE
                || value.getDataCase() == Value.DataCase.UINT8_VALUE;
    }

    /**
     * judge the Value is Float type
     *
     * @return boolean
     */
    public boolean isFloat() {
        return value.getDataCase() == Value.DataCase.FLOAT_VALUE;
    }


    /**
     * judge the Value is Double type
     *
     * @return boolean
     */
    public boolean isDouble() {
        return value.getDataCase() == Value.DataCase.DOUBLE_VALUE;
    }

    /**
     * judge the Value is String type
     *
     * @return boolean
     */
    public boolean isString() {
        return value.getDataCase() == Value.DataCase.STRING_VALUE;
    }

    /**
     * judge the Value is List type, the List type is the nebula's type
     *
     * @return boolean
     */
    public boolean isList() {
        return value.getDataCase() == Value.DataCase.LIST_VALUE;
    }


    public boolean isNode() {
        return value.getDataCase() == Value.DataCase.NODE_VALUE;
    }

    public boolean isEdge() {
        return value.getDataCase() == Value.DataCase.EDGE_VALUE;
    }

    public boolean isLocalTime() {
        return value.getDataCase() == Value.DataCase.LOCAL_TIME_VALUE;
    }

    public boolean isZonedTime() {
        return value.getDataCase() == Value.DataCase.ZONED_TIME_VALUE;
    }

    public boolean isLocalDateTime() {
        return value.getDataCase() == Value.DataCase.LOCAL_DATETIME_VALUE;
    }

    public boolean isZonedDateTime() {
        return value.getDataCase() == Value.DataCase.ZONED_DATETIME_VALUE;
    }

    public boolean isDate() {
        return value.getDataCase() == Value.DataCase.DATE_VALUE;
    }

    public boolean isRecord() {
        return value.getDataCase() == Value.DataCase.RECORD_VALUE;
    }

    public boolean isDuration() {
        return value.getDataCase() == Value.DataCase.DURATION_VALUE;
    }

    public boolean isPath() {
        return value.getDataCase() == Value.DataCase.PATH_VALUE;
    }


    /**
     * Convert the original data type Value to boolean
     *
     * @return boolean
     * @throws InvalidValueException if the value type is not boolean
     */
    public boolean asBoolean() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.BOOL_VALUE) {
            return value.getBoolValue();
        }
        throw new InvalidValueException(
                "Cannot get field `boolean` because value's type is " + getDataType());
    }


    /**
     * Convert the original data type Value to int
     *
     * @return int
     * @throws InvalidValueException if the value type is not int32
     */
    public int asInt() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.INT8_VALUE) {
            return value.getInt8Value();
        } else if (value.getDataCase() == Value.DataCase.INT16_VALUE) {
            return value.getInt16Value();
        } else if (value.getDataCase() == Value.DataCase.INT32_VALUE) {
            return value.getInt32Value();
        } else if (value.getDataCase() == Value.DataCase.UINT8_VALUE) {
            return value.getUint8Value();
        } else if (value.getDataCase() == Value.DataCase.UINT16_VALUE) {
            return value.getUint16Value();
        } else if (value.getDataCase() == Value.DataCase.UINT32_VALUE) {
            return value.getUint32Value();
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
        if (value.getDataCase() == Value.DataCase.INT64_VALUE) {
            return value.getInt64Value();
        } else if (value.getDataCase() == Value.DataCase.UINT64_VALUE) {
            return value.getUint64Value();
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
        if (value.getDataCase() == Value.DataCase.STRING_VALUE) {
            return value.getStringValue().toString(charset);
        }
        throw new InvalidValueException(
                "Cannot get field `string` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to float
     *
     * @return float
     * @throws InvalidValueException if the value type is not float
     */
    public float asFloat() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.FLOAT_VALUE) {
            return value.getFloatValue();
        }
        throw new InvalidValueException(
                "Cannot get field `float` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to double
     *
     * @return double
     * @throws InvalidValueException if the value type is not double
     */
    public double asDouble() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.DOUBLE_VALUE) {
            return value.getDoubleValue();
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
        if (value.getDataCase() == Value.DataCase.LIST_VALUE) {
            List<ValueWrapper> values = new ArrayList<>();
            for (Value value : value.getListValue().getValuesList()) {
                values.add(new ValueWrapper(value));
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
        if (value.getDataCase() == Value.DataCase.NODE_VALUE) {
            return new Vertex(value.getNodeValue());
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
        if (value.getDataCase() == Value.DataCase.EDGE_VALUE) {
            return new Relationship(value.getEdgeValue());
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
        if (value.getDataCase() == Value.DataCase.LOCAL_TIME_VALUE) {
            return new NTime(value.getLocalTimeValue());
        }
        throw new InvalidValueException(
                "cannot get field `localtime` because value's type is " + getDataType());
    }


    /**
     * Convert the original data type Value to {@link NZonedTime}
     *
     * @return NTime
     * @throws InvalidValueException if the value type is not localtime
     */
    public NZonedTime asZonedTime() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.ZONED_TIME_VALUE) {
            return new NZonedTime(value.getZonedTimeValue());
        }
        throw new InvalidValueException(
                "cannot get field `zonedtime` because value's type is " + getDataType());
    }


    /**
     * Convert the original data type Value to {@link NDate}
     *
     * @return NDate
     * @throws InvalidValueException if the value type is not localtime
     */
    public NDate asDate() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.DATE_VALUE) {
            return new NDate(value.getDateValue());
        }
        throw new InvalidValueException(
                "cannot get field `date` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link NDateTime}
     *
     * @return NDateTime
     * @throws InvalidValueException if the value type is not localDatetime
     */
    public NDateTime asLocalDateTime() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.LOCAL_DATETIME_VALUE) {
            return new NDateTime(value.getLocalDatetimeValue());
        }
        throw new InvalidValueException(
                "cannot get field `localdatetime` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link NZonedDateTime}
     *
     * @return NZonedDateTime
     * @throws InvalidValueException if the value type is not localDatetime
     */
    public NZonedDateTime asZonedDateTime() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.ZONED_DATETIME_VALUE) {
            return new NZonedDateTime(value.getZonedDatetimeValue());
        }
        throw new InvalidValueException(
                "cannot get field `zoneddatetime` because value's type is " + getDataType());
    }

    /**
     * Convert the original data type Value to {@link NDuration}
     *
     * @return NDuration
     * @throws InvalidValueException if the value type is not duration
     */
    public NDuration asDuration() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.DURATION_VALUE) {
            return new NDuration(value.getDurationValue());
        }
        throw new InvalidValueException(
                "cannot get field `duration` because value's type is " + getDataType());
    }

    public NRecord asRecord() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.RECORD_VALUE) {
            return new NRecord(value.getRecordValue());
        }
        throw new InvalidValueException(
                "cannot get field `record` because value's type is " + getDataType());
    }

    public NPath asPath() throws InvalidValueException {
        if (value.getDataCase() == Value.DataCase.PATH_VALUE) {
            return new NPath(value.getPathValue());
        }
        throw new InvalidValueException(
                "cannot get field `path` because value's type is " + getDataType());
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
        return Objects.equals(value, that.value);
    }

    @Override
    public int hashCode() {
        return Objects.hash(value);
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
        } else if (isInt()) {
            return String.valueOf(asInt());
        } else if (isLong()) {
            return String.valueOf(asLong());
        } else if (isFloat()) {
            return String.valueOf(asFloat());
        } else if (isDouble()) {
            return String.valueOf(asDouble());
        } else if (isString()) {
            return asString();
        } else if (isList()) {
            return asList().toString();
        } else if (isRecord()) {
            return asRecord().toString();
        } else if (isNode()) {
            return asNode().toString();
        } else if (isEdge()) {
            return asEdge().toString();
        } else if (isLocalTime()) {
            return asLocalTime().toString();
        } else if (isZonedTime()) {
            return asZonedTime().toString();
        } else if (isLocalDateTime()) {
            return asLocalDateTime().toString();
        } else if (isZonedDateTime()) {
            return asZonedDateTime().toString();
        } else if (isDate()) {
            return asDate().toString();
        } else if (isDuration()) {
            return asDuration().toString();
        } else if (isPath()) {
            return asPath().toString();
        }
        return "Unknown type: " + getDataType();

    }
}
