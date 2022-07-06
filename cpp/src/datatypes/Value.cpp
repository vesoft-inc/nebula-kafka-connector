// Copyright (c) 2022 vesoft inc. All rights reserved.


#include "common/datatype/Value.h"

#include <folly/Conv.h>
#include <glog/logging.h>

#include <memory_resource>

#include "common/datatype/Edge.h"
#include "common/datatype/List.h"
#include "common/datatype/Map.h"
#include "common/datatype/Node.h"

namespace nebula::client {

const Value Value::kNullValue{};

Value::allocator_type Value::get_allocator() const noexcept {
    switch (type_) {
        case Type::kNull:
            // TODO(jie): Check the consistency of the alloc copy.
            return {data_.mr_};
        case Type::kBool:
        case Type::kInt8:
        case Type::kInt16:
        case Type::kInt32:
        case Type::kInt64:
        case Type::kFloat:
        case Type::kDouble:
            return {};
        case Type::kString:
            return data_.string_->get_allocator();
        case Type::kList:
            return data_.list_->get_allocator();
        case Type::kMap:
            return data_.map_->get_allocator();
        case Type::kNode:
            return data_.node_->get_allocator();
        case Type::kEdge:
            return data_.edge_->get_allocator();
    }
    LOG(FATAL) << "Unknown value type: " << type_;
}

Value::Value(allocator_type alloc) : type_(Type::kNull) {
    data_.mr_ = alloc.resource();
}
Value::Value(bool val, allocator_type) : type_(Type::kBool) {
    data_.bool_ = val;
}
Value::Value(int8_t val, allocator_type) : type_(Type::kInt8) {
    data_.int8_ = val;
}
Value::Value(int16_t val, allocator_type) : type_(Type::kInt16) {
    data_.int16_ = val;
}
Value::Value(int32_t val, allocator_type) : type_(Type::kInt32) {
    data_.int32_ = val;
}
Value::Value(int64_t val, allocator_type) : type_(Type::kInt64) {
    data_.int64_ = val;
}
Value::Value(float val, allocator_type) : type_(Type::kFloat) {
    data_.float_ = val;
}
Value::Value(double val, allocator_type) : type_(Type::kDouble) {
    data_.double_ = val;
}

Value::Value(const char* val, allocator_type alloc) : type_(Type::kString) {
    data_.string_ = newObject<std::pmr::string>(alloc, val);
}

Value::Value(std::string_view val, allocator_type alloc) : type_(Type::kString) {
    data_.string_ = newObject<std::pmr::string>(alloc, val);
}

Value::Value(const std::string& val, allocator_type alloc) : type_(Type::kString) {
    data_.string_ = newObject<std::pmr::string>(alloc, val);
}

Value::Value(const std::pmr::string& val)
        : Value(val,
                std::allocator_traits<allocator_type>::select_on_container_copy_construction(
                        val.get_allocator())) {}

Value::Value(const std::pmr::string& val, allocator_type alloc) : type_(Type::kString) {
    data_.string_ = newObject<std::pmr::string>(alloc, val);
}

Value::Value(std::pmr::string&& val) noexcept : Value(std::move(val), val.get_allocator()) {}

Value::Value(std::pmr::string&& val, allocator_type alloc) : type_(Type::kString) {
    if (val.get_allocator() == alloc) {
        data_.string_ = newObject<std::pmr::string>(alloc, std::move(val));
    } else {
        data_.string_ = newObject<std::pmr::string>(alloc, val);
    }
}

Value::Value(const List& val, allocator_type alloc) : type_(Type::kList) {
    data_.list_ = newObject<List>(alloc, val);
}

Value::Value(List&& val) noexcept : Value(std::move(val), val.get_allocator()) {}

Value::Value(List&& val, allocator_type alloc) : type_(Type::kList) {
    if (val.get_allocator() == alloc) {
        data_.list_ = newObject<List>(alloc, std::move(val));
    } else {
        data_.list_ = newObject<List>(alloc, val);
    }
}

Value::Value(const Map& val, allocator_type alloc) : type_(Type::kMap) {
    data_.map_ = newObject<Map>(alloc, val);
}

Value::Value(Map&& val) noexcept : Value(std::move(val), val.get_allocator()) {}

Value::Value(Map&& val, allocator_type alloc) : type_(Type::kMap) {
    if (val.get_allocator() == alloc) {
        data_.map_ = newObject<Map>(alloc, std::move(val));
    } else {
        data_.map_ = newObject<Map>(alloc, val);
    }
}

Value::Value(const Node& val, allocator_type alloc) : type_(Type::kNode) {
    data_.node_ = newObject<Node>(alloc, val);
}

Value::Value(Node&& val) noexcept : Value(std::move(val), val.get_allocator()) {}

Value::Value(Node&& val, allocator_type alloc) : type_(Type::kNode) {
    if (val.get_allocator() == alloc) {
        data_.node_ = newObject<Node>(alloc, std::move(val));
    } else {
        data_.node_ = newObject<Node>(alloc, val);
    }
}

Value::Value(const Edge& val, allocator_type alloc) : type_(Type::kEdge) {
    data_.edge_ = newObject<Edge>(alloc, val);
}

Value::Value(Edge&& val) noexcept : Value(std::move(val), val.get_allocator()) {}

Value::Value(Edge&& val, allocator_type alloc) : type_(Type::kEdge) {
    if (val.get_allocator() == alloc) {
        data_.edge_ = newObject<Edge>(alloc, std::move(val));
    } else {
        data_.edge_ = newObject<Edge>(alloc, val);
    }
}

// std::allocator_traits<allocator_type>::select_on_container_copy_construction(other.get_allocator())
Value::Value(const Value& other) : Value(other, allocator_type()) {}

Value::Value(const Value& other, allocator_type alloc) : type_(other.getType()) {
    copyValue(other, alloc);
}

Value::Value(Value&& other) noexcept : Value(std::move(other), other.get_allocator()) {}

Value::Value(Value&& other, allocator_type alloc) : type_(other.getType()) {
    if (other.get_allocator() == alloc) {
        moveValue(std::move(other));
    } else {
        copyValue(other, alloc);
    }
}

Value::~Value() {
    clear();
}

int64_t Value::getInteger() const {
    switch (type_) {
        case Type::kInt8:
            return data_.int8_;
        case Type::kInt16:
            return data_.int16_;
        case Type::kInt32:
            return data_.int32_;
        case Type::kInt64:
            return data_.int64_;
        default:
            LOG(FATAL) << "Maleformed value type: " << static_cast<int>(type_);
    }
}

// pmr allocator does not propagate on container copy assignment.
Value& Value::operator=(const Value& other) {
    if (this != &other) {
        // Still use the alloc of this
        auto alloc = get_allocator();
        clear();
        copyValue(other, alloc);
    }
    return *this;
}

// pmr allocator does not propagate on container move assignment.
Value& Value::operator=(Value&& other) {
    if (this != &other) {
        // Still use the alloc of this
        auto alloc = get_allocator();
        clear();
        if (alloc == other.get_allocator()) {
            moveValue(std::move(other));
        } else {
            copyValue(other, alloc);
        }
    }
    return *this;
}

std::string Value::toString() const {
    switch (type_) {
        case Value::Type::kNull: {
            return "NULL";
        }
        case Value::Type::kBool: {
            return getBool() ? "true" : "false";
        }
        case Value::Type::kInt8: {
            return folly::to<std::string>(getInt8());
        }
        case Value::Type::kInt16: {
            return folly::to<std::string>(getInt16());
        }
        case Value::Type::kInt32: {
            return folly::to<std::string>(getInt32());
        }
        case Value::Type::kInt64: {
            return folly::to<std::string>(getInt64());
        }
        case Value::Type::kFloat: {
            return folly::to<std::string>(getFloat());
        }
        case Value::Type::kDouble: {
            return folly::to<std::string>(getDouble());
        }
        case Type::kString: {
            return "\"" + std::string(getString()) + "\"";
        }
        case Type::kList: {
            return data_.list_->toString();
        }
        case Type::kMap: {
            return data_.map_->toString();
        }
        case Type::kNode: {
            return data_.node_->toString();
        }
        case Type::kEdge: {
            return data_.edge_->toString();
        }
    }
    LOG(FATAL) << "Unknown value type: " << type_;
}

void Value::clear() {
    switch (type_) {
        case Type::kNull:
            // TODO(jie): Clear the mr_?
        case Type::kBool:
        case Type::kInt8:
        case Type::kInt16:
        case Type::kInt32:
        case Type::kInt64:
        case Type::kFloat:
        case Type::kDouble:
            break;
        case Type::kString: {
            if (data_.string_) {
                deleteObject(data_.string_->get_allocator(), data_.string_);
            }
            break;
        }
        case Type::kList: {
            if (data_.list_) {
                deleteObject(data_.list_->get_allocator(), data_.list_);
            }
            break;
        }
        case Type::kMap: {
            if (data_.map_) {
                deleteObject(data_.map_->get_allocator(), data_.map_);
            }
            break;
        }
        case Type::kNode: {
            if (data_.node_) {
                deleteObject(data_.node_->get_allocator(), data_.node_);
            }
            break;
        }
        case Type::kEdge: {
            if (data_.edge_) {
                deleteObject(data_.edge_->get_allocator(), data_.edge_);
            }
            break;
        }
    }
    type_ = Type::kNull;
}


void Value::copyValue(const Value& other, const allocator_type& alloc) {
    type_ = other.type_;
    switch (other.getType()) {
        case Type::kNull:
            // TODO(jie): Check the consistency of the allocator copy.
            data_.mr_ = alloc.resource();
            break;
        case Type::kBool:
            data_.bool_ = other.getBool();
            break;
        case Type::kInt8:
            data_.int8_ = other.getInt8();
            break;
        case Type::kInt16:
            data_.int16_ = other.getInt16();
            break;
        case Type::kInt32:
            data_.int32_ = other.getInt32();
            break;
        case Type::kInt64:
            data_.int64_ = other.getInt64();
            break;
        case Type::kFloat:
            data_.float_ = other.getFloat();
            break;
        case Type::kDouble:
            data_.double_ = other.getDouble();
            break;
        case Type::kString: {
            data_.string_ = newObject<std::pmr::string>(alloc, other.getString());
            break;
        }
        case Type::kList: {
            data_.list_ = newObject<List>(alloc, other.getList());
            break;
        }
        case Type::kMap: {
            data_.map_ = newObject<Map>(alloc, other.getMap());
            break;
        }
        case Type::kNode: {
            data_.node_ = newObject<Node>(alloc, other.getNode());
            break;
        }
        case Type::kEdge: {
            data_.edge_ = newObject<Edge>(alloc, other.getEdge());
            break;
        }
    }
}

void Value::moveValue(Value&& other) noexcept {
    type_ = other.type_;
    other.type_ = Type::kNull;
    std::swap(data_, other.data_);
    std::memset(&other.data_, 0, sizeof(data_));
}

std::string enum2String(const Value::Type type) {
    switch (type) {
        case Value::Type::kInt8:
            return "INT8";
        case Value::Type::kInt16:
            return "INT16";
        case Value::Type::kInt32:
            return "INT32";
        case Value::Type::kInt64:
            return "INT64";
        case Value::Type::kBool:
            return "BOOL";
        case Value::Type::kFloat:
            return "FLOAT";
        case Value::Type::kDouble:
            return "DOUBLE";
        case Value::Type::kNode:
            return "NODE";
        case Value::Type::kEdge:
            return "EDGE";
        case Value::Type::kList:
            return "LIST";
        case Value::Type::kMap:
            return "MAP";
        case Value::Type::kString:
            return "STRING";
        case Value::Type::kNull:
            return "NULL";
    }

    return "UNKNOWN VALUE TYPE";
}

bool operator==(const Value& lhs, const Value& rhs) {
    auto lType = lhs.getType();
    auto rType = rhs.getType();
    if (lType != rType) {
        return false;
    }
    // TODO(jie): null is equal to null
    if (lType == Value::Type::kNull) {
        return true;
    }
    // TODO(jie): Consider the equality of compatible types
    switch (lType) {
        case Value::Type::kBool:
            return lhs.getBool() == rhs.getBool();
        case Value::Type::kInt8:
            return lhs.getInt8() == rhs.getInt8();
        case Value::Type::kInt16:
            return lhs.getInt16() == rhs.getInt16();
        case Value::Type::kInt32:
            return lhs.getInt32() == rhs.getInt32();
        case Value::Type::kInt64:
            return lhs.getInt64() == rhs.getInt64();
        case Value::Type::kFloat:
            return lhs.getFloat() == rhs.getFloat();
        case Value::Type::kDouble:
            return lhs.getDouble() == rhs.getDouble();
        case Value::Type::kString:
            return lhs.getString() == rhs.getString();
        case Value::Type::kList:
            return lhs.getList() == rhs.getList();
        case Value::Type::kMap:
            return lhs.getMap() == rhs.getMap();
        case Value::Type::kNode:
            return lhs.getNode() == rhs.getNode();
        case Value::Type::kEdge:
            return lhs.getEdge() == rhs.getEdge();
        default:
            DLOG(FATAL) << "Unknown value type " << static_cast<int>(lType);
    }

    return false;
}

bool operator!=(const Value& lhs, const Value& rhs) {
    return !(lhs == rhs);
}

bool operator<(const Value& lhs, const Value& rhs) {
    auto lType = lhs.getType();
    auto rType = rhs.getType();
    if (lType != rType) {
        LOG(ERROR) << "Cannot compare different types " << enum2String(lType) << " and "
                   << enum2String(rType);
        return false;
    }
    // TODO(Aiee): null is equal to null
    if (lType == Value::Type::kNull) {
        return true;
    }
    // TODO(Aiee): Consider the equality of compatible types
    switch (lType) {
        case Value::Type::kBool:
            return lhs.getBool() < rhs.getBool();
        case Value::Type::kInt8:
            return lhs.getInt8() < rhs.getInt8();
        case Value::Type::kInt16:
            return lhs.getInt16() < rhs.getInt16();
        case Value::Type::kInt32:
            return lhs.getInt32() < rhs.getInt32();
        case Value::Type::kInt64:
            return lhs.getInt64() < rhs.getInt64();
        case Value::Type::kFloat:
            return lhs.getFloat() < rhs.getFloat();
        case Value::Type::kDouble:
            return lhs.getDouble() < rhs.getDouble();
        case Value::Type::kString:
            return lhs.getString() < rhs.getString();
        case Value::Type::kList:
            return lhs.getList() < rhs.getList();
        case Value::Type::kMap:
            return lhs.getMap() < rhs.getMap();
        case Value::Type::kNode:
            return lhs.getNode() < rhs.getNode();
        case Value::Type::kEdge:
            return lhs.getEdge() < rhs.getEdge();
        default:
            DLOG(FATAL) << "Unknown value type " << static_cast<int>(lType);
    }

    return false;
}

bool operator>(const Value& lhs, const Value& rhs) {
    return rhs < lhs;
}

bool operator<=(const Value& lhs, const Value& rhs) {
    return !(rhs < lhs);
}

bool operator>=(const Value& lhs, const Value& rhs) {
    return !(lhs < rhs);
}

Value operator-(const Value& rhs) {
    if (rhs.isNullValue()) {
        return rhs;
    }

    switch (rhs.getType()) {
        case Value::Type::kInt8: {
            int8_t rVal = rhs.getInt8();
            if (rVal == INT8_MIN) {
                // overflow
                return Value();
            }
            int8_t val = -rVal;
            return val;
        }
        case Value::Type::kInt16: {
            int16_t rVal = rhs.getInt16();
            if (rVal == INT16_MIN) {
                // overflow
                return Value();
            }
            int16_t val = -rVal;
            return val;
        }
        case Value::Type::kInt32: {
            int32_t rVal = rhs.getInt32();
            if (rVal == INT32_MIN) {
                // overflow
                return Value();
            }
            int32_t val = -rVal;
            return val;
        }
        case Value::Type::kInt64: {
            int64_t rVal = rhs.getInt64();
            if (rVal == INT64_MIN) {
                // overflow
                return Value();
            }
            int64_t val = -rVal;
            return val;
        }
        case Value::Type::kFloat: {
            auto val = -rhs.getFloat();
            return val;
        }
        case Value::Type::kDouble: {
            auto val = -rhs.getDouble();
            return val;
        }
        default: {
            // Unsupported type
            return Value();
        }
    }
}

}  // namespace nebula::client


