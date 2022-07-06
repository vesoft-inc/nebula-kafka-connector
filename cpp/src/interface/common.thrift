// Copyright (c) 2022 vesoft inc. All rights reserved.

namespace cpp nebula
namespace java com.vesoft.nebula
namespace go nebula
namespace js nebula
namespace csharp nebula
namespace py nebula3.common

cpp_include "common/utils/Types.h"
cpp_include "common/datatype/thriftSerialization/ValueOps-inl.h"
cpp_include "common/datatype/thriftSerialization/ListOps-inl.h"
cpp_include "common/datatype/thriftSerialization/MapOps-inl.h"
cpp_include "common/datatype/thriftSerialization/NodeOps-inl.h"
cpp_include "common/datatype/thriftSerialization/EdgeOps-inl.h"
cpp_include "common/datatype/thriftSerialization/RecordOps-inl.h"
cpp_include "common/datatype/thriftSerialization/BindingTableOps-inl.h"

//  Note: In order to support multiple languages, all strings
//        have to be defined as **binary** in the thrift file

// Version of the thrift interface.
// This is used to detect incompatible changes in the thrift interface.
const binary (cpp.type = "char const *") version = "5.0.0"

typedef i32 (cpp.type = "nebula::client::NodeTypeID") NodeTypeID
typedef i32 (cpp.type = "nebula::client::EdgeTypeID") EdgeTypeID
typedef i64 (cpp.type = "nebula::client::InternalID") InternalID
typedef i64 (cpp.type = "nebula::client::EdgeRank") EdgeRank

// The type to hold any supported values during the query
union Value {
    1: bool boolVal;
    2: byte int8Val;
    3: i16  int16Val;
    4: i32  int32Val;
    5: i64  int64Val;
    // Thrift does not support float type for generating js:node code, so we use double here instead
    6: double floatVal;
    7: double doubleVal;
    8: binary stringVal;
    9: NList listVal;
    10: NMap mapVal;
    11: Node nodeVal;
    12: Edge edgeVal;
} (cpp.type = "nebula::client::Value")

struct NList {
    1: list<Value> (cpp.template = "std::pmr::vector") values;
} (cpp.type = "nebula::client::List")

// The key of map can be non-string type
struct NMap {
    // TODO(Aiee) Temprarily comment this so the generated client code can compile.
    // Still need to discuss if we should use string or value for the key.
    # 1: map<Value, Value> (cpp.template = "std::pmr::unordered_map") values;
} (cpp.type = "nebula::client::Map")

struct Node {
    1: InternalID nodeID;
    2: NodeTypeID nodeTypeID;
    3: map<binary, Value> (cpp.template = "std::pmr::unordered_map") properties;
} (cpp.type = "nebula::client::Node")

struct Edge {
    1: InternalID srcID;
    2: InternalID dstID;
    3: EdgeTypeID edgeTypeID;
    4: EdgeRank rank;
    5: map<binary, Value> (cpp.template = "std::pmr::unordered_map") properties;
} (cpp.type = "nebula::client::Edge")

enum ValueType {
    kNull     = 0,        
    kBool     = 1,   
    kInt8     = 2,       
    kInt16    = 3,       
    kInt32    = 4,        
    kInt64    = 5,      
    kFloat    = 6,      
    kDouble   = 7,       
    kString   = 8,       
    kList     = 9,     
    kMap      = 10,    
    kNode     = 11,    
    kEdge     = 12,   
} (cpp.enum_strict, cpp.type = "nebula::client::ValueType")

struct FieldType {
    1: binary filedName;
    2: ValueType (cpp.type = "nebula::client::ValueType") valueType;
} (cpp.type = "nebula::client::FieldType")

struct RecordType {
    1: list<FieldType> (cpp.template = "std::pmr::vector") fieldType;
    2: map<binary, i32> (cpp.template = "std::pmr::unordered_map") fieldNameIndexMap;
} (cpp.type = "nebula::client::RecordType")

