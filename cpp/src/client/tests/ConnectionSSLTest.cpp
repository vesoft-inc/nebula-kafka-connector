
#include <common/Init.h>
#include <folly/json.h>
#include <folly/synchronization/Baton.h>
#include <glog/logging.h>
#include <gtest/gtest.h>
#include <nebula/client/Connection.h>

#include "./ClientTest.h"

// Require a nebula server could access

static constexpr char kServerHost[] = "graphd";

class ConnectionTest : public ClientTest {};

TEST_F(ConnectionTest, SSL) {
    nebula::client::Connection c;

    ASSERT_TRUE(c.open(kServerHost, 9669, 10, true, ""));

    // auth
    auto authResp = c.authenticate("root", "nebula");
    ASSERT_EQ(authResp.errorCode, nebula::client::ErrorCode::SUCCEEDED) << *authResp.errorMsg;

    // execute
    auto resp = c.execute(*authResp.sessionId, "YIELD 1");
    ASSERT_EQ(resp.errorCode, nebula::client::ErrorCode::SUCCEEDED);
    nebula::client::DataSet expected({"1"});
    expected.emplace_back(nebula::client::List({1}));
    EXPECT_TRUE(verifyResultWithoutOrder(*resp.data, expected));
}

int main(int argc, char **argv) {
    testing::InitGoogleTest(&argc, argv);
    nebula::client::init(&argc, &argv);
    google::SetStderrLogging(google::INFO);

    return RUN_ALL_TESTS();
}
