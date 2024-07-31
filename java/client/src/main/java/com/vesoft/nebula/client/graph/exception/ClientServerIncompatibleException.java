package com.vesoft.nebula.client.graph.exception;

/**
 *
 */
public class ClientServerIncompatibleException extends Exception {
    public ClientServerIncompatibleException(String message) {
        super("Current client is not compatible with the remote server, please check the "
              + "version: "  + message);
    }
}

