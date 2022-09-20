// Copyright (c) 2022 vesoft inc. All rights reserved.

#include <common/Init.h>
#include <nebula/client/Config.h>
#include <nebula/client/ConnectionPool.h>

#include <atomic>
#include <chrono>
#include <iostream>
#include <thread>

int main(int argc, char* argv[]) {
    nebula::client::init(&argc, &argv);
    nebula::client::ConnectionPool pool;
    pool.init({"127.0.0.1:22347"}, nebula::client::Config{});
    auto session =
            pool.getSession(nebula::client::AuthReq{"root", "nebula", "graph", "v5.0.0"});
    assert(session.valid());
    auto resp = session.execute("malformed query.");
    std::cout << "result: " << resp.executionOutcome_.gqlStatus_.status_ << std::endl;
    return 0;
}
