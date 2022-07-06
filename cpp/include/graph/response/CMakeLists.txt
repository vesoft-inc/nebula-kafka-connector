# Copyright (c) 2022 vesoft inc. All rights reserved.

nebula_add_library(
    graph_response_obj OBJECT
    AuthResponse.cpp
    ExecutionOutcome.cpp
    ExecutionResponse.cpp
    GqlStatus.cpp
)

nebula_add_subdirectory(test)
