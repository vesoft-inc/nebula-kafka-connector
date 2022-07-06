// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_EDGEOPS_H_
#define COMMON_DATATYPE_EDGEOPS_H_

#include <thrift/lib/cpp/protocol/TType.h>
#include <thrift/lib/cpp2/GeneratedCodeHelper.h>
#include <thrift/lib/cpp2/TypeClass.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>
#include <thrift/lib/cpp2/protocol/ProtocolReaderStructReadState.h>

#include "common/datatype/Edge.h"
#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"
#include "common/utils/Types.h"

namespace apache {
namespace thrift {
/**************************************
 *
 * Ops for class Edge
 *
 *************************************/
namespace detail {

template <>
struct TccStructTraits<nebula::client::Edge> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "srcID") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_I64;
        } else if (_fname == "dstID") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_I64;
        } else if (_fname == "edgeTypeID") {
            fid = 3;
            _ftype = apache::thrift::protocol::T_I32;
        } else if (_fname == "rank") {
            fid = 4;
            _ftype = apache::thrift::protocol::T_I64;
        } else if (_fname == "properties") {
            fid = 5;
            _ftype = apache::thrift::protocol::T_MAP;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::Edge>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Edge>::write(Protocol* proto,
                                              nebula::client::Edge const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("Edge");

    // Write field srcID (required)
    xfer += proto->writeFieldBegin("srcID", apache::thrift::protocol::T_I64, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, int64_t>::write(*proto,
                                                                               obj->getSrcID());
    xfer += proto->writeFieldEnd();

    // Write field dstID (required)
    xfer += proto->writeFieldBegin("srcID", apache::thrift::protocol::T_I64, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, int64_t>::write(*proto,
                                                                               obj->getDstID());
    xfer += proto->writeFieldEnd();

    // Write field edgeTypeID (required)
    xfer += proto->writeFieldBegin("edgeTypeID", apache::thrift::protocol::T_I32, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, int32_t>::write(
            *proto, obj->getEdgeTypeID());
    xfer += proto->writeFieldEnd();

    // Write field rank (required)
    xfer += proto->writeFieldBegin("rank", apache::thrift::protocol::T_I64, 4);
    xfer += detail::pm::protocol_methods<type_class::integral, int64_t>::write(
            *proto, obj->getEdgeRank());
    xfer += proto->writeFieldEnd();

    // Write field properties (required)
    xfer += proto->writeFieldBegin("properties", apache::thrift::protocol::T_MAP, 5);
    xfer += detail::pm::protocol_methods<
            type_class::map<type_class::string, type_class::structure>,
            std::pmr::unordered_map<std::pmr::string, nebula::client::Value>>::
            write(*proto, obj->getProperties());
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::Edge>::read(Protocol* proto, nebula::client::Edge* obj) {
    detail::ProtocolReaderStructReadState<Protocol> readState;

    readState.readStructBegin(proto);

    using apache::thrift::protocol::TProtocolException;

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_I64))) {
        goto _loop;
    }

_readField_srcID : {
    obj->srcID_ = 0;
    detail::pm::protocol_methods<type_class::integral, int64_t>::read(*proto, obj->srcID_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_I64))) {
        goto _loop;
    }

_readField_dstID : {
    obj->dstID_ = 0;
    detail::pm::protocol_methods<type_class::integral, int64_t>::read(*proto, obj->dstID_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 2, 3, protocol::T_I32))) {
        goto _loop;
    }

_readField_edgeTypeID : {
    obj->edgeTypeID_ = 0;
    detail::pm::protocol_methods<type_class::integral, int32_t>::read(*proto, obj->edgeTypeID_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 3, 4, protocol::T_I64))) {
        goto _loop;
    }

_readField_rank : {
    obj->rank_ = 0;
    detail::pm::protocol_methods<type_class::integral, int64_t>::read(*proto, obj->rank_);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 4, 5, protocol::T_MAP))) {
        goto _loop;
    }

_readField_properties : {
    obj->properties_.clear();
    detail::pm::protocol_methods<
            type_class::map<type_class::binary, type_class::structure>,
            std::pmr::unordered_map<std::pmr::string,
                                    nebula::client::Value>>::read(*proto, obj->properties_);
}
    if (UNLIKELY(!readState.advanceToNextField(proto, 5, 0, protocol::T_STOP))) {
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
        detail::TccStructTraits<nebula::client::Edge>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I64)) {
                goto _readField_srcID;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I64)) {
                goto _readField_dstID;
            } else {
                goto _skip;
            }
        }
        case 3: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I32)) {
                goto _readField_edgeTypeID;
            } else {
                goto _skip;
            }
        }
        case 4: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_I64)) {
                goto _readField_rank;
            } else {
                goto _skip;
            }
        }
        case 5: {
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
uint32_t Cpp2Ops<nebula::client::Edge>::serializedSize(Protocol const* proto,
                                                       nebula::client::Edge const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Edge");

    xfer += proto->serializedFieldSize("srcID", apache::thrift::protocol::T_I64, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::InternalID>::
            serializedSize<false>(*proto, obj->srcID_);
    xfer += proto->serializedFieldSize("dstID", apache::thrift::protocol::T_I64, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::InternalID>::
            serializedSize<false>(*proto, obj->dstID_);
    xfer += proto->serializedFieldSize("edgeTypeID", apache::thrift::protocol::T_I32, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::EdgeTypeID>::
            serializedSize<false>(*proto, obj->edgeTypeID_);
    xfer += proto->serializedFieldSize("rank", apache::thrift::protocol::T_I64, 4);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::EdgeRank>::
            serializedSize<false>(*proto, obj->rank_);
    xfer += proto->serializedFieldSize("properties", apache::thrift::protocol::T_MAP, 5);
    xfer += detail::pm::protocol_methods<
            type_class::map<type_class::string, type_class::structure>,
            std::pmr::unordered_map<std::pmr::string, nebula::client::Value>>::
            serializedSize<false>(*proto, obj->properties_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::Edge>::serializedSizeZC(Protocol const* proto,
                                                         nebula::client::Edge const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("Edge");

    xfer += proto->serializedFieldSize("srcID", apache::thrift::protocol::T_I64, 1);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::InternalID>::
            serializedSize<false>(*proto, obj->srcID_);
    xfer += proto->serializedFieldSize("dstID", apache::thrift::protocol::T_I64, 2);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::InternalID>::
            serializedSize<false>(*proto, obj->dstID_);
    xfer += proto->serializedFieldSize("edgeTypeID", apache::thrift::protocol::T_I32, 3);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::EdgeTypeID>::
            serializedSize<false>(*proto, obj->edgeTypeID_);
    xfer += proto->serializedFieldSize("rank", apache::thrift::protocol::T_I64, 4);
    xfer += detail::pm::protocol_methods<type_class::integral, nebula::client::EdgeRank>::
            serializedSize<false>(*proto, obj->rank_);
    xfer += proto->serializedFieldSize("properties", apache::thrift::protocol::T_MAP, 5);
    xfer += detail::pm::protocol_methods<
            type_class::map<type_class::string, type_class::structure>,
            std::pmr::unordered_map<std::pmr::string, nebula::client::Value>>::
            serializedSize<false>(*proto, obj->properties_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
#endif  // COMMON_DATATYPE_EDGEOPS_H_
