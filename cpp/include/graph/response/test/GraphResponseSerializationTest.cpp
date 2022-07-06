// Copyright (c) 2022 vesoft inc. All rights reserved.

#include <gtest/gtest.h>

#include "common/base/Base.h"
#include "common/datatype/thriftSerialization/CommonCpp2Ops.h"
#include "common/nrpc/BufferReaderWriter.h"
#include "graph/response/AuthResponse.h"
#include "graph/response/ExecutionOutcome.h"
#include "graph/response/ExecutionResponse.h"
#include "graph/response/GqlStatus.h"
#include "interface/GraphCpp2Ops.h"

namespace nebula::client {

using serializer = apache::thrift::CompactSerializer;

template <class ObjType>
testing::AssertionResult checkSerialization(const ObjType& obj) {
    std::string buf;
    buf.reserve(128);
    nebula::client::serializer::serialize(obj, &buf);
    ObjType deserializedObj;
    std::size_t s = nebula::client::serializer::deserialize(buf, deserializedObj);
    if (s != buf.size()) {
        return testing::AssertionFailure() << "deserialized size " << s << " != " << buf.size();
    }
    if (obj != deserializedObj) {
        return testing::AssertionFailure()
               << "The deserialized object is not equal to the original one";
    }
    return testing::AssertionSuccess();
}

BindingTable newBindingTable() {
    auto mr = std::pmr::new_delete_resource();
    auto table = BindingTable(mr);
    auto& record1 = table.emplace();
    record1.append(Value(true));
    record1.append(Value(88));
    record1.append(Value(3.14));
    record1.append(Value("hello", mr));

    return table;
}

TEST(GqlStatus, thriftSerialization) {
    auto status = GQLStatus(std::string("SUCCEEDED"));
    EXPECT_TRUE(checkSerialization(status));
}

TEST(ExecutionOutcome, thriftSerialization) {
    auto outcome = ExecutionOutcome();
    outcome.gqlStatus_ = GQLStatus(std::string("SUCCEEDED"));
    EXPECT_TRUE(checkSerialization(outcome));

    outcome.result_ = newBindingTable();
    EXPECT_TRUE(checkSerialization(outcome));
}

TEST(AuthResponse, thriftSerialization) {
    auto authResponse = AuthResponse();
    authResponse.gqlStatus_ = GQLStatus(std::string("SUCCEEDED"));
    EXPECT_TRUE(checkSerialization(authResponse));

    authResponse.identifier_ = 100;
    EXPECT_TRUE(checkSerialization(authResponse));
}

TEST(ExecutionResponse, thriftSerialization) {
    auto outcome = ExecutionOutcome();
    outcome.gqlStatus_ = GQLStatus(std::string("SUCCEEDED"));
    outcome.result_ = newBindingTable();

    auto executionResponse = ExecutionResponse();
    executionResponse.executionOutcome_ = outcome;
    executionResponse.latencyInUs_ = 100;
    EXPECT_TRUE(checkSerialization(executionResponse));
}


}  // namespace nebula::client

int main(int argc, char** argv) {
    testing::InitGoogleTest(&argc, argv);
    folly::init(&argc, &argv, true);
    google::SetStderrLogging(google::INFO);

    return RUN_ALL_TESTS();
}
