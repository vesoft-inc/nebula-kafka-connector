// Copyright (c) 2022 vesoft inc. All rights reserved.

#ifndef COMMON_CATALOG_ELEMENT_TYPE_H_
#define COMMON_CATALOG_ELEMENT_TYPE_H_

#include <cstdint>
#include <string>

namespace nebula::client::catalog {

using DirectoryID = int64_t;
using NodeTypeID = int32_t;
using EdgeTypeID = int32_t;
using SchemaID = int64_t;
using GraphTypeID = int64_t;
using GraphID = int32_t;
using PropertyID = size_t;
using ElementTypeVersion = size_t;
using NodeTypeVersion = size_t;
using EdgeTypeVersion = size_t;
using CatalogVersion = size_t;
using Label = std::string;

static_assert(std::is_same_v<ElementTypeVersion, NodeTypeVersion>);
static_assert(std::is_same_v<ElementTypeVersion, EdgeTypeVersion>);

// forward declaration
class Graph;
class GraphType;
class Directory;
class Schema;
class ElementType;
class NodeType;
class NodeTypeStub;
class EdgeType;
class EdgeTypeStub;
class PatternMacro;
class Procedure;
class Constant;
class Library;
class Property;
class IndexProperty;

}  // namespace nebula::client::catalog

#endif
