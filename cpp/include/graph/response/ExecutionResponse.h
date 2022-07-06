// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_GRAPH_EXECUTIONRESPONSE_H
#define COMMON_GRAPH_EXECUTIONRESPONSE_H

#include "graph/response/ExecutionOutcome.h"

namespace nebula::client {

// TODO(Aiee) Add profiling information
struct ProfilingStats {};
struct PlanNodeDescription {};
struct PlanDescription {};


// ExecutionResponse is the response returned to the user after the execution of a query
// finished.
// Besides the ExecutionOutcome, it contains extra information such as profiling information,
// latency, etc.
struct ExecutionResponse {
    void clear() {
        executionOutcome_.clear();
        latencyInUs_ = 0;
        // planDesc_.reset();
    }

    void __clear() {
        clear();
    }

    bool operator==(const ExecutionResponse& rhs) const {
        return executionOutcome_ == rhs.executionOutcome_ && latencyInUs_ == rhs.latencyInUs_;
    }

    bool operator!=(const ExecutionResponse& rhs) const {
        return !(*this == rhs);
    }

    ExecutionOutcome executionOutcome_;
    int64_t latencyInUs_{0};
    // std::optional<PlanDescription> planDesc_{};

    // TODO(Aiee) Returns the response as a JSON string
};

}  // namespace nebula::client

#endif  // COMMON_GRAPH_EXECUTIONRESPONSE_H
