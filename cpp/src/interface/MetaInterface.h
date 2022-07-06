// Copyright (c) 2022 vesoft inc. All rights reserved.
#ifndef INTERFACE_METAINTERFACE_H
#define INTERFACE_METAINTERFACE_H

#include <folly/io/IOBuf.h>

#include <cstdint>
#include <string>

#include "common/base/ErrorCode.h"
#include "common/nrpc/BufferReaderWriter.h"
#include "common/nrpc/CommonDefine.h"
#include "common/nrpc/Context.h"
#include "common/service/Service.h"
#include "common/utils/Types.h"

namespace nebula::client {
namespace meta {

enum class HostRole : uint8_t {
    UNKNOWN = 0,
    GRAPH = 1,
    META = 2,
    STORAGE = 3,
};

enum class HostState : uint8_t {
    UNKNOWN = 0,
    NORMAL = 1,
};


struct LeaderInfo {
    PartID partId;
    TermID termId;
};

class MetaRequestName final {
public:
    MetaRequestName() = delete;
    // space related
    static constexpr const char *createSpace = "createSpace";
    static constexpr const char *getSpace = "getSpace";
    static constexpr const char *dropSpace = "dropSpace";
    static constexpr const char *listSpaces = "listSpaces";
    static constexpr const char *partAllocation = "partAllocation";

    // catalog related
    static constexpr const char *catalogHeartbeat = "catalogHeartbeat";
    static constexpr const char *ddl = "ddl";
};

struct RequestHeader {
    std::string type;
};

struct ResponseHeader {
    ErrorCode code;
    std::string msg;
    Service host;
};

struct NullResponse {
    NullResponse() = default;
};

// TODO(@SuperYoko): currently, we only on version
struct CatalogHeartbeatReq {
    uint64_t currentVersion;
};

struct CatalogHeartbeatResp {
    CatalogHeartbeatResp() = default;
    bool result;
    int64_t heartbeatInterval;
    int64_t maxheartbeatFailedCount;
    std::string catalogJsonStr;
};

struct DDLReq {
    std::string operation;
};

// empty allowed.
struct DDLResp {
    DDLResp() = default;
    bool result;
};

/**
 * @brief SpaceDesc used to pass space information to meta server
 * @Description: stub now.
 */
struct SpaceDesc {
    std::string space_name;
    int32_t partition_num;
    int32_t replica_factor;
    std::string charset;
    std::string collate;
};

// Stubed now.
struct CreateSpaceRequest {
    SpaceDesc spaceDesc;
    bool ifNotExists;
};

struct CreateSpaceResponse {
    CreateSpaceResponse() = default;
    SpaceID spaceId;
};

enum class RaftRole {
    FOLLOWER = 1,
    LEADER = 2,
    LEANER = 3,
};

struct Host {
    Service address;
    bool isLeaner;
    Host() = default;

    Host(Service &&addr, bool leaner) {
        address = std::move(addr);
        isLeaner = leaner;
    }

    explicit Host(Service &&addr) {
        address = std::move(addr);
        isLeaner = false;
    }

    bool operator==(const Host &other) const {
        return isLeaner == other.isLeaner && address == other.address;
    }

    meta::Host &operator=(const Service &rhs) {
        address = rhs;
        return *this;
    }

    meta::Host &operator=(Service &&rhs) {
        address = std::move(rhs);
        return *this;
    }
};

struct PartInfo {
    PartID partId;
    std::vector<Host> hosts;

    PartInfo() = default;

    bool operator==(const PartInfo &other) const {
        return partId == other.partId && hosts == other.hosts;
    }
};

struct PartRequest {
    Service host;
};

struct PartResponse {
    PartResponse() = default;
    std::map<SpaceID, std::vector<PartInfo>> spaceParts;
};

struct GetSpaceRequest {
    std::string spaceName;
};

struct GetSpaceResponse {
    GetSpaceResponse() = default;
    SpaceID id;
    SpaceDesc desc;
};


struct ListSpacesRequest {};


struct ListSpacesResponse {
    ListSpacesResponse() = default;
    std::map<SpaceID, SpaceDesc> spaces;
};


using PartsMap = std::unordered_map<SpaceID, std::unordered_map<PartID, PartInfo>>;

}  // namespace meta


SERIALIZE_EACH_MEMBER(meta::RequestHeader, type)
SERIALIZE_EACH_MEMBER(meta::ResponseHeader, code, msg, host)
SERIALIZE_EACH_MEMBER(meta::CatalogHeartbeatReq, currentVersion)
SERIALIZE_EACH_MEMBER(meta::CatalogHeartbeatResp,
                      result,
                      heartbeatInterval,
                      maxheartbeatFailedCount,
                      catalogJsonStr)
SERIALIZE_EACH_MEMBER(meta::DDLReq, operation)
SERIALIZE_EACH_MEMBER(meta::DDLResp, result)
SERIALIZE_EACH_MEMBER(
        meta::SpaceDesc, space_name, partition_num, replica_factor, charset, collate)
SERIALIZE_EACH_MEMBER(meta::CreateSpaceRequest, spaceDesc, ifNotExists)
SERIALIZE_CONST_LENGTH_MEMBER(meta::CreateSpaceResponse, spaceId)
SERIALIZE_EACH_MEMBER(meta::Host, address, isLeaner)
SERIALIZE_EACH_MEMBER(meta::GetSpaceRequest, spaceName)
SERIALIZE_EACH_MEMBER(meta::GetSpaceResponse, id, desc)
SERIALIZE_EACH_MEMBER(meta::ListSpacesResponse, spaces)
SERIALIZE_EACH_MEMBER(meta::PartInfo, partId, hosts)
SERIALIZE_EACH_MEMBER(meta::PartRequest, host)
SERIALIZE_EACH_MEMBER(meta::PartResponse, spaceParts)

namespace meta {

template <typename Request>
nrpc::ClientContext::BufferPtr ConstructRequest(const Request &req,
                                                const std::string &requestType) {
    nrpc::BufferReaderWriter<RequestHeader> headerWriter;
    nrpc::BufferReaderWriter<Request> requestWriter;

    RequestHeader header{requestType};
    auto size = headerWriter.encodedSize(header) + requestWriter.encodedSize(req);
    auto buf = folly::IOBuf::create(size);

    headerWriter.write(buf.get(), header);
    requestWriter.write(buf.get(), req);
    return buf;
}

template <typename Response>
nrpc::ServerContext::BufferPtr ConstructResponse(const Response &resp,
                                                 ErrorCode code,
                                                 const std::string &msg,
                                                 bool skipResp = false) {
    nrpc::BufferReaderWriter<ResponseHeader> headerWriter;
    nrpc::BufferReaderWriter<Response> responseWriter;
    ResponseHeader header{code, msg, Service()};

    if (skipResp) {
        auto size = headerWriter.encodedSize(header);
        auto buf = folly::IOBuf::create(size);
        headerWriter.write(buf.get(), header);
        return buf;
    }

    auto size = headerWriter.encodedSize(header) + responseWriter.encodedSize(resp);
    auto buf = folly::IOBuf::create(size);
    headerWriter.write(buf.get(), header);
    responseWriter.write(buf.get(), resp);
    return buf;
}
}  // namespace meta

}  // namespace nebula::client

#endif  // INTERFACE_METAINTERFACE_H
