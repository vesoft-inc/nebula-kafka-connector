// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_LIST_H_
#define COMMON_DATATYPE_LIST_H_

#include <algorithm>
#include <memory_resource>

#include "common/datatype/Value.h"
namespace nebula::client {
class List final {
public:
    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept {
        return values_.get_allocator();
    }

    explicit List(const allocator_type& alloc = allocator_type()) : values_(alloc) {}
    explicit List(const std::pmr::vector<Value>& values,
                  const allocator_type& alloc = allocator_type())
            : values_(values.begin(), values.end(), alloc) {}

    List(const List& other, const allocator_type& alloc) : values_(other.values_, alloc) {}
    List(List&& other, const allocator_type& alloc) noexcept
            : values_(std::move(other.values_), alloc) {}

    const std::pmr::vector<Value>& getValues() const {
        return values_;
    }

    size_t size() const {
        return values_.size();
    }

    const Value& operator[](size_t i) const {
        return values_[i];
    }

    template <typename... Args>
    void emplaceBack(Args&&... args) {
        values_.emplace_back(std::forward<Args>(args)...);
    }

    std::string toString() const;

private:
    // Serialization using thrift
    friend class apache::thrift::Cpp2Ops<List, void>;
    std::pmr::vector<Value> values_;
};

inline std::ostream& operator<<(std::ostream& os, const List& list) {
    return os << list.toString();
}

bool operator==(const List& lhs, const List& rhs);
bool operator!=(const List& lhs, const List& rhs);


}  // namespace nebula::client

namespace std {

template <>
struct hash<nebula::client::List> {
    std::size_t operator()(const nebula::client::List& list) const noexcept;
};

}  // namespace std

#endif  // COMMON_DATATYPE_LIST_H_
