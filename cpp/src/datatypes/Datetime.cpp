// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "common/datatype/Datetime.h"

#include "common/datatype/Duration.h"

namespace nebula::client {

void LocalDatetime::addDuration(const Duration& duration) {
    // The origin fields of Datetime is unsigned, but there will be some negative intermediate
    // results so I define some variable(field, tYear, tMonth) for this.
    int64_t tmp{0}, carry{0}, field{0};
    tmp = month + duration.months;
    if (std::abs(tmp) > 12) {
        carry = tmp / 12;
        /*month*/ field = tmp % 12;
    } else {
        /*month*/ field = tmp;
    }
    if (/*month*/ field <= 0) {
        carry--;
        /*month*/ field += 12;
    }
    year += carry;
    carry = 0;
    month = field;
    field = 0;

    tmp = microsec + duration.microsecondsInSecond();
    if (std::abs(tmp) >= 1000000) {
        carry = tmp / 1000000;
        /*microsec*/ field = tmp % 1000000;
    } else {
        /*microsec*/ field = tmp;
    }
    if (/*microsec*/ field < 0) {
        carry--;
        /*microsec*/ field += 1000000;
    }
    microsec = field;
    field = 0;
    tmp = sec + duration.seconds + carry;
    carry = 0;
    if (std::abs(tmp) >= 60) {
        carry = tmp / 60;
        /*sec*/ field = tmp % 60;
    } else {
        /*sec*/ field = tmp;
    }
    if (/*sec*/ field < 0) {
        carry--;
        /*sec*/ field += 60;
    }
    sec = field;
    field = 0;
    tmp = minute + carry;
    carry = 0;
    if (std::abs(tmp) >= 60) {
        carry = tmp / 60;
        /*minute*/ field = tmp % 60;
    } else {
        /*minute*/ field = tmp;
    }
    if (/*minute*/ field < 0) {
        carry--;
        /*minute*/ field += 60;
    }
    minute = field;
    field = 0;
    tmp = hour + carry;
    carry = 0;
    if (std::abs(tmp) >= 24) {
        carry = tmp / 24;
        /*hour*/ field = tmp % 24;
    } else {
        /*hour*/ field = tmp;
    }
    if (/*hour*/ field < 0) {
        carry--;
        /*hour*/ field += 24;
    }
    hour = field;
    field = 0;

    tmp = day + carry;
    carry = 0;
    if (tmp > 0) {
        int8_t dom = dayOfMonth(year, month);
        while (tmp > dom) {
            tmp -= dom;
            month += 1;
            if (month > 12) {
                year += 1;
                month = 1;
            }
            dom = dayOfMonth(year, month);
        }
    } else {
        int8_t dom = (month == 1 ? dayOfMonth(year - 1, 12) : dayOfMonth(year, month - 1));
        int64_t tMonth = month;
        int64_t tYear = year;
        while (tmp <= 0) {
            tmp += dom;
            /*month*/ tMonth -= 1;
            if (/*month*/ tMonth <= 0) {
                /*year*/ tYear--;
                /*month*/ tMonth += 12;
            }
            dom = (/*month*/ tMonth == 1 ? dayOfMonth(/*year*/ tYear - 1, 12)
                                         : dayOfMonth(/*year*/ tYear, /*month*/ tMonth - 1));
        }
        month = tMonth;
        year = tYear;
    }
    day = tmp;
}

void LocalDatetime::subDuration(const Duration& duration) {
    return addDuration(-duration);
}

// TODO(Aiee) +00:00 is the timezone offset for datetime instead of local datetime.
// Remove when we support timezone.
std::string LocalDatetime::toString() const {
    // It's in current timezone already
    return folly::sformat("DATETIME {}-{:0>2}-{:0>2}T{:0>2}:{:0>2}:{:0>2}.{:0>6}+00:00",
                          static_cast<int16_t>(year),
                          static_cast<uint8_t>(month),
                          static_cast<uint8_t>(day),
                          static_cast<uint8_t>(hour),
                          static_cast<uint8_t>(minute),
                          static_cast<uint8_t>(sec),
                          static_cast<uint32_t>(microsec));
}

}  // namespace nebula::client

namespace std {

// Inject a customized hash function
std::size_t hash<nebula::client::LocalDatetime>::operator()(
        const nebula::client::LocalDatetime& h) const noexcept {
    return h.qword;
}

// std::size_t hash<nebula::client::Datetime>::operator()(
//         const nebula::client::Datetime& h) const noexcept {
//     std::size_t hv =
//             folly::hash::fnv64_buf(reinterpret_cast<const void*>(&h.dt), sizeof(h.dt));
//     return folly::hash::fnv64_buf(
//             reinterpret_cast<const void*>(&h.timezoneOffset), sizeof(h.timezoneOffset), hv);
//     return 1;
// }

}  // namespace std
