// Copyright (c) 2022 vesoft inc. All rights reserved.

#pragma once

#include <string>
#include <vector>

#include "common/datatypes/Geography.h"

namespace nebula::client {
namespace geo {

class WKTWriter {
public:
    WKTWriter() {}

    ~WKTWriter() {}

    std::string write(const Geography& geog) const;

    void writeCoordinate(std::string& wkt, const Coordinate& coord) const;

    void writeCoordinateList(std::string& wkt, const std::vector<Coordinate>& coordList) const;

    void writeCoordinateListList(
            std::string& wkt, const std::vector<std::vector<Coordinate>>& coordListList) const;

    void writeDouble(std::string& wkt, double v) const;
};

}  // namespace geo
}  // namespace nebula::client
