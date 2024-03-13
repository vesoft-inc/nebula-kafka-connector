/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import java.io.Serializable;
import java.util.List;
import org.apache.commons.pool2.BasePooledObjectFactory;
import org.apache.commons.pool2.PooledObject;
import org.apache.commons.pool2.impl.DefaultPooledObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class SessionPoolFactory extends BasePooledObjectFactory<Session> implements Serializable {

    private final Logger log = LoggerFactory.getLogger(this.getClass());

    private SessionPoolConfig sessionPoolConfig;
    private final LoadBalancer loadBalancer;

    public SessionPoolFactory(SessionPoolConfig sessionPoolConfig, LoadBalancer loadBalancer) {
        this.sessionPoolConfig = sessionPoolConfig;
        this.loadBalancer = loadBalancer;
    }


    @Override
    public Session create() throws Exception {
        List<HostAddress> goodHosts = loadBalancer.getGoodAddresses();
        if (goodHosts.size() == 0) {
            throw new IOErrorException(IOErrorException.E_ALL_BROKEN, "No servers host is "
                    + "available, please check your servers is up and network between client and "
                    + "server is connected.");
        }

        GrpcConnection connection = new GrpcConnection();
        int tryConnect = goodHosts.size();
        while (tryConnect-- > 0) {
            try {
                connection.open(loadBalancer.getAddress(), sessionPoolConfig.getRequestTimeout());
                break;
            } catch (Exception e) {
                if (tryConnect == 0 || !sessionPoolConfig.isReconnect()) {
                    throw e;
                } else {
                    log.warn("connect failed, " + e.getMessage());
                }
            }
        }
        AuthResult authResult;
        try {
            authResult = connection.authenticate(sessionPoolConfig.getUsername(),
                    sessionPoolConfig.getPassword());
        } catch (AuthFailedException e) {
            log.error(e.getMessage());
            throw e;
        }

        Session session = new Session(connection, sessionPoolConfig.getRequestTimeout(),
                authResult, sessionPoolConfig.isReconnect(), loadBalancer);
        return session;
    }

    @Override
    public PooledObject<Session> wrap(Session session) {
        return new DefaultPooledObject<>(session);
    }

    @Override
    public void destroyObject(PooledObject<Session> sessionPooledObject) throws Exception {
        Session session = sessionPooledObject.getObject();
        try {
            session.release();
        } catch (Exception e) {
            log.warn("session release failed ", e);
        }
        super.destroyObject(sessionPooledObject);
    }

    @Override
    public boolean validateObject(PooledObject<Session> sessionPooledObject) {
        Session session = sessionPooledObject.getObject();
        try {
            return session.pingSession();
        } catch (IOErrorException e) {
            log.warn("session for server {} is invalid, remove it",
                    session.getGraphHost().toString());
            return false;
        }
    }
}
