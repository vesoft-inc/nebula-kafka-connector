// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "common/datatype/Node.h"

#include <folly/hash/Hash.h>

#include <sstream>

namespace nebula::client {

std::string Node::toString() const {
    std::stringstream os;
    os << "({";
    auto n = properties_.size();
    for (const auto& item : properties_) {
        os << item.first << ":" << item.second;
        if (--n > 0) {
            os << ", ";
        }
    }
    os << "})";
    return os.str();
}

bool operator==(const Node& lhs, const Node& rhs) {
    return lhs.getNodeID() == rhs.getNodeID() && lhs.getNodeTypeID() == rhs.getNodeTypeID() &&
           lhs.getProperties() == rhs.getProperties();
}

bool compareWithoutDynamicId(const Node& lhs, const Node& rhs) {
    return lhs.getProperties() == rhs.getProperties();
}

bool operator!=(const Node& lhs, const Node& rhs) {
    return !(lhs == rhs);
}

}  // namespace nebula::client

namespace std {

std::size_t hash<nebula::client::Node>::operator()(
        const nebula::client::Node& node) const noexcept {
    return node.getNodeID();
}

}  // namespace std
