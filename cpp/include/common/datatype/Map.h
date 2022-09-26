// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_MAP_H_
#define COMMON_DATATYPE_MAP_H_

#include <unordered_map>

#include "common/datatype/Value.h"

namespace nebula::client {
class Map final {
public:
    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept {
        return values_.get_allocator();
    }

    explicit Map(const allocator_type& alloc = allocator_type()) : values_(alloc) {}
    explicit Map(const std::pmr::unordered_map<Value, Value>& values,
                 const allocator_type& alloc = allocator_type())
            : values_(values, alloc) {}

    Map(const Map& other, const allocator_type& alloc) : values_(other.getValues(), alloc) {}
    Map(Map&& other, const allocator_type& alloc) noexcept
            : values_(std::move(other.values_), alloc) {}

    const std::pmr::unordered_map<Value, Value>& getValues() const {
        return values_;
    }

    size_t size() const {
        return values_.size();
    }

    const Value& operator[](const Value& key) const {
        return values_.at(key);
    }

    template <typename... Args>
    auto emplace(Args&&... args) {
        values_.emplace(std::forward<Args>(args)...);
    }

    std::string toString() const;

private:
    friend class apache::thrift::Cpp2Ops<Map, void>;
    // The key of gql's map can be other types except string.
    std::pmr::unordered_map<Value, Value> values_;
};

inline std::ostream& operator<<(std::ostream& os, const Map& map) {
    return os << map.toString();
}

bool operator==(const Map& lhs, const Map& rhs);
bool compareWithoutDynamicId(const Map& lhs, const Map& rhs);
bool operator!=(const Map& lhs, const Map& rhs);


}  // namespace nebula::client

namespace std {

template <>
struct hash<nebula::client::Map> {
    std::size_t operator()(const nebula::client::Map& map) const noexcept;
};

}  // namespace std

#endif  // COMMON_DATATYPE_MAP_H_
