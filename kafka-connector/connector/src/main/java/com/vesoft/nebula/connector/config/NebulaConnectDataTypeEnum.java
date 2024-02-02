/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.config;


public enum NebulaConnectDataTypeEnum {
    NODE,
    EDGE,
    // there
    BOTH;

    public static NebulaConnectDataTypeEnum getDataType(String type) {
        switch (type.toUpperCase().trim()) {
            case "NODE":
                return NODE;
            case "EDGE":
                return EDGE;
            case "BOTH":
                return BOTH;
            default:
                throw new IllegalArgumentException(
                        String.format("data type %s is not supported.", type));
        }
    }
}
