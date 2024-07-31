
package com.vesoft.nebula.connector.exceptions;

import org.apache.kafka.connect.errors.ConnectException;

public class RecordConversionException extends ConnectException {
    public RecordConversionException(String msg) {
        super(msg);
    }
}
