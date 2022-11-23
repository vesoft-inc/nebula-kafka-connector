// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_VALUEOPS_H_
#define COMMON_DATATYPE_VALUEOPS_H_

#include "common/datatype/Value.h"
#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"
#include "common/datatype/thriftSerialization/DateOps-inl.h"
#include "common/datatype/thriftSerialization/DatetimeOps-inl.h"
#include "common/datatype/thriftSerialization/DurationOps-inl.h"
#include "common/datatype/thriftSerialization/EdgeOps-inl.h"
#include "common/datatype/thriftSerialization/ListOps-inl.h"
#include "common/datatype/thriftSerialization/MapOps-inl.h"
#include "common/datatype/thriftSerialization/NodeOps-inl.h"
#include "common/datatype/thriftSerialization/TimeOps-inl.h"


namespace apache {
namespace thrift {

namespace detail {

template <>
struct TccStructTraits<nebula::client::Value> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "boolVal") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_BOOL;
        } else if (_fname == "int8Val") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_I08;
        } else if (_fname == "int16Val") {
            fid = 3;
            _ftype = apache::thrift::protocol::T_I16;
        } else if (_fname == "int32Val") {
            fid = 4;
            _ftype = apache::thrift::protocol::T_I32;
        } else if (_fname == "int64Val") {
            fid = 5;
            _ftype = apache::thrift::protocol::T_I64;
        } else if (_fname == "floatVal") {
            fid = 6;
            _ftype = apache::thrift::protocol::T_DOUBLE;
        } else if (_fname == "doubleVal") {
            fid = 7;
            _ftype = apache::thrift::protocol::T_DOUBLE;
        } else if (_fname == "stringVal") {
            fid = 8;
            _ftype = apache::thrift::protocol::T_STRING;
        } else if (_fname == "listVal") {
            fid = 9;
            _ftype = apache::thrift::protocol::T_STRUCT;
        } else if (_fname == "mapVal") {
            fid = 10;
            _ftype = apache::thrift::protocol::T_STRUCT;
        } else if (_fname == "nodeVal") {
            fid = 11;
            _ftype = apache::thrift::protocol::T_STRUCT;
        } else if (_fname == "edgeVal") {
            fid = 12;
            _ftype = apache::thrift::protocol::T_STRUCT;
        } else if (_fname == "durationVal") {
            fid = 13;
            _ftype = apache::thrift::protocol::T_STRUCT;
        } else if (_fname == "localTimeVal") {
            fid = 14;
            _ftype = apache::thrift::protocol::T_STRUCT;
        } else if (_fname == "dateVal") {
            fid = 15;
            _ftype = apache::thrift::protocol::T_STRUCT;
        } else if (_fname == "localDatetimeVal") {
            fid = 16;
            _ftype = apache::thrift::protocol::T_STRUCT;
        }
    }
};

}  // namespace detail

