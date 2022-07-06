// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "graph/response/ExecutionResponse.h"

namespace nebula::client {
bool operator==(const ExecutionResponse &lhs, ExecutionResponse &rhs) {
    return lhs.executionOutcome_ == rhs.executionOutcome_ &&
           lhs.latencyInUs_ == rhs.latencyInUs_;
    // TODO(Aiee) Add later
    //    lhs.planDesc_ == rhs.planDesc_;
}

bool operator!=(const ExecutionResponse &lhs, ExecutionResponse &rhs) {
    return !(lhs == rhs);
}
}  // namespace nebula::client
