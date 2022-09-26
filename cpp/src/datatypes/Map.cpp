// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "common/datatype/Map.h"

#include <folly/String.h>

#include <sstream>

#include "common/datatype/Value.h"

namespace nebula::client {

std::string Map::toString() const {
    std::vector<std::string> strs(values_.size());
    std::transform(
            values_.begin(), values_.end(), strs.begin(), [](const auto& iter) -> std::string {
                std::stringstream out;
                out << iter.first << ":" << iter.second;
                return out.str();
            });

    std::stringstream os;
    os << "{" << folly::join(",", strs) << "}";
    return os.str();
}

bool operator==(const Map& lhs, const Map& rhs) {
    return lhs.getValues() == rhs.getValues();
}

bool compareWithoutDynamicId(const Map& lhs, const Map& rhs) {
    if (lhs.size() != rhs.size()) {
        return false;
    }
    for (const auto& iter : lhs.getValues()) {
        auto rhsIter = rhs.getValues().find(iter.first);
        if (rhsIter == rhs.getValues().end()) {
            return false;
        }
        if (!compareWithoutDynamicId(iter.second, rhsIter->second)) {
            return false;
        }
    }
    return true;
}

bool operator!=(const Map& lhs, const Map& rhs) {
    return !(lhs == rhs);
}

}  // namespace nebula::client


namespace std {

std::size_t hash<nebula::client::Map>::operator()(
        const nebula::client::Map& map) const noexcept {
    size_t seed = 0;
    for (auto& v : map.getValues()) {
        seed ^= hash<nebula::client::Value>()(v.first) + 0x9e3779b9 + (seed << 6) + (seed >> 2);
        seed ^= hash<nebula::client::Value>()(v.second) + 0x9e3779b9 + (seed << 6) +
                (seed >> 2);
    }
    return seed;
}

}  // namespace std
