// Copyright (c) 2022 vesoft inc. All rights reserved.

#pragma once

#include "nebula/client/Connection.h"

namespace nebula::client {

class ConnectionPool;

class Session {
public:
    using ExecuteCallback = std::function<void(ExecutionResponse &&)>;
    using ExecuteJsonCallback = std::function<void(std::string &&)>;

    Session() = default;
    Session(int64_t sessionId,
            Connection &&conn,
            ConnectionPool *pool,
            const AuthReq &req,
            const std::string &timezoneName,
            int32_t offsetSecs)
            : sessionId_(sessionId),
              conn_(std::move(conn)),
              pool_(pool),
              req_(req),
              timezoneName_(timezoneName),
              offsetSecs_(offsetSecs) {}
    Session(Session &&session)
            : sessionId_(session.sessionId_),
              conn_(std::move(session.conn_)),
              pool_(session.pool_),
              req_(std::move(session.req_)),
              timezoneName_(std::move(session.timezoneName_)),
              offsetSecs_(session.offsetSecs_) {
        session.sessionId_ = -1;
        session.pool_ = nullptr;
        session.offsetSecs_ = 0;
    }
    ~Session() {
        release();
    }

    ExecutionResponse execute(const std::string &stmt);

    void asyncExecute(const std::string &stmt, ExecuteCallback cb);

    // ExecutionResponse executeWithParameter(const std::string &stmt,
    // const std::unordered_map<std::string, Value> &parameters);

    // void asyncExecuteWithParameter(const std::string &stmt,
    // const std::unordered_map<std::string, Value> &parameters,
    // ExecuteCallback cb);

    // std::string executeJson(const std::string &stmt);

    // void asyncExecuteJson(const std::string &stmt, ExecuteJsonCallback cb);

    // std::string executeJsonWithParameter(const std::string &stmt,
    // const std::unordered_map<std::string, Value> &parameters);

    // void asyncExecuteJsonWithParameter(const std::string &stmt,
    // const std::unordered_map<std::string, Value> &parameters,
    // ExecuteJsonCallback cb);

    bool ping();

    GQLStatus retryConnect();

    void release();

    bool valid() const {
        return sessionId_ > 0;
    }

    const std::string &timeZoneName() const {
        return timezoneName_;
    }

    int32_t timeZoneOffsetSecs() const {
        return offsetSecs_;
    }

    // // convert the time to server time zone
    // void toLocal(DataSet &data) {
    // return toLocal(data, offsetSecs_);
    // }

    // // convert the time to specific time zone
    // static void toLocal(DataSet &data, int32_t offsetSecs);

private:
    int64_t sessionId_{-1};
    Connection conn_;
    ConnectionPool *pool_{nullptr};
    AuthReq req_;
    // empty means not a named timezone
    std::string timezoneName_;
    int32_t offsetSecs_;
};

}  // namespace nebula::client
