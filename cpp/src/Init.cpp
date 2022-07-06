// Copyright (c) 2022 vesoft inc. All rights reserved.

#include <folly/init/Init.h>

namespace nebula::client {

void init(int *argc, char **argv[]) {
    folly::init(argc, argv, true);
}

}  // namespace nebula::client
