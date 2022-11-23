// Copyright (c) 2022 vesoft inc. All rights reserved.
#include "common/datatype/Date.h"

namespace nebula::client {

int8_t dayOfMonth(int16_t year, int8_t month) {
    return isLeapYear(year) ? kLeapDaysSoFar[month] - kLeapDaysSoFar[month - 1]
                            : kDaysSoFar[month] - kDaysSoFar[month - 1];
}

Date::Date(uint64_t days) {
    fromInt(days);
}

int64_t Date::toInt() const {
    // Year
    int64_t yearsPassed = year + 32768L;
    int64_t days = yearsPassed * 365L;
    // Add one day per leap year
    if (yearsPassed > 0) {
        days += (yearsPassed - 1) / 4 + 1;
    }

    // Month
    if (yearsPassed % 4 == 0) {
        // Leap year
        days += kLeapDaysSoFar[month - 1];
    } else {
        days += kDaysSoFar[month - 1];
    }

    // Day
    days += day;

    // Since we start from -32768/1/1, we need to reduce one day
    return days - 1;
}

void Date::fromInt(int64_t days) {
    // year
    int64_t yearsPassed = (days + 1) / 365;
    year = yearsPassed - 32768;
    int64_t daysInYear = (days + 1) % 365;

    // Deduce the number of days for leap years
    if (yearsPassed > 0) {
        daysInYear -= (yearsPassed - 1) / 4 + 1;
    }

    // Adjust the year if necessary
    while (daysInYear <= 0) {
        year = year - 1;
        if (year % 4 == 0) {
            // Leap year
            daysInYear += 366;
        } else {
            daysInYear += 365;
        }
    }

    // Month and day
    month = 1;
    while (true) {
        if (year % 4 == 0) {
            // Leap year
            if (daysInYear <= kLeapDaysSoFar[month]) {
                day = daysInYear - kLeapDaysSoFar[month - 1];
                break;
            }
        } else {
            if (daysInYear <= kDaysSoFar[month]) {
                day = daysInYear - kDaysSoFar[month - 1];
                break;
            }
        }
        month++;
    }
}

Date Date::operator+(int64_t days) const {
    int64_t daysSince = toInt();
    return Date(daysSince + days);
}

Date Date::operator-(int64_t days) const {
    int64_t daysSince = toInt();
    return Date(daysSince - days);
}

void Date::addDuration(const Duration& duration) {
    int64_t tmp{0}, carry{0};
    tmp = month + duration.months;
    if (std::abs(tmp) > 12) {
        carry = tmp / 12;
        month = tmp % 12;
    } else {
        month = tmp;
    }
    if (month <= 0) {
        carry--;
        month += 12;
    }
    year += carry;

    tmp = day + duration.days();
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
        while (tmp <= 0) {
            tmp += dom;
            month -= 1;
            if (month <= 0) {
                year--;
                month += 12;
            }
            dom = (month == 1 ? dayOfMonth(year - 1, 12) : dayOfMonth(year, month - 1));
        }
    }
    day = tmp;
}

void Date::subDuration(const Duration& duration) {
    return addDuration(-duration);
}

std::string Date::toString() const {
    // It's in current timezone already
    return folly::stringPrintf("DATE %d-%02d-%02d", year, month, day);
}

}  // namespace nebula::client

namespace std {

// Inject a customized hash function
std::size_t hash<nebula::client::Date>::operator()(const nebula::client::Date& h) const noexcept {
    size_t hv = folly::hash::fnv64_buf(reinterpret_cast<const void*>(&h.year), sizeof(h.year));
    hv = folly::hash::fnv64_buf(reinterpret_cast<const void*>(&h.month), sizeof(h.month), hv);
    return folly::hash::fnv64_buf(reinterpret_cast<const void*>(&h.day), sizeof(h.day), hv);
}

}  // namespace std
