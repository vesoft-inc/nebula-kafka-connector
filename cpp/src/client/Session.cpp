// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "nebula/client/Session.h"

//#include "common/time/TimeConversion.h"
#include "nebula/client/ConnectionPool.h"

namespace nebula::client {

ExecutionResponse Session::execute(const std::string &stmt) {
    return ExecutionResponse(conn_.execute(sessionId_, stmt));
}

void Session::asyncExecute(const std::string &stmt, ExecuteCallback cb) {
    conn_.asyncExecute(sessionId_, stmt, [cb = std::move(cb)](auto &&resp) {
        cb(ExecutionResponse(std::move(resp)));
    });
}

// ExecutionResponse Session::executeWithParameter(
// const std::string &stmt, const std::unordered_map<std::string, Value> &parameters) {
// return ExecutionResponse(conn_.executeWithParameter(sessionId_, stmt, parameters));
// }

// void Session::asyncExecuteWithParameter(const std::string &stmt,
// const std::unordered_map<std::string, Value> &parameters,
// ExecuteCallback cb) {
// conn_.asyncExecuteWithParameter(sessionId_, stmt, parameters, [cb = std::move(cb)](auto
// &&resp) { cb(ExecutionResponse(std::move(resp)));
// });
// }

// std::string Session::executeJson(const std::string &stmt) {
// return conn_.executeJson(sessionId_, stmt);
// }

// void Session::asyncExecuteJson(const std::string &stmt, ExecuteJsonCallback cb) {
// conn_.asyncExecuteJson(
// sessionId_, stmt, [cb = std::move(cb)](auto &&json) { cb(std::move(json)); });
// }

// std::string Session::executeJsonWithParameter(
// const std::string &stmt, const std::unordered_map<std::string, Value> &parameters) {
// return conn_.executeJsonWithParameter(sessionId_, stmt, parameters);
// }

// void Session::asyncExecuteJsonWithParameter(
// const std::string &stmt,
// const std::unordered_map<std::string, Value> &parameters,
// ExecuteJsonCallback cb) {
// conn_.asyncExecuteJsonWithParameter(
// sessionId_, stmt, parameters, [cb = std::move(cb)](auto &&json) { cb(std::move(json)); });
// }

bool Session::ping() {
    return conn_.ping();
}

GQLStatus Session::retryConnect() {
    pool_->giveBack(std::move(conn_));
    conn_ = pool_->getConnection();
    auto resp = conn_.authenticate(req_);
    sessionId_ = resp.identifier_.has_value() ? resp.identifier_.value() : -1;
    return resp.gqlStatus_;
}

void Session::release() {
    if (valid()) {
        conn_.signout(sessionId_);
        pool_->giveBack(std::move(conn_));
        sessionId_ = -1;
    }
}

// void Session::toLocal(DataSet &data, int32_t offsetSecs) {
// for (auto &row : data.rows) {
// for (auto &col : row.values) {
// if (col.isTime()) {
// col.setTime(time::TimeConversion::timeShift(col.getTime(), offsetSecs));
// } else if (col.isDateTime()) {
// col.setDateTime(time::TimeConversion::dateTimeShift(col.getDateTime(), offsetSecs));
// }
// }
// }
// }

}  // namespace nebula::client
