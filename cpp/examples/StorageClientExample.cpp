// Copyright (c) 2022 vesoft inc. All rights reserved.

#include <common/Init.h>
#include <nebula/sclient/ScanEdgeIter.h>
#include <nebula/sclient/StorageClient.h>

#include <atomic>
#include <chrono>
#include <limits>
#include <thread>

int main(int argc, char* argv[]) {
    nebula::client::init(&argc, &argv);

    nebula::client::StorageClient c({"127.0.0.1:9559"});

    nebula::client::ScanEdgeIter scanEdgeIter =
            c.scanEdgeWithPart("nba",
                               1,
                               "like",
                               std::vector<std::string>{"likeness"},
                               10,
                               0,
                               std::numeric_limits<int64_t>::max(),
                               "",
                               true,
                               true);
    std::cout << "scan edge..." << std::endl;
    while (scanEdgeIter.hasNext()) {
        std::cout << "-------------------------" << std::endl;
        nebula::client::DataSet ds = scanEdgeIter.next();
        std::cout << ds << std::endl;
        std::cout << "+++++++++++++++++++++++++" << std::endl;
    }

    return 0;
}
