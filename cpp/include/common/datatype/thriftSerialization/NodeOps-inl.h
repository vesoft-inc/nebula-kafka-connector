// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_NODEOPS_H_
#define COMMON_DATATYPE_NODEOPS_H_

#include <thrift/lib/cpp2/GeneratedCodeHelper.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>
#include <thrift/lib/cpp2/protocol/ProtocolReaderStructReadState.h>

#include "common/datatype/Node.h"
#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"


namespace apache {
namespace thrift {
/**************************************
 *
 * Ops for class Node
 *
 *************************************/
namespace detail {

template <>
struct TccStructTraits<nebula::client::Node> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "nodeId") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_I64;
        } else if (_fname == "nodeTypeID") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_I32;
        } else if (_fname == "properties") {
            fid = 3;
            _ftype = apache::thrift::protocol::T_MAP;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::Node>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Node>::write(Protocol* proto,
                                              nebula::client::Node const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("Node");
    // Write field nodeID (required)
    xfer += proto->writeFieldBegin("nodeID", apache::thrift::protocol::T_I64, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int64_t>::write(
            *proto, obj->getNodeID());
    xfer += proto->writeFieldEnd();

    // Write field nodeTypeID (required)
    xfer += proto->writeFieldBegin("nodeTypeID", apache::thrift::protocol::T_I32, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::write(
            *proto, obj->getNodeTypeID());
    xfer += proto->writeFieldEnd();

    // Write field properties (required)
    xfer += proto->writeFieldBegin("properties", apache::thrift::protocol::T_MAP, 3);
    xfer += detail::pm::protocol_methods<
            type_class::map<type_class::string, type_class::structure>,
            std::pmr::unordered_map<std::pmr::string,
                                    nebula::client::Value>>::write(*proto, obj->properties_);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::Node>::read(Protocol* proto, nebula::client::Node* obj) {
    detail::ProtocolReaderStructReadState<Protocol> readState;

    readState.readStructBegin(proto);

    using apache::thrift::protocol::TProtocolException;

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_I64))) {
        goto _loop;
    }

_readField_nodeID : {
    obj->nodeID_ = 0;
    detail::pm::protocol_methods<type_class::integral, int64_t>::read(*proto, obj->nodeID_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_I32))) {
        goto _loop;
    }

_readField_nodeTypeID : {
    obj->nodeTypeID_ = 0;
    detail::pm::protocol_methods<type_class::integral, int32_t>::read(*proto, obj->nodeTypeID_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 2, 3, protocol::T_MAP))) {
        goto _loop;
    }

_readField_properties : {
    obj->properties_.clear();
    detail::pm::protocol_methods<
            type_class::map<type_class::binary, type_class::structure>,
            std::pmr::unordered_map<std::pmr::string,
                                    nebula::client::Value>>::read(*proto, obj->properties_);
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
        detail::TccStructTraits<nebula::client::Node>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I64)) {
                goto _readField_nodeID;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I32)) {
                goto _readField_nodeTypeID;
            } else {
                goto _skip;
            }
        }
        case 3: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_MAP)) {
                goto _readField_properties;
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
uint32_t Cpp2Ops<nebula::client::Node>::serializedSize(Protocol const* proto,
                                                       nebula::client::Node const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Node");

    xfer += proto->serializedFieldSize("nodeID", apache::thrift::protocol::T_I64, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::InternalID>::
            serializedSize<false>(*proto, obj->getNodeID());

    xfer += proto->serializedFieldSize("nodeTypeID", apache::thrift::protocol::T_I32, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::NodeTypeID>::
            serializedSize<false>(*proto, obj->getNodeTypeID());

    xfer += proto->serializedFieldSize("properties", apache::thrift::protocol::T_MAP, 3);
    xfer += detail::pm::protocol_methods<
            type_class::map<type_class::string, type_class::structure>,
            std::pmr::unordered_map<std::pmr::string, nebula::client::Value>>::
            serializedSize<false>(*proto, obj->properties_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Node>::serializedSizeZC(Protocol const* proto,
                                                         nebula::client::Node const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Node");

    xfer += proto->serializedFieldSize("nodeID", apache::thrift::protocol::T_I64, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::InternalID>::
            serializedSize<false>(*proto, obj->getNodeID());

    xfer += proto->serializedFieldSize("nodeTypeID", apache::thrift::protocol::T_I32, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::NodeTypeID>::
            serializedSize<false>(*proto, obj->getNodeTypeID());

    xfer += proto->serializedFieldSize("properties", apache::thrift::protocol::T_MAP, 3);
    xfer += detail::pm::protocol_methods<
            type_class::map<type_class::string, type_class::structure>,
            std::pmr::unordered_map<std::pmr::string, nebula::client::Value>>::
            serializedSize<false>(*proto, obj->properties_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
#endif  // COMMON_DATATYPE_NODEOPS_H_
