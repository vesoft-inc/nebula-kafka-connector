// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef INTERFACE_STORAGEINTERFACE_H_
#define INTERFACE_STORAGEINTERFACE_H_

#include <glog/logging.h>

#include <memory_resource>

#include "common/datatype/BindingTable.h"
#include "common/datatype/Edge.h"
#include "common/datatype/List.h"
#include "common/datatype/Node.h"
#include "common/nrpc/BufferReaderWriter.h"
#include "common/nrpc/CommonDefine.h"
#include "common/nrpc/Context.h"
#include "graph/plan/Query.h"
#include "kvstore/KVStore.h"

namespace nebula::client {
namespace storage {

using PartID = int32_t;

enum class StorageReqType : uint8_t {
    kUnknown = 0,
    kQueryPlan = 1,
    kAddNode = 2,
    kAddEdge = 3,
};

struct StorageReqHeader {
    StorageReqType type{StorageReqType::kUnknown};
};

struct StorageRespHeader {
    ErrorCode code{ErrorCode::SUCCESS};
    std::string msg;
    Service leader;
};


struct AddNodesRequest : StorageReqHeader {
    AuthIdentifier authId;
    GraphID graphId;
    NodeTypeID nodeTypeId;
    std::unordered_map<PartID, std::vector<nebula::client::Node>> nodeValues;
    // TODO Whether use DAG physical storage plan
    // Exec graphPlan;
};

struct AddEdgesRequest : StorageReqHeader {
    AuthIdentifier authId;
    GraphID graphId;
    EdgeTypeID edgeTypeId;
    std::unordered_map<PartID, std::vector<Edge>> edgeValues;
    // TODO Whether use DAG physical storage plan
    // Exec graphPlan;
};

struct CommonResponse : StorageRespHeader {
    std::unordered_map<PartID, StorageRespHeader> failedParts;
};


struct QueryPlanRequest : StorageReqHeader {
    AuthIdentifier authId;
    GraphID graphId;
    std::vector<List> inputs;

    // TODO(spw): may extend to multi-roots in the future
    // DAG physical storage plan
    std::shared_ptr<graph::StoragePlanNode> graphPlan;
};

struct QueryPlanResponse : StorageRespHeader {
    using allocator_type = std::pmr::polymorphic_allocator<std::byte>;
    allocator_type get_allocator() const noexcept {
        return {mr};
    }

    explicit QueryPlanResponse(allocator_type alloc = allocator_type()) {
        mr = alloc.resource();
    }
    nebula::client::BindingTable *result{nullptr};
    std::pmr::memory_resource *mr{nullptr};
};

}  // namespace storage

SERIALIZE_EACH_MEMBER(storage::StorageReqHeader, type)
SERIALIZE_EACH_MEMBER(storage::StorageRespHeader, code, leader)

SERIALIZE_EACH_MEMBER_BASE(storage::AddNodesRequest,
                           storage::StorageReqHeader,
                           authId,
                           graphId,
                           nodeTypeId,
                           nodeValues)

SERIALIZE_EACH_MEMBER_BASE(storage::AddEdgesRequest,
                           storage::StorageReqHeader,
                           authId,
                           graphId,
                           edgeTypeId,
                           edgeValues)

SERIALIZE_EACH_MEMBER_BASE(storage::CommonResponse, storage::StorageRespHeader, failedParts)


namespace nrpc {

template <>
struct BufferReaderWriter<storage::QueryPlanRequest> {
    static void write(folly::IOBuf *buf, const storage::QueryPlanRequest &msg) {
        BufferReaderWriter<storage::StorageReqHeader>::write(buf, msg);
        BufferReaderWriter<AuthIdentifier>::write(buf, msg.authId);
        BufferReaderWriter<GraphID>::write(buf, msg.graphId);
        BufferReaderWriter<std::vector<nebula::client::List>>::write(buf, msg.inputs);
        BufferReaderWriter<graph::PlanNodeType>::write(buf, msg.graphPlan->type());
        switch (msg.graphPlan->type()) {
            case graph::PlanNodeType::kNodeScanById: {
                BufferReaderWriter<graph::NodeScanById>::write(
                        buf, static_cast<const graph::NodeScanById &>(*msg.graphPlan));
                break;
            }
            case graph::PlanNodeType::kEdgeScanById: {
                BufferReaderWriter<graph::EdgeScanById>::write(
                        buf, static_cast<const graph::EdgeScanById &>(*msg.graphPlan));
                break;
            }
            case graph::PlanNodeType::kNEJoin: {
                BufferReaderWriter<graph::NEJoin>::write(
                        buf, static_cast<const graph::NEJoin &>(*msg.graphPlan));
                break;
            }
            default: {
                break;
            }
        }
    }

