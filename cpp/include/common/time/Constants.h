// Copyright (c) 2022 vesoft inc. All rights reserved.

#pragma once

namespace nebula::client {
namespace time {

static constexpr int kDayOfLeapYear = 366;
static constexpr int kDayOfCommonYear = 365;

static constexpr int64_t kSecondsOfMinute = 60;
static constexpr int64_t kSecondsOfHour = 60 * kSecondsOfMinute;
static constexpr int64_t kSecondsOfDay = 24 * kSecondsOfHour;

}  // namespace time
}  // namespace nebula::client
