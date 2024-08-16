package com.vesoft.nebula.driver.graph.net;

import java.io.Serializable;

public class AuthResult implements Serializable {

    private static final long serialVersionUID = 8795815613377375650L;

    private final long sessionId;

    public AuthResult(long sessionId) {
        this.sessionId = sessionId;
    }

    public long getSessionId() {
        return sessionId;
    }

}
