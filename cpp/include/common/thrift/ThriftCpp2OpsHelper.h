// Copyright (c) 2022 vesoft inc. All rights reserved.
#ifndef COMMON_THRIFT_THRIFTCPP2OPS_HELPER_H_
#define COMMON_THRIFT_THRIFTCPP2OPS_HELPER_H_

#include <thrift/lib/cpp2/GeneratedCodeHelper.h>

/**
 * Comment of serializedSize() from
 * https://github.com/facebook/fbthrift/blob/main/thrift/lib/cpp2/protocol/BinaryProtocol.h
 *
 * Functions that return the [estimated] serialized size
 * Notes:
 * * Serialized size estimates for Binary protocol are generally accurate,
 *   but this is not the case for other protocols, e.g. Compact.
 *   Don't use these values as more than an estimate.
 *
 * * ZC versions are the preallocated estimate if any IOBufs are shared (i.e.
 *   there are IOBuf fields, and their sizes aren't too small to be packed),
 *   and won't count in the ZC estimate.
 */

#define SPECIALIZE_CPP2OPS(X)                                                     \
    template <>                                                                   \
    class Cpp2Ops<X> {                                                            \
    public:                                                                       \
        using Type = X;                                                           \
        inline static constexpr protocol::TType thriftType();                     \
        template <class Protocol>                                                 \
        static uint32_t write(Protocol *proto, const Type *obj);                  \
        template <class Protocol>                                                 \
        static void read(Protocol *proto, Type *obj);                             \
        template <class Protocol>                                                 \
        static uint32_t serializedSize(const Protocol *proto, const Type *obj);   \
        template <class Protocol>                                                 \
        static uint32_t serializedSizeZC(const Protocol *proto, const Type *obj); \
    }

#endif  // COMMON_THRIFT_THRIFTCPP2OPS_HELPER_H_
