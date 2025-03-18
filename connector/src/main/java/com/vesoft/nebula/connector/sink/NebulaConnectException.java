
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
