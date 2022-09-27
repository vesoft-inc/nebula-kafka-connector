// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_EDGE_H_
#define COMMON_DATATYPE_EDGE_H_

#include <cstdint>
#include <unordered_map>

#include "common/datatype/Value.h"
#include "common/utils/Types.h"

namespace nebula::client {

class Edge final {
public:
    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept {
        return properties_.get_allocator();
    }

    explicit Edge(const allocator_type& alloc = allocator_type()) : properties_(alloc) {}
    Edge(InternalID srcID,
         InternalID dstID,
         EdgeTypeID edgeTypeID,
         EdgeRank rank,
         const allocator_type& alloc = allocator_type())
            : srcID_(srcID),
              dstID_(dstID),
              edgeTypeID_(edgeTypeID),
              rank_(rank),
              properties_(alloc) {}

    Edge(InternalID srcID,
         InternalID dstID,
         EdgeTypeID edgeTypeID,
         EdgeRank rank,
         const std::pmr::unordered_map<std::pmr::string, Value>& properties,
         const allocator_type& alloc = allocator_type())
            : srcID_(srcID),
              dstID_(dstID),
              edgeTypeID_(edgeTypeID),
              rank_(rank),
              properties_(properties, alloc) {}

    Edge(const Edge& other, const allocator_type& alloc)
            : srcID_(other.getSrcID()),
              dstID_(other.getDstID()),
              edgeTypeID_(other.getEdgeTypeID()),
              rank_(other.getEdgeRank()),
              properties_(other.getProperties(), alloc) {}

    Edge(Edge&& other, const allocator_type& alloc) noexcept
            : srcID_(std::move(other.srcID_)),
              dstID_(std::move(other.dstID_)),
              edgeTypeID_(std::move(other.edgeTypeID_)),
              rank_(std::move(other.rank_)),
              properties_(std::move(other.properties_), alloc) {}

    InternalID getSrcID() const {
        return srcID_;
    }

    void setSrcID(InternalID srcID) {
        srcID_ = srcID;
    }

    InternalID getDstID() const {
        return dstID_;
    }

    void setDstID(InternalID dstID) {
        dstID_ = dstID;
    }

    EdgeTypeID getEdgeTypeID() const {
        return edgeTypeID_;
    }

    void setEdgeTypeID(EdgeTypeID edgeTypeID) {
        edgeTypeID_ = edgeTypeID;
    }

    EdgeRank getEdgeRank() const {
        return rank_;
    }

    void setEdgeRank(EdgeRank rank) {
        rank_ = rank;
    }

    const std::pmr::unordered_map<std::pmr::string, Value>& getProperties() const {
        return properties_;
    }

    const Value& getPropertyValue(const char* propName) const {
        return getPropertyValue(std::pmr::string(propName));
    }

    const Value& getPropertyValue(const std::string& propName) const {
        return getPropertyValue(std::pmr::string(propName));
    }

    const Value& getPropertyValue(const std::pmr::string& propName) const {
        auto it = properties_.find(propName.c_str());
        if (it == properties_.end()) {
            return Value::kNullValue;
        }
        return it->second;
    }

    bool setProperty(std::string_view propName, const Value& value) {
        auto ret = properties_.emplace(propName, value);
        return ret.second;
    }

    std::string toString() const;

private:
    // Serialization using fbthrift
    friend class apache::thrift::Cpp2Ops<Edge, void>;

    InternalID srcID_;
    InternalID dstID_;
    EdgeTypeID edgeTypeID_;
    EdgeRank rank_;
    std::pmr::unordered_map<std::pmr::string, Value> properties_;
};

inline std::ostream& operator<<(std::ostream& os, const Edge& edge) {
    return os << edge.toString();
}

bool operator==(const Edge& lhs, const Edge& rhs);
bool compareWithoutDynamicId(const Edge& lhs, const Edge& rhs);
bool operator!=(const Edge& lhs, const Edge& rhs);

}  // namespace nebula::client

namespace std {

template <>
struct hash<nebula::client::Edge> {
    std::size_t operator()(const nebula::client::Edge& edge) const noexcept;
};

}  // namespace std

#endif  // COMMON_DATATYPE_EDGE_H_
