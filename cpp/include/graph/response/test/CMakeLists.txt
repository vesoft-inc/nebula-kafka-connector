# Copyright (c) 2022 vesoft inc. All rights reserved.

nebula_add_test(
    NAME graph_response_serialization_test
    SOURCES
        GraphResponseSerializationTest.cpp
    OBJECTS
        $<TARGET_OBJECTS:base_obj>
        $<TARGET_OBJECTS:datatype_obj>
        $<TARGET_OBJECTS:graph_response_obj>
        $<TARGET_OBJECTS:common_thrift_obj>
        $<TARGET_OBJECTS:graph_thrift_obj>
    LIBRARIES
        gtest
        ${THRIFT_LIBRARIES}
)
