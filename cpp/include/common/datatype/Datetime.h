// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_DATETIME_H_
#define COMMON_DATATYPE_DATETIME_H_

#include "common/datatype/Date.h"
#include "common/datatype/Time.h"
namespace nebula::client {

// In nebula only store UTC time, and the interpretation of time value based on
// the timezone configuration in current system.

int8_t dayOfMonth(int16_t year, int8_t month);

// An instant capturing the date and the time, but not the time zone.
//
// ISO 8601-1:2019 is the standard for date and time representation specified in the doc.
// https://en.wikipedia.org/wiki/ISO_8601
struct LocalDatetime {
#if defined(__GNUC__)
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wpedantic"
#endif  // defined(__GNUC__)
    union {
        // Using bit fields to specify the memory layout so datetime is 8 bytes
        struct {
            int64_t year : 16;
            uint64_t month : 4;
            uint64_t day : 5;
            uint64_t hour : 5;
            uint64_t minute : 6;
            uint64_t sec : 6;
            uint32_t microsec : 22;
        };
        uint64_t qword;
    };
#if defined(__GNUC__)
#pragma GCC diagnostic pop
#endif  // defined(__GNUC__)

    LocalDatetime() : year{0}, month{1}, day{1}, hour{0}, minute{0}, sec{0}, microsec{0} {}
    LocalDatetime(int16_t y, int8_t m, int8_t d, int8_t h, int8_t min, int8_t s, int32_t us) {
        year = y;
        month = m;
        day = d;
        hour = h;
        minute = min;
        sec = s;
        microsec = us;
    }
    explicit LocalDatetime(const Date& date) {
        year = date.year;
        month = date.month;
        day = date.day;
        hour = 0;
        minute = 0;
        sec = 0;
        microsec = 0;
    }
    LocalDatetime(const Date& date, const LocalTime& time) {
        year = date.year;
        month = date.month;
        day = date.day;
        hour = time.hour;
        minute = time.minute;
        sec = time.sec;
        microsec = time.microsec;
    }

    Date date() const {
        return Date(year, month, day);
    }

    LocalTime time() const {
        return LocalTime(hour, minute, sec, microsec);
    }

    void clear() {
        year = 0;
        month = 1;
        day = 1;
        hour = 0;
        minute = 0;
        sec = 0;
        microsec = 0;
    }

    void __clear() {
        clear();
    }

    bool operator==(const LocalDatetime& rhs) const {
        return year == rhs.year && month == rhs.month && day == rhs.day && hour == rhs.hour &&
               minute == rhs.minute && sec == rhs.sec && microsec == rhs.microsec;
    }
    bool operator<(const LocalDatetime& rhs) const {
        if (!(year == rhs.year)) {
            return year < rhs.year;
        }
        if (!(month == rhs.month)) {
            return month < rhs.month;
        }
        if (!(day == rhs.day)) {
            return day < rhs.day;
        }
        if (!(hour == rhs.hour)) {
            return hour < rhs.hour;
        }
        if (!(minute == rhs.minute)) {
            return minute < rhs.minute;
        }
        if (!(sec == rhs.sec)) {
            return sec < rhs.sec;
        }
        if (!(microsec == rhs.microsec)) {
            return microsec < rhs.microsec;
        }
        return false;
    }

    void addDuration(const Duration& duration);
    void subDuration(const Duration& duration);

    std::string toString() const;
    // 'Z' representing UTC timezone
    folly::dynamic toJson() const {
        return toString() + "Z";
    }
};

inline std::ostream& operator<<(std::ostream& os, const LocalDatetime& d) {
    os << d.toString();
    return os;
}

inline LocalDatetime operator+(const LocalDatetime& l, const Duration& r) {
    LocalDatetime dt = l;
    dt.addDuration(r);
    return dt;
}

inline LocalDatetime operator-(const LocalDatetime& l, const Duration& r) {
    LocalDatetime dt = l;
    dt.subDuration(r);
    return dt;
}

}  // namespace nebula::client

namespace std {

// Inject a customized hash function
template <>
struct hash<nebula::client::LocalDatetime> {
    std::size_t operator()(const nebula::client::LocalDatetime& h) const noexcept;
};

}  // namespace std

#endif  // COMMON_DATATYPE_DATETIME_H_