struct RawRecord {
    1: list<Value> (cpp.template = "std::pmr::vector") values;
} (cpp.type = "nebula::client::RawRecord")

// TODO(Aiee) Unimplemented
# struct Record {}

struct BindingTable {
    1: list<RawRecord> (cpp.template = "std::pmr::deque") records;
} (cpp.type = "nebula::client::BindingTable")


/*
 * ErrorCode for graphd, metad, storaged,raftd
 * -1xxx for graphd
 * -2xxx for metad
 * -3xxx for storaged
 */
enum ErrorCode {
    // for common code
    SUCCEEDED                         = 0,

    E_DISCONNECTED                    = -1,        // RPC Failure
    E_FAIL_TO_CONNECT                 = -2,
    E_RPC_FAILURE                     = -3,
    E_LEADER_CHANGED                  = -4,


    // only unify metad and storaged error code
    E_SPACE_NOT_FOUND                 = -5,
    E_TAG_NOT_FOUND                   = -6,
    E_EDGE_NOT_FOUND                  = -7,
    E_INDEX_NOT_FOUND                 = -8,
    E_EDGE_PROP_NOT_FOUND             = -9,
    E_TAG_PROP_NOT_FOUND              = -10,
    E_ROLE_NOT_FOUND                  = -11,
    E_CONFIG_NOT_FOUND                = -12,
    E_MACHINE_NOT_FOUND               = -13,
    E_ZONE_NOT_FOUND                  = -14,
    E_LISTENER_NOT_FOUND              = -15,
    E_PART_NOT_FOUND                  = -16,
    E_KEY_NOT_FOUND                   = -17,
    E_USER_NOT_FOUND                  = -18,
    E_STATS_NOT_FOUND                 = -19,
    E_SERVICE_NOT_FOUND               = -20,
    E_DRAINER_NOT_FOUND               = -21,
    E_DRAINER_CLIENT_NOT_FOUND        = -22,

    // backup failed
    E_BACKUP_FAILED                   = -24,
    E_BACKUP_EMPTY_TABLE              = -25,
    E_BACKUP_TABLE_FAILED             = -26,
    E_PARTIAL_RESULT                  = -27,
    E_REBUILD_INDEX_FAILED            = -28,
    E_INVALID_PASSWORD                = -29,
    E_FAILED_GET_ABS_PATH             = -30,


    // 1xxx for graphd
    E_BAD_USERNAME_PASSWORD           = -1001,     // Authentication error
    E_SESSION_INVALID                 = -1002,     // Execution errors
    E_SESSION_TIMEOUT                 = -1003,
    E_SYNTAX_ERROR                    = -1004,
    E_EXECUTION_ERROR                 = -1005,
    E_STATEMENT_EMPTY                 = -1006,     // Nothing is executed When command is comment

    E_BAD_PERMISSION                  = -1008,
    E_SEMANTIC_ERROR                  = -1009,     // semantic error
    E_TOO_MANY_CONNECTIONS            = -1010,     // Exceeding the maximum number of connections
    E_PARTIAL_SUCCEEDED               = -1011,


    // 2xxx for metad
    E_NO_HOSTS                        = -2001,     // Operation Failure
    E_EXISTED                         = -2002,
    E_INVALID_HOST                    = -2003,
    E_UNSUPPORTED                     = -2004,
    E_NOT_DROP                        = -2005,
    E_BALANCER_RUNNING                = -2006,
    E_CONFIG_IMMUTABLE                = -2007,
    E_CONFLICT                        = -2008,
    E_INVALID_PARM                    = -2009,
    E_WRONGCLUSTER                    = -2010,
    E_LISTENER_CONFLICT               = -2011,
    E_ZONE_NOT_ENOUGH                 = -2012,
    E_ZONE_IS_EMPTY                   = -2013,

