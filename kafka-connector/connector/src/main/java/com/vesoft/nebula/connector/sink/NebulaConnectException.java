/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import org.apache.kafka.connect.errors.ConnectException;

public class NebulaConnectException extends ConnectException {
    public NebulaConnectException(String reason) {
        super(reason);
    }

    public NebulaConnectException(String reason, Throwable throwable) {
        super(reason, throwable);
    }
}
