/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.google.common.net.InetAddresses;
import com.google.common.net.InternetDomainName;
import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import java.io.Serializable;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;
import org.apache.commons.pool2.impl.GenericObjectPool;
import org.apache.commons.pool2.impl.GenericObjectPoolConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class NebulaClient implements Serializable {

    private final Logger log = LoggerFactory.getLogger(this.getClass());

    private SessionPoolConfig sessionPoolConfig;
    private GenericObjectPool<Session> pool;
    private LoadBalancer loadBalancer;

    private final AtomicBoolean hasInit = new AtomicBoolean(false);
    private final AtomicBoolean isClosed = new AtomicBoolean(false);
    private int retryTimes;
    private int intervalTime;
    private long maxWaitMills;

    private NebulaClient(Builder builder) throws IOErrorException {
        if (hasInit.get()) {
            return;
        }
        this.retryTimes = builder.retryTimes;
        this.intervalTime = builder.intervalTime;
        this.maxWaitMills = builder.maxWaitMills;

        this.loadBalancer = new RoundRobinLoadBalancer(builder.address, builder.connectTimeout,
                builder.requestTimeout, builder.strictlyServerHealthy, builder.healthCheckTime);
        if (!loadBalancer.isServersOK()) {
            loadBalancer.close();
            log.error("servers status is not ok, please check the server status or network.");
            throw new IOErrorException(IOErrorException.E_SERVER_BAD, "Servers status is not ok.");
        }

        sessionPoolConfig = new SessionPoolConfig(builder.address, builder.userName,
                builder.password)
                .setMaxSessionSize(builder.maxSessionSize)
                .setMinSessionSize(builder.minSessionSize)
                .setConnTimeout(builder.connectTimeout)
                .setRequestTimeout(builder.requestTimeout)
                .setRetryTimes(builder.retryTimes)
                .setIntervalTime(builder.intervalTime)
                .setReconnect(builder.reconnect)
                .setHealthCheckTime(builder.healthCheckTime)
                .setBlockWhenExhausted(builder.blockWhenExhausted)
                .setMaxWaitMills(builder.maxWaitMills)
                .setStrictlyServerHealthy(builder.strictlyServerHealthy);
        // pool config
        GenericObjectPoolConfig objConfig = new GenericObjectPoolConfig();
        objConfig.setMaxIdle(builder.maxSessionSize);
        objConfig.setMinIdle(builder.minSessionSize);
        objConfig.setMaxTotal(builder.maxSessionSize);
        objConfig.setBlockWhenExhausted(builder.blockWhenExhausted);
        objConfig.setMaxWaitMillis(builder.maxWaitMills);
        // just test the validation when session is idle. When session execute failed,
        // there's retry mechanism to get new Session, so no need to test when borrow or return.
        if (builder.healthCheckTime > 0) {
            objConfig.setTestWhileIdle(true);
            objConfig.setTimeBetweenEvictionRunsMillis(builder.healthCheckTime);
        }

        SessionPoolFactory factory = new SessionPoolFactory(sessionPoolConfig, loadBalancer);
        pool = new GenericObjectPool<>(factory, objConfig);
        hasInit.compareAndSet(false, true);
    }

    /**
     * execute the graph statement
     *
     * @param stmt graph statement
     * @return {@link ResultSet}
     */
    public ResultSet execute(String stmt) throws IOErrorException, NoValidSessionException {
        checkClosed();
        int tryTimes = 0;
        Session session = getSession();
        ResultSet resultSet = null;

        // execute times will be (retryTimes + 1)
        while (tryTimes++ < retryTimes + 1) {
            try {
                resultSet = session.execute(stmt);
                if (resultSet.isSucceeded() || "E_SEMANTIC_ERROR".equals(resultSet.getGqlStatus())
                        || "E_SYNTAX_ERROR".equals(resultSet.getGqlStatus())) {
                    pool.returnObject(session);
                    return resultSet;
                }
                log.debug(String.format("execute error,  message: %s, retry: %d",
                        resultSet.getGqlStatus(), tryTimes));

                // destory invalid session and re-execute
                if ("E_SESSION_INVALID".equals(resultSet.getGqlStatus())
                        || "E_SESSION_TIMEOUT".equals(resultSet.getGqlStatus())) {
                    invalidSession(session);
                } else {
                    pool.returnObject(session);
                }
                if (tryTimes <= retryTimes) {
                    sleep();
                    session = getSession();
                }
            } catch (IOErrorException e) {
                loadBalancer.updateServersStatus();
                // still has retry time. update the connection for session.
                if (tryTimes <= retryTimes) {
                    log.warn(String.format("execute failed for IOErrorException, "
                            + "message: %s, tryTime: %d", e.getMessage(), tryTimes));
                    sleep();
                    session.retryConnect();
                } else {
                    // retry time is exhausted
                    pool.returnObject(session);
                    throw e;
                }
            }
        }
        return resultSet;
    }


    /**
     * close the client
     */
    public void close() {
        if (isClosed.get()) {
            return;
        }
        isClosed.compareAndSet(false, true);
        loadBalancer.close();
        if (pool != null && !pool.isClosed()) {
            pool.close();
        }
    }


    /**
     * get available session
     *
     * @return Session
     */
    private Session getSession() throws NoValidSessionException {
        checkClosed();
        try {
            return pool.borrowObject(maxWaitMills);
        } catch (Exception e) {
            log.error("get session from pool failed.", e);
            throw new NoValidSessionException(e.getMessage());
        }
    }

    /**
     * invalid broken session
     */
    private void invalidSession(Session session) {
        try {
            pool.invalidateObject(session);
        } catch (Exception e) {
            log.warn("set session to invalidate object failed.", e);
        }
    }


    /**
     * sleep interval time
     */
    private void sleep() {
        try {
            Thread.sleep(intervalTime);
        } catch (InterruptedException e) {
            // ignore
        }
    }

    /**
     * check if the client has been closed
     *
     * @throws RuntimeException if client has been closed.
     */
    private void checkClosed() {
        if (isClosed.get()) {
            throw new RuntimeException("NebulaClient has closed. Couldn't use again.");
        }
    }

    /**
     * internal class to build a {@link NebulaClient}
     */
    public static class Builder {
        private final List<HostAddress> address;
        private final String userName;
        private final String password;
        // define default value

        // The max sessions in pool
        private int maxSessionSize = 10;

        // The min sessions in pool
        private int minSessionSize = 1;

        // socket timeout for connection, unit: millisecond
        private int connectTimeout = 0;

        // socket timeout for request, unit: millisecond
        private int requestTimeout = 0;

        // retry times for failed execute
        private int retryTimes = 0;

        // interval time for retry, unit: millisecond
        private int intervalTime = 0;

        // if reconnect for broken graphd server or network jitter
        private boolean reconnect = false;

        // The healthCheckTime for schedule check the health of session, unit: millisecond
        private int healthCheckTime = 3000;

        // if block when session is exhausted, if false, throw exception.
        private boolean blockWhenExhausted = false;

        // the max wait time if blockWhenExhausted is true. if value is less than 0, always wait.
        // unit: millisecond
        private long maxWaitMills = -1;

        // if need all servers are strictly healthy.
        // if true, all addresses must be available, if false, at least one address is available.
        private boolean strictlyServerHealthy = false;


        public Builder(String address, String userName, String password)
                throws UnknownHostException {
            this.address = validateAddress(address);
            this.userName = userName;
            this.password = password;
        }

        public Builder setMaxSessionSize(int maxSessionSize) {
            if (maxSessionSize < 1) {
                throw new IllegalArgumentException("maxSessionSize cannot be less than 1.");
            }
            this.maxSessionSize = maxSessionSize;
            return this;
        }

        public Builder setMinSessionSize(int minSessionSize) {
            if (minSessionSize < 0) {
                throw new IllegalArgumentException("minSessionSize cannot be less than 0.");
            }
            this.minSessionSize = minSessionSize;
            return this;
        }

        public Builder setConnectTimeoutMills(int connectTimeout) {
            if (connectTimeout < 0) {
                throw new IllegalArgumentException("connect timeout cannot be less than 0.");
            }
            this.connectTimeout = connectTimeout;
            return this;
        }

        public Builder setRequestTimeoutMills(int requestTimeout) {
            if (requestTimeout < 0) {
                throw new IllegalArgumentException("request timeout cannot be less than 0.");
            }
            this.requestTimeout = requestTimeout;
            return this;
        }

        public Builder setRetryTimes(int retryTimes) {
            if (retryTimes < 0) {
                throw new IllegalArgumentException("retryTimes cannot be less than 0.");
            }
            this.retryTimes = retryTimes;
            return this;
        }

        public Builder setIntervalTimeMills(int intervalTime) {
            if (intervalTime < 0) {
                throw new IllegalArgumentException("intervalTime cannot be less than 0.");
            }
            this.intervalTime = intervalTime;
            return this;
        }

        public Builder setReconnect(boolean reconnect) {
            this.reconnect = reconnect;
            return this;
        }

        public Builder setHealthCheckTimeMills(int healthCheckTime) {
            if (healthCheckTime < 0) {
                throw new IllegalArgumentException("healthCheckTime cannot be less than 0.");
            }
            this.healthCheckTime = healthCheckTime;
            return this;
        }

        public Builder setBlockWhenExhausted(boolean blockWhenExhausted) {
            this.blockWhenExhausted = blockWhenExhausted;
            return this;
        }

        public Builder setMaxWaitMills(long maxWaitMills) {
            this.maxWaitMills = maxWaitMills;
            return this;
        }

        public Builder setStrictlyServerHealthy(boolean strictlyServerHealthy) {
            this.strictlyServerHealthy = strictlyServerHealthy;
            return this;
        }

        public void check() {
            if (address == null) {
                throw new IllegalArgumentException("Graph addresses cannot be empty.");
            }
            if (userName == null || userName.trim().isEmpty()) {
                throw new IllegalArgumentException("user name cannot be blank.");
            }
            if (password == null || password.trim().isEmpty()) {
                throw new IllegalArgumentException("password cannot be blank.");
            }
        }

        /**
         * construct a NebulaClient with configs
         */
        public NebulaClient build() throws IOErrorException {
            check();
            return new NebulaClient(this);
        }

        /**
         * validate the graph addresses
         *
         * @param addresses graph server addresses, multiple addresses are split by comma
         * @return List of HostAddress
         * @throws IllegalArgumentException if address id not split by comma or port is beyond range
         * @throws UnknownHostException     if address host is wrong
         */
        private List<HostAddress> validateAddress(String addresses) throws UnknownHostException {
            List<HostAddress> newAddrs = new ArrayList<>();
            for (String addr : addresses.split(",")) {
                String[] hostAndPort = addr.split(":");
                if (hostAndPort.length < 2) {
                    throw new IllegalArgumentException("wrong server address " + addr);
                }
                String host = hostAndPort[0];
                int port = Integer.parseInt(hostAndPort[1]);

                // get all host name
                InetAddress[] inetAddresses = InetAddress.getAllByName(host);
                for (InetAddress inetAddress : inetAddresses) {
                    String ip = inetAddress.getHostAddress();
                    if (!(InetAddresses.isInetAddress(ip)
                            || InetAddresses.isUriInetAddress(ip)
                            || InternetDomainName.isValid(ip))
                            || (port <= 0 || port >= 65535)) {
                        throw new IllegalArgumentException(
                                String.format("host %s and port %d is illegal.", ip, port));
                    }
                    newAddrs.add(new HostAddress(ip, port));
                }
            }
            return newAddrs;
        }
    }

    public static Builder builder(String address, String userName, String password)
            throws UnknownHostException {
        return new Builder(address, userName, password);
    }


    /**
     * get configs of NebulaClient
     */
    public SessionPoolConfig getConfig() {
        return sessionPoolConfig;
    }

    /**
     * get the active sessions of NebulaClient
     *
     * @return number of sessions which are being used
     */
    public int getActiveSessions() {
        return pool.getNumActive();
    }

    /**
     * get the idle sessions of NebulaClient
     *
     * @return number of sessions which are idle
     */
    public int getIdleSessions() {
        return pool.getNumIdle();
    }

    /**
     * get the execute requests waiting to get session
     *
     * @return number of execute requests waiting to get session
     */
    public int getWaiters() {
        return pool.getNumWaiters();
    }
}
