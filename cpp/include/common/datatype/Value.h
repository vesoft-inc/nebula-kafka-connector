// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_VALUE_H_
#define COMMON_DATATYPE_VALUE_H_

#include <memory_resource>

namespace apache {
namespace thrift {

template <class, class>
class Cpp2Ops;

}  // namespace thrift
}  // namespace apache

namespace nebula::client {

class List;
class Map;
class Node;
class Edge;

class Value final {
public:
    static const Value kNullValue;

    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept;

    // TODO(jie): Consider the null value.
    // TODO(jie): Redesign the any type and other abstract types
    enum class Type : uint16_t {
        kNull = 0,
        kBool = 1,
        kInt8 = 2,
        kInt16 = 3,
        kInt32 = 4,
        kInt64 = 5,
        kFloat = 6,
        kDouble = 7,
        kString = 8,
        kList = 9,
        kMap = 10,
        kNode = 11,
        kEdge = 12,
    };
    union Data {
        // Fixed length data types less or equal than 8 bytes
        bool bool_;
        int8_t int8_;
        int16_t int16_;
        int32_t int32_;
        int64_t int64_;
        float float_;
        double double_;
        // Variable length data types or fixed length data types greater than 8 bytes
        // TODO: use
        std::pmr::string* string_;
        List* list_;
        Map* map_;
        Node* node_;
        Edge* edge_;

        // A little tricky here.
        // For an allocator-aware `Value`, we need 8 bytes to store the pmr allocator.
        // Considering that the primitive types don't need the allocator while the
        // non-primitive types store the allocator themselves, so we can use the union to store
        // the allocator temporarily.
        std::pmr::memory_resource* mr_;
    };

    Value(allocator_type alloc = allocator_type());                          // NOLINT
    Value(bool val, allocator_type alloc = allocator_type());                // NOLINT
    Value(int8_t val, allocator_type alloc = allocator_type());              // NOLINT
    Value(int16_t val, allocator_type alloc = allocator_type());             // NOLINT
    Value(int32_t val, allocator_type alloc = allocator_type());             // NOLINT
    Value(int64_t val, allocator_type alloc = allocator_type());             // NOLINT
    Value(float val, allocator_type alloc = allocator_type());               // NOLINT
    Value(double val, allocator_type alloc = allocator_type());              // NOLINT
    Value(const char* val, allocator_type alloc = allocator_type());         // NOLINT
    Value(std::string_view val, allocator_type alloc = allocator_type());    // NOLINT
    Value(const std::string& val, allocator_type alloc = allocator_type());  // NOLINT
    Value(const std::pmr::string& val);                                      // NOLINT
    Value(const std::pmr::string& val, allocator_type alloc);                // NOLINT
    Value(std::pmr::string&& val) noexcept;                                  // NOLINT
    Value(std::pmr::string&& val, allocator_type alloc);                     // NOLINT
    Value(const List& val, allocator_type alloc = allocator_type());         // NOLINT
    Value(List&& val) noexcept;                                              // NOLINT
    Value(List&& val, allocator_type alloc);                                 // NOLINT
    Value(const Map& val, allocator_type alloc = allocator_type());          // NOLINT
    Value(Map&& val) noexcept;                                               // NOLINT
    Value(Map&& val, allocator_type alloc);                                  // NOLINT
    Value(const Node& val, allocator_type alloc = allocator_type());         // NOLINT
    Value(Node&& val) noexcept;                                              // NOLINT
    Value(Node&& val, allocator_type alloc);                                 // NOLINT
    Value(const Edge& val, allocator_type alloc = allocator_type());         // NOLINT
    Value(Edge&& val) noexcept;                                              // NOLINT
    Value(Edge&& val, allocator_type alloc);                                 // NOLINT

    Value(const Value& other);
    Value(const Value& other, allocator_type alloc);
    Value(Value&& other) noexcept;
    Value(Value&& other, allocator_type alloc);

    ~Value();

    Value& operator=(const Value& other);
    // TODO(jie): A copy may happen here, so cannot marked as noexcept?
    Value& operator=(Value&& other);

    Type getType() const {
        return type_;
    }

    bool isAny() const {
        return true;
    }

    bool isNumeric() const {
        // return type_ & kNumericMask == Numeric;
        return isInt8() || isInt16() || isInt32() || isInt64() || isFloat() || isDouble();
    }

    bool isNullValue() const {
        return type_ == Type::kNull;
    }

    bool isBool() const {
        return type_ == Type::kBool;
    }

    bool isInt8() const {
        return type_ == Type::kInt8;
    }

    bool isInt16() const {
        return type_ == Type::kInt16;
    }

    bool isInt32() const {
        return type_ == Type::kInt32;
    }

    bool isInt64() const {
        return type_ == Type::kInt64;
    }

