// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_THRIFTSERIALIZATION_DATEOPS_INL_H_
#define COMMON_DATATYPE_THRIFTSERIALIZATION_DATEOPS_INL_H_

#include <thrift/lib/cpp/protocol/TType.h>
#include <thrift/lib/cpp2/TypeClass.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>

#include "common/datatype/Date.h"

namespace apache {
namespace thrift {

/**************************************
 *
 * Ops for class Date
 *
 *************************************/
namespace detail {

template <>
struct TccStructTraits<nebula::client::Date> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "year") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_I16;
        } else if (_fname == "month") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_BYTE;
        } else if (_fname == "day") {
            fid = 3;
            _ftype = apache::thrift::protocol::T_BYTE;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::Date>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Date>::write(Protocol* proto, nebula::client::Date const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("Date");

    xfer += proto->writeFieldBegin("year", apache::thrift::protocol::T_I16, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int16_t>::write(*proto,
                                                                               obj->year);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("month", apache::thrift::protocol::T_BYTE, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(*proto,
                                                                              obj->month);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("day", apache::thrift::protocol::T_BYTE, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(*proto, obj->day);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();

    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::Date>::read(Protocol* proto, nebula::client::Date* obj) {
    detail::ProtocolReaderStructReadState<Protocol> readState;

    readState.readStructBegin(proto);

    using apache::thrift::protocol::TProtocolException;

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_I16))) {
        goto _loop;
    }

_readField_year : {
    detail::pm::protocol_methods<type_class::integral, int16_t>::read(*proto, obj->year);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_BYTE))) {
        goto _loop;
    }

_readField_month : {
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, obj->month);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 2, 3, protocol::T_BYTE))) {
        goto _loop;
    }

_readField_day : {
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, obj->day);
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
        detail::TccStructTraits<nebula::client::Date>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I16)) {
                goto _readField_year;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_BYTE)) {
                goto _readField_month;
            } else {
                goto _skip;
            }
        }
        case 3: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_BYTE)) {
                goto _readField_day;
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
uint32_t Cpp2Ops<nebula::client::Date>::serializedSize(Protocol const* proto, nebula::client::Date const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Date");

    xfer += proto->serializedFieldSize("year", apache::thrift::protocol::T_I16, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int16_t>::serializedSize<false>(
            *proto, obj->year);

    xfer += proto->serializedFieldSize("month", apache::thrift::protocol::T_BYTE, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->month);

    xfer += proto->serializedFieldSize("day", apache::thrift::protocol::T_BYTE, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->day);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Date>::serializedSizeZC(Protocol const* proto,
                                                 nebula::client::Date const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Date");

    xfer += proto->serializedFieldSize("year", apache::thrift::protocol::T_I16, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int16_t>::serializedSize<false>(
            *proto, obj->year);

    xfer += proto->serializedFieldSize("month", apache::thrift::protocol::T_BYTE, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->month);

    xfer += proto->serializedFieldSize("day", apache::thrift::protocol::T_BYTE, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->day);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
#endif  // COMMON_DATATYPE_THRIFTSERIALIZATION_DATEOPS_INL_H_
