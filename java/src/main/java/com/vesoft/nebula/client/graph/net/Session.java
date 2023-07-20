/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.NList;
import com.vesoft.nebula.Value;
import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import java.io.Serializable;
import java.util.ArrayList;
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
    private SyncConnection connection;
    private final int connTimeout;
    private final int requestTimeout;
    private final LoadBalancer loadBalancer;
    private final Boolean retryConnect;

    /**
     * Constructor
     *
     * @param connection   the connection from the pool
     * @param authResult   the auth result from graph service
     * @param loadBalancer the loadBalancer
     * @param retryConnect whether to retry after the connection is disconnected
     */
    protected Session(SyncConnection connection,
                      int connTimeout,
                      int requestTimeout,
                      AuthResult authResult,
                      Boolean retryConnect,
                      LoadBalancer loadBalancer) {
        this.connection = connection;
        this.sessionID = authResult.getSessionId();
        this.loadBalancer = loadBalancer;
        this.retryConnect = retryConnect;
        this.connTimeout = connTimeout;
        this.requestTimeout = requestTimeout;
    }

    /**
     * Execute the nGql sentence.
     *
     * @param stmt The nGql sentence.
     *             such as insert ngql `INSERT VERTEX person(name) VALUES "Tom":("Tom");`
     * @return The ResultSet
     */
    protected ResultSet execute(String stmt) throws IOErrorException {
        return new ResultSet(connection.execute(sessionID, stmt));
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
            log.warn("Release session or return object to pool failed:" + e.getMessage());
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

    /**
     * set current connection is invalid, and get a new connection from the pool,
     * if get connection failed, return false, else return true
     */
    protected void retryConnect() throws IOErrorException {
        connection.close();
        List<HostAddress> goodHosts = loadBalancer.getGoodAddresses();
        int tryConnect = goodHosts.size();
        SyncConnection newConnection = new SyncConnection();
        while (tryConnect-- > 0) {
            try {
                newConnection.open(loadBalancer.getAddress(), connTimeout, requestTimeout);
                connection = newConnection;
                break;
            } catch (IOErrorException e) {
                if (tryConnect == 0 || !retryConnect) {
                    throw e;
                } else {
                    log.warn("connect failed, " + e.getMessage());
                }
            }
        }
    }

    /**
     * convert java list to nebula thrift list
     *
     * @param list java list
     * @return nebula list
     */
    private static NList list2Nlist(List<Object> list) throws UnsupportedOperationException {
        NList nlist = new NList(new ArrayList<Value>());
        for (Object item : list) {
            nlist.values.add(value2Nvalue(item));
        }
        return nlist;
    }

    /**
     * convert java value type to nebula thrift value type
     *
     * @param value java obj
     * @return nebula value
     */
    private static Value value2Nvalue(Object value) throws UnsupportedOperationException {
        Value nvalue = new Value();
        if (value == null) {
            // TODO nvalue.set(NullType.__NULL__);
        } else if (value instanceof Boolean) {
            boolean bval = (Boolean) value;
            nvalue.setBoolVal(bval);
        } else if (value instanceof Integer) {
            int ival = (Integer) value;
            nvalue.setInt32Val(ival);
        } else if (value instanceof Short) {
            short ival = (Short) value;
            nvalue.setInt16Val(ival);
        } else if (value instanceof Byte) {
            byte ival = (Byte) value;
            nvalue.setInt8Val(ival);
        } else if (value instanceof Long) {
            long ival = (Long) value;
            nvalue.setInt64Val(ival);
        } else if (value instanceof Float) {
            float fval = (Float) value;
            nvalue.setFloatVal(fval);
        } else if (value instanceof Double) {
            double dval = (Double) value;
            nvalue.setDoubleVal(dval);
        } else if (value instanceof String) {
            byte[] sval = ((String) value).getBytes();
            nvalue.setStringVal(sval);
        } else if (value instanceof List) {
            nvalue.setListVal(list2Nlist((List<Object>) value));
        } else if (value instanceof Value) {
            return (Value) value;
        } else {
            // unsupport other Value type, use this function carefully
            throw new UnsupportedOperationException(
                    "Only support convert boolean/float/int/string/map/list to nebula.Value but was"
                            + value.getClass().getTypeName());
        }
        return nvalue;
    }
}