    E_STORE_FAILURE                   = -2021,
    E_STORE_SEGMENT_ILLEGAL           = -2022,
    E_BAD_BALANCE_PLAN                = -2023,
    E_BALANCED                        = -2024,
    E_NO_RUNNING_BALANCE_PLAN         = -2025,
    E_NO_VALID_HOST                   = -2026,
    E_CORRUPTED_BALANCE_PLAN          = -2027,
    E_NO_INVALID_BALANCE_PLAN         = -2028,
    E_NO_VALID_DRAINER                = -2029,

    // Authentication Failure
    E_IMPROPER_ROLE                   = -2030,
    E_INVALID_PARTITION_NUM           = -2031,
    E_INVALID_REPLICA_FACTOR          = -2032,
    E_INVALID_CHARSET                 = -2033,
    E_INVALID_COLLATE                 = -2034,
    E_CHARSET_COLLATE_NOT_MATCH       = -2035,

    // Admin Failure
    E_SNAPSHOT_FAILURE                = -2040,
    E_BLOCK_WRITE_FAILURE             = -2041,
    E_REBUILD_INDEX_FAILURE           = -2042,
    E_INDEX_WITH_TTL                  = -2043,
    E_ADD_JOB_FAILURE                 = -2044,
    E_STOP_JOB_FAILURE                = -2045,
    E_SAVE_JOB_FAILURE                = -2046,
    E_BALANCER_FAILURE                = -2047,
    E_JOB_NOT_FINISHED                = -2048,
    E_TASK_REPORT_OUT_DATE            = -2049,
    E_JOB_NOT_IN_SPACE                = -2050,
    E_JOB_NEED_RECOVER                = -2051,
    E_INVALID_JOB                     = -2065,

    // Backup Failure
    E_BACKUP_BUILDING_INDEX           = -2066,
    E_BACKUP_SPACE_NOT_FOUND          = -2067,

    // RESTORE Failure
    E_RESTORE_FAILURE                 = -2068,

    E_SESSION_NOT_FOUND               = -2069,

    // ListClusterInfo Failure
    E_LIST_CLUSTER_FAILURE              = -2070,
    E_LIST_CLUSTER_GET_ABS_PATH_FAILURE = -2071,
    E_LIST_CLUSTER_NO_AGENT_FAILURE     = -2072,

    E_QUERY_NOT_FOUND                 = -2073,
    E_AGENT_HB_FAILUE                 = -2074,

    E_INVALID_VARIABLE                = -2080,
    E_VARIABLE_TYPE_VALUE_MISMATCH    = -2081,

    // 3xxx for storaged
    E_CONSENSUS_ERROR                 = -3001,
    E_KEY_HAS_EXISTS                  = -3002,
    E_DATA_TYPE_MISMATCH              = -3003,
    E_INVALID_FIELD_VALUE             = -3004,
    E_INVALID_OPERATION               = -3005,
    E_NOT_NULLABLE                    = -3006,     // Not allowed to be null
    // The field neither can be NULL, nor has a default value
    E_FIELD_UNSET                     = -3007,
    // Value exceeds the range of type
    E_OUT_OF_RANGE                    = -3008,
    E_DATA_CONFLICT_ERROR             = -3010,     // data conflict, for index write without toss.

    E_WRITE_STALLED                   = -3011,

    // meta failures
    E_IMPROPER_DATA_TYPE              = -3021,
    E_INVALID_SPACEVIDLEN             = -3022,

    // Invalid request
    E_INVALID_FILTER                  = -3031,
    E_INVALID_UPDATER                 = -3032,
    E_INVALID_STORE                   = -3033,
    E_INVALID_PEER                    = -3034,
    E_RETRY_EXHAUSTED                 = -3035,
    E_TRANSFER_LEADER_FAILED          = -3036,
    E_INVALID_STAT_TYPE               = -3037,
    E_INVALID_VID                     = -3038,
    E_NO_TRANSFORMED                  = -3039,

