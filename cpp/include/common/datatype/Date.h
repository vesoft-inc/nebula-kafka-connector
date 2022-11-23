// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_DATE_H_
#define COMMON_DATATYPE_DATE_H_

#include <cstdint>
#include <ostream>

#include "common/datatype/Duration.h"

namespace nebula::client {

// https://en.wikipedia.org/wiki/Leap_year#Leap_day
static inline bool isLeapYear(int16_t year) {
    if (year % 4 != 0) {
        return false;
    } else if (year % 100 != 0) {
        return true;
    } else if (year % 400 != 0) {
        return false;
    } else {
        return true;
    }
}

// In nebula only store UTC time, and the interpretation of time value based on
// the timezone configuration in current system.

const int64_t kDaysSoFar[] = {0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334, 365};
const int64_t kLeapDaysSoFar[] = {0, 31, 60, 91, 121, 152, 182, 213, 244, 274, 305, 335, 366};

// An instant capturing the date, but not the time, nor the time zone.
struct Date {
    int16_t year;  // Any integer
    int8_t month;  // 1 - 12
    int8_t day;    // 1 - 31

    Date() : year{0}, month{1}, day{1} {}
    Date(int16_t y, int8_t m, int8_t d) : year{y}, month{m}, day{d} {}
    // Tak the number of days since -32768/1/1, and convert to the real date
    explicit Date(uint64_t days);

    void clear() {
        year = 0;
        month = 1;
        day = 1;
    }

    void __clear() {
        clear();
    }

    void reset(int16_t y, int8_t m, int8_t d) {
        year = y;
        month = m;
        day = d;
    }

    Date date() const {
        return {year, month, day};
    }

    bool operator==(const Date& rhs) const {
        return year == rhs.year && month == rhs.month && day == rhs.day;
    }

    bool operator<(const Date& rhs) const {
        if (!(year == rhs.year)) {
            return year < rhs.year;
        }
        if (!(month == rhs.month)) {
            return month < rhs.month;
        }
        if (!(day == rhs.day)) {
            return day < rhs.day;
        }
        return false;
    }

    Date operator+(int64_t days) const;
    Date operator-(int64_t days) const;

    void addDuration(const Duration& duration);
    void subDuration(const Duration& duration);

    std::string toString() const;
    folly::dynamic toJson() const {
        return toString();
    }

    // Return the number of days since -32768/1/1
    int64_t toInt() const;
    // Convert the number of days since -32768/1/1 to the real date
    void fromInt(int64_t days);
};

inline Date operator+(const Date& l, const Duration& r) {
    Date d = l;
    d.addDuration(r);
    return d;
}

inline Date operator-(const Date& l, const Duration& r) {
    Date d = l;
    d.subDuration(r);
    return d;
}

inline std::ostream& operator<<(std::ostream& os, const Date& d) {
    os << d.toString();
    return os;
}

}  // namespace nebula::client

namespace std {

// Inject a customized hash function
template <>
struct hash<nebula::client::Date> {
    std::size_t operator()(const nebula::client::Date& h) const noexcept;
};

}  // namespace std
#endif  // COMMON_DATATYPE_DATE_H_
