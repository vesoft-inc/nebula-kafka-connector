// Copyright (c) 2022 vesoft inc. All rights reserved.

#pragma once

#ifndef MAYBE_UNUSED
#if (__cplusplus >= 201703L)  // c++17
#include <folly/CppAttributes.h>
#define MAYBE_UNUSED FOLLY_MAYBE_UNUSED
#else
#define MAYBE_UNUSED __attribute__((unused))
#endif
#endif
