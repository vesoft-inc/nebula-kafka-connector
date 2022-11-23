// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_DURATION_H
#define COMMON_DATATYPE_DURATION_H

#include <folly/dynamic.h>

#include <sstream>

#include "common/time/Constants.h"
namespace nebula::client {

// TODO(Aiee) refactor this class according ISO 8601-1:2019, Date and time
//
// Duration equals to months + seconds + microseconds
// The base between months and days is not fixed, so we store years and months
// separately.
//
// The valid duration type size is 16 bytes: 8bytes for days and 8 bytes for secs, 4 bytes for
// ms, 4 bytes for months
struct Duration {
    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept {
        return alloc_;
    }

    // day + hours + minutes + seconds + microseconds
    int64_t seconds{0};
    int32_t microseconds{0};
    // years + months
    int32_t months{0};
    // TODO(Aiee) this alloc is used in newObject()
    allocator_type alloc_{};

    explicit Duration(const allocator_type& alloc = allocator_type()) : alloc_(alloc) {}
    Duration(int32_t m, int64_t s, int32_t us, allocator_type alloc = allocator_type())
            : seconds(s), microseconds(us), months(m), alloc_(alloc) {}
    Duration(const Duration& other, const allocator_type& alloc)
            : seconds(other.seconds),
              microseconds(other.microseconds),
              months(other.months),
              alloc_(alloc) {}


    Duration(Duration&& other, const allocator_type& alloc) noexcept
            : seconds(std::move(other.seconds)),
              microseconds(std::move(other.microseconds)),
              months(std::move(other.months)),
              alloc_(alloc) {}


    int64_t years() const {
        return months / 12;
    }

    int64_t monthsInYear() const {
        return months % 12;
    }

    int64_t days() const {
        return seconds / time::kSecondsOfDay;
    }

    int64_t hours() const {
        return seconds % time::kSecondsOfDay / time::kSecondsOfHour;
    }

    int64_t minutes() const {
        return seconds % time::kSecondsOfHour / time::kSecondsOfMinute;
    }

    int64_t secondsInMinute() const {
        return seconds % time::kSecondsOfMinute;
    }

    int64_t microsecondsInSecond() const {
        return microseconds;
    }

    Duration operator-() const {
        return {-months, -seconds, -microseconds};
    }

    Duration operator+(const Duration& rhs) const {
        return {months + rhs.months, seconds + rhs.seconds, microseconds + rhs.microseconds};
    }

    Duration operator-(const Duration& rhs) const {
        return {months - rhs.months, seconds - rhs.seconds, microseconds - rhs.microseconds};
    }

    Duration& addYears(int32_t y) {
        months += y * 12;
        return *this;
    }

    Duration& addQuarters(int32_t q) {
        months += q * 3;
        return *this;
    }

    Duration& addMonths(int32_t m) {
        months += m;
        return *this;
    }

    Duration& addWeeks(int32_t w) {
        seconds += (w * 7 * time::kSecondsOfDay);
        return *this;
    }

    Duration& addDays(int64_t d) {
        seconds += d * time::kSecondsOfDay;
        return *this;
    }

    Duration& addHours(int64_t h) {
        seconds += h * time::kSecondsOfHour;
        return *this;
    }

    Duration& addMinutes(int64_t minutes) {
        seconds += minutes * time::kSecondsOfMinute;
        return *this;
    }

    Duration& addSeconds(int64_t s) {
        seconds += s;
        return *this;
    }

    Duration& addMilliseconds(int64_t ms) {
        seconds += ms / 1000;
        microseconds += ((ms % 1000) * 1000);
        return *this;
    }

    Duration& addMicroseconds(int32_t us) {
        microseconds += us;
        return *this;
    }

    // can't compare
    bool operator<(const Duration& rhs) const {
        (void)(rhs);
        return false;
    }

    bool operator==(const Duration& rhs) const {
        return months == rhs.months && seconds == rhs.seconds &&
               microseconds == rhs.microseconds;
    }

    void addDuration(const Duration& duration);
    void subDuration(const Duration& duration);

    std::string toString() const {
        return folly::sformat("P{}MT{}.{:0>6}000S",
                              months,
                              seconds + microseconds / 1000000,
                              microseconds % 1000000);
    }

    folly::dynamic toJson() const {
        return toString();
    }
};

}  // namespace nebula::client

namespace std {

// Inject a customized hash function
template <>
struct hash<nebula::client::Duration> {
    std::size_t operator()(const nebula::client::Duration& d) const noexcept {
        size_t hv = folly::hash::fnv64_buf(reinterpret_cast<const void*>(&d.months),
                                           sizeof(d.months));
        hv = folly::hash::fnv64_buf(
                reinterpret_cast<const void*>(d.seconds), sizeof(d.seconds), hv);
        return folly::hash::fnv64_buf(
                reinterpret_cast<const void*>(d.microseconds), sizeof(d.microseconds), hv);
    }
};

}  // namespace std
#endif
