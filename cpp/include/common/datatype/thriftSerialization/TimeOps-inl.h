// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_THRIFTSERIALIZATION_TIMEOPS_INL_H_
#define COMMON_DATATYPE_THRIFTSERIALIZATION_TIMEOPS_INL_H_

#include <thrift/lib/cpp/protocol/TType.h>
#include <thrift/lib/cpp2/TypeClass.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>

#include "common/datatype/Time.h"
namespace apache {
namespace thrift {

/**************************************
 *
 * Ops for class LocalTime
 *
 *************************************/
namespace detail {

template <>
struct TccStructTraits<nebula::client::LocalTime> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "hour") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_BYTE;
        } else if (_fname == "minute") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_BYTE;
        } else if (_fname == "sec") {
            fid = 3;
            _ftype = apache::thrift::protocol::T_BYTE;
        } else if (_fname == "microsec") {
            fid = 4;
            _ftype = apache::thrift::protocol::T_I32;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::LocalTime>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::LocalTime>::write(Protocol* proto, nebula::client::LocalTime const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("LocalTime");

    xfer += proto->writeFieldBegin("hour", apache::thrift::protocol::T_BYTE, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(*proto,
                                                                              obj->hour);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("minute", apache::thrift::protocol::T_BYTE, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(*proto,
                                                                              obj->minute);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("sec", apache::thrift::protocol::T_BYTE, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(*proto, obj->sec);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("microsec", apache::thrift::protocol::T_I32, 4);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::write(*proto,
                                                                               obj->microsec);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::LocalTime>::read(Protocol* proto, nebula::client::LocalTime* obj) {
    detail::ProtocolReaderStructReadState<Protocol> readState;

    readState.readStructBegin(proto);

    using apache::thrift::protocol::TProtocolException;

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_I16))) {
        goto _loop;
    }

_readField_hour : {
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, obj->hour);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_BYTE))) {
        goto _loop;
    }

_readField_minute : {
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, obj->minute);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 2, 3, protocol::T_BYTE))) {
        goto _loop;
    }

_readField_sec : {
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, obj->sec);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 3, 4, protocol::T_I32))) {
        goto _loop;
    }

_readField_microsec : {
    detail::pm::protocol_methods<type_class::integral, int32_t>::read(*proto, obj->microsec);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 4, 5, protocol::T_I32))) {
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
        detail::TccStructTraits<nebula::client::LocalTime>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_BYTE)) {
                goto _readField_hour;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_BYTE)) {
                goto _readField_minute;
            } else {
                goto _skip;
            }
        }
        case 3: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_BYTE)) {
                goto _readField_sec;
            } else {
                goto _skip;
            }
        }
        case 4: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I32)) {
                goto _readField_microsec;
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
uint32_t Cpp2Ops<nebula::client::LocalTime>::serializedSize(Protocol const* proto,
                                                    nebula::client::LocalTime const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("LocalTime");

    xfer += proto->serializedFieldSize("hour", apache::thrift::protocol::T_BYTE, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->hour);

    xfer += proto->serializedFieldSize("minute", apache::thrift::protocol::T_BYTE, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->minute);

    xfer += proto->serializedFieldSize("sec", apache::thrift::protocol::T_BYTE, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->sec);

    xfer += proto->serializedFieldSize("microsec", apache::thrift::protocol::T_I32, 4);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<false>(
            *proto, obj->microsec);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::LocalTime>::serializedSizeZC(Protocol const* proto,
                                                      nebula::client::LocalTime const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("LocalTime");

    xfer += proto->serializedFieldSize("hour", apache::thrift::protocol::T_BYTE, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->hour);

    xfer += proto->serializedFieldSize("minute", apache::thrift::protocol::T_BYTE, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->minute);

    xfer += proto->serializedFieldSize("sec", apache::thrift::protocol::T_BYTE, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->sec);

    xfer += proto->serializedFieldSize("microsec", apache::thrift::protocol::T_I32, 4);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<false>(
            *proto, obj->microsec);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
#endif  // COMMON_DATATYPE_THRIFTSERIALIZATION_TIMEOPS_INL_H_