namespace std {

std::size_t hash<nebula::client::Value>::operator()(
        const nebula::client::Value& v) const noexcept {
    switch (v.getType()) {
        case nebula::client::Value::Type::kNull: {
            return 0;
        }
        case nebula::client::Value::Type::kBool: {
            return hash<bool>()(v.getBool());
        }
        case nebula::client::Value::Type::kInt8: {
            return hash<int8_t>()(v.getInt8());
        }
        case nebula::client::Value::Type::kInt16: {
            return hash<int16_t>()(v.getInt16());
        }
        case nebula::client::Value::Type::kInt32: {
            return hash<int32_t>()(v.getInt32());
        }
        case nebula::client::Value::Type::kInt64: {
            return hash<int64_t>()(v.getInt64());
        }
        case nebula::client::Value::Type::kFloat: {
            return hash<float>()(v.getFloat());
        }
        case nebula::client::Value::Type::kDouble: {
            return hash<double>()(v.getDouble());
        }
        case nebula::client::Value::Type::kString: {
            return hash<std::string_view>()(v.getStringView());
        }
        case nebula::client::Value::Type::kList: {
            return hash<nebula::client::List>()(v.getList());
        }
        case nebula::client::Value::Type::kMap: {
            return hash<nebula::client::Map>()(v.getMap());
        }
        case nebula::client::Value::Type::kNode: {
            return hash<nebula::client::Node>()(v.getNode());
        }
        case nebula::client::Value::Type::kEdge: {
            return hash<nebula::client::Edge>()(v.getEdge());
        }
    }
    LOG(FATAL) << "Unknown value type: " << v.getType();
}

}  // namespace std
