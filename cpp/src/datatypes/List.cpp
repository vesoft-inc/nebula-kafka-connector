// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "common/datatype/List.h"

#include <folly/String.h>

#include <algorithm>
#include <sstream>

#include "common/datatype/Value.h"

namespace nebula::client {

std::string List::toString() const {
    std::vector<std::string> strs(values_.size());
    std::transform(
            values_.begin(), values_.end(), strs.begin(), [](const auto& v) -> std::string {
                return v.toString();
            });
    std::stringstream os;
    os << "[" << folly::join(",", strs) << "]";
    return os.str();
}

bool operator==(const List& lhs, const List& rhs) {
    return lhs.getValues() == rhs.getValues();
}

bool operator!=(const List& lhs, const List& rhs) {
    return !(lhs == rhs);
}

}  // namespace nebula::client


namespace std {

std::size_t hash<nebula::client::List>::operator()(
        const nebula::client::List& list) const noexcept {
    size_t seed = 0;
    // FIXME(jie)
    for (auto& v : list.getValues()) {
        seed ^= hash<nebula::client::Value>()(v) + 0x9e3779b9 + (seed << 6) + (seed >> 2);
    }
    return seed;
}

}  // namespace std
