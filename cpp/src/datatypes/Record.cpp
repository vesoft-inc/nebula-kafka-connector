// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "common/datatype/Record.h"

namespace nebula::client {

bool operator==(const FieldType& lhs, const FieldType& rhs) {
    return lhs.getFieldName() == rhs.getFieldName() && lhs.getType() == rhs.getType();
}

bool operator!=(const FieldType& lhs, const FieldType& rhs) {
    return !(lhs == rhs);
}

bool operator==(const RecordType& lhs, const RecordType& rhs) {
    // return lhs.getFieldTypes() == rhs.getFieldTypes();
    return lhs.getFieldTypes() == rhs.getFieldTypes() &&
           lhs.getFieldNameIndexMap() == rhs.getFieldNameIndexMap();
}

bool operator!=(const RecordType& lhs, const RecordType& rhs) {
    return !(lhs == rhs);
}

bool operator==(const RawRecord& lhs, const RawRecord& rhs) {
    return lhs.getValues() == rhs.getValues();
}

bool operator!=(const RawRecord& lhs, const RawRecord& rhs) {
    return !(lhs == rhs);
}

bool operator<(const RawRecord& lhs, const RawRecord& rhs) {
    return lhs.getValues() < rhs.getValues();
}

std::ostream& operator<<(std::ostream& os, const RawRecord& record) {
    os << "[";
    for (auto& value : record.getValues()) {
        os << value << ", ";
    }
    os << "]";
    return os;
}

}  // namespace nebula::client
