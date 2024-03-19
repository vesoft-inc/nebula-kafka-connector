/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import java.io.Serializable;
import java.util.List;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * The Session is an object that operates with nebula-graph.
 * It provides an interface  `execute` to send any NGQL.
 * The returned data result `ResultSet` include wrapped string encoding
 * and time zone calculations and Node and Relationship and PathWrapper
 * and DateWrapper and TimeWrapper and DateTimeWrapper.
 * The data type obtained by the user is `ValueWrapper`,
 * which is the wrapper of the original data structure Value returned by the server.
 * The user can directly read the data using the interface of ValueWrapper.
 */

public class Session implements Serializable {

    private static final long serialVersionUID = -8855886967097862376L;

    private static final Logger log = LoggerFactory.getLogger(Session.class);

    private final long sessionID;
    private GrpcConnection connection;
    private final long requestTimeout;
    private final LoadBalancer loadBalancer;

    /**
     * Constructor
     *
     * @param connection   the connection from the pool
     * @param authResult   the auth result from graph service
     * @param loadBalancer the loadBalancer
     */
    protected Session(GrpcConnection connection,
                      long requestTimeout,
                      AuthResult authResult,
                      LoadBalancer loadBalancer) {
        this.connection = connection;
        this.sessionID = authResult.getSessionId();
        this.loadBalancer = loadBalancer;
        this.requestTimeout = requestTimeout;
    }

    /**
     * Execute the nGql sentence.
     *
     * @param stmt The nGql sentence.
     *             such as insert ngql `INSERT VERTEX person(name) VALUES "Tom":("Tom");`
     * @return The ResultSet
     */
    public ResultSet execute(String stmt, long timeoutMs) throws IOErrorException {
        return new ResultSet(connection.execute(sessionID, stmt, timeoutMs));
    }


    /**
     * check current session is ok
     *
     * @return boolean
     */
    protected boolean pingSession() throws IOErrorException {
        if (connection == null) {
            return false;
        }
        return connection.ping(sessionID);
    }

    /**
     * Notifies the server that the session is no longer needed
     * and returns the connection to the pool,
     * and the connection will be reuse.
     * This function is called if the user is no longer using the session.
     */
    protected void release() {
        if (connection == null) {
            return;
        }
        try {
            connection.signout(sessionID);
        } catch (Exception e) {
            log.warn("Release session failed:" + e.getMessage());
        }
        connection = null;
    }

    /**
     * Gets the service address of the current connection
     *
     * @return HostAddress the graph service address
     */
    protected HostAddress getGraphHost() {
        if (connection == null) {
            return null;
        }
        return connection.getServerAddress();
    }

}
