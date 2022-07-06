// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef GRAPH_RESPONSE_EXECUTIONOPS_INL_H_
#define GRAPH_RESPONSE_EXECUTIONOPS_INL_H_

#include <thrift/lib/cpp2/GeneratedCodeHelper.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>
#include <thrift/lib/cpp2/protocol/ProtocolReaderStructReadState.h>

#include "common/datatype/BindingTable.h"
#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"
#include "common/datatype/thriftSerialization/ValueOps-inl.h"
#include "graph/response/ExecutionOutcome.h"
#include "graph/response/thriftSerialization/GqlStatusOps-inl.h"

namespace apache {
namespace thrift {

namespace detail {

template <>
struct TccStructTraits<nebula::client::ExecutionOutcome> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "gqlStatus") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_STRUCT;
        } else if (_fname == "result") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_STRUCT;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::ExecutionOutcome>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::ExecutionOutcome>::write(
        Protocol* proto, nebula::client::ExecutionOutcome const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("ExecutionOutcome");

    // Write field status (required)

    xfer += proto->writeFieldBegin("gqlStatus", protocol::T_STRUCT, 1);
    xfer += Cpp2Ops<nebula::client::GQLStatus>::write(proto, &obj->gqlStatus_);
    xfer += proto->writeFieldEnd();

    // Write field result (optional)
    if (obj->result_.has_value()) {
        xfer += proto->writeFieldBegin("result", protocol::T_STRUCT, 2);
        xfer += Cpp2Ops<nebula::client::BindingTable>::write(proto, &obj->result_.value());
        xfer += proto->writeFieldEnd();
    }

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::ExecutionOutcome>::read(Protocol* proto,
                                                     nebula::client::ExecutionOutcome* obj) {
    using apache::thrift::protocol::TProtocolException;
    detail::ProtocolReaderStructReadState<Protocol> readState;
    readState.readStructBegin(proto);

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_STRUCT))) {
        goto _loop;
    }

_readField_gqlStatus : { Cpp2Ops<nebula::client::GQLStatus>::read(proto, &obj->gqlStatus_); }

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_STRUCT))) {
        goto _loop;
    }

_readField_result : {
    obj->result_ = nebula::client::BindingTable();
    Cpp2Ops<nebula::client::BindingTable>::read(proto, &obj->result_.value());
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
        detail::TccStructTraits<nebula::client::ExecutionOutcome>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_STRUCT)) {
                goto _readField_gqlStatus;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_STRUCT)) {
                goto _readField_result;
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
uint32_t Cpp2Ops<nebula::client::ExecutionOutcome>::serializedSize(
        Protocol const* proto, nebula::client::ExecutionOutcome const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("ExecutionOutcome");

    xfer += proto->serializedFieldSize("gqlStatus", apache::thrift::protocol::T_STRUCT, 1);
    xfer += Cpp2Ops<nebula::client::GQLStatus>::serializedSize(proto, &obj->gqlStatus_);

    if (obj->result_.has_value()) {
        xfer += proto->serializedFieldSize("result", apache::thrift::protocol::T_STRUCT, 2);
        xfer += Cpp2Ops<nebula::client::BindingTable>::serializedSize(proto,
                                                                      &obj->result_.value());
    }
    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::ExecutionOutcome>::serializedSizeZC(
        Protocol const* proto, nebula::client::ExecutionOutcome const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("ExecutionOutcome");

    xfer += proto->serializedFieldSize("gqlStatus", apache::thrift::protocol::T_STRUCT, 1);
    xfer += Cpp2Ops<nebula::client::GQLStatus>::serializedSizeZC(proto, &obj->gqlStatus_);

    if (obj->result_.has_value()) {
        xfer += proto->serializedFieldSize("result", apache::thrift::protocol::T_STRUCT, 2);
        xfer += Cpp2Ops<nebula::client::BindingTable>::serializedSizeZC(proto,
                                                                        &obj->result_.value());
    }
    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache

#endif  // GRAPH_RESPONSE_EXECUTIONOPS_INL_H_
