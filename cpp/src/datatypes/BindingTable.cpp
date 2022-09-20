// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "common/datatype/BindingTable.h"

#include <glog/logging.h>

namespace nebula::client {

void BindingTableDescriptor::setColumnNames(const std::vector<std::string>& colNames) {
    recordType_.getFieldTypes().resize(colNames.size());
    for (std::size_t i = 0; i < colNames.size(); ++i) {
        recordType_.getFieldTypes()[i].setFieldName(colNames[i]);
    }
}

bool operator==(const BindingTableDescriptor& lhs, const BindingTableDescriptor& rhs) {
    return lhs.getRecordType() == rhs.getRecordType() && lhs.isOrdered() == rhs.isOrdered() &&
           lhs.isDuplicateFree() == rhs.isDuplicateFree();
}

bool operator!=(const BindingTableDescriptor& lhs, const BindingTableDescriptor& rhs) {
    return !(lhs == rhs);
}

bool operator==(const BindingTable& lhs, const BindingTable& rhs) {
    return lhs.getDescriptor() == rhs.getDescriptor() && lhs.getRecords() == rhs.getRecords();
}

bool operator!=(const BindingTable& lhs, const BindingTable& rhs) {
    return !(lhs == rhs);
}

}  // namespace nebula::client
