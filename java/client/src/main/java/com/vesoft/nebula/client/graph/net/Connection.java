package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import java.io.Serializable;

public abstract class Connection implements Serializable {

    private static final long serialVersionUID = -8425216612015802331L;

    protected HostAddress serverAddr = null;

    public HostAddress getServerAddress() {
        return this.serverAddr;
    }

    public abstract void open(HostAddress address, long requestTimeout)
            throws IOErrorException;

    public abstract void close();

    public abstract boolean ping(long sessionID) throws IOErrorException;
}
