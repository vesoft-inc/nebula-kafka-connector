// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_GRAPH_AUTHRESPONSE_H
#define COMMON_GRAPH_AUTHRESPONSE_H

#include <optional>

#include "graph/response/GqlStatus.h"

namespace nebula::client {

// The response of the authentication of a principal.
struct AuthResponse {
    void __clear() {
        gqlStatus_.clear();
        identifier_ = {};
    }

    void clear() {
        __clear();
    }

    bool operator==(const AuthResponse& rhs) const;
    bool operator!=(const AuthResponse& rhs) const;

    GQLStatus gqlStatus_;
    std::optional<int64_t> identifier_{};
    // May be not required in v5.0
    // std::optional<int32_t> timeZoneOffsetSeconds{};
    // std::optional<std::string> timeZoneName{};
};

}  // namespace nebula::client

#endif  // COMMON_GRAPH_AUTHRESPONSE_H
