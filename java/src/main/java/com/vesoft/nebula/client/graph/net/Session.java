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
import com.vesoft.nebula.graph.ExecutionResponse;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicBoolean;
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

    private final long sessionID;
    private SyncConnection connection;
    private final NebulaPool pool;
    private final Boolean retryConnect;
    private final AtomicBoolean connectionIsBroken = new AtomicBoolean(false);
    private final Logger log = LoggerFactory.getLogger(getClass());

    /**
     * Constructor
     *
     * @param connection   the connection from the pool
     * @param authResult   the auth result from graph service
     * @param connPool     the connection pool
     * @param retryConnect whether to retry after the connection is disconnected
     */
    public Session(SyncConnection connection,
                   AuthResult authResult,
                   NebulaPool connPool,
                   Boolean retryConnect) {
        this.connection = connection;
        this.sessionID = authResult.getSessionId();
        this.pool = connPool;
        this.retryConnect = retryConnect;
    }

    /**
     * Execute the nGql sentence.
     *
     * @param stmt The nGql sentence.
     *             such as insert ngql `INSERT VERTEX person(name) VALUES "Tom":("Tom");`
     * @return The ResultSet
     */
    public synchronized ResultSet execute(
            String stmt)
            throws IOErrorException {
        if (connection == null) {
            throw new IOErrorException(IOErrorException.E_CONNECT_BROKEN,
                    "The session was released, couldn't use again.");
        }

        if (connectionIsBroken.get() && retryConnect) {
            if (retryConnect()) {
                ExecutionResponse resp =
                        connection.execute(sessionID, stmt);
                return new ResultSet(resp);
            } else {
                throw new IOErrorException(IOErrorException.E_ALL_BROKEN,
                        "All servers are broken.");
            }
        }

        try {
            ExecutionResponse resp = connection.execute(sessionID, stmt);
            return new ResultSet(resp);
        } catch (IOErrorException ie) {
            if (ie.getType() == IOErrorException.E_CONNECT_BROKEN) {
                connectionIsBroken.set(true);
                pool.updateServerStatus();

                if (retryConnect) {
                    if (retryConnect()) {
                        connectionIsBroken.set(false);
                        ExecutionResponse resp =
                                connection.execute(sessionID, stmt);
                        return new ResultSet(resp);
                    } else {
                        connectionIsBroken.set(true);
                        throw new IOErrorException(IOErrorException.E_ALL_BROKEN,
                                "All servers are broken.");
                    }
                }
            }
            throw ie;
        }
    }

    /**
     * Check current connection is ok
     *
     * @return boolean
     */
    public synchronized boolean ping() {
        if (connection == null) {
            return false;
        }
        return connection.ping();
    }

    /**
     * check current session is ok
     */
    public synchronized boolean pingSession() {
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
    public synchronized void release() {
        if (connection == null) {
            return;
        }
        try {
            connection.signout(sessionID);
            pool.returnConnection(connection);
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
    public synchronized HostAddress getGraphHost() {
        if (connection == null) {
            return null;
        }
        return connection.getServerAddress();
    }

    /**
     * set current connection is invalid, and get a new connection from the pool,
     * if get connection failed, return false, else return true
     *
     * @return true or false
     */
    private boolean retryConnect() {
        try {
            pool.setInvalidateConnection(connection);
            SyncConnection newConn = pool.getConnection();
            if (newConn == null) {
                log.error("Get connection object failed.");
                return false;
            }
            connection = newConn;
            return true;
        } catch (Exception e) {
            log.error("Reconnected failed: " + e);
            return false;
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
    public static Value value2Nvalue(Object value) throws UnsupportedOperationException {
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
