// Copyright (c) 2022 vesoft inc. All rights reserved.

#include <common/Init.h>
#include <nebula/client/Config.h>
#include <nebula/client/ConnectionPool.h>

#include <atomic>
#include <chrono>
#include <thread>

int main(int argc, char* argv[]) {
    nebula::client::init(&argc, &argv);
    auto address = "127.0.0.1:9669";
    if (argc == 2) {
        address = argv[1];
    }
    std::cout << "Current address: " << address << std::endl;
    nebula::client::ConnectionPool pool;
    pool.init({address}, nebula::client::Config{});
    auto session = pool.getSession("root", "nebula");
    if (!session.valid()) {
        return -1;
    }

    auto result = session.execute("SHOW HOSTS");
    if (result.errorCode != nebula::client::ErrorCode::SUCCEEDED) {
        std::cout << "Exit with error code: " << static_cast<int>(result.errorCode)
                  << std::endl;
        return static_cast<int>(result.errorCode);
    }
    std::cout << *result.data;

    std::atomic_bool complete{false};
    session.asyncExecute("SHOW HOSTS",
                         [&complete](nebula::client::ExecutionResponse&& cbResult) {
                             if (cbResult.errorCode != nebula::client::ErrorCode::SUCCEEDED) {
                                 std::cout << "Exit with error code: "
                                           << static_cast<int>(cbResult.errorCode) << std::endl;
                                 std::exit(static_cast<int>(cbResult.errorCode));
                             }
                             std::cout << *cbResult.data;
                             complete.store(true);
                         });

    while (!complete.load()) {
        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    session.release();
    return 0;
}
