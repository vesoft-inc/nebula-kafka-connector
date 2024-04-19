/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_BLOCK_WHEN_EXHAUSTED;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_HEALTH_CHECK_TIME_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_IDLE_EVICT_SCHEDULE_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_MAX_CLIENT_SIZE;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_MAX_WAIT_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_MIN_CLIENT_SIZE;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_MIN_EVICTABLE_IDLE_TIME_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_REQUEST_TIMEOUT;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_SCAN_PARALLEL;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_STRICT_SERVER_HEALTHY;

import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.utils.AddressUtil;
import java.io.Serializable;
import java.net.UnknownHostException;
import java.time.ZoneId;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicBoolean;
import org.apache.commons.pool2.impl.GenericObjectPool;
import org.apache.commons.pool2.impl.GenericObjectPoolConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class NebulaPool implements Serializable {
    private final Logger logger = LoggerFactory.getLogger(this.getClass());

    private GenericObjectPool<NebulaClient> pool;
    private LoadBalancer loadBalancer;
    private final AtomicBoolean hasInit = new AtomicBoolean(false);
    private final AtomicBoolean isClosed = new AtomicBoolean(false);
    private long maxWaitMills;


    public static NebulaPool.Builder builder(String addresses, String userName) {
        return builder(addresses, userName, null);
    }

    public static NebulaPool.Builder builder(String addresses, String userName, String password) {
        return new NebulaPool.Builder(addresses, userName, password);
    }

    private NebulaPool(Builder builder) throws IOErrorException {
        if (hasInit.get()) {
            return;
        }
        this.maxWaitMills = builder.maxWaitMills;

        this.loadBalancer = new RoundRobinLoadBalancer(
                builder.address,
                builder.userName,
                builder.authOptions,
                builder.strictlyServerHealthy,
                builder.healthCheckTimeMills);

        if (!loadBalancer.isServersOK()) {
            loadBalancer.close();
            logger.error("servers status is not ok, please check the server status or network.");
            throw new IOErrorException(IOErrorException.E_SERVER_BAD, "Servers status is not ok.");
        }
        // pool config
        GenericObjectPoolConfig objConfig = new GenericObjectPoolConfig();
        objConfig.setMaxIdle(builder.maxClientSize);
        objConfig.setMinIdle(builder.minClientSize);
        objConfig.setMaxTotal(builder.maxClientSize);
        objConfig.setBlockWhenExhausted(builder.blockWhenExhausted);
        objConfig.setMaxWaitMillis(builder.maxWaitMills);
        objConfig.setTimeBetweenEvictionRunsMillis(builder.idleEvictScheduleMills);
        objConfig.setMinEvictableIdleTimeMillis(builder.minEvictableIdleTimeMillis);
        // just test the validation when session is idle.
        if (builder.healthCheckTimeMills > 0) {
            objConfig.setTestWhileIdle(true);
            objConfig.setTimeBetweenEvictionRunsMillis(builder.healthCheckTimeMills);
        }

        ClientPoolFactory factory = new ClientPoolFactory(
                loadBalancer,
                builder.userName,
                builder.authOptions,
                builder.requestTimeoutMills,
                builder.scanParallel,
                builder.workingGraph,
                builder.timeZone);
        pool = new GenericObjectPool<>(factory, objConfig);
        hasInit.compareAndSet(false, true);
    }


    /**
     * get NebulaClient from pool
     */
    public NebulaClient getClient() throws Exception {
        return pool.borrowObject(maxWaitMills);
    }

    /**
     * return the client to object pool
     *
     * @param client NebulaClient
     */
    public void returnClient(NebulaClient client) {
        if (client.isClosed()) {
            try {
                pool.invalidateObject(client);
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        } else {
            pool.returnObject(client);
        }
    }


    public void close() {
        loadBalancer.close();
        pool.close();
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


    public static class Builder {
        private final List<HostAddress> address;
        private final String userName;
        private final String password;
        private Map<String, Object> authOptions = new HashMap<>();

        private int maxClientSize = DEFAULT_MAX_CLIENT_SIZE;
        private int minClientSize = DEFAULT_MIN_CLIENT_SIZE;
        private long requestTimeoutMills = DEFAULT_REQUEST_TIMEOUT;

        // The healthCheckTime for schedule check the health of session, unit: millisecond
        private long healthCheckTimeMills = DEFAULT_HEALTH_CHECK_TIME_MS;

        // if block when session is exhausted, if false, throw exception.
        private boolean blockWhenExhausted = DEFAULT_BLOCK_WHEN_EXHAUSTED;

        // the max wait time if blockWhenExhausted is true. if value is less than 0, always wait.
        // unit: millisecond
        private long maxWaitMills = DEFAULT_MAX_WAIT_MS;

        // the schedule time for test the idle session and evict it. if value is less than 0,
        // never evict the idle sessions.
        private long idleEvictScheduleMills = DEFAULT_IDLE_EVICT_SCHEDULE_MS;

        // the min idle time for idle session
        private long minEvictableIdleTimeMillis = DEFAULT_MIN_EVICTABLE_IDLE_TIME_MS;

        // if need all servers are strictly healthy.
        // if true, all addresses must be available, if false, at least one address is available.
        private boolean strictlyServerHealthy = DEFAULT_STRICT_SERVER_HEALTHY;

        private String workingGraph = null;

        // the time zone, used to parse ZonedTime and ZonedDatetime
        private ZoneId timeZone = null;
        private int scanParallel = DEFAULT_SCAN_PARALLEL;

        public Builder(String addresses, String userName, String password) {
            try {
                this.address = AddressUtil.validateAddress(addresses);
            } catch (UnknownHostException e) {
                throw new RuntimeException(e);
            }
            this.userName = userName;
            this.password = password;
        }

        public Builder withAuthOptions(Map<String, Object> authOptions) {
            if (authOptions != null) {
                this.authOptions.putAll(authOptions);
            }
            return this;
        }

        public Builder withMaxClientSize(int maxClientSize) {
            if (maxClientSize < 1) {
                throw new IllegalArgumentException("maxClientSize cannot be less than 1.");
            }
            this.maxClientSize = maxClientSize;
            return this;
        }

        public Builder withMinClientSize(int minClientSize) {
            if (minClientSize < 0) {
                throw new IllegalArgumentException("minClientSize cannot be less than 0.");
            }
            this.minClientSize = minClientSize;
            return this;
        }

        public Builder withRequestTimeoutMills(long requestTimeoutMills) {
            this.requestTimeoutMills =
                    requestTimeoutMills < 0 ? Long.MAX_VALUE : requestTimeoutMills;
            return this;
        }

        public Builder withHealthCheckTimeMills(long healthCheckTimeMills) {
            this.healthCheckTimeMills = Math.max(healthCheckTimeMills, 0);
            return this;
        }

        public Builder withBlockWhenExhausted(boolean blockWhenExhausted) {
            this.blockWhenExhausted = blockWhenExhausted;
            return this;
        }

        public Builder withMaxWaitMills(long maxWaitMills) {
            this.maxWaitMills = maxWaitMills < 0 ? Long.MAX_VALUE : maxWaitMills;
            return this;
        }

        public Builder withIdleEvictScheduleMills(long idleEvictScheduleMills) {
            this.idleEvictScheduleMills = idleEvictScheduleMills;
            return this;
        }

        public Builder withMinEvictableIdleTimeMillis(long minEvictableIdleTimeMillis) {
            this.minEvictableIdleTimeMillis = minEvictableIdleTimeMillis;
            return this;
        }

        public Builder withStrictlyServerHealthy(boolean strictlyServerHealthy) {
            this.strictlyServerHealthy = strictlyServerHealthy;
            return this;
        }

        public Builder withWorkingGraph(String graphName) {
            this.workingGraph = workingGraph;
            return this;
        }

        public Builder withTimeZone(ZoneId zoneId) {
            this.timeZone = zoneId;
            return this;
        }

        public Builder setScanParallel(int scanParallel) {
            this.scanParallel = scanParallel;
            return this;
        }

        public void check() {
            if (address == null) {
                throw new IllegalArgumentException("Graph addresses cannot be empty.");
            }
            if (userName == null || userName.trim().isEmpty()) {
                throw new IllegalArgumentException("user name cannot be empty.");
            }
            if (authOptions.isEmpty() && (password == null || password.trim().isEmpty())) {
                throw new IllegalArgumentException(
                        "auth options and password cannot be empty at the same time.");
            }
        }

        /**
         * construct a NebulaClient with configs
         */
        public NebulaPool build() throws IOErrorException {
            check();
            if (password != null) {
                this.authOptions.put("password", password);
            }
            return new NebulaPool(this);
        }
    }
}
