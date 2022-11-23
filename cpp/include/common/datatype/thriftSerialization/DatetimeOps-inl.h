// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_THRIFTSERIALIZATION_DATETIMEOPS_INL_H_
#define COMMON_DATATYPE_THRIFTSERIALIZATION_DATETIMEOPS_INL_H_

#include <thrift/lib/cpp/protocol/TType.h>
#include <thrift/lib/cpp2/TypeClass.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>

#include "common/datatype/Datetime.h"
namespace apache {
namespace thrift {

/**************************************
 *
 * Ops for class LocalDatetime
 *
 *************************************/
namespace detail {

template <>
struct TccStructTraits<nebula::client::LocalDatetime> {
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
        } else if (_fname == "hour") {
            fid = 4;
            _ftype = apache::thrift::protocol::T_BYTE;
        } else if (_fname == "minute") {
            fid = 5;
            _ftype = apache::thrift::protocol::T_BYTE;
        } else if (_fname == "sec") {
            fid = 6;
            _ftype = apache::thrift::protocol::T_BYTE;
        } else if (_fname == "microsec") {
            fid = 7;
            _ftype = apache::thrift::protocol::T_I32;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::LocalDatetime>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::LocalDatetime>::write(
        Protocol* proto, nebula::client::LocalDatetime const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("LocalDatetime");

    xfer += proto->writeFieldBegin("year", apache::thrift::protocol::T_I16, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int16_t>::write(*proto,
                                                                               obj->year);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("month", apache::thrift::protocol::T_BYTE, 2);
    xfer += detail::pm::protocol_methods<thrift::type_class::integral, int8_t>::write(
            *proto, obj->month);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("day", apache::thrift::protocol::T_BYTE, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(*proto, obj->day);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("hour", apache::thrift::protocol::T_BYTE, 4);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(*proto,
                                                                              obj->hour);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("minute", apache::thrift::protocol::T_BYTE, 5);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(*proto,
                                                                              obj->minute);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("sec", apache::thrift::protocol::T_BYTE, 6);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::write(*proto, obj->sec);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("microsec", apache::thrift::protocol::T_I32, 7);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::write(*proto,
                                                                               obj->microsec);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::LocalDatetime>::read(Protocol* proto,
                                                  nebula::client::LocalDatetime* obj) {
    detail::ProtocolReaderStructReadState<Protocol> readState;

    readState.readStructBegin(proto);

    using apache::thrift::protocol::TProtocolException;

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_I16))) {
        goto _loop;
    }

_readField_year : {
    int16_t year;
    detail::pm::protocol_methods<type_class::integral, int16_t>::read(*proto, year);
    obj->year = year;
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_BYTE))) {
        goto _loop;
    }

_readField_month : {
    int8_t month;
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, month);
    obj->month = month;
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 2, 3, protocol::T_BYTE))) {
        goto _loop;
    }

_readField_day : {
    int8_t day;
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, day);
    obj->day = day;
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 3, 4, protocol::T_BYTE))) {
        goto _loop;
    }

_readField_hour : {
    int8_t hour;
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, hour);
    obj->hour = hour;
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 4, 5, protocol::T_BYTE))) {
        goto _loop;
    }

_readField_minute : {
    int8_t minute;
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, minute);
    obj->minute = minute;
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 5, 6, protocol::T_BYTE))) {
        goto _loop;
    }

_readField_sec : {
    int8_t sec;
    detail::pm::protocol_methods<type_class::integral, int8_t>::read(*proto, sec);
    obj->sec = sec;
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 6, 7, protocol::T_I32))) {
        goto _loop;
    }

_readField_microsec : {
    int32_t microsec;
    detail::pm::protocol_methods<type_class::integral, int32_t>::read(*proto, microsec);
    obj->microsec = microsec;
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 7, 8, protocol::T_I32))) {
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
        detail::TccStructTraits<nebula::client::LocalDatetime>::translateFieldName(
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
        case 4: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_BYTE)) {
                goto _readField_hour;
            } else {
                goto _skip;
            }
        }
        case 5: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_BYTE)) {
                goto _readField_minute;
            } else {
                goto _skip;
            }
        }
        case 6: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_BYTE)) {
                goto _readField_sec;
            } else {
                goto _skip;
            }
        }
        case 7: {
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
uint32_t Cpp2Ops<nebula::client::LocalDatetime>::serializedSize(
        Protocol const* proto, nebula::client::LocalDatetime const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("LocalDatetime");

    xfer += proto->serializedFieldSize("year", apache::thrift::protocol::T_I16, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int16_t>::serializedSize<false>(
            *proto, obj->year);

    xfer += proto->serializedFieldSize("month", apache::thrift::protocol::T_BYTE, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->month);

    xfer += proto->serializedFieldSize("day", apache::thrift::protocol::T_BYTE, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->day);

    xfer += proto->serializedFieldSize("hour", apache::thrift::protocol::T_BYTE, 4);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->hour);

    xfer += proto->serializedFieldSize("minute", apache::thrift::protocol::T_BYTE, 5);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->minute);

    xfer += proto->serializedFieldSize("sec", apache::thrift::protocol::T_BYTE, 6);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->sec);

    xfer += proto->serializedFieldSize("microsec", apache::thrift::protocol::T_I32, 7);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<false>(
            *proto, obj->microsec);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::LocalDatetime>::serializedSizeZC(
        Protocol const* proto, nebula::client::LocalDatetime const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("LocalDatetime");

    xfer += proto->serializedFieldSize("year", apache::thrift::protocol::T_I16, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int16_t>::serializedSize<false>(
            *proto, obj->year);

    xfer += proto->serializedFieldSize("month", apache::thrift::protocol::T_BYTE, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->month);

    xfer += proto->serializedFieldSize("day", apache::thrift::protocol::T_BYTE, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->day);

    xfer += proto->serializedFieldSize("hour", apache::thrift::protocol::T_BYTE, 4);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->hour);

    xfer += proto->serializedFieldSize("minute", apache::thrift::protocol::T_BYTE, 5);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->minute);

    xfer += proto->serializedFieldSize("sec", apache::thrift::protocol::T_BYTE, 6);
    xfer += detail::pm::protocol_methods<type_class::integral, int8_t>::serializedSize<false>(
            *proto, obj->sec);

    xfer += proto->serializedFieldSize("microsec", apache::thrift::protocol::T_I32, 7);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::serializedSize<false>(
            *proto, obj->microsec);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
#endif  // COMMON_DATATYPE_THRIFTSERIALIZATION_DATETIMEOPS_INL_H_
