package com.vesoft.nebula.driver.graph.data;

import java.io.Serializable;

public abstract class SSLParam implements Serializable {

    private static final long serialVersionUID = 7410233298826490747L;

    public enum SignMode {
        NONE,
        SELF_SIGNED,
        CA_SIGNED
    }

    private SignMode signMode;

    public SSLParam(SignMode signMode) {
        this.signMode = signMode;
    }

    public SignMode getSignMode() {
        return signMode;
    }
}
