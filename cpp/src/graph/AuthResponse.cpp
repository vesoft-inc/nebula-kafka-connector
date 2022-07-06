// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "graph/response/AuthResponse.h"

namespace nebula::client {
bool AuthResponse::operator==(const AuthResponse& rhs) const {
    return gqlStatus_ == rhs.gqlStatus_ && identifier_ == rhs.identifier_;
}

bool AuthResponse::operator!=(const AuthResponse& rhs) const {
    return !(*this == rhs);
}

}  // namespace nebula::client
