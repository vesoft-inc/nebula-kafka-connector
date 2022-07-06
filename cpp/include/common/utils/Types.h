#ifndef COMMON_UTILS_TYPES_H_
#define COMMON_UTILS_TYPES_H_

#include <sys/types.h>

#include <cstdint>
#include <vector>

#include "common/catalog/element/type.h"

namespace nebula::client {

using ClusterID = int64_t;

/**
 * @brief
 *
 * GraphID: It is **NOT** same as spaceId. Multiple graph with different graphId could be saved
 *          in same space. Besides, cloned graph, subgraph and view which origins from a graph
 *          will be saved in the same graph space with different GraphID as well.
 * SpaceID: A space is a resource domain where all resources in the same space share the same
 *          configurations, e.g., partition numbers, replica numbers.
 * PartID: partition id when a graph is distributed to multiple shard
 * NodeTypeID/EdgeTypeID: schema id from catalog
 * InternalID: the primary key will be mapped to a unique integer id
 * EdgeRank: edge rank in big endian to support multi graph, encoding in big endian
 * Version: version number to support mvcc, encoding in big endian, and the order of version is
 *          same as lexicographical order
 */

using GraphID =
        ::nebula::client::catalog::GraphID;  // Only least three significant bytes will be used
using SpaceID = uint32_t;
using PartID = uint32_t;
using InternalID = int64_t;
using NodeTypeID = ::nebula::client::catalog::NodeTypeID;
using EdgeTypeID = ::nebula::client::catalog::EdgeTypeID;
using PropertyID = ::nebula::client::catalog::PropertyID;
using EdgeRank = int64_t;
using Version = uint64_t;
using Port = uint32_t;
using CatalogVersion = ::nebula::client::catalog::CatalogVersion;
using AuthIdentifier = uint64_t;  // sth like session id, from session.h

// raft related
using TermID = int64_t;
using LogID = int64_t;
using ClusterID = int64_t;

static_assert(sizeof(NodeTypeID) == sizeof(EdgeTypeID));

enum class NebulaKeyType : uint8_t {
    kNode = 0x01,     // For node
    kEdge = 0x02,     // For edge
    kPrimary = 0x03,  // For mapping from primary key to internal id
    kIndex = 0x04,    // For index
    kSystem = 0x05,   // For internal use
};

enum class SystemKeySubType : uint8_t {
    kSystemPart = 0x01,    // Partition id in storage
    kSystemCommit = 0x02,  // Raft commit id of each partitions
    kSystemPeers = 0x03,   // Raft peers of each partitions
};

struct PartitionList {
    std::vector<PartID> partList;
};

struct LogEntry {
    ClusterID cluster;
    std::string logStr;
};

}  // namespace nebula::client

#endif  // COMMON_UTILS_TYPES_H_
