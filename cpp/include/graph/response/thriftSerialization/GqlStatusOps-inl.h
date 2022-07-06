// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef GRAPH_RESPONSE_GQLSTATUSOPS_INL_H_
#define GRAPH_RESPONSE_GQLSTATUSOPS_INL_H_

#include <thrift/lib/cpp2/GeneratedCodeHelper.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>
#include <thrift/lib/cpp2/protocol/ProtocolReaderStructReadState.h>

#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"
#include "graph/response/GqlStatus.h"
#include "interface/GraphCpp2Ops.h"

namespace apache {
namespace thrift {

namespace detail {

template <>
struct TccStructTraits<nebula::client::GQLStatus> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "status") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_STRING;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::GQLStatus>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::GQLStatus>::write(Protocol* proto,
                                                   nebula::client::GQLStatus const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("GQLStatus");

    // Write field status (required)
    xfer += proto->writeFieldBegin("stringVal", protocol::T_STRING, 1);
    xfer += proto->writeBinary(obj->status_);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::GQLStatus>::read(Protocol* proto, nebula::client::GQLStatus* obj) {
    using apache::thrift::protocol::TProtocolException;
    detail::ProtocolReaderStructReadState<Protocol> readState;
    readState.readStructBegin(proto);

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_STRING))) {
        goto _loop;
    }

_readField_status : {
    obj->status_.clear();
    proto->readBinary(obj->status_);
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
        detail::TccStructTraits<nebula::client::GQLStatus>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_STRING)) {
                goto _readField_status;
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
uint32_t Cpp2Ops<nebula::client::GQLStatus>::serializedSize(
        Protocol const* proto, nebula::client::GQLStatus const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("GQLStatus");

    xfer += proto->serializedFieldSize("status", apache::thrift::protocol::T_STRING, 1);
    xfer += proto->serializedSizeBinary(obj->status_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::GQLStatus>::serializedSizeZC(
        Protocol const* proto, nebula::client::GQLStatus const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("GQLStatus");

    xfer += proto->serializedFieldSize("status", apache::thrift::protocol::T_STRING, 1);
    xfer += proto->serializedSizeZCBinary(obj->status_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache

#endif  // GRAPH_RESPONSE_GQLSTATUSOPS_INL_H_
