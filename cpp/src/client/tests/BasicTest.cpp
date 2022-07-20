// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "common/Init.h"
#include "gtest/gtest.h"
#include "nebula/client/Config.h"
#include "nebula/client/ConnectionPool.h"
#include "nebula/client/Session.h"

TEST(Session, Basic) {
    nebula::client::ConnectionPool pool;
    pool.init({"127.0.0.1:22347"}, nebula::client::Config{});
    auto session =
            pool.getSession(nebula::client::AuthReq{"root", "nebula", "graph", "v5.0.0"});
    ASSERT_TRUE(session.valid());
    auto resp = session.execute("malformed query.");
    std::cout << "result: " << resp.executionOutcome_.gqlStatus_.status_ << std::endl;
}

int main(int argc, char *argv[]) {
    testing::InitGoogleTest(&argc, argv);
    nebula::client::init(&argc, &argv);
    return RUN_ALL_TESTS();
}