    // meta client failed
    E_LOAD_META_FAILED                = -3040,

    // checkpoint failed
    E_FAILED_TO_CHECKPOINT            = -3041,
    E_CHECKPOINT_BLOCKED              = -3042,

    // Filter out
    E_FILTER_OUT                      = -3043,
    E_INVALID_DATA                    = -3044,

    E_MUTATE_EDGE_CONFLICT            = -3045,
    E_MUTATE_TAG_CONFLICT             = -3046,

    // transaction
    E_OUTDATED_LOCK                   = -3047,

    // task manager failed
    E_INVALID_TASK_PARA               = -3051,
    E_USER_CANCEL                     = -3052,
    E_TASK_EXECUTION_FAILED           = -3053,

    E_PLAN_IS_KILLED                  = -3060,
    // toss
    E_NO_TERM                         = -3070,
    E_OUTDATED_TERM                   = -3071,
    E_OUTDATED_EDGE                   = -3072,
    E_WRITE_WRITE_CONFLICT            = -3073,

    E_CLIENT_SERVER_INCOMPATIBLE      = -3061,
    // get id failed
    E_ID_FAILED                       = -3062,

    // 35xx for storaged raft
    E_RAFT_UNKNOWN_PART               = -3500,
    // Raft consensus errors
    E_RAFT_LOG_GAP                    = -3501,
    E_RAFT_LOG_STALE                  = -3502,
    E_RAFT_TERM_OUT_OF_DATE           = -3503,
    E_RAFT_UNKNOWN_APPEND_LOG         = -3504,
    // Raft state errors
    E_RAFT_WAITING_SNAPSHOT           = -3511,
    E_RAFT_SENDING_SNAPSHOT           = -3512,
    E_RAFT_INVALID_PEER               = -3513,
    E_RAFT_NOT_READY                  = -3514,
    E_RAFT_STOPPED                    = -3515,
    E_RAFT_BAD_ROLE                   = -3516,
    // Local errors
    E_RAFT_WAL_FAIL                   = -3521,
    E_RAFT_HOST_STOPPED               = -3522,
    E_RAFT_TOO_MANY_REQUESTS          = -3523,
    E_RAFT_PERSIST_SNAPSHOT_FAILED    = -3524,
    E_RAFT_RPC_EXCEPTION              = -3525,
    E_RAFT_NO_WAL_FOUND               = -3526,
    E_RAFT_HOST_PAUSED                = -3527,
    E_RAFT_WRITE_BLOCKED              = -3528,
    E_RAFT_BUFFER_OVERFLOW            = -3529,
    E_RAFT_ATOMIC_OP_FAILED           = -3530,
    E_LEADER_LEASE_FAILED             = -3531,
    E_RAFT_CAUGHT_UP                  = -3532,
    
    // 4xxx for drainer
    E_LOG_GAP                         = -4001,
    E_LOG_STALE                       = -4002,
    E_INVALID_DRAINER_STORE           = -4003,
    E_SPACE_MISMATCH                  = -4004,
    E_PART_MISMATCH                   = -4005,
    E_DATA_CONFLICT                   = -4006,
    E_REQ_CONFLICT                    = -4007,
    E_DATA_ILLEGAL                    = -4008,

    // 5xxx for cache
    E_CACHE_CONFIG_ERROR              = -5001,
    E_NOT_ENOUGH_SPACE                = -5002,
    E_CACHE_MISS                      = -5003,
    E_POOL_NOT_FOUND                  = -5004,
    E_CACHE_WRITE_FAILURE             = -5005,

    // 7xxx for nebula enterprise
    // license related
    E_NODE_NUMBER_EXCEED_LIMIT        = -7001,
    E_PARSING_LICENSE_FAILURE         = -7002,

    E_UNKNOWN                         = -8000,
} (cpp.enum_strict)
