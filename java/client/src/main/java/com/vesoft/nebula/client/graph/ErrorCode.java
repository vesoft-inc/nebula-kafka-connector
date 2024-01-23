package com.vesoft.nebula.client.graph;

public enum ErrorCode {
    SUCCESSFUL_COMPLETION("00000");

    public final String code;

    ErrorCode(String c) {
        code = c;
    }
}
