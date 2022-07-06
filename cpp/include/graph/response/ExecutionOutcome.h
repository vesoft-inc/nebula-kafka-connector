// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef GRAPH_RESPONSE_EXECUTIONOUTCOME_H_
#define GRAPH_RESPONSE_EXECUTIONOUTCOME_H_

#include <optional>

#include "common/datatype/BindingTable.h"
// #include "common/datatype/Value.h"
#include "graph/response/GqlStatus.h"

namespace nebula::client {

// An execution outcome is a component of an execution context representing the outcome of an
// execution and comprises:
//  1) A GQL-status object.
//  2) An optional result that is always a value. Every result is required to be a supported
//  result, which is a valid result of a successful outcome. Supported results are direct values
//  or indirect values that possibly reference GQL-objects.
struct ExecutionOutcome {
    ExecutionOutcome() = default;

    void clear() {
        gqlStatus_.clear();
        result_.reset();
    }

    void __clear() {
        clear();
    }

    bool operator==(const ExecutionOutcome& rhs) const {
        return gqlStatus_ == rhs.gqlStatus_ && result_ == rhs.result_;
    }

    bool operator!=(const ExecutionOutcome& rhs) const {
        return !(*this == rhs);
    }

    // TODO(Aiee) add default value
    GQLStatus gqlStatus_;
    std::optional<BindingTable> result_{};  // optional
};


}  // namespace nebula::client

#endif  // GRAPH_RESPONSE_EXECUTIONOUTCOME_H_
