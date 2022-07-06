// Copyright (c) 2022 vesoft inc. All rights reserved.

#pragma once

#include <folly/io/async/SSLContext.h>

namespace nebula::client {

std::shared_ptr<folly::SSLContext> createSSLContext(const std::string &CAPath);

}  // namespace nebula::client
