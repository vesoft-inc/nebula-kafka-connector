# Copyright (c) 2022 vesoft inc. All rights reserved.

nebula_add_test(
    NAME
        binding_table_test
    SOURCES
        BindingTableTest.cpp
    OBJECTS
        ${NEBULA_CLIENT_OBJS}
        $<TARGET_OBJECTS:nebula_common_obj>
        $<TARGET_OBJECTS:nebula_graph_client_obj>
    LIBRARIES
        gtest
        ${NEBULA_THIRD_PARTY_LIBRARIES}
)
