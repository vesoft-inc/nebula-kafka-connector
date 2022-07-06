// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "nebula/client/ConnectionPool.h"

#include <folly/String.h>

#include <atomic>

#include "graph/response/AuthReq.h"

namespace nebula::client {

ConnectionPool::~ConnectionPool() {
    close();
}

void ConnectionPool::init(const std::vector<std::string> &addresses, const Config &config) {
    for (const auto &addr : addresses) {
        std::vector<std::string> splits;
        folly::split(':', addr, splits, true);
        if (splits.size() != 2) {
            // ignore error
            continue;
        }
        address_.emplace_back(std::make_pair(splits[0], folly::to<int32_t>(splits[1])));
    }
    if (address_.empty()) {
        // no valid address
        return;
    }
    config_ = config;
    newConnection(0, config.maxConnectionPoolSize_);
}

void ConnectionPool::close() {
    std::lock_guard<std::mutex> l(lock_);
    for (auto &conn : conns_) {
        conn.close();
    }
}

Session ConnectionPool::getSession(const AuthReq &req, bool retryConnect) {
    (void)retryConnect;
    Connection conn = getConnection();
    auto resp = conn.authenticate(req);
    // TODO check GQLStatus
    if (!resp.identifier_.has_value()) {
        return Session();
    }
    return Session(resp.identifier_.value(),
                   std::move(conn),
                   this,
                   req,
                   /*TODO*/ "",
                   /*TODO*/ 0);
}

Connection ConnectionPool::getConnection() {
    std::lock_guard<std::mutex> l(lock_);
    // check connection
    for (auto c = conns_.begin(); c != conns_.end();) {
        if (!c->isOpen()) {
            c = conns_.erase(c);
            newConnection(nextCursor(), 1);
        } else {
            c++;
        }
    }
    if (conns_.empty()) {
        return Connection();
    }
    Connection conn = std::move(conns_.front());
    conns_.pop_front();
    return conn;
}

void ConnectionPool::giveBack(Connection &&conn) {
    std::lock_guard<std::mutex> l(lock_);
    conns_.emplace_back(std::move(conn));
}

void ConnectionPool::newConnection(std::size_t cursor, std::size_t count) {
    for (std::size_t connectionCount = 0, addrCursor = cursor, loopCount = 0;
         connectionCount < count;
         ++addrCursor, ++loopCount) {
        if (loopCount > count * address_.size()) {
            // Can't get so many connections, return to avoid dead loop
            return;
        }
        if (addrCursor >= address_.size()) {
            addrCursor = 0;
        }
        Connection conn;
        if (conn.open(address_[addrCursor].first,
                      address_[addrCursor].second,
                      config_.timeout_,
                      config_.enableSSL_,
                      config_.CAPath_)) {
            ++connectionCount;
            conns_.emplace_back(std::move(conn));
        }
        // ignore error
    }
}

}  // namespace nebula::client
