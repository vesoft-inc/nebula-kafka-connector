// Copyright (c) 2022 vesoft inc. All rights reserved.
#ifndef COMMON_DATATYPE_THRIFTSERIALIZATION_DURATIONOPS_INL_H_
#define COMMON_DATATYPE_THRIFTSERIALIZATION_DURATIONOPS_INL_H_

#include <thrift/lib/cpp/protocol/TType.h>
#include <thrift/lib/cpp2/TypeClass.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>

// #include "common/datatype/Duration.h"
#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"
#include "common/utils/Types.h"

namespace apache {
namespace thrift {

/**************************************
 *
 * Ops for class Duration
 *
 *************************************/
namespace detail {

template <>
struct TccStructTraits<nebula::client::Duration> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "seconds") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_I64;
        } else if (_fname == "microseconds") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_I32;
        } else if (_fname == "months") {
            fid = 3;
            _ftype = apache::thrift::protocol::T_I32;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::Duration>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Duration>::write(Protocol* proto, nebula::client::Duration const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("Duration");

    xfer += proto->writeFieldBegin("seconds", apache::thrift::protocol::T_I64, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int64_t>::write(*proto,
                                                                               obj->seconds);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("microseconds", apache::thrift::protocol::T_I32, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::write(
            *proto, obj->microseconds);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("months", apache::thrift::protocol::T_I32, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::write(*proto,
                                                                               obj->months);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();

    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::Duration>::read(Protocol* proto, nebula::client::Duration* obj) {
    detail::ProtocolReaderStructReadState<Protocol> readState;

    readState.readStructBegin(proto);

    using apache::thrift::protocol::TProtocolException;

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_I64))) {
        goto _loop;
    }

_readField_seconds : {
    detail::pm::protocol_methods<type_class::integral, int64_t>::read(*proto, obj->seconds);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_I32))) {
        goto _loop;
    }

_readField_microseconds : {
    detail::pm::protocol_methods<type_class::integral, int32_t>::read(*proto,
                                                                      obj->microseconds);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 2, 3, protocol::T_I32))) {
        goto _loop;
    }

_readField_months : {
    detail::pm::protocol_methods<type_class::integral, int32_t>::read(*proto, obj->months);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 3, 0, protocol::T_STOP))) {
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
        detail::TccStructTraits<nebula::client::Duration>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I64)) {
                goto _readField_seconds;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I32)) {
                goto _readField_microseconds;
            } else {
                goto _skip;
            }
        }
        case 3: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I32)) {
                goto _readField_months;
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
uint32_t Cpp2Ops<nebula::client::Duration>::serializedSize(Protocol const* proto,
                                                   nebula::client::Duration const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Duration");

    xfer += proto->serializedFieldSize("seconds", apache::thrift::protocol::T_I64, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int64_t>::serializedSize<false>(
            *proto, obj->seconds);

    xfer += proto->serializedFieldSize("microseconds", apache::thrift::protocol::T_I32, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<false>(
            *proto, obj->microseconds);

    xfer += proto->serializedFieldSize("months", apache::thrift::protocol::T_I32, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<false>(
            *proto, obj->months);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Duration>::serializedSizeZC(Protocol const* proto,
                                                     nebula::client::Duration const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Duration");

    xfer += proto->serializedFieldSize("seconds", apache::thrift::protocol::T_I64, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int64_t>::serializedSize<false>(
            *proto, obj->seconds);

    xfer += proto->serializedFieldSize("microseconds", apache::thrift::protocol::T_I32, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<false>(
            *proto, obj->microseconds);

    xfer += proto->serializedFieldSize("months", apache::thrift::protocol::T_I32, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<false>(
            *proto, obj->months);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
#endif  // COMMON_DATATYPE_THRIFTSERIALIZATION_DURATIONOPS_INL_H_
