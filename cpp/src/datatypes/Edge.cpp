// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "common/datatype/Edge.h"

#include <folly/hash/Hash.h>

#include <sstream>

namespace nebula::client {

std::string Edge::toString() const {
    std::stringstream os;
    os << "[{";
    auto n = properties_.size();
    for (const auto& item : properties_) {
        os << item.first << ":" << item.second;
        if (--n>0) {
            os << ", ";
        }
    }
    os << "}]";
    return os.str();
}


bool Edge::compareWithoutId(const Node& rhs) const {
    return properties_ == rhs.properties_;
}

bool operator==(const Edge& lhs, const Edge& rhs) {
    return lhs.getSrcID() == rhs.getSrcID() && lhs.getDstID() == rhs.getDstID() &&
           lhs.getEdgeRank() == rhs.getEdgeRank() &&
           lhs.getEdgeTypeID() == rhs.getEdgeTypeID() &&
           lhs.getProperties() == rhs.getProperties();
}

bool operator!=(const Edge& lhs, const Edge& rhs) {
    return !(lhs == rhs);
}

}  // namespace nebula::client

namespace std {

std::size_t hash<nebula::client::Edge>::operator()(
        const nebula::client::Edge& edge) const noexcept {
    size_t seed = 0;
    seed ^= edge.getSrcID() + 0x9e3779b9 + (seed << 6) + (seed >> 2);
    seed ^= edge.getDstID() + 0x9e3779b9 + (seed << 6) + (seed >> 2);
    seed ^= edge.getEdgeTypeID() + 0x9e3779b9 + (seed << 6) + (seed >> 2);
    return seed;
}

}  // namespace std