    static Status read(folly::IOBuf *buf, storage::QueryPlanRequest *msgPtr) {
        NG_RETURN_IF_ERROR(BufferReaderWriter<storage::StorageReqHeader>::read(buf, msgPtr));
        NG_RETURN_IF_ERROR(BufferReaderWriter<AuthIdentifier>::read(buf, &msgPtr->authId));
        NG_RETURN_IF_ERROR(BufferReaderWriter<GraphID>::read(buf, &msgPtr->graphId));
        NG_RETURN_IF_ERROR(BufferReaderWriter<std::vector<nebula::client::List>>::read(
                buf, &msgPtr->inputs));
        graph::PlanNodeType type = graph::PlanNodeType::kUnknown;
        NG_RETURN_IF_ERROR(BufferReaderWriter<graph::PlanNodeType>::read(buf, &type));
        switch (type) {
            case graph::PlanNodeType::kNodeScanById: {
                auto root = std::make_shared<graph::NodeScanById>();
                NG_RETURN_IF_ERROR(
                        BufferReaderWriter<graph::NodeScanById>::read(buf, root.get()));
                msgPtr->graphPlan = root;
                break;
            }
            case graph::PlanNodeType::kEdgeScanById: {
                auto root = std::make_shared<graph::EdgeScanById>();
                NG_RETURN_IF_ERROR(
                        BufferReaderWriter<graph::EdgeScanById>::read(buf, root.get()));
                msgPtr->graphPlan = root;
                break;
            }
            case graph::PlanNodeType::kNEJoin: {
                auto root = std::make_shared<graph::NEJoin>();
                NG_RETURN_IF_ERROR(BufferReaderWriter<graph::NEJoin>::read(buf, root.get()));
                msgPtr->graphPlan = root;
                break;
            }
            default: {
                break;
            }
        }
        return Status::OK();
    }

    static size_t encodedSize(const storage::QueryPlanRequest &msg) {
        size_t sz = BufferReaderWriter<storage::StorageReqHeader>::encodedSize(msg);
        sz += BufferReaderWriter<AuthIdentifier>::encodedSize(msg.authId);
        sz += BufferReaderWriter<GraphID>::encodedSize(msg.graphId);
        sz += BufferReaderWriter<std::vector<nebula::client::List>>::encodedSize(msg.inputs);
        sz += BufferReaderWriter<graph::PlanNodeType>::encodedSize(msg.graphPlan->type());
        switch (msg.graphPlan->type()) {
            case graph::PlanNodeType::kNodeScanById: {
                sz += BufferReaderWriter<graph::NodeScanById>::encodedSize(
                        static_cast<const graph::NodeScanById &>(*msg.graphPlan));
                break;
            }
            case graph::PlanNodeType::kEdgeScanById: {
                sz += BufferReaderWriter<graph::EdgeScanById>::encodedSize(
                        static_cast<const graph::EdgeScanById &>(*msg.graphPlan));
                break;
            }
            case graph::PlanNodeType::kNEJoin: {
                sz += BufferReaderWriter<graph::NEJoin>::encodedSize(
                        static_cast<const graph::NEJoin &>(*msg.graphPlan));
                break;
            }
            default: {
                break;
            }
        }
        return sz;
    }
};

template <>
struct BufferReaderWriter<storage::QueryPlanResponse> {
    static void write(folly::IOBuf *buf, const storage::QueryPlanResponse &msg) {
        BufferReaderWriter<storage::StorageRespHeader>::write(buf, msg);
        BufferReaderWriter<bool>::write(buf, msg.result != nullptr);
        if (msg.result) {
            BufferReaderWriter<nebula::client::BindingTable>::write(buf, *msg.result);
        }
    }

    static Status read(folly::IOBuf *buf, storage::QueryPlanResponse *msgPtr) {
        NG_RETURN_IF_ERROR(BufferReaderWriter<storage::StorageRespHeader>::read(buf, msgPtr));
        bool b = false;
        BufferReaderWriter<bool>::read(buf, &b);
        if (b) {
            auto *mr = msgPtr->mr;
            msgPtr->result = BindingTable::make(CHECK_NOTNULL(mr));
            BufferReaderWriter<BindingTable>::read(buf, msgPtr->result);
        }
        return Status::OK();
    }

    static size_t encodedSize(const storage::QueryPlanResponse &msg) {
        size_t sz =
                BufferReaderWriter<storage::StorageRespHeader>::encodedSize(msg) + sizeof(bool);
        if (msg.result) {
            sz += BufferReaderWriter<nebula::client::BindingTable>::encodedSize(*msg.result);
        }
        return sz;
    }
};

}  // namespace nrpc

}  // namespace nebula::client

#endif
