// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_DATATYPE_COMMON_CPP2OPS_H_
#define COMMON_DATATYPE_COMMON_CPP2OPS_H_

#include "common/thrift/ThriftCpp2OpsHelper.h"
#include "common/utils/utils.h"

namespace nebula::client {
struct Service;
class Value;
class Map;
class List;
class Node;
class Edge;
class FieldType;
class BindingTable;
class RecordType;
class RawRecord;
class BindingTable;
struct Duration;
struct LocalTime;
struct LocalDatetime;
struct Date;

struct GQLStatus;
struct ExecutionOutcome;
struct ExecutionResponse;
struct AuthResponse;
struct AuthReq;

}  // namespace nebula::client

namespace apache::thrift {

SPECIALIZE_CPP2OPS(nebula::client::Value);
SPECIALIZE_CPP2OPS(nebula::client::Service);
SPECIALIZE_CPP2OPS(nebula::client::Map);
SPECIALIZE_CPP2OPS(nebula::client::List);
SPECIALIZE_CPP2OPS(nebula::client::Node);
SPECIALIZE_CPP2OPS(nebula::client::Edge);
SPECIALIZE_CPP2OPS(nebula::client::FieldType);
SPECIALIZE_CPP2OPS(nebula::client::RecordType);
SPECIALIZE_CPP2OPS(nebula::client::RawRecord);
SPECIALIZE_CPP2OPS(nebula::client::BindingTable);
SPECIALIZE_CPP2OPS(nebula::client::Duration);
SPECIALIZE_CPP2OPS(nebula::client::LocalTime);
SPECIALIZE_CPP2OPS(nebula::client::Date);
SPECIALIZE_CPP2OPS(nebula::client::LocalDatetime);

SPECIALIZE_CPP2OPS(nebula::client::GQLStatus);
SPECIALIZE_CPP2OPS(nebula::client::ExecutionOutcome);
SPECIALIZE_CPP2OPS(nebula::client::ExecutionResponse);
SPECIALIZE_CPP2OPS(nebula::client::AuthResponse);
SPECIALIZE_CPP2OPS(nebula::client::AuthReq);


}  // namespace apache::thrift

#endif  // COMMON_DATATYPE_COMMON_CPP2OPS_H_
