package com.vesoft.nebula.client.graph.exception;

/**
 *
 */
public class AuthFailedException extends Exception {
    public AuthFailedException(String message) {
        super(String.format("Auth failed: %s", message));
    }
}

