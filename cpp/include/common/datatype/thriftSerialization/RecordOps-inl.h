// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_RECORDOPS_H_
#define COMMON_DATATYPE_RECORDOPS_H_

#include <thrift/lib/cpp2/GeneratedCodeHelper.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>
#include <thrift/lib/cpp2/protocol/ProtocolReaderStructReadState.h>

#include "common/datatype/Record.h"
#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"

// TODO(Aiee) The implementation of Record is not complete.
namespace apache {
namespace thrift {

// Ops for class FieldType
namespace detail {

template <>
struct TccStructTraits<nebula::client::FieldType> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED protocol::TType& _ftype) {
        if (_fname == "filedName") {
            fid = 1;
            _ftype = protocol::T_STRING;
        } else if (_fname == "valueType") {
            fid = 2;
            _ftype = protocol::T_I32;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::FieldType>::thriftType() {
    return protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::FieldType>::write(Protocol* proto,
                                                   nebula::client::FieldType const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("FieldType");

    // Write field filedName (required)
    xfer += proto->writeFieldBegin("filedName", protocol::T_STRING, 1);
    xfer += proto->writeBinary(obj->getFieldName());
    xfer += proto->writeFieldEnd();

    // Write field valueType (required)
    xfer += proto->writeFieldBegin("valueType", protocol::T_I32, 2);
    xfer += detail::pm::protocol_methods<type_class::enumeration,
                                         nebula::client::ValueType>::write(*proto,
                                                                           obj->getType());
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::FieldType>::read(Protocol* proto, nebula::client::FieldType* obj) {
    using protocol::TProtocolException;
    detail::ProtocolReaderStructReadState<Protocol> readState;
    readState.readStructBegin(proto);

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_STRING))) {
        goto _loop;
    }

_readField_filedName : {
    obj->fieldName_.clear();
    proto->readBinary(const_cast<std::pmr::string&>(obj->fieldName_));
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_I32))) {
        goto _loop;
    }

_readField_valueType : {
    obj->valueType_ = nebula::client::Value::Type::kNull;
    detail::pm::protocol_methods<::apache::thrift::type_class::enumeration,
                                 ::nebula::client::Value::Type>::read(*proto, obj->valueType_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 2, 0, protocol::T_STOP))) {
        goto _loop;
    }

_end:
    readState.readStructEnd(proto);

    return;

_loop:
    if (readState.fieldType == protocol::T_STOP) {
        goto _end;
    }

    if (proto->kUsesFieldNames()) {
        detail::TccStructTraits<nebula::client::FieldType>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == protocol::T_STRING)) {
                goto _readField_filedName;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == protocol::T_I32)) {
                goto _readField_valueType;
            } else {
                goto _skip;
            }
        }

        default: {
_skip:
            proto->skip(readState.fieldType);
            readState.readFieldEnd(proto);
            readState.readFieldBeginNoInline(proto);
            goto _loop;
        }
    }
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::FieldType>::serializedSize(
        Protocol const* proto, nebula::client::FieldType const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("FieldType");

    xfer += proto->serializedFieldSize("filedName", protocol::T_STRING, 1);
    xfer += proto->serializedSizeBinary(obj->getFieldName());

    xfer += proto->serializedFieldSize("valueType", protocol::T_I32, 2);
    xfer += detail::pm::protocol_methods<type_class::enumeration, nebula::client::ValueType>::
            serializedSize<false>(*proto, obj->getType());

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::FieldType>::serializedSizeZC(
        Protocol const* proto, nebula::client::FieldType const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("FieldType");

    xfer += proto->serializedFieldSize("filedName", protocol::T_STRING, 1);
    xfer += proto->serializedSizeZCBinary(obj->getFieldName());

    xfer += proto->serializedFieldSize("valueType", protocol::T_I32, 2);
    xfer += detail::pm::protocol_methods<type_class::enumeration, nebula::client::ValueType>::
            serializedSize<false>(*proto, obj->getType());

    xfer += proto->serializedSizeStop();
    return xfer;
}

// Ops for nebula::client::RawRecord
namespace detail {

template <>
struct TccStructTraits<nebula::client::RawRecord> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "values") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_LIST;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::RawRecord>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::RawRecord>::write(Protocol* proto,
                                                   nebula::client::RawRecord const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("RawRecord");

    xfer += proto->writeFieldBegin("values", apache::thrift::protocol::T_LIST, 1);
    xfer += detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::pmr::vector<nebula::client::Value>>::write(*proto, obj->values_);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::RawRecord>::read(Protocol* proto, nebula::client::RawRecord* obj) {
    detail::ProtocolReaderStructReadState<Protocol> readState;

    readState.readStructBegin(proto);

    using apache::thrift::protocol::TProtocolException;

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_LIST))) {
        goto _loop;
    }

_readField_values : {
    obj->values_ = std::pmr::vector<nebula::client::Value>();
    detail::pm::protocol_methods<type_class::list<type_class::structure>,
                                 std::pmr::vector<nebula::client::Value>>::read(*proto,
                                                                                obj->values_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 0, protocol::T_STOP))) {
        goto _loop;
    }

_end:
    readState.readStructEnd(proto);

    return;

_loop:
    if (readState.fieldType == apache::thrift::protocol::T_STOP) {
        goto _end;
    }

    if (proto->kUsesFieldNames()) {
        detail::TccStructTraits<nebula::client::RawRecord>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_LIST)) {
                goto _readField_values;
            } else {
                goto _skip;
            }
        }
        default: {
_skip:
            proto->skip(readState.fieldType);
            readState.readFieldEnd(proto);
            readState.readFieldBeginNoInline(proto);
            goto _loop;
        }
    }
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::RawRecord>::serializedSize(
        Protocol const* proto, nebula::client::RawRecord const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("RawRecord");

    xfer += proto->serializedFieldSize("values", apache::thrift::protocol::T_LIST, 1);
    xfer += detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::pmr::vector<nebula::client::Value>>::serializedSize<false>(*proto,
                                                                            obj->values_);
    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::RawRecord>::serializedSizeZC(
        Protocol const* proto, nebula::client::RawRecord const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("RawRecord");

    xfer += proto->serializedFieldSize("values", apache::thrift::protocol::T_LIST, 1);
    xfer += detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::pmr::vector<nebula::client::Value>>::serializedSize<false>(*proto,
                                                                            obj->values_);
    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
#endif  // COMMON_DATATYPE_RECORDOPS_H_
