// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_TIME_H_
#define COMMON_DATATYPE_TIME_H_

#include <thrift/lib/cpp2/protocol/Cpp2Ops.h>

#include <cstdint>
#include <ostream>

#include "common/datatype/Date.h"
#include "common/datatype/Duration.h"
namespace nebula::client {

struct LocalTime {
    int8_t hour;
    int8_t minute;
    int8_t sec;
    int32_t microsec;

    LocalTime() : hour{0}, minute{0}, sec{0}, microsec{0} {}
    LocalTime(int8_t h, int8_t min, int8_t s, int32_t us)
            : hour{h}, minute{min}, sec{s}, microsec{us} {}

    void clear() {
        hour = 0;
        minute = 0;
        sec = 0;
        microsec = 0;
    }

    void __clear() {
        clear();
    }

    bool operator==(const LocalTime& rhs) const {
        return hour == rhs.hour && minute == rhs.minute && sec == rhs.sec &&
               microsec == rhs.microsec;
    }

    bool operator<(const LocalTime& rhs) const {
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

    // 'Z' representing UTC timezone
    std::string toString() const;

private:
    // Serialization using fbthrift
    friend class apache::thrift::Cpp2Ops<LocalTime>;
};

inline std::ostream& operator<<(std::ostream& os, const LocalTime& d) {
    os << d.toString();
    return os;
}

inline LocalTime operator+(const LocalTime& l, const Duration& r) {
    LocalTime t = l;
    t.addDuration(r);
    return t;
}

inline LocalTime operator-(const LocalTime& l, const Duration& r) {
    LocalTime t = l;
    t.subDuration(r);
    return t;
}

// struct Time {
//     Time() = default;

//     Time time;
//     int32_t timezoneOffset;
// };

}  // namespace nebula::client

namespace std {

// Inject a customized hash function
template <>
struct hash<nebula::client::LocalTime> {
    std::size_t operator()(const nebula::client::LocalTime& h) const noexcept;
};

// template <>
// struct hash<nebula::client::Time> {
//     std::size_t operator()(const nebula::client::Time& h) const noexcept;
// };

}  // namespace std

#endif  // COMMON_DATATYPE_TIME_H_
