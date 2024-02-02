/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.util;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_MODE;
import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import java.util.Map;

public class ConfigUtils {
    public static String connectorName(Map<?, ?> connectorProps) {
        Object nameValue = connectorProps.get("name");
        return nameValue == null ? null : nameValue.toString();
    }

    public static NebulaSinkConnectConfig.InsertMode sinkMode(Map<?, ?> connectorProps) {
        Object modeValue = connectorProps.get(CONNECT_SINK_MODE);
        if (modeValue == null) {
            return NebulaSinkConnectConfig.InsertMode.INSERT;
        } else {
            switch (modeValue.toString().trim().toUpperCase()) {
                case "INSERT":
                    return NebulaSinkConnectConfig.InsertMode.INSERT;
                case "UPDATE":
                    return NebulaSinkConnectConfig.InsertMode.UPDATE;
                case "DELETE":
                    return NebulaSinkConnectConfig.InsertMode.DELETE;
            }
        }
        throw new IllegalArgumentException("wrong config for sink mode, optional mode are: " +
                "INSERT, UPDATE, DELETE.");
    }
}
