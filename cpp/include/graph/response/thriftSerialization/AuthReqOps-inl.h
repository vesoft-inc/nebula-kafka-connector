// Copyright (c) 2022 vesoft inc. All rights reserved.

#pragma once

#include <thrift/lib/cpp2/GeneratedCodeHelper.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>
#include <thrift/lib/cpp2/protocol/ProtocolReaderStructReadState.h>

#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"
#include "graph/response/AuthReq.h"

namespace apache {
namespace thrift {

namespace detail {

template <>
struct TccStructTraits<nebula::client::AuthReq> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "username") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_STRING;
        } else if (_fname == "password") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_STRING;
        } else if (_fname == "clientType") {
            fid = 3;
            _ftype = apache::thrift::protocol::T_STRING;
        } else if (_fname == "clientVersion") {
            fid = 4;
            _ftype = apache::thrift::protocol::T_STRING;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::AuthReq>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::AuthReq>::write(Protocol* proto,
                                                 nebula::client::AuthReq const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("AuthReq");

    // Write field status (required)

    xfer += proto->writeFieldBegin("username", protocol::T_STRING, 1);
    xfer += proto->writeBinary(obj->username);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("password", protocol::T_STRING, 2);
    xfer += proto->writeBinary(obj->password);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("clientType", protocol::T_STRING, 3);
    xfer += proto->writeBinary(obj->clientType);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("clientVersion", protocol::T_STRING, 4);
    xfer += proto->writeBinary(obj->clientVersion);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::AuthReq>::read(Protocol* proto, nebula::client::AuthReq* obj) {
    using apache::thrift::protocol::TProtocolException;
    detail::ProtocolReaderStructReadState<Protocol> readState;
    readState.readStructBegin(proto);

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_STRING))) {
        goto _loop;
    }

_readField_username : {
    obj->username.clear();
    proto->readBinary(obj->username);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_STRING))) {
        goto _loop;
    }

_readField_password : {
    obj->password.clear();
    proto->readBinary(obj->password);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 2, 3, protocol::T_STRING))) {
        goto _loop;
    }

_readField_clientType : {
    obj->clientType.clear();
    proto->readBinary(obj->clientType);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 3, 4, protocol::T_STRING))) {
        goto _loop;
    }

_readField_clientVersion : {
    obj->clientVersion.clear();
    proto->readBinary(obj->clientVersion);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 4, 0, protocol::T_STOP))) {
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
        detail::TccStructTraits<nebula::client::AuthReq>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_STRING)) {
                goto _readField_username;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_STRING)) {
                goto _readField_password;
            } else {
                goto _skip;
            }
        }
        case 3: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_STRING)) {
                goto _readField_clientType;
            } else {
                goto _skip;
            }
        }
        case 4: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_STRING)) {
                goto _readField_clientVersion;
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
uint32_t Cpp2Ops<nebula::client::AuthReq>::serializedSize(Protocol const* proto,
                                                          nebula::client::AuthReq const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("AuthReq");

    xfer += proto->serializedfieldsize("username", apache::thrift::protocol::T_STRING, 1);
    xfer += proto->serializedsizebinary(obj->username);

    xfer += proto->serializedfieldsize("password", apache::thrift::protocol::T_STRING, 2);
    xfer += proto->serializedsizebinary(obj->password);

    xfer += proto->serializedfieldsize("clientType", apache::thrift::protocol::T_STRING, 3);
    xfer += proto->serializedsizebinary(obj->clientType);

    xfer += proto->serializedfieldsize("clientVersion", apache::thrift::protocol::T_STRING, 4);
    xfer += proto->serializedsizebinary(obj->clientVersion);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::AuthReq>::serializedSizeZC(
        Protocol const* proto, nebula::client::AuthReq const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("AuthReq");

    xfer += proto->serializedFieldSize("username", apache::thrift::protocol::T_STRING, 1);
    xfer += proto->serializedSizeZCBinary(obj->username);

    xfer += proto->serializedFieldSize("password", apache::thrift::protocol::T_STRING, 2);
    xfer += proto->serializedSizeZCBinary(obj->password);

    xfer += proto->serializedFieldSize("clientType", apache::thrift::protocol::T_STRING, 3);
    xfer += proto->serializedSizeZCBinary(obj->clientType);

    xfer += proto->serializedFieldSize("clientVersion", apache::thrift::protocol::T_STRING, 4);
    xfer += proto->serializedSizeZCBinary(obj->clientVersion);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