    bool isInteger() const {
        switch (type_) {
            case Type::kInt8:
            case Type::kInt16:
            case Type::kInt32:
            case Type::kInt64:
                return true;
            default:
                return false;
        }
    }

    bool isFloat() const {
        return type_ == Type::kFloat;
    }

    bool isDouble() const {
        return type_ == Type::kDouble;
    }

    bool isString() const {
        return type_ == Type::kString;
    }

    bool isList() const {
        return type_ == Type::kList;
    }

    bool isMap() const {
        return type_ == Type::kMap;
    }

    bool isNode() const {
        return type_ == Type::kNode;
    }

    bool isEdge() const {
        return type_ == Type::kEdge;
    }

    bool getBool() const {
        return data_.bool_;
    }

    int8_t getInt8() const {
        return data_.int8_;
    }

    int16_t getInt16() const {
        return data_.int16_;
    }

    int32_t getInt32() const {
        return data_.int32_;
    }

    int64_t getInt64() const {
        return data_.int64_;
    }

    int64_t getInteger() const;

    float getFloat() const {
        return data_.float_;
    }

    double getDouble() const {
        return data_.double_;
    }

    const std::pmr::string& getString() const {
        return *data_.string_;
    }

    std::string_view getStringView() const {
        return {data_.string_->data(), data_.string_->size()};
    }

    const List& getList() const {
        return *data_.list_;
    }

    const Map& getMap() const {
        return *data_.map_;
    }

    const Node& getNode() const {
        return *data_.node_;
    }

    const Edge& getEdge() const {
        return *data_.edge_;
    }

    bool& mutableBool() {
        return data_.bool_;
    }

    int8_t& mutableInt8() {
        return data_.int8_;
    }

    int16_t& mutableInt16() {
        return data_.int16_;
    }

    int32_t& mutableInt32() {
        return data_.int32_;
    }

    int64_t& mutableInt64() {
        return data_.int64_;
    }

    float& mutableFloat() {
        return data_.float_;
    }

    double& mutableDouble() {
        return data_.double_;
    }

    std::pmr::string& mutableString() {
        return *data_.string_;
    }

    List& mutableList() {
        return *data_.list_;
    }

    Map& mutableMap() {
        return *data_.map_;
    }

    Node& mutableNode() {
        return *data_.node_;
    }

    Edge& mutableEdge() {
        return *data_.edge_;
    }

    std::string toString() const;

    void clear();

    // To be compatible with thrift serialization.
    void __clear() {
        clear();
    }

private:
    void copyValue(const Value& other, const allocator_type& alloc);
    void moveValue(Value&& other) noexcept;

private:
    template <typename T>
    T* allocate(const allocator_type& alloc) {
        return static_cast<T*>(alloc.resource()->allocate(sizeof(T), alignof(T)));
    }

    template <typename T>
    void deallocate(const allocator_type& alloc, T* p) noexcept {
        alloc.resource()->deallocate(p, sizeof(T), alignof(T));
    }

    template <typename T, typename... Args>
    void construct(T* p, Args&&... args) {
        // Is it need to handle the exception of constructor here?
        ::new (p) T(std::forward<Args>(args)...);
    }

    template <typename T>
    void destroy(T* p) {
        p->~T();
    }

    template <typename T, typename... Args>
    T* newObject(const allocator_type& alloc, Args&&... args) {
        T* p = allocate<T>(alloc);
        // Don't forget to forward the alloc to the constructor
        construct(p, std::forward<Args>(args)..., alloc);
        return p;
    }

    template <typename T>
    void deleteObject(const allocator_type& alloc, T* p) {
        destroy(p);
        deallocate(alloc, p);
    }

private:
    // Serialization using fbthrift
    friend class apache::thrift::Cpp2Ops<Value, void>;
    Type type_;
    Data data_;
};

using ValueType = Value::Type;

std::string enum2String(const ValueType);
inline std::ostream& operator<<(std::ostream& os, ValueType type) {
    return os << enum2String(type);
}
inline std::ostream& operator<<(std::ostream& os, const Value& value) {
    return os << value.toString();
}

// Comparison operations
bool operator==(const Value& lhs, const Value& rhs);
bool operator!=(const Value& lhs, const Value& rhs);
bool operator<(const Value& lhs, const Value& rhs);
bool operator>(const Value& lhs, const Value& rhs);
bool operator<=(const Value& lhs, const Value& rhs);
bool operator>=(const Value& lhs, const Value& rhs);
// unary operations
Value operator-(const Value& rhs);

}  // namespace nebula::client

namespace std {

template <>
struct hash<nebula::client::Value> {
    std::size_t operator()(const nebula::client::Value& v) const noexcept;
};

}  // namespace std

#endif  // COMMON_DATATYPE_VALUE_H_
