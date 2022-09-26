// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_RECORD_H_
#define COMMON_DATATYPE_RECORD_H_

#include <cassert>
#include <memory_resource>
#include <unordered_map>

#include "common/datatype/Value.h"

namespace nebula::client {

// FieldType describes the type of a field in a record. It contains the name of the field and
// the value type of the field.
class FieldType final {
public:
    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept {
        return fieldName_.get_allocator();
    }

    explicit FieldType(const allocator_type& alloc = allocator_type()) : fieldName_(alloc) {}
    FieldType(const std::string& name,
              const ValueType& type,
              const allocator_type& alloc = allocator_type())
            : fieldName_(name, alloc), valueType_(type) {}

    FieldType(const FieldType& other, const allocator_type& alloc)
            : fieldName_(other.fieldName_, alloc), valueType_(other.valueType_) {}

    FieldType(FieldType&& other, const allocator_type& alloc) noexcept
            : fieldName_(std::move(other.fieldName_), alloc), valueType_(other.valueType_) {}

    const std::pmr::string& getFieldName() const {
        return fieldName_;
    }

    void setFieldName(const std::string& name) {
        fieldName_ = name;
    }

    ValueType getType() const {
        return valueType_;
    }

private:
    // Serialization using fbthrift
    friend class apache::thrift::Cpp2Ops<FieldType, void>;

    std::pmr::string fieldName_;
    // TODO(jie): Think again about whether it's appropriate to use kNull here.
    ValueType valueType_{ValueType::kNull};
};

bool operator==(const FieldType& lhs, const FieldType& rhs);
bool operator!=(const FieldType& lhs, const FieldType& rhs);

// RecordType describes the set of fields of a record in terms of their field type.
class RecordType final {
public:
    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept {
        return fieldTypes_.get_allocator();
    }

    explicit RecordType(const allocator_type& alloc = allocator_type())
            : fieldTypes_(alloc), fieldNameIndexMap_(alloc) {}

    explicit RecordType(const std::vector<FieldType>& fieldTypes,
                        const allocator_type& alloc = allocator_type())
            : fieldTypes_(alloc), fieldNameIndexMap_(alloc) {
        fieldTypes_.assign(fieldTypes.begin(), fieldTypes.end());
        buildFieldNameIndexMap();
    }

    explicit RecordType(const std::vector<std::string>& colNames,
                        const allocator_type& alloc = allocator_type())
            : fieldTypes_(alloc), fieldNameIndexMap_(alloc) {
        for (const auto& colName : colNames) {
            // TODO(jie): default type should be kAny
            // TODO(jie): Test whether the alloc of RecordType is passed
            fieldTypes_.emplace_back(colName, ValueType::kNull);
        }
        buildFieldNameIndexMap();
    }

    RecordType(const std::pmr::vector<FieldType>& fieldTypes,
               const std::pmr::unordered_map<std::pmr::string, int32_t>& fieldNameIndexMap,
               const allocator_type& alloc = allocator_type())
            : fieldTypes_(fieldTypes, alloc), fieldNameIndexMap_(fieldNameIndexMap, alloc) {}

    RecordType(const RecordType& other, const allocator_type& alloc)
            : fieldTypes_(other.getFieldTypes(), alloc),
              fieldNameIndexMap_(other.getFieldNameIndexMap(), alloc) {}

    RecordType(RecordType&& other, const allocator_type& alloc) noexcept
            : fieldTypes_(std::move(other.fieldTypes_), alloc),
              fieldNameIndexMap_(std::move(other.fieldNameIndexMap_), alloc) {}

    const std::pmr::vector<FieldType>& getFieldTypes() const {
        return fieldTypes_;
    }

    std::pmr::vector<FieldType>& getFieldTypes() {
        return fieldTypes_;
    }

    size_t getNumFieldTypes() const {
        return fieldTypes_.size();
    }

    const std::pmr::unordered_map<std::pmr::string, int32_t>& getFieldNameIndexMap() const {
        return fieldNameIndexMap_;
    }

    size_t getFieldIndexByName(const std::string& name) const {
        auto it = fieldNameIndexMap_.find(name.c_str());
        if (it == fieldNameIndexMap_.end()) {
            return -1;
        }
        return it->second;
    }

    void clear() {
        fieldTypes_.clear();
        fieldNameIndexMap_.clear();
    }

private:
    void buildFieldNameIndexMap() {
        for (size_t i = 0; i < fieldTypes_.size(); i++) {
            fieldNameIndexMap_.emplace(fieldTypes_[i].getFieldName(), i);
        }
    }

private:
    // Serialization using fbthrift
    friend class apache::thrift::Cpp2Ops<RecordType, void>;

    std::pmr::vector<FieldType> fieldTypes_;
    std::pmr::unordered_map<std::pmr::string, int32_t> fieldNameIndexMap_;
};

bool operator==(const RecordType& lhs, const RecordType& rhs);
bool operator!=(const RecordType& lhs, const RecordType& rhs);

class RawRecord final {
public:
    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept {
        return values_.get_allocator();
    }

    explicit RawRecord(const allocator_type& alloc = allocator_type()) : values_(alloc) {}
    explicit RawRecord(const std::pmr::vector<Value>& values,
                       const allocator_type& alloc = allocator_type())
            : values_(values.begin(), values.end(), alloc) {}

    RawRecord(const RawRecord& other, const allocator_type& alloc)
            : values_(other.values_, alloc) {}

    RawRecord(RawRecord&& other, const allocator_type& alloc) noexcept
            : values_(std::move(other.values_), alloc) {}

    // Append the value to the end of the record.
    void append(const Value& value) {
        values_.emplace_back(value);
    }
    void append(Value&& value) {
        values_.emplace_back(std::move(value));
    }

    // Get the field value at the given index with bounding check.
    const Value& at(size_t fieldIndex) const {
        // CHECK_LT(fieldIndex, values_.size());
        assert(fieldIndex == values_.size());
        return values_[fieldIndex];
    }

    // Get the field value at the given index.
    const Value& operator[](size_t fieldIndex) const {
        return values_[fieldIndex];
    }

    const std::pmr::vector<Value>& getValues() const {
        return values_;
    }

    size_t getNumValues() const {
        return values_.size();
    }

    // Reserve capacity for the number of fields.
    void reserve(size_t size) {
        values_.reserve(size);
    }

private:
    // Serialization using fbthrift
    friend class apache::thrift::Cpp2Ops<RawRecord, void>;

    // TODO(jie): Replace `values_` with:
    // Value* values_;
    // size_t size_;   // The size_ may also be omitted.
    std::pmr::vector<Value> values_;
};

bool operator==(const RawRecord& lhs, const RawRecord& rhs);
bool compareWithoutDynamicId(const RawRecord& lhs, const RawRecord& rhs);
bool operator!=(const RawRecord& lhs, const RawRecord& rhs);
bool operator<(const RawRecord& lhs, const RawRecord& rhs);
std::ostream& operator<<(std::ostream& os, const RawRecord& record);

// Record is not used currently.
class Record {};

}  // namespace nebula::client

#endif  // COMMON_DATATYPE_RECORD_H_