inline constexpr apache::thrift::protocol::TType Cpp2Ops<nebula::client::Value>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Value>::write(Protocol* proto,
                                               nebula::client::Value const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("Value");

    switch (obj->getType()) {
        // Nothing to write for a Null value
        case nebula::client::Value::Type::kNull: {
            break;
        }
        case nebula::client::Value::Type::kBool: {
            xfer += proto->writeFieldBegin("boolVal", protocol::T_BOOL, 1);
            xfer += proto->writeBool(obj->getBool());
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kInt8: {
            xfer += proto->writeFieldBegin("int8Val", protocol::T_I08, 2);
            xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(
                    *proto, obj->getInt8());
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kInt16: {
            xfer += proto->writeFieldBegin("int16Val", protocol::T_I16, 3);
            xfer += detail::pm::protocol_methods<type_class::integral, int16_t>::write(
                    *proto, obj->getInt16());
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kInt32: {
            xfer += proto->writeFieldBegin("int32Val", protocol::T_I32, 4);
            xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::write(
                    *proto, obj->getInt32());
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kInt64: {
            xfer += proto->writeFieldBegin("int64Val", protocol::T_I64, 5);
            xfer += detail::pm::protocol_methods<type_class::integral, int64_t>::write(
                    *proto, obj->getInt64());
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kFloat: {
            xfer += proto->writeFieldBegin("floatVal", protocol::T_DOUBLE, 6);
            // Float type is not supported so cast to double
            xfer += proto->writeDouble(static_cast<double>(obj->getFloat()));
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kDouble: {
            xfer += proto->writeFieldBegin("doubleVal", protocol::T_DOUBLE, 7);
            xfer += proto->writeDouble(obj->getDouble());
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kString: {
            xfer += proto->writeFieldBegin("stringVal", protocol::T_STRING, 8);
            xfer += proto->writeBinary(obj->getString());
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kList: {
            xfer += proto->writeFieldBegin("listVal", protocol::T_STRUCT, 9);
            // If the type is a list, there is always a list object
            xfer += Cpp2Ops<nebula::client::List>::write(proto, obj->data_.list_);
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kMap: {
            xfer += proto->writeFieldBegin("mapVal", protocol::T_STRUCT, 10);
            // If the type is a map, there is always a map object
            xfer += Cpp2Ops<nebula::client::Map>::write(proto, obj->data_.map_);
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kNode: {
            xfer += proto->writeFieldBegin("nodeVal", protocol::T_STRUCT, 11);
            xfer += Cpp2Ops<nebula::client::Node>::write(proto, obj->data_.node_);
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kEdge: {
            xfer += proto->writeFieldBegin("edgeVal", protocol::T_STRUCT, 12);
            xfer += Cpp2Ops<nebula::client::Edge>::write(proto, obj->data_.edge_);
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kDuration: {
            xfer += proto->writeFieldBegin("durationVal", protocol::T_STRUCT, 13);
            xfer += Cpp2Ops<nebula::client::Duration>::write(proto, obj->data_.duration_);
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kLocalTime: {
            xfer += proto->writeFieldBegin("localTimeVal", protocol::T_STRUCT, 14);
            xfer += Cpp2Ops<nebula::client::LocalTime>::write(proto, &obj->getLocalTime());
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kDate: {
            xfer += proto->writeFieldBegin("dateVal", protocol::T_STRUCT, 15);
            xfer += Cpp2Ops<nebula::client::Date>::write(proto, &obj->getDate());
            xfer += proto->writeFieldEnd();
            break;
        }
        case nebula::client::Value::Type::kLocalDatetime: {
            xfer += proto->writeFieldBegin("localDatetimeVal", protocol::T_STRUCT, 16);
            xfer += Cpp2Ops<nebula::client::LocalDatetime>::write(proto,
                                                                  &obj->getLocalDatetime());
            xfer += proto->writeFieldEnd();
            break;
        }
        default: {
            LOG(FATAL) << "Unknown type: " << static_cast<int>(obj->getType());
        }
    }

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::Value>::read(Protocol* proto, nebula::client::Value* obj) {
    apache::thrift::detail::ProtocolReaderStructReadState<Protocol> readState;
    readState.fieldId = 0;

    readState.readStructBegin(proto);

    using apache::thrift::protocol::TProtocolException;

    readState.readFieldBegin(proto);
    if (readState.fieldType == apache::thrift::protocol::T_STOP) {
        obj->clear();
    } else {
        if (proto->kUsesFieldNames()) {
            detail::TccStructTraits<nebula::client::Value>::translateFieldName(
                    readState.fieldName(), readState.fieldId, readState.fieldType);
        }

        // Allocator
        auto mr = obj->get_allocator();

        switch (readState.fieldId) {
            case 1: {
                if (readState.fieldType == apache::thrift::protocol::T_BOOL) {
                    obj->type_ = nebula::client::Value::Type::kBool;
                    obj->data_.bool_ = false;
                    proto->readBool(obj->data_.bool_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            case 2: {
                if (readState.fieldType == apache::thrift::protocol::T_I08) {
                    obj->type_ = nebula::client::Value::Type::kInt8;
                    obj->data_.int8_ = 0;
                    detail::pm::protocol_methods<type_class::integral, int8_t>::read(
                            *proto, obj->data_.int8_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            case 3: {
                if (readState.fieldType == apache::thrift::protocol::T_I16) {
                    obj->type_ = nebula::client::Value::Type::kInt16;
                    obj->data_.int16_ = 0;
                    detail::pm::protocol_methods<type_class::integral, int16_t>::read(
                            *proto, obj->data_.int16_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            case 4: {
                if (readState.fieldType == apache::thrift::protocol::T_I32) {
                    obj->type_ = nebula::client::Value::Type::kInt32;
                    obj->data_.int32_ = 0;
                    detail::pm::protocol_methods<type_class::integral, int32_t>::read(
                            *proto, obj->data_.int32_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            case 5: {
                if (readState.fieldType == apache::thrift::protocol::T_I64) {
                    obj->type_ = nebula::client::Value::Type::kInt64;
                    obj->data_.int64_ = 0;
                    detail::pm::protocol_methods<type_class::integral, int64_t>::read(
                            *proto, obj->data_.int64_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            case 6: {
                if (readState.fieldType == apache::thrift::protocol::T_DOUBLE) {
                    obj->type_ = nebula::client::Value::Type::kFloat;
                    obj->data_.float_ = 0;

                    //  Float type is not supported in thrift, so use a temporary value to
                    //  accept the value from the protocol and cast to a float value then write
                    //  into the object
                    double temp = 0;
                    detail::pm::protocol_methods<type_class::floating_point, double>::read(
                            *proto, temp);
                    obj->data_.float_ = static_cast<float>(temp);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            case 7: {
                if (readState.fieldType == apache::thrift::protocol::T_DOUBLE) {
                    obj->type_ = nebula::client::Value::Type::kDouble;
                    obj->data_.double_ = 0;
                    detail::pm::protocol_methods<type_class::floating_point, double>::read(
                            *proto, obj->data_.double_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            case 8: {
                if (readState.fieldType == apache::thrift::protocol::T_STRING) {
                    obj->type_ = nebula::client::Value::Type::kString;
                    obj->data_.string_ = new std::pmr::string("", mr);
                    proto->readBinary(const_cast<std::pmr::string&>(obj->getString()));
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            case 9: {
                if (readState.fieldType == apache::thrift::protocol::T_STRUCT) {
                    obj->type_ = nebula::client::Value::Type::kList;
                    obj->data_.list_ = new nebula::client::List(mr);
                    Cpp2Ops<nebula::client::List>::read(proto, obj->data_.list_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            case 10: {
                if (readState.fieldType == apache::thrift::protocol::T_STRUCT) {
                    obj->type_ = nebula::client::Value::Type::kMap;
                    obj->data_.map_ = new nebula::client::Map(mr);
                    Cpp2Ops<nebula::client::Map>::read(proto, obj->data_.map_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            // Node type
            case 11: {
                if (readState.fieldType == apache::thrift::protocol::T_STRUCT) {
                    obj->type_ = nebula::client::Value::Type::kNode;
                    obj->data_.node_ = new nebula::client::Node(mr);
                    Cpp2Ops<nebula::client::Node>::read(proto, obj->data_.node_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            // Edge type
            case 12: {
                if (readState.fieldType == apache::thrift::protocol::T_STRUCT) {
                    obj->type_ = nebula::client::Value::Type::kEdge;
                    obj->data_.edge_ = new nebula::client::Edge(mr);
                    Cpp2Ops<nebula::client::Edge>::read(proto, obj->data_.edge_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            // Duration type
            case 13: {
                if (readState.fieldType == apache::thrift::protocol::T_STRUCT) {
                    obj->type_ = nebula::client::Value::Type::kDuration;
                    obj->data_.duration_ = new nebula::client::Duration(mr);
                    Cpp2Ops<nebula::client::Duration>::read(proto, obj->data_.duration_);
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            // LocalTime type
            case 14: {
                if (readState.fieldType == apache::thrift::protocol::T_STRUCT) {
                    obj->type_ = nebula::client::Value::Type::kLocalTime;
                    obj->data_.localTime_ = nebula::client::LocalTime();
                    Cpp2Ops<nebula::client::LocalTime>::read(proto, &obj->mutableLocalTime());
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            // Date type
            case 15: {
                if (readState.fieldType == apache::thrift::protocol::T_STRUCT) {
                    obj->type_ = nebula::client::Value::Type::kDate;
                    obj->data_.date_ = nebula::client::Date();
                    Cpp2Ops<nebula::client::Date>::read(proto, &obj->mutableDate());
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
            // LocalDatetime type
            case 16: {
                if (readState.fieldType == apache::thrift::protocol::T_STRUCT) {
                    obj->type_ = nebula::client::Value::Type::kLocalDatetime;
                    obj->data_.localDatetime_ = nebula::client::LocalDatetime();
                    Cpp2Ops<nebula::client::LocalDatetime>::read(proto,
                                                                 &obj->mutableLocalDatetime());
                } else {
                    proto->skip(readState.fieldType);
                }
                break;
            }
        }

        readState.readFieldEnd(proto);
        readState.readFieldBegin(proto);
        if (UNLIKELY(readState.fieldType != apache::thrift::protocol::T_STOP)) {
            using apache::thrift::protocol::TProtocolException;
            TProtocolException::throwUnionMissingStop();
        }
    }
    readState.readStructEnd(proto);
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Value>::serializedSize(Protocol const* proto,
                                                        nebula::client::Value const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Value");
    switch (obj->getType()) {
        case nebula::client::Value::Type::kNull: {
            break;
        }
        case nebula::client::Value::Type::kBool: {
            xfer += proto->serializedFieldSize("nVal", protocol::T_BOOL, 1);
            xfer += proto->serializedSizeBool(obj->getBool());
            break;
        }
        case nebula::client::Value::Type::kInt8: {
            xfer += proto->serializedFieldSize("bVal", protocol::T_I08, 2);
            detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
                    *proto, obj->getInt8());
            break;
        }
        case nebula::client::Value::Type::kInt16: {
            xfer += proto->serializedFieldSize("i16Val", protocol::T_I16, 3);
            detail::pm::protocol_methods<type_class::integral, int16_t>::serializedSize<false>(
                    *proto, obj->getInt16());
            break;
        }
        case nebula::client::Value::Type::kInt32: {
            xfer += proto->serializedFieldSize("i32Val", protocol::T_I32, 4);
            detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<false>(
                    *proto, obj->getInt32());
            break;
        }
        case nebula::client::Value::Type::kInt64: {
            xfer += proto->serializedFieldSize("i64Val", protocol::T_I64, 5);
            detail::pm::protocol_methods<type_class::integral, int64_t>::serializedSize<false>(
                    *proto, obj->getInt64());
            break;
        }
        case nebula::client::Value::Type::kFloat: {
            // Float type is not supported so cast to double
            xfer += proto->serializedFieldSize("floatVal", protocol::T_DOUBLE, 6);
            xfer += proto->serializedSizeDouble(static_cast<double>(obj->getFloat()));
            break;
        }
        case nebula::client::Value::Type::kDouble: {
            xfer += proto->serializedFieldSize("doubleVal", protocol::T_DOUBLE, 7);
            xfer += proto->serializedSizeDouble(obj->getDouble());
            break;
        }
        case nebula::client::Value::Type::kString: {
            xfer += proto->serializedFieldSize("stringVal", protocol::T_STRING, 8);
            xfer += proto->serializedSizeBinary(obj->getString());
            break;
        }
        case nebula::client::Value::Type::kList: {
            xfer += proto->serializedFieldSize("listVal", protocol::T_STRUCT, 9);
            // If the type is a list, there is always a list object
            xfer += Cpp2Ops<nebula::client::List>::serializedSize(proto, obj->data_.list_);
            break;
        }
        case nebula::client::Value::Type::kMap: {
            xfer += proto->serializedFieldSize("mapVal", protocol::T_STRUCT, 10);
            // If the type is a map, there is always a map object
            xfer += Cpp2Ops<nebula::client::Map>::serializedSize(proto, obj->data_.map_);
            break;
        }
        case nebula::client::Value::Type::kNode: {
            xfer += proto->serializedFieldSize("nodeVal", protocol::T_STRUCT, 11);
            xfer += Cpp2Ops<nebula::client::Node>::serializedSize(proto, obj->data_.node_);
            break;
        }
        case nebula::client::Value::Type::kEdge: {
            xfer += proto->serializedFieldSize("edgeVal", protocol::T_STRUCT, 12);
            xfer += Cpp2Ops<nebula::client::Edge>::serializedSize(proto, obj->data_.edge_);
            break;
        }
        case nebula::client::Value::Type::kDuration: {
            xfer += proto->serializedFieldSize("durationVal", protocol::T_STRUCT, 13);
            xfer += Cpp2Ops<nebula::client::Duration>::serializedSize(proto,
                                                                      obj->data_.duration_);
            break;
        }
        case nebula::client::Value::Type::kLocalTime: {
            xfer += proto->serializedFieldSize("localTimeVal", protocol::T_STRUCT, 14);
            xfer += Cpp2Ops<nebula::client::LocalTime>::serializedSize(proto,
                                                                       &obj->getLocalTime());
            break;
        }
        case nebula::client::Value::Type::kDate: {
            xfer += proto->serializedFieldSize("dateVal", protocol::T_STRUCT, 15);
            xfer += Cpp2Ops<nebula::client::Date>::serializedSize(proto, &obj->getDate());
            break;
        }
        case nebula::client::Value::Type::kLocalDatetime: {
            xfer += proto->serializedFieldSize("localDatetimeVal", protocol::T_STRUCT, 16);
            xfer += Cpp2Ops<nebula::client::LocalDatetime>::serializedSize(
                    proto, &obj->getLocalDatetime());
            break;
        }
        default: {
            LOG(FATAL) << "Unknown type " << static_cast<int>(obj->getType());
            break;
        }
    }

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Value>::serializedSizeZC(Protocol const* proto,
                                                          nebula::client::Value const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Value");

    switch (obj->getType()) {
        case nebula::client::Value::Type::kNull: {
            break;
        }
        case nebula::client::Value::Type::kBool: {
            xfer += proto->serializedFieldSize("nVal", protocol::T_BOOL, 1);
            xfer += proto->serializedSizeBool(obj->getBool());
            break;
        }
        case nebula::client::Value::Type::kInt8: {
            xfer += proto->serializedFieldSize("bVal", protocol::T_I08, 2);
            detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<true>(
                    *proto, obj->getInt8());
            break;
        }
        case nebula::client::Value::Type::kInt16: {
            xfer += proto->serializedFieldSize("i16Val", protocol::T_I16, 3);
            detail::pm::protocol_methods<type_class::integral, int16_t>::serializedSize<true>(
                    *proto, obj->getInt16());
            break;
        }
        case nebula::client::Value::Type::kInt32: {
            xfer += proto->serializedFieldSize("i32Val", protocol::T_I32, 4);
            detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<true>(
                    *proto, obj->getInt32());
            break;
        }
        case nebula::client::Value::Type::kInt64: {
            xfer += proto->serializedFieldSize("i64Val", protocol::T_I64, 5);
            detail::pm::protocol_methods<type_class::integral, int64_t>::serializedSize<true>(
                    *proto, obj->getInt64());
            break;
        }
        case nebula::client::Value::Type::kFloat: {
            xfer += proto->serializedFieldSize("floatVal", protocol::T_DOUBLE, 6);
            // Float type is not supported so cast to double
            xfer += proto->serializedSizeDouble(static_cast<double>(obj->getFloat()));
            break;
        }
        case nebula::client::Value::Type::kDouble: {
            xfer += proto->serializedFieldSize("doubleVal", protocol::T_DOUBLE, 7);
            xfer += proto->serializedSizeDouble(obj->getDouble());
            break;
        }
        case nebula::client::Value::Type::kString: {
            xfer += proto->serializedFieldSize("stringVal", protocol::T_STRING, 8);
            xfer += proto->serializedSizeZCBinary(obj->getString());
            break;
        }
        case nebula::client::Value::Type::kList: {
            xfer += proto->serializedFieldSize("listVal", protocol::T_STRUCT, 9);
            // If the type is a list, there is always a list object
            xfer += Cpp2Ops<nebula::client::List>::serializedSizeZC(proto, obj->data_.list_);
            break;
        }
        case nebula::client::Value::Type::kMap: {
            xfer += proto->serializedFieldSize("mapVal", protocol::T_STRUCT, 10);
            // If the type is a map, there is always a map object
            xfer += Cpp2Ops<nebula::client::Map>::serializedSizeZC(proto, obj->data_.map_);
            break;
        }
        case nebula::client::Value::Type::kNode: {
            xfer += proto->serializedFieldSize("nodeVal", protocol::T_STRUCT, 11);
            xfer += Cpp2Ops<nebula::client::Node>::serializedSizeZC(proto, obj->data_.node_);
            break;
        }
        case nebula::client::Value::Type::kEdge: {
            xfer += proto->serializedFieldSize("edgeVal", protocol::T_STRUCT, 12);
            xfer += Cpp2Ops<nebula::client::Edge>::serializedSizeZC(proto, obj->data_.edge_);
            break;
        }
        case nebula::client::Value::Type::kDuration: {
            xfer += proto->serializedFieldSize("durationVal", protocol::T_STRUCT, 13);
            xfer += Cpp2Ops<nebula::client::Duration>::serializedSizeZC(proto,
                                                                        obj->data_.duration_);
            break;
        }
        case nebula::client::Value::Type::kLocalTime: {
            xfer += proto->serializedFieldSize("localTimeVal", protocol::T_STRUCT, 14);
            xfer += Cpp2Ops<nebula::client::LocalTime>::serializedSizeZC(proto,
                                                                         &obj->getLocalTime());
            break;
        }
        case nebula::client::Value::Type::kDate: {
            xfer += proto->serializedFieldSize("dateVal", protocol::T_STRUCT, 15);
            xfer += Cpp2Ops<nebula::client::Date>::serializedSizeZC(proto, &obj->getDate());
            break;
        }
        case nebula::client::Value::Type::kLocalDatetime: {
            xfer += proto->serializedFieldSize("localDatetimeVal", protocol::T_STRUCT, 16);
            xfer += Cpp2Ops<nebula::client::LocalDatetime>::serializedSizeZC(
                    proto, &obj->getLocalDatetime());
            break;
        }
        default: {
            LOG(FATAL) << "Unknown type " << static_cast<int>(obj->getType());
            break;
        }
    }

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
#endif  // COMMON_DATATYPE_VALUEOPS_H_
