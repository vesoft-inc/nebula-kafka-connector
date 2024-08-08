package com.vesoft.nebula.client.graph.net;

import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_BLOCK_WHEN_EXHAUSTED;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_CONNECT_TIMEOUT_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_HEALTH_CHECK_TIME_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_IDLE_EVICT_SCHEDULE_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_MAX_CLIENT_SIZE;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_MAX_TIMEOUT_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_MAX_WAIT_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_MIN_CLIENT_SIZE;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_MIN_EVICTABLE_IDLE_TIME_MS;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_REQUEST_TIMEOUT_MS;
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

    private       GenericObjectPool<NebulaClient> pool;
    private       LoadBalancer                    loadBalancer;
    private final AtomicBoolean                   hasInit  = new AtomicBoolean(false);
    private final AtomicBoolean                   isClosed = new AtomicBoolean(false);
    private       long                            maxWaitMills;


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
                builder.connectTimeoutMills,
                builder.requestTimeoutMills,
                builder.scanParallel,
                builder.workingGraph,
                builder.timeZone,
                builder.schemaName,
                builder.parameters);
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
        private final List<HostAddress>   address;
        private final String              userName;
        private final String              password;
        private       Map<String, Object> authOptions = new HashMap<>();

        private int maxClientSize = DEFAULT_MAX_CLIENT_SIZE;
        private int minClientSize = DEFAULT_MIN_CLIENT_SIZE;

        private long connectTimeoutMills = DEFAULT_CONNECT_TIMEOUT_MS;
        private long requestTimeoutMills = DEFAULT_REQUEST_TIMEOUT_MS;

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
        private ZoneId              timeZone     = null;
        private String              schemaName   = null;
        private Map<String, String> parameters   = new HashMap<>();
        private int                 scanParallel = DEFAULT_SCAN_PARALLEL;

        /**
         * Builder for {@link NebulaPool}
         *
         * @param addresses graphd servers address
         * @param userName  username
         * @param password  password for user
         */
        public Builder(String addresses, String userName, String password) {
            try {
                this.address = AddressUtil.validateAddress(addresses);
            } catch (UnknownHostException e) {
                throw new RuntimeException(e);
            }
            this.userName = userName;
            this.password = password;
        }

        /**
         * config the auth options for user
         *
         * @param authOptions map of auth options
         * @return NebulaPool.Builder
         */
        public Builder withAuthOptions(Map<String, Object> authOptions) {
            if (authOptions != null) {
                this.authOptions.putAll(authOptions);
            }
            return this;
        }

        /**
         * config the max client size for pool
         *
         * @param maxClientSize max client size
         * @return NebulaPool.Builder
         */
        public Builder withMaxClientSize(int maxClientSize) {
            if (maxClientSize < 1) {
                throw new IllegalArgumentException("maxClientSize cannot be less than 1.");
            }
            this.maxClientSize = maxClientSize;
            return this;
        }

        /**
         * config the min client size for pool
         *
         * @param minClientSize min client size
         * @return NebulaPool.Builder
         */
        public Builder withMinClientSize(int minClientSize) {
            if (minClientSize < 0) {
                throw new IllegalArgumentException("minClientSize cannot be less than 0.");
            }
            this.minClientSize = minClientSize;
            return this;
        }

        /**
         * config the timeout for tcp connect, unit: ms
         * the value must be larger than 0 and smaller than Integer.MAX_VALUE in jdk 8.
         *
         * @param connectTimeoutMills timeout ms
         * @return NebulaPool.Builder
         */
        public Builder withConnectTimeoutMills(long connectTimeoutMills) {
            if (connectTimeoutMills <= 0 || connectTimeoutMills > DEFAULT_MAX_TIMEOUT_MS) {
                this.connectTimeoutMills = DEFAULT_MAX_TIMEOUT_MS;
            } else {
                this.connectTimeoutMills = connectTimeoutMills;
            }
            return this;
        }

        /**
         * config the timeout for rpc request, unit: ms
         * the value should be larger than 0 and smaller than Integer.MAX_VALUE in jdk 8.
         *
         * @param requestTimeoutMills timeout ms
         * @return NebulaPool.Builder
         */
        public Builder withRequestTimeoutMills(long requestTimeoutMills) {
            if (requestTimeoutMills <= 0 || requestTimeoutMills > DEFAULT_MAX_TIMEOUT_MS) {
                this.requestTimeoutMills = DEFAULT_MAX_TIMEOUT_MS;
            } else {
                this.connectTimeoutMills = requestTimeoutMills;
            }
            return this;
        }

        /**
         * config time to periodically check the health of graphd servers
         *
         * @param healthCheckTimeMills health check time for graphd servers
         * @return NebulaPool.Builder
         */
        public Builder withHealthCheckTimeMills(long healthCheckTimeMills) {
            this.healthCheckTimeMills = Math.max(healthCheckTimeMills, 0);
            return this;
        }

        /**
         * config if block and wait when object in NebulaPool is exhausted.
         * if false, then throw exception immediately when there's no idle object in NebulaPool.
         *
         * @param blockWhenExhausted if block when NebulaPool is exhausted
         * @return NebulaPool.Builder
         */
        public Builder withBlockWhenExhausted(boolean blockWhenExhausted) {
            this.blockWhenExhausted = blockWhenExhausted;
            return this;
        }

        /**
         * config the maximum wait time that the getClient should block before throwing
         * exception when the NebulaPool is exhausted and blockWhenExhausted is true. unit: ms
         * if the value is not positive, then always wait.
         *
         * @param maxWaitMills maximum time
         * @return NebulaPool.Builder
         */
        public Builder withMaxWaitMills(long maxWaitMills) {
            this.maxWaitMills = maxWaitMills <= 0 ? Long.MAX_VALUE : maxWaitMills;
            return this;
        }

        /**
         * config the schedule time to evict the idle object in NebulaPool.
         *
         * @param idleEvictScheduleMills sleep time between runs of the idle object evict task,
         *                               if the value is not positive, do not run the evict task.
         * @return NebulaPool.Builder
         */
        public Builder withIdleEvictScheduleMills(long idleEvictScheduleMills) {
            this.idleEvictScheduleMills = idleEvictScheduleMills;
            return this;
        }

        /**
         * config the minimum idle time for object in NebulaPool before it is eligible for
         * eviction by the evict task. unit: ms
         *
         * @param minEvictableIdleTimeMillis minimum idle for object in pool, if the value is
         *                                   not positive, do not evict any idle object
         * @return NebulaPool.Builder
         */
        public Builder withMinEvictableIdleTimeMillis(long minEvictableIdleTimeMillis) {
            this.minEvictableIdleTimeMillis = minEvictableIdleTimeMillis;
            return this;
        }

        /**
         * config whether to require all graphd servers are all strictly available.
         *
         * @param strictlyServerHealthy whether the servers are strictly healthy.
         *                              if true, all servers must be available,
         *                              if false, at least one server must be available.
         * @return NebulaPool.Builder
         */
        public Builder withStrictlyServerHealthy(boolean strictlyServerHealthy) {
            this.strictlyServerHealthy = strictlyServerHealthy;
            return this;
        }

        /**
         * config the initial working graph for NebulaClient in NebulaPool
         *
         * @param workingGraph working graph name
         * @return NebulaPool.Builder
         */
        public Builder withWorkingGraph(String workingGraph) {
            this.workingGraph = workingGraph;
            return this;
        }

        /**
         * config the initial ZonedId for NebulaClient in NebulaPool
         *
         * @param zoneId zone id
         * @return NebulaPool.Builder
         */
        public Builder withTimeZone(ZoneId zoneId) {
            this.timeZone = zoneId;
            return this;
        }

        /**
         * config the initial schema for NebulaClient in NebulaPool
         *
         * @param schemaName schema name
         * @return NebulaPool.Builder
         */
        public Builder withSchemaName(String schemaName) {
            this.schemaName = schemaName;
            return this;
        }

        /**
         * config the parameters for NebulaClient in NebulaPool
         * session set value $key=value
         *
         * @param parameters map of parameter key and value
         * @return NebulaPool.Builder
         */
        public Builder withParameters(Map<String, String> parameters) {
            if (parameters != null) {
                this.parameters = parameters;
            }
            return this;
        }


        /**
         * add the parameter into parameters for NebulaClient in NebulaPool
         * session set value $key=value
         *
         * @param paramName map of parameter key and value
         * @return NebulaPool.Builder
         */
        public Builder addParameter(String paramName, String value) {
            if (paramName != null) {
                this.parameters.put(paramName, value);
            }
            return this;
        }

        /**
         * config the parallel for data scan
         *
         * @param scanParallel number of the concurrency for data scan
         * @return NebulaClient.Builder
         */
        public Builder withScanParallel(int scanParallel) {
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
         * build a new {@link NebulaPool} with configs
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
