/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.exceptions;

import org.apache.kafka.connect.errors.ConnectException;

public class RecordConversionException extends ConnectException {
    public RecordConversionException(String msg) {
        super(msg);
    }
}
