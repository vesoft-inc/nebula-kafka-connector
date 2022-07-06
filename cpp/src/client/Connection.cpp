// Copyright (c) 2022 vesoft inc. All rights reserved.

#include "nebula/client/Connection.h"

#include <folly/SocketAddress.h>
#include <folly/io/async/AsyncSSLSocket.h>
#include <folly/io/async/AsyncSocket.h>
#include <folly/io/async/ScopedEventBaseThread.h>
#include <thrift/lib/cpp2/async/HeaderClientChannel.h>

#include <memory>
#include <stdexcept>

#include "../SSLConfig.h"
#include "interface/gen-cpp2/GraphServiceAsyncClient.h"

namespace nebula::client {

class NebulaConnectionErrMessageCallback : public folly::AsyncSocket::ErrMessageCallback {
public:
    /**
     * errMessage() will be invoked when kernel puts a message to
     * the error queue associated with the socket.
     *
     * @param cmsg      Reference to cmsghdr structure describing
     *                  a message read from error queue associated
     *                  with the socket.
     */
    void errMessage(const cmsghdr &) noexcept override {
        DLOG(ERROR) << "Connection error.";
    }

    /**
     * errMessageError() will be invoked if an error occurs reading a message
     * from the socket error stream.
     *
     * @param ex        An exception describing the error that occurred.
     */
    void errMessageError(const folly::AsyncSocketException &ex) noexcept override {
        DLOG(ERROR) << "Connection error: " << ex.what();
    }

