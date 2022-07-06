// Copyright (c) 2022 vesoft inc. All rights reserved.

#pragma once

#include <gtest/gtest.h>

#include "common/datatypes/DataSet.h"
#include "common/graph/Response.h"
class SClientTest : public ::testing::Test {
protected:
    static ::testing::AssertionResult verifyResultWithoutOrder(
            const nebula::client::DataSet& resp, const nebula::client::DataSet& expect) {
        nebula::client::DataSet respCopy = resp;
        nebula::client::DataSet expectCopy = expect;
        std::sort(respCopy.rows.begin(), respCopy.rows.end());
        std::sort(expectCopy.rows.begin(), expectCopy.rows.end());
        if (respCopy != expectCopy) {
            return ::testing::AssertionFailure() << "Resp is : " << resp << std::endl
                                                 << "Expect : " << expect;
        }
        return ::testing::AssertionSuccess();
    }

    static ::testing::AssertionResult verifyResultWithoutOrder(
            const nebula::client::ExecutionResponse& resp,
            const nebula::client::DataSet& expect) {
        auto result = succeeded(resp);
        if (!result) {
            return result;
        }
        const auto& data = *resp.data;
        return verifyResultWithoutOrder(data, expect);
    }

    static ::testing::AssertionResult succeeded(const nebula::client::ExecutionResponse& resp) {
        if (resp.errorCode != nebula::client::ErrorCode::SUCCEEDED) {
            return ::testing::AssertionFailure()
                   << "Execution Failed with error: " << resp.errorCode << ","
                   << *resp.errorMsg;
        }
        return ::testing::AssertionSuccess();
    }
};
