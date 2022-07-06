// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef GRAPH_RESPONSE_RESPONSEOPS_INL_H_
#define GRAPH_RESPONSE_RESPONSEOPS_INL_H_

#include <thrift/lib/cpp/protocol/TType.h>
#include <thrift/lib/cpp2/GeneratedCodeHelper.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>
#include <thrift/lib/cpp2/protocol/ProtocolReaderStructReadState.h>

#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"
#include "graph/response/ExecutionOutcome.h"
#include "graph/response/ExecutionResponse.h"
#include "graph/response/thriftSerialization/ExecutionOutcomeOps-inl.h"

namespace apache {
namespace thrift {

namespace detail {

template <>
struct TccStructTraits<nebula::client::ExecutionResponse> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "executionOutcome") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_STRUCT;
        } else if (_fname == "latencyInUs") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_I64;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::ExecutionResponse>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::ExecutionResponse>::write(
        Protocol* proto, nebula::client::ExecutionResponse const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("ExecutionResponse");

    // Write field status (required)
    xfer += proto->writeFieldBegin("executionOutcome", protocol::T_STRUCT, 1);
    xfer += Cpp2Ops<nebula::client::ExecutionOutcome>::write(proto, &obj->executionOutcome_);
    xfer += proto->writeFieldEnd();

    // Write field latencyInUs (required)
    xfer += proto->writeFieldBegin("latencyInUs", protocol::T_I64, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int64_t>::write(
            *proto, obj->latencyInUs_);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::ExecutionResponse>::read(Protocol* proto,
                                                      nebula::client::ExecutionResponse* obj) {
    using apache::thrift::protocol::TProtocolException;
    detail::ProtocolReaderStructReadState<Protocol> readState;
    readState.readStructBegin(proto);

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_STRUCT))) {
        goto _loop;
    }

_readField_executionOutcome : {
    Cpp2Ops<nebula::client::ExecutionOutcome>::read(proto, &obj->executionOutcome_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_I64))) {
        goto _loop;
    }

_readField_latencyInUs : {
    obj->latencyInUs_ = 0;
    detail::pm::protocol_methods<type_class::integral, int64_t>::read(*proto,
                                                                      obj->latencyInUs_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 2, 0, protocol::T_STOP))) {
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
        detail::TccStructTraits<nebula::client::ExecutionResponse>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_STRUCT)) {
                goto _readField_executionOutcome;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I64)) {
                goto _readField_latencyInUs;
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
uint32_t Cpp2Ops<nebula::client::ExecutionResponse>::serializedSize(
        Protocol const* proto, nebula::client::ExecutionResponse const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("ExecutionResponse");

    xfer += proto->serializedFieldSize(
            "executionOutcome", apache::thrift::protocol::T_STRUCT, 1);
    xfer += Cpp2Ops<nebula::client::ExecutionOutcome>::serializedSize(proto,
                                                                      &obj->executionOutcome_);

    xfer += proto->serializedFieldSize("latencyInUs", apache::thrift::protocol::T_I64, 2);
    detail::pm::protocol_methods<type_class::integral, int64_t>::serializedSize<false>(
            *proto, obj->latencyInUs_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::ExecutionResponse>::serializedSizeZC(
        Protocol const* proto, nebula::client::ExecutionResponse const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("ExecutionResponse");

    xfer += proto->serializedFieldSize(
            "executionOutcome", apache::thrift::protocol::T_STRUCT, 1);
    xfer += Cpp2Ops<nebula::client::ExecutionOutcome>::serializedSizeZC(
            proto, &obj->executionOutcome_);

    xfer += proto->serializedFieldSize("latencyInUs", apache::thrift::protocol::T_I64, 2);
    detail::pm::protocol_methods<type_class::integral, int64_t>::serializedSize<false>(
            *proto, obj->latencyInUs_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache

#endif  // GRAPH_RESPONSE_RESPONSEOPS_INL_H_
