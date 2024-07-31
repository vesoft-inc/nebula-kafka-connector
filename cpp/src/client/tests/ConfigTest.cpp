
#include <common/Init.h>
#include <glog/logging.h>
#include <gtest/gtest.h>
#include <gtest/gtest_prod.h>
#include <nebula/client/ConnectionPool.h>
#include <nebula/client/Session.h>

#include "./ClientTest.h"

// Require a nebula server could access

#define kServerHost "graphd"

class ConfigTest : public ClientTest {};

TEST_F(ConfigTest, Timeout) {
    {
        // don't set
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        ASSERT_EQ(c.timeout_, 0);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_EQ(pool.size(), 10);
    }
    {
        // set to 0
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        c.timeout_ = 0;
        ASSERT_EQ(c.timeout_, 0);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_EQ(pool.size(), 10);
    }
    {
        // set to positive integer
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        c.timeout_ = 3;
        ASSERT_EQ(c.timeout_, 3);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_EQ(pool.size(), 10);
    }
}

TEST_F(ConfigTest, IdleTime) {
    {
        // don't set
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        ASSERT_EQ(c.idleTime_, 0);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_EQ(pool.size(), 10);
    }
    {
        // set to 0
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        c.idleTime_ = 0;
        ASSERT_EQ(c.idleTime_, 0);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_EQ(pool.size(), 10);
    }
    {
        // set to positive number
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        c.idleTime_ = 3;
        ASSERT_EQ(c.idleTime_, 3);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_EQ(pool.size(), 10);
    }
}

TEST_F(ClientTest, maxConnectionPoolSize) {
    {
        // don't set
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        ASSERT_EQ(c.maxConnectionPoolSize_, 10);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_LE(pool.size(), c.maxConnectionPoolSize_);
        EXPECT_GE(pool.size(), c.minConnectionPoolSize_);
    }
    {
        // set to positive number
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        c.maxConnectionPoolSize_ = 5;
        ASSERT_EQ(c.maxConnectionPoolSize_, 5);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_LE(pool.size(), c.maxConnectionPoolSize_);
        EXPECT_GE(pool.size(), c.minConnectionPoolSize_);
    }
    {
        // set to zero
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        c.maxConnectionPoolSize_ = 0;
        ASSERT_EQ(c.maxConnectionPoolSize_, 0);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_LE(pool.size(), c.maxConnectionPoolSize_);
        EXPECT_GE(pool.size(), c.minConnectionPoolSize_);
    }
}

TEST_F(ClientTest, minConnectionPoolSize) {
    {
        // don't set
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        ASSERT_EQ(c.minConnectionPoolSize_, 0);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_LE(pool.size(), c.maxConnectionPoolSize_);
        EXPECT_GE(pool.size(), c.minConnectionPoolSize_);
    }
    {
        // set to zero
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        c.minConnectionPoolSize_ = 0;
        ASSERT_EQ(c.minConnectionPoolSize_, 0);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_LE(pool.size(), c.maxConnectionPoolSize_);
        EXPECT_GE(pool.size(), c.minConnectionPoolSize_);
    }
    {
        // set to positive number
        nebula::client::ConnectionPool pool;
        nebula::client::Config c{};
        c.minConnectionPoolSize_ = 4;
        ASSERT_EQ(c.minConnectionPoolSize_, 4);
        pool.init({kServerHost ":9669"}, c);
        EXPECT_LE(pool.size(), c.maxConnectionPoolSize_);
        EXPECT_GE(pool.size(), c.minConnectionPoolSize_);
    }
}

int main(int argc, char** argv) {
    testing::InitGoogleTest(&argc, argv);
    nebula::client::init(&argc, &argv);
    google::SetStderrLogging(google::INFO);

    return RUN_ALL_TESTS();
}
