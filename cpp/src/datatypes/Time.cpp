// Copyright (c) 2022 vesoft inc. All rights reserved.


#include "common/datatype/Time.h"

namespace nebula::client {

void LocalTime::addDuration(const Duration& duration) {
    int64_t tmp{0}, carry{0};
    tmp = microsec + duration.microsecondsInSecond();
    if (std::abs(tmp) >= 1000000) {
        carry = tmp / 1000000;
        microsec = tmp % 1000000;
    } else {
        microsec = tmp;
    }
    if (microsec < 0) {
        carry--;
        microsec += 1000000;
    }
    tmp = sec + duration.seconds + carry;
    carry = 0;
    if (std::abs(tmp) >= 60) {
        carry = tmp / 60;
        sec = tmp % 60;
    } else {
        sec = tmp;
    }
    if (sec < 0) {
        carry--;
        sec += 60;
    }
    tmp = minute + carry;
    carry = 0;
    if (std::abs(tmp) >= 60) {
        carry = tmp / 60;
        minute = tmp % 60;
    } else {
        minute = tmp;
    }
    if (minute < 0) {
        carry--;
        minute += 60;
    }
    tmp = hour + carry;
    carry = 0;
    if (std::abs(tmp) >= 24) {
        carry = tmp / 24;
        hour = tmp % 24;
    } else {
        hour = tmp;
    }
    if (hour < 0) {
        carry--;
        hour += 24;
    }
}

void LocalTime::subDuration(const Duration& duration) {
    addDuration(-duration);
}

std::string LocalTime::toString() const {
    return folly::sformat("TIME {:0>2}:{:0>2}:{:0>2}.{:0>6}",
                          static_cast<uint8_t>(hour),
                          static_cast<uint8_t>(minute),
                          static_cast<uint8_t>(sec),
                          static_cast<uint32_t>(microsec));
}


}  // namespace nebula::client

namespace std {

// Inject a customized hash function
std::size_t hash<nebula::client::LocalTime>::operator()(
        const nebula::client::LocalTime& h) const noexcept {
    std::size_t hv =
            folly::hash::fnv64_buf(reinterpret_cast<const void*>(&h.hour), sizeof(h.hour));
    hv = folly::hash::fnv64_buf(reinterpret_cast<const void*>(&h.minute), sizeof(h.minute), hv);
    hv = folly::hash::fnv64_buf(reinterpret_cast<const void*>(&h.sec), sizeof(h.sec), hv);
    return folly::hash::fnv64_buf(
            reinterpret_cast<const void*>(&h.microsec), sizeof(h.microsec), hv);
}

}  // namespace std
