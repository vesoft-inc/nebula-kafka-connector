// Copyright (c) 2022 vesoft inc. All rights reserved.

#pragma once

#include <cstdint>
#include <string>

namespace nebula::client {

struct Config {
    std::uint32_t timeout_{0};   // in ms
    std::uint32_t idleTime_{0};  // in ms
    std::uint32_t maxConnectionPoolSize_{10};
    std::uint32_t minConnectionPoolSize_{0};
    std::string CAPath_;
    bool enableSSL_{false};
};

}  // namespace nebula::client
