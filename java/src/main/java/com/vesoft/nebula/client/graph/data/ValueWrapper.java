/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.Value;
import com.vesoft.nebula.client.graph.exception.InvalidValueException;
import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
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
    private int timezoneOffset = 0;

    private String descType() {
        switch (value.getSetField()) {
            case Value.BOOLVAL:
                return "BOOLEAN";
            case Value.INT8VAL | Value.INT16VAL | Value.INT32VAL | Value.INT64VAL:
                return "INT";
            case Value.FLOATVAL:
                return "FLOAT";
            case Value.STRINGVAL:
                return "STRING";
            case Value.NODEVAL:
                return "NODE";
            case Value.EDGEVAL:
                return "EDGE";
            case Value.LISTVAL:
                return "LIST";
            case Value.MAPVAL:
                return "MAP";
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
        this.timezoneOffset = 0;
    }

    /**
     * @param value          the Value get from service
     * @param decodeType     the decodeType get from the service to decode the byte array,
     *                       but now the service no return the decodeType, so use the utf-8
     * @param timezoneOffset the timezone offset get from the service to calculate local time
     */
    public ValueWrapper(Value value, String decodeType, int timezoneOffset) {
        this.value = value;
        this.decodeType = decodeType;
        this.timezoneOffset = timezoneOffset;
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


    /**
     * judge the Value is Map type, the Map type is the nebula's type
     *
     * @return boolean
     */
    public boolean isMap() {
        return value.getSetField() == Value.MAPVAL;
    }


    public boolean isNode() {
        return value.getSetField() == Value.NODEVAL;
    }

    public boolean isEdge() {
        return value.getSetField() == Value.EDGEVAL;
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
                "Cannot get field boolean because value's type is " + descType());
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
        } else if (value.getSetField() == Value.INT32VAL) {
            return ((Integer) value.getInt32Val()).longValue();
        } else if (value.getSetField() == Value.INT16VAL) {
            return ((Short) value.getInt16Val()).longValue();
        } else if (value.getSetField() == Value.INT8VAL) {
            return ((Byte) value.getInt8Val()).longValue();
        } else {
            throw new InvalidValueException(
                    "Cannot get field long because value's type is " + descType());
        }
    }


    /**
     * Convert the original data type Value to String
     *
     * @return String
     * @throws InvalidValueException        if the value type is not string
     * @throws UnsupportedEncodingException if decode bianry failed
     */
    public String asString() throws InvalidValueException, UnsupportedEncodingException {
        if (value.getSetField() == Value.STRINGVAL) {
            return new String((byte[]) value.getFieldValue(), decodeType);
        }
        throw new InvalidValueException(
                "Cannot get field string because value's type is " + descType());
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
                "Cannot get field double because value's type is " + descType());
    }

    /**
     * Convert the original data type Value to ArrayList
     *
     * @return ArrayList of ValueWrapper
     * @throws InvalidValueException if the value type is not list
     */
    public ArrayList<ValueWrapper> asList() throws InvalidValueException {
        if (value.getSetField() != Value.LISTVAL) {
            throw new InvalidValueException(
                    "Cannot get field type `list' because value's type is " + descType());
        }
        ArrayList<ValueWrapper> values = new ArrayList<>();
        for (Value value : value.getListVal().getValues()) {
            values.add(new ValueWrapper(value, decodeType, timezoneOffset));
        }
        return values;
    }

    public Vertex asNode() throws InvalidValueException, UnsupportedEncodingException {
        if (value.getSetField() == Value.NODEVAL) {
            return (Vertex) new Vertex(value.getNodeVal())
                    .setDecodeType(decodeType);
        }
        throw new InvalidValueException(
                "Cannot get field Node because value's type is " + descType());
    }

    public Relationship asEdge() throws InvalidValueException {
        if (value.getSetField() == Value.EDGEVAL) {
            return (Relationship) new Relationship(value.getEdgeVal()).setDecodeType(decodeType);
        }
        throw new InvalidValueException(
                "Cannot get field Edge because value's type is " + descType());
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
        try {
            if (isEmpty()) {
                return "__EMPTY__";
            } else if (isBoolean()) {
                return String.valueOf(asBoolean());
            } else if (isLong() || isInt() || isShort() || isByte()) {
                return String.valueOf(asLong());
            } else if (isDouble()) {
                return String.valueOf(asDouble());
            } else if (isString()) {
                return "\"" + asString() + "\"";
            } else if (isList()) {
                return asList().toString();
            } else if (isNode()) {
                return asNode().toString();
            } else if (isEdge()) {
                return asEdge().toString();
            }
            return "Unknown type: " + descType();
        } catch (UnsupportedEncodingException e) {
            return e.getMessage();
        }
    }
}
