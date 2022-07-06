// Copyright (c) 2022 vesoft inc. All rights reserved.

#pragma once

namespace nebula::client {

// The response of the authentication of a principal.
struct AuthReq {
    void __clear() {
        username.clear();
        password.clear();
        clientType.clear();
        clientVersion.clear();
    }

    void clear() {
        __clear();
    }

    std::string username;
    std::string password;
    std::string clientType;
    std::string clientVersion;
};

}  // namespace nebula::client
