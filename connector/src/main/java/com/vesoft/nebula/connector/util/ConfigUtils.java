
package com.vesoft.nebula.connector.util;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_AUTH_OPTIONS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_BATCH_SIZE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_CONNECT_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_REQUEST_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_MODE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_PARTITION;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_SINK_RETRY_TIMES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_BATCH_SIZE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECTOR_CONNECT_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECTOR_REQUEST_TIMEOUT;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECT_RETRY_TIMES;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_CONNECT_SINK_PARTITION;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_INSERT_MODE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.DEFAULT_INTERVAL_BETWEEN_RETRY;

import com.alibaba.fastjson.JSON;
import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import java.util.HashMap;
import java.util.Map;

public class ConfigUtils {
    public static String connectorName(Map<?, ?> connectorProps) {
        Object nameValue = connectorProps.get("name");
        return nameValue == null ? null : nameValue.toString();
    }

    public static Map<String, Object> authOptions(Map<?, ?> connectorProps) {
        if (!connectorProps.containsKey(CONNECT_AUTH_OPTIONS)) {
            return new HashMap<>();
        }
        String authOptionsString = (String) connectorProps.get(CONNECT_AUTH_OPTIONS);
        if (authOptionsString == null || authOptionsString.isEmpty()) {
            return new HashMap<>();
        }
        Map<String, Object> authOptions = (Map<String, Object>) JSON.parse(authOptionsString);
        return authOptions;
    }

    public static int retryTimes(Map<?, ?> connectorProps) {
        if (connectorProps.containsKey(CONNECT_SINK_RETRY_TIMES)) {
            return (Integer) connectorProps.get(CONNECT_SINK_RETRY_TIMES);
        } else {
            return DEFAULT_CONNECT_RETRY_TIMES;
        }
    }

    public static long intervalTimeMill(Map<?, ?> connectorProps) {
        if (connectorProps.containsKey(CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY)) {
            return (Long) connectorProps.get(CONNECT_SINK_INTERVAL_TIME_MILL_BETWEEN_RETRY);
        } else {
            return DEFAULT_INTERVAL_BETWEEN_RETRY;
        }
    }

    public static NebulaSinkConnectConfig.InsertMode sinkMode(Map<?, ?> connectorProps) {
        if (!connectorProps.containsKey(CONNECT_SINK_MODE)) {
            return DEFAULT_INSERT_MODE;
        }
        Object modeValue = connectorProps.get(CONNECT_SINK_MODE);
        if (modeValue == null) {
            return NebulaSinkConnectConfig.InsertMode.INSERTIGNORE;
        } else {
            switch (String.valueOf(modeValue).trim().toUpperCase()) {
                case "INSERT":
                    return NebulaSinkConnectConfig.InsertMode.INSERT;
                case "INSERTIGNORE":
                    return NebulaSinkConnectConfig.InsertMode.INSERTIGNORE;
                case "INSERTREPLACE":
                    return NebulaSinkConnectConfig.InsertMode.INSERTREPLACE;
                case "UPDATE":
                    return NebulaSinkConnectConfig.InsertMode.UPDATE;
                case "DELETE":
                    return NebulaSinkConnectConfig.InsertMode.DELETE;
                case "DETACHDELETE":
                    return NebulaSinkConnectConfig.InsertMode.DETACHDELETE;
                default:
            }
        }
        throw new IllegalArgumentException("wrong config for sink mode, optional mode are: "
                                                   + "INSERT, UPDATE, DELETE.");
    }


    public static int connectTimeout(Map<?, ?> connectorProps) {
        if (connectorProps.containsKey(CONNECT_CONNECT_TIMEOUT)) {
            return (Integer) connectorProps.get(CONNECT_CONNECT_TIMEOUT);
        } else {
            return DEFAULT_CONNECTOR_CONNECT_TIMEOUT;
        }
    }

    public static int requestTimeout(Map<?, ?> connectorProps) {
        if (connectorProps.containsKey(CONNECT_REQUEST_TIMEOUT)) {
            return (Integer) connectorProps.get(CONNECT_REQUEST_TIMEOUT);
        } else {
            return DEFAULT_CONNECTOR_REQUEST_TIMEOUT;
        }
    }

    public static int sinkPartition(Map<?, ?> connectorProps) {
        if (connectorProps.containsKey(CONNECT_SINK_PARTITION)) {
            return (Integer) connectorProps.get(CONNECT_SINK_PARTITION);
        } else {
            return DEFAULT_CONNECT_SINK_PARTITION;
        }
    }

    public static int batchSize(Map<?, ?> connectorProps) {
        if (connectorProps.containsKey(CONNECT_BATCH_SIZE)) {
            return (Integer) connectorProps.get(CONNECT_BATCH_SIZE);
        } else {
            return DEFAULT_BATCH_SIZE;
        }
    }
}
