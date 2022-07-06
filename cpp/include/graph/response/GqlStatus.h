// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef GRAPH_RESPONSE_GQLSTATUS_H_
#define GRAPH_RESPONSE_GQLSTATUS_H_

#include <string>

namespace nebula::client {

struct Diagnostics {
    // COMMAND_FUNCTION : CF;
    // COMMAND_FUNCTION_CODE : CFC;
    // NUMBER : N;
    // CURRENT_SCHEMA : CS;
    // HOME_GRAPH : HG;
    // CURRENT_GRAPH : CG;
};

// GQLStatus comprises a condition code and additional diagnostic information.
//
// Every GQL-program returns some diagnostic information to the GQL-client that originated the
// GQL-request of which the GQL-program was a part.
struct GQLStatus {
    // TODO(Aiee): This is a temp constructor used to quickly mock a GQLStatus.
    GQLStatus() = default;
    explicit GQLStatus(std::string status) : status_(status) {}

    void clear() {
        status_.clear();
        // description_.clear();
    }

    void __clear() {
        return clear();
    }

    bool operator==(const GQLStatus& rhs) const {
        return status_ == rhs.status_;
    }
    bool operator!=(const GQLStatus& rhs) const {
        return !(*this == rhs);
    }

    // TODO(Aiee) fix generated code compilation
    bool operator<(const GQLStatus&) const {
        return false;
    }
    bool operator>(const GQLStatus&) const {
        return false;
    }
    // A GQLSTATUS character string.
    std::string status_;

    // TODO(Aiee) Omitted now for simplicity.
    // // A character string describing the GQLSTATUS character string.
    // //  Table 7, “GQLSTATUS class and subclass codes”
    // std::string description_;
    // // A map value with diagnostics information as defined in Clause 23, “Diagnostics”.
    // Diagnostics diagnostics_;
    // // A chain of nested GQL-status objects.
    // // GQLStatus* statusChain_{nullptr};
    // // An optional map of implementation-defined diagnostic information.
};

}  // namespace nebula::client
#endif  // GRAPH_RESPONSE_GQLSTATUS_H_
