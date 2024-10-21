/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.driver.graph.decode.datatype;

import static com.vesoft.nebula.driver.graph.decode.struct.SizeConstant.EDGE_TYPE_ID_SIZE;
import static com.vesoft.nebula.driver.graph.decode.struct.SizeConstant.GRAPH_ELEMENT_TYPE_NUM_SIZE;
import static com.vesoft.nebula.driver.graph.decode.struct.SizeConstant.NODE_TYPE_ID_SIZE;
import static com.vesoft.nebula.driver.graph.decode.struct.SizeConstant.PATH_ELEMENT_NUM_SIZE;
import static com.vesoft.nebula.driver.graph.decode.struct.SizeConstant.PROPERTY_NUM_SIZE;
import static com.vesoft.nebula.driver.graph.decode.struct.SizeConstant.RECORD_FIELD_NUM_SIZE;
import static com.vesoft.nebula.driver.graph.decode.struct.SizeConstant.VALUE_TYPE_SIZE;

import com.vesoft.nebula.driver.graph.decode.BytesReader;
import com.vesoft.nebula.driver.graph.decode.ColumnType;
import com.vesoft.nebula.driver.graph.decode.DecodeUtils;
import java.nio.ByteOrder;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ValueTypeParser {
    private final ByteOrder byteOrder;

    public ValueTypeParser(ByteOrder byteOrder) {
        this.byteOrder = byteOrder;
    }

    public DataType getDataType(BytesReader reader) {
        return decodeValueType(reader);
    }

    public DataType decodeValueType(BytesReader reader) {
        ColumnType type = ColumnType.getColumnType(
                DecodeUtils.bytesToUInt8(reader.read(VALUE_TYPE_SIZE)));
        switch (type) {
            case COLUMN_TYPE_NULL:
            case COLUMN_TYPE_BOOL:
            case COLUMN_TYPE_INT8:
            case COLUMN_TYPE_UINT8:
            case COLUMN_TYPE_INT16:
            case COLUMN_TYPE_UINT16:
            case COLUMN_TYPE_INT32:
            case COLUMN_TYPE_UINT32:
            case COLUMN_TYPE_INT64:
            case COLUMN_TYPE_UINT64:
            case COLUMN_TYPE_FLOAT32:
            case COLUMN_TYPE_FLOAT64:
                // not support
                //case COLUMN_TYPE_BYTES:
            case COLUMN_TYPE_STRING:
            case COLUMN_TYPE_DATE:
            case COLUMN_TYPE_LOCALTIME:
            case COLUMN_TYPE_ZONEDTIME:
            case COLUMN_TYPE_LOCALDATETIME:
            case COLUMN_TYPE_ZONEDDATETIME:
            case COLUMN_TYPE_DURATION:
            case COLUMN_TYPE_REFERENCE:
            case COLUMN_TYPE_ANY:
            case COLUMN_TYPE_INVALID:
                return new BasicType(type);
            case COLUMN_TYPE_DECIMAL:
                // skip two value for decimal: 2 byte precision and 2 byte scale
                // these two value is useless for client, just read skip them
                reader.read(4);
                return new BasicType(type);
            case COLUMN_TYPE_NODE:
                Map<Integer, Map<String, DataType>> nodeTypes =
                        getPropertyNameAndTypeFromValueType(reader, NODE_TYPE_ID_SIZE, byteOrder);
                return new NodeType(nodeTypes);
            case COLUMN_TYPE_EDGE:
                Map<Integer, Map<String, DataType>> edgeTypes =
                        getPropertyNameAndTypeFromValueType(reader, EDGE_TYPE_ID_SIZE, byteOrder);
                return new EdgeType(edgeTypes);
            case COLUMN_TYPE_PATH:
                int elementNum = DecodeUtils.bytesToInt32(
                        reader.read(PATH_ELEMENT_NUM_SIZE), byteOrder);
                List<DataType> dataTypes = new ArrayList<>();
                for (int i = 0; i < elementNum; i++) {
                    dataTypes.add(decodeValueType(reader));
                }
                return new PathType(dataTypes);
            case COLUMN_TYPE_LIST:
                DataType dataType = decodeValueType(reader);
                return new ListType(dataType);
            case COLUMN_TYPE_RECORD:
                int fieldNum = DecodeUtils.bytesToInt32(
                        reader.read(RECORD_FIELD_NUM_SIZE), byteOrder);
                Map<String, DataType> fieldTypes = new HashMap<>();
                for (int i = 0; i < fieldNum; i++) {
                    String fieldName = reader.readSizedString(byteOrder);
                    fieldTypes.put(fieldName, decodeValueType(reader));
                }
                return new RecordType(fieldTypes);
            default:
                throw new RuntimeException("unsupported type:" + type);
        }
    }


    /**
     * decode value type for node and edge, get each type's property name and property data type
     * typeId -> (propName -> prop DataType）
     */
    private Map<Integer, Map<String, DataType>> getPropertyNameAndTypeFromValueType(
            BytesReader reader,
            int typeIdSize,
            ByteOrder byteOrder) {
        // 1-5: node or edge type number, 4 bytes
        int typeNum = DecodeUtils.bytesToInt32(reader.read(GRAPH_ELEMENT_TYPE_NUM_SIZE), byteOrder);

        Map<Integer, Map<String, DataType>> typeFields = new HashMap<>(typeNum);
        for (int i = 0; i < typeNum; i++) {
            // node type ID or edge type ID
            final int typeId =
                    typeIdSize == NODE_TYPE_ID_SIZE
                            ? DecodeUtils.bytesToInt16(reader.read(typeIdSize), byteOrder)
                            : DecodeUtils.bytesToInt32(reader.read(typeIdSize), byteOrder);
            // node or edge type property number, 4 bytes
            int typePropertyNum = DecodeUtils.bytesToInt32(
                    reader.read(PROPERTY_NUM_SIZE), byteOrder);
            Map<String, DataType> propertyAndType = new HashMap<>();
            //read the property name and data type for node or edge type, property name end with \0.
            for (int j = 0; j < typePropertyNum; j++) {
                String propertyName = reader.readSizedString(byteOrder);
                // data type, 1 byte
                DataType dataType = decodeValueType(reader);
                propertyAndType.put(propertyName, dataType);
            }
            typeFields.put(typeId, propertyAndType);
        }
        return typeFields;
    }

}
