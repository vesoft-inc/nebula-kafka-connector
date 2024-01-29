package com.vesoft.nebula.client.graph;

public enum ErrorCode {
    SUCCESSFUL_COMPLETION("00000"),
    SEMANTIC_ERROR_PREFIX("NS"),
    SYNTAX_ERROR_PREFIX("42"),
    SESSION_ERROR_PREFIX("NE"),


    ;

    public final String code;

    ErrorCode(String c) {
        code = c;
    }
}
