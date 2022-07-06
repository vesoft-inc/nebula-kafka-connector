// Copyright (c) 2022 vesoft inc. All rights reserved.

namespace cpp nebula.graph
namespace java com.vesoft.nebula.graph
namespace go nebula.graph
namespace js nebula.graph
namespace csharp nebula.graph
namespace py nebula3.graph

include "common.thrift"

cpp_include "graph/response/thriftSerialization/AuthResponseOps-inl.h"
cpp_include "graph/response/thriftSerialization/AuthReqOps-inl.h"
cpp_include "graph/response/thriftSerialization/GqlStatusOps-inl.h"
cpp_include "graph/response/thriftSerialization/ExecutionOutcomeOps-inl.h"
cpp_include "graph/response/thriftSerialization/ExecutionResponseOps-inl.h"

//   Note: In order to support multiple languages, all string
//         have to be defined as **binary** in the thrift file

struct ProfilingStats {
    // How many rows being processed in an executor.
    1: required i64  rows;
    // Duration spent in an executor.
    2: required i64  exec_duration_in_us;
    // Total duration spent in an executor, contains schedule time
    3: required i64  total_duration_in_us;
    // Other profiling stats data map
    4: optional map<binary, binary>
        (cpp.template = "std::unordered_map") other_stats;
}

struct Pair {
    1: required binary key;
    2: required binary value;
}

struct PlanNodeDescription {
    1: required binary                          name;
    2: required i64                             id;
    3: required binary                          output_var;
    // other description of an executor
    4: optional list<Pair>                      description;
    // If an executor would be executed multi times,
    // the profiling statistics should be multi-versioned.
    5: optional list<ProfilingStats>            profiles;
    6: optional list<i64>                       dependencies;
}

struct PlanDescription {
    1: required list<PlanNodeDescription>     plan_node_descs;
    // map from node id to index of list
    2: required map<i64, i64>
        (cpp.template = "std::unordered_map") node_index_map;
    // the print format of exec plan, lowercase string like `dot'
    3: required binary                        format;
    // the time optimizer spent
    4: required i32                           optimize_time_in_us;
}

// TODO: Repleace with complete implementation

// Response
struct GQLStatus {
    // A GQLSTATUS character string.
    // TODO(Aiee) Unimplemented. For simplicity, we use a string to represent the status of the execution.
    1: required binary status;
    // A character string describing the GQLSTATUS character string.
    //  Table 7, “GQLSTATUS class and subclass codes”
    // 2: required binary description;
    // A map value with diagnostics information as defined in Clause 23, “Diagnostics”.
    // 3: required Diagnostics diagnostics;
    // A chain of nested GQL-status objects.
    // Note sure how to represent here
} (cpp.type = "nebula::client::GQLStatus")

// The outcome of the execution of a GQL request.
struct ExecutionOutcome {
    1: required GQLStatus (cpp.type = "nebula::client::GQLStatus")               gqlStatus;
    2: optional common.BindingTable (cpp.type = "nebula::client::BindingTable")  result;  // optional
} (cpp.type = "nebula::client::ExecutionOutcome")

struct ExecutionResponse {
    1: required ExecutionOutcome executionOutcome;  // The outcome of the execution of the GQL request
    2: required i64              latencyInUs;  // Execution time on server
    # 3: optional PlanDescription         plan_desc;      // Description of the execution plan
} (cpp.type = "nebula::client::ExecutionResponse")

//(TODO) This is a draft for now, we should try to determine the version compatibility
// when open a connection.
struct AuthReq {
    1: required binary username;
    2: required binary password;
    3: required binary client_type;
    4: required binary client_version;
} (cpp.type = "nebula::client::AuthReq")

struct AuthResponse {
    1: required GQLStatus   gqlStatus;
    2: optional i64         identifier;
} (cpp.type = "nebula::client::AuthResponse")

service GraphService {
    // Authenticates a principle
    AuthResponse authenticate(1: AuthReq authReq)
    
    // Logs out the principle by its identifier
    oneway void signout(1: i64 sessionId)

    // Exeuctes the query
    ExecutionResponse execute(1: i64 sessionId, 2: binary stmt)

    // // Same as execute(), but the response will be a json string
    // binary executeJson(1: i64 sessionId, 2: binary stmt)
    // binary executeJsonWithParameter(1: i64 sessionId, 2: binary stmt)   
}
