// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "common/datatype/BindingTable.h"

namespace nebula::client {

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
