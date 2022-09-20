// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_BINDINGTABLEOPS_H_
#define COMMON_DATATYPE_BINDINGTABLEOPS_H_

#include <thrift/lib/cpp2/GeneratedCodeHelper.h>
#include <thrift/lib/cpp2/gen/module_types_tcc.h>
#include <thrift/lib/cpp2/protocol/ProtocolReaderStructReadState.h>

#include "common/datatype/BindingTable.h"
#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"

namespace apache {
namespace thrift {

namespace detail {

template <>
struct TccStructTraits<nebula::client::BindingTable> {
    static void translateFieldName(MAYBE_UNUSED folly::StringPiece _fname,
                                   MAYBE_UNUSED int16_t& fid,
                                   MAYBE_UNUSED apache::thrift::protocol::TType& _ftype) {
        if (_fname == "columnNames") {
            fid = 1;
            _ftype = apache::thrift::protocol::T_LIST;
        } else if (_fname == "records") {
            fid = 2;
            _ftype = apache::thrift::protocol::T_LIST;
        }
    }
};

}  // namespace detail

inline constexpr protocol::TType Cpp2Ops<nebula::client::BindingTable>::thriftType() {
    return apache::thrift::protocol::T_STRUCT;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::BindingTable>::write(Protocol* proto,
                                                      nebula::client::BindingTable const* obj) {
    uint32_t xfer = 0;
    xfer += proto->writeStructBegin("BindingTable");

    xfer += proto->writeFieldBegin("columnNames", apache::thrift::protocol::T_LIST, 1);
    xfer += detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::vector<std::string>>::write(*proto, obj->desc_.getColumnNames());
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldBegin("records", apache::thrift::protocol::T_LIST, 2);
    xfer += detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::pmr::deque<nebula::client::RawRecord>>::write(*proto, obj->records_);
    xfer += proto->writeFieldEnd();

    xfer += proto->writeFieldStop();
    xfer += proto->writeStructEnd();
    return xfer;
}

template <class Protocol>
void Cpp2Ops<nebula::client::BindingTable>::read(Protocol* proto,
                                                 nebula::client::BindingTable* obj) {
    detail::ProtocolReaderStructReadState<Protocol> readState;

    readState.readStructBegin(proto);

    using apache::thrift::protocol::TProtocolException;

    if (UNLIKELY(!readState.advanceToNextField(proto, 0, 1, protocol::T_LIST))) {
        goto _loop;
    }

_readField_columnNames : {
    obj->desc_.clear();
    // TODO refactor table descriptor
    std::vector<std::string> colNames;
    detail::pm::protocol_methods<type_class::list<type_class::structure>,
                                 std::vector<std::string>>::read(*proto, colNames);
    obj->setColumnNames(colNames);
}

    if (UNLIKELY(!readState.advanceToNextField(proto, 1, 2, protocol::T_LIST))) {
        goto _loop;
    }

_readField_values : {
    obj->records_.clear();
    detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::pmr::deque<nebula::client::RawRecord>>::read(*proto, obj->records_);
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
        detail::TccStructTraits<nebula::client::BindingTable>::translateFieldName(
                readState.fieldName(), readState.fieldId, readState.fieldType);
    }

    switch (readState.fieldId) {
        case 1: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_LIST)) {
                goto _readField_columnNames;
            } else {
                goto _skip;
            }
        }
        case 2: {
            if (LIKELY(readState.fieldType == apache::thrift::protocol::T_LIST)) {
                goto _readField_values;
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
uint32_t Cpp2Ops<nebula::client::BindingTable>::serializedSize(
        Protocol const* proto, nebula::client::BindingTable const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("BindingTable");

    xfer += proto->serializedFieldSize("columnNames", apache::thrift::protocol::T_LIST, 1);
    xfer += detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::vector<std::string>>::serializedSize<false>(*proto,
                                                             obj->desc_.getColumnNames());

    xfer += proto->serializedFieldSize("records", apache::thrift::protocol::T_LIST, 2);
    xfer += detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::pmr::deque<nebula::client::RawRecord>>::serializedSize<false>(*proto,
                                                                               obj->records_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

template <class Protocol>
uint32_t Cpp2Ops<nebula::client::BindingTable>::serializedSizeZC(
        Protocol const* proto, nebula::client::BindingTable const* obj) {
    uint32_t xfer = 0;
    xfer += proto->serializedStructSize("BindingTable");

    xfer += proto->serializedFieldSize("columnNames", apache::thrift::protocol::T_LIST, 1);
    xfer += detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::vector<std::string>>::serializedSize<false>(*proto,
                                                             obj->desc_.getColumnNames());

    xfer += proto->serializedFieldSize("records", apache::thrift::protocol::T_LIST, 2);
    xfer += detail::pm::protocol_methods<
            type_class::list<type_class::structure>,
            std::pmr::deque<nebula::client::RawRecord>>::serializedSize<false>(*proto,
                                                                               obj->records_);

    xfer += proto->serializedSizeStop();
    return xfer;
}

}  // namespace thrift
}  // namespace apache
#endif  // COMMON_DATATYPE_BINDINGTABLEOPS_H_
