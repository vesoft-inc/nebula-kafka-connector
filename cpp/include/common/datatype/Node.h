// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_NODE_H_
#define COMMON_DATATYPE_NODE_H_

#include <cstdint>
#include <unordered_map>

#include "common/datatype/Value.h"
#include "common/utils/Types.h"

namespace nebula::client {

class Node final {
public:
    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept {
        return properties_.get_allocator();
    }

    explicit Node(const allocator_type& alloc = allocator_type()) : properties_(alloc) {}
    Node(InternalID nodeID,
         NodeTypeID nodeTypeID,
         const allocator_type& alloc = allocator_type())
            : nodeID_(nodeID), nodeTypeID_(nodeTypeID), properties_(alloc) {}
    Node(InternalID nodeID,
         NodeTypeID nodeTypeID,
         const std::pmr::unordered_map<std::pmr::string, Value>& properties,
         const allocator_type& alloc = allocator_type())
            : nodeID_(nodeID), nodeTypeID_(nodeTypeID), properties_(properties, alloc) {}

    Node(const Node& other, const allocator_type& alloc)
            : nodeID_(other.getNodeID()),
              nodeTypeID_(other.getNodeTypeID()),
              properties_(other.getProperties(), alloc) {}

    Node(Node&& other, const allocator_type& alloc) noexcept
            : nodeID_(std::move(other.nodeID_)),
              nodeTypeID_(std::move(other.nodeTypeID_)),
              properties_(std::move(other.properties_), alloc) {}

    InternalID getNodeID() const {
        return nodeID_;
    }

    void setNodeID(InternalID nodeID) {
        nodeID_ = nodeID;
    }

    NodeTypeID getNodeTypeID() const {
        return nodeTypeID_;
    }

    void setNodeTypeID(NodeTypeID nodeTypeID) {
        nodeTypeID_ = nodeTypeID;
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
        auto it = properties_.find(propName);
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
    friend class apache::thrift::Cpp2Ops<Node, void>;

    // using NodeID = uInternalID;
    InternalID nodeID_;
    NodeTypeID nodeTypeID_;
    // TODO: avoid to use std::pmr::vector here
    std::pmr::unordered_map<std::pmr::string, Value> properties_;
};

inline std::ostream& operator<<(std::ostream& os, const Node& node) {
    return os << node.toString();
}

bool operator==(const Node& lhs, const Node& rhs);
bool operator!=(const Node& lhs, const Node& rhs);

}  // namespace nebula::client

namespace std {

template <>
struct hash<nebula::client::Node> {
    std::size_t operator()(const nebula::client::Node& node) const noexcept;
};

}  // namespace std

#endif  // COMMON_DATATYPE_NODE_H_
