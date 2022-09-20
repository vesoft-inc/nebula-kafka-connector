// Copyright (c) 2022 vesoft inc. All rights reserved.

#include <common/datatype/BindingTable.h>
#include <common/datatype/thriftSerialization/BindingTableOps-inl.h>
#include <folly/Random.h>
#include <folly/init/Init.h>
#include <glog/logging.h>
#include <gtest/gtest.h>
#include <thrift/lib/cpp2/protocol/Serializer.h>

#include <string>

namespace nebula {

using serializer = apache::thrift::CompactSerializer;

namespace client {

static const int seed = folly::randomNumberSeed();
using RandomT = std::mt19937;
static RandomT rng(seed);

template <class Integral1, class Integral2>
Integral2 random(Integral1 low, Integral2 up) {
    std::uniform_int_distribution<> range(low, up);
    return range(rng);
}

std::string randomString(size_t size = 15) {
    std::string str(size, ' ');
    for (size_t p = 0; p < size; p++) {
        str[p] = random('a', 'z');
    }
    return str;
}

void randomStrBindingTable(BindingTable* table,
                           size_t numRecords,
                           size_t numColumns,
                           size_t strSize = 15) {
    std::vector<std::string> colNames;
    for (size_t i = 0; i < numColumns; ++i) {
        colNames.emplace_back(std::to_string(i));
    }
    table->setColumnNames(colNames);
    for (size_t i = 0; i < numRecords; ++i) {
        auto& newRecord = table->emplace();
        newRecord.reserve(numColumns);
        for (size_t j = 0; j < numColumns; ++j) {
            // Note this tmp value's allocator.
            // nebula::Value tmp(randomString(1000), table->get_allocator());
            Value tmp(randomString(strSize), std::pmr::new_delete_resource());
            newRecord.append(tmp);
        }
    }
}

template <class ObjType>
testing::AssertionResult checkSerialization(const ObjType& obj) {
    std::string buf;
    buf.reserve(128);
    nebula::serializer::serialize(obj, &buf);
    ObjType deserializedObj;
    std::size_t s = nebula::serializer::deserialize(buf, deserializedObj);
    if (s != buf.size()) {
        return testing::AssertionFailure() << "deserialized size " << s << " != " << buf.size();
    }
    if (obj != deserializedObj) {
        return testing::AssertionFailure()
               << "The deserialized object is not equal to the original one";
    }
    return testing::AssertionSuccess();
}

// Thrift serialization
TEST(BindingTable, FieldTypeThriftSerialization) {
    auto ft = FieldType("name", ValueType::kString);
    EXPECT_TRUE(checkSerialization(ft));
}

TEST(BindingTable, BindingTableThriftSerialization) {
    BindingTable table(std::pmr::new_delete_resource());
    randomStrBindingTable(&table, 2, 4);
    EXPECT_TRUE(checkSerialization(table));
}

}  // namespace client
}  // namespace nebula

int main(int argc, char** argv) {
    testing::InitGoogleTest(&argc, argv);
    folly::init(&argc, &argv, true);
    google::SetStderrLogging(google::INFO);

    return RUN_ALL_TESTS();
}
