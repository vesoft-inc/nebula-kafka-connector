
#include <common/Init.h>
#include <glog/logging.h>
#include <gtest/gtest.h>
#include <gtest/gtest_prod.h>
#include <nebula/client/ConnectionPool.h>
#include <nebula/client/Session.h>

#include "./ClientTest.h"

// Require a nebula server could access

#define kServerHost "graphd"

class AddressTest : public ClientTest {};

TEST_F(AddressTest, One) {
    nebula::client::ConnectionPool pool;
    pool.init({kServerHost ":9669"}, nebula::client::Config{});
    EXPECT_EQ(pool.size(), 10);
}

TEST_F(AddressTest, Multiple) {
    nebula::client::ConnectionPool pool;
    pool.init({"graphd:9669", "graphd1:9669", "graphd2:9669"}, nebula::client::Config{});
    EXPECT_EQ(pool.size(), 10);
}

int main(int argc, char** argv) {
    testing::InitGoogleTest(&argc, argv);
    nebula::client::init(&argc, &argv);
    google::SetStderrLogging(google::INFO);

    return RUN_ALL_TESTS();
}