    static auto &instance() {
        return cb_;
    }

private:
    static NebulaConnectionErrMessageCallback cb_;
};

NebulaConnectionErrMessageCallback NebulaConnectionErrMessageCallback::cb_;

Connection::Connection()
        : client_{nullptr}, clientLoopThread_(new folly::ScopedEventBaseThread()) {}

Connection::~Connection() {
    close();
    delete clientLoopThread_;
}

Connection &Connection::operator=(Connection &&c) {
    close();
    client_ = c.client_;
    c.client_ = nullptr;

    delete clientLoopThread_;
    clientLoopThread_ = c.clientLoopThread_;
    c.clientLoopThread_ = nullptr;

    return *this;
}

bool Connection::open(const std::string &address,
                      int32_t port,
                      uint32_t timeout,
                      bool enableSSL,
                      const std::string &CAPath) {
    if (address.empty()) {
        return false;
    }
    bool complete{false};
    //  std::shared_ptr<folly::AsyncSocket> socket;
    folly::AsyncTransport::UniquePtr socket;
    folly::SocketAddress socketAddr;
    try {
        socketAddr = folly::SocketAddress(address, port, true);
    } catch (const std::exception &e) {
        DLOG(ERROR) << "Invalid address: " << address << ":" << port << ": " << e.what();
        return false;
    }
    clientLoopThread_->getEventBase()->runImmediatelyOrRunInEventBaseThreadAndWait(
            [this, &complete, &socket, timeout, &socketAddr, enableSSL, &CAPath]() {
                try {
                    if (enableSSL) {
                        auto s = folly::AsyncSSLSocket::newSocket(
                                nebula::client::createSSLContext(CAPath),
                                clientLoopThread_->getEventBase());
                        s->connect(nullptr, std::move(socketAddr), timeout);
                        socket = std::move(s);
                    } else {
                        socket =
                                folly::AsyncSocket::newSocket(clientLoopThread_->getEventBase(),
                                                              std::move(socketAddr),
                                                              timeout);
                    }
                    complete = true;
                } catch (const std::exception &e) {
                    DLOG(ERROR) << "Connect failed: " << e.what();
                    complete = false;
                }
            });
    if (!complete) {
        return false;
    }
    if (!socket->good()) {
        return false;
    }
    // TODO workaround for issue #72
    // socket->setErrMessageCB(&NebulaConnectionErrMessageCallback::instance());
    auto channel = apache::thrift::HeaderClientChannel::newChannel(std::move(socket));
    channel->setTimeout(timeout);
    client_ = new graph::cpp2::GraphServiceAsyncClient(std::move(channel));
    return true;
}

AuthResponse Connection::authenticate(const AuthReq &req) {
    if (client_ == nullptr) {
        return AuthResponse{GQLStatus{"Not open connection."}, {}};
    }

    AuthResponse resp;
    try {
        resp = client_->future_authenticate(req).get();
    } catch (const std::exception &ex) {
        resp = AuthResponse{GQLStatus{std::string("RPC failure: ") + ex.what()}, {}};
    }
    return resp;
}

ExecutionResponse Connection::execute(int64_t sessionId, const std::string &stmt) {
    if (client_ == nullptr) {
        return ExecutionResponse{ExecutionOutcome{GQLStatus{"Not open connection."}, {}}, 0};
    }

    ExecutionResponse resp;
    try {
        resp = client_->future_execute(sessionId, stmt).get();
    } catch (const std::exception &ex) {
        resp = ExecutionResponse{
                ExecutionOutcome{GQLStatus{std::string("RPC failure: ") + ex.what()}, {}}, 0};
    }

    return resp;
}

void Connection::asyncExecute(int64_t sessionId, const std::string &stmt, ExecuteCallback cb) {
    if (client_ == nullptr) {
        cb(ExecutionResponse{ExecutionOutcome{GQLStatus{"Not open connection."}, {}}, 0});
        return;
    }
    client_->future_execute(sessionId, stmt).thenValue([cb = std::move(cb)](auto &&resp) {
        cb(std::move(resp));
    });
}

// ExecutionResponse Connection::executeWithParameter(
// int64_t sessionId,
// const std::string &stmt,
// const std::unordered_map<std::string, Value> &parameters) {
// if (client_ == nullptr) {
// return ExecutionResponse{ExecutionOutcome{GQLStatus{"Not open connection."}, {}}, 0};
// }

// ExecutionResponse resp;
// try {
// resp = client_->future_executeWithParameter(sessionId, stmt, parameters).get();
// } catch (const std::exception &ex) {
// resp = ExecutionResponse{ExecutionOutcome{GQLStatus{std::string("RPC failure: ") +
// ex.what()}, {}}, 0}, 0};
// }

// return resp;
// }

// void Connection::asyncExecuteWithParameter(int64_t sessionId,
// const std::string &stmt,
// const std::unordered_map<std::string, Value> &parameters,
// ExecuteCallback cb) {
// if (client_ == nullptr) {
// cb(ExecutionResponse{ExecutionOutcome{GQLStatus{"Not open connection."}, {}}, 0});
// return;
// }
// client_->future_executeWithParameter(sessionId, stmt, parameters)
// .thenValue([cb = std::move(cb)](auto &&resp) { cb(std::move(resp)); });
// }

// std::string Connection::executeJson(int64_t sessionId, const std::string &stmt) {
// if (client_ == nullptr) {
// // TODO handle error
// return "";
// }

// std::string json;
// try {
// json = client_->future_executeJson(sessionId, stmt).get();
// } catch (const std::exception &ex) {
// // TODO handle error
// json = "";
// }

// return json;
// }

// void Connection::asyncExecuteJson(int64_t sessionId,
// const std::string &stmt,
// ExecuteJsonCallback cb) {
// if (client_ == nullptr) {
// // TODO handle error
// cb("");
// return;
// }
// client_->future_executeJson(sessionId, stmt).thenValue(std::move(cb));
// }

// std::string Connection::executeJsonWithParameter(
// int64_t sessionId,
// const std::string &stmt,
// const std::unordered_map<std::string, Value> &parameters) {
// if (client_ == nullptr) {
// // TODO handle error
// return "";
// }

// std::string json;
// try {
// json = client_->future_executeJsonWithParameter(sessionId, stmt, parameters).get();
// } catch (const std::exception &ex) {
// // TODO handle error
// json = "";
// }

// return json;
// }

// void Connection::asyncExecuteJsonWithParameter(
// int64_t sessionId,
// const std::string &stmt,
// const std::unordered_map<std::string, Value> &parameters,
// ExecuteJsonCallback cb) {
// if (client_ == nullptr) {
// // TODO handle error
// cb("");
// return;
// }
// client_->future_executeJsonWithParameter(sessionId, stmt,
// parameters).thenValue(std::move(cb));
// }

bool Connection::isOpen() {
    return ping();
}

void Connection::close() {
    if (client_ != nullptr) {
        clientLoopThread_->getEventBase()->runImmediatelyOrRunInEventBaseThreadAndWait(
                [this]() { delete client_; });
        client_ = nullptr;
    }
}

bool Connection::ping() {
    auto resp = execute(0 /*Only check connection*/, "YIELD 1");
    // TODO check GQLStatus
    // if (resp.errorCode == ErrorCode::E_RPC_FAILURE || resp.errorCode ==
    // ErrorCode::E_DISCONNECTED) { DLOG(ERROR) << "Ping failed: " << *resp.errorMsg; return
    // false;
    // }
    return true;
}

void Connection::signout(int64_t sessionId) {
    if (client_ != nullptr) {
        client_->future_signout(sessionId).wait();
    }
}

}  // namespace nebula::client
