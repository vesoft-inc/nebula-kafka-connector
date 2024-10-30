package com.vesoft.nebula.driver.graph.net;

import static com.vesoft.nebula.driver.graph.net.Constants.DEFAULT_ENABLE_TLS;
import static com.vesoft.nebula.driver.graph.net.Constants.DEFAULT_TLS_PEER_NAME_VERIFY;

import com.vesoft.nebula.driver.graph.data.HostAddress;
import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import com.vesoft.nebula.driver.graph.utils.AddressUtil;
import java.io.Serializable;
import java.net.UnknownHostException;
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

        this.loadBalancer = new RoundRobinLoadBalancer(builder);

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
        objConfig.setTestOnBorrow(builder.testOnBorrow);
        // just test the validation when session is idle.
        if (builder.healthCheckTimeMills > 0) {
            objConfig.setTestWhileIdle(true);
            objConfig.setTimeBetweenEvictionRunsMillis(builder.healthCheckTimeMills);
        }

        ClientPoolFactory factory = new ClientPoolFactory(
                loadBalancer,
                builder);
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
        protected final List<HostAddress>   address;
        protected final String              userName;
        protected final String              password;
        protected       Map<String, Object> authOptions = new HashMap<>();

        protected int maxClientSize = Constants.DEFAULT_MAX_CLIENT_SIZE;
        protected int minClientSize = Constants.DEFAULT_MIN_CLIENT_SIZE;

        protected long connectTimeoutMills = Constants.DEFAULT_CONNECT_TIMEOUT_MS;
        protected long requestTimeoutMills = Constants.DEFAULT_REQUEST_TIMEOUT_MS;

        // The healthCheckTime for schedule check the health of session, unit: millisecond
        protected long healthCheckTimeMills = Constants.DEFAULT_HEALTH_CHECK_TIME_MS;

        // whether to check the client when borrow from poll
        private boolean testOnBorrow = Constants.DEFAULT_TEST_ON_BORROW;

        // if block when session is exhausted, if false, throw exception.
        protected boolean blockWhenExhausted = Constants.DEFAULT_BLOCK_WHEN_EXHAUSTED;

        // the max wait time if blockWhenExhausted is true. if value is less than 0, always wait.
        // unit: millisecond
        protected long maxWaitMills = Constants.DEFAULT_MAX_WAIT_MS;

        // the schedule time for test the idle session and evict it. if value is less than 0,
        // never evict the idle sessions.
        protected long idleEvictScheduleMills = Constants.DEFAULT_IDLE_EVICT_SCHEDULE_MS;

        // the min idle time for idle session
        protected long minEvictableIdleTimeMillis = Constants.DEFAULT_MIN_EVICTABLE_IDLE_TIME_MS;

        // if need all servers are strictly healthy.
        // if true, all addresses must be available, if false, at least one address is available.
        protected boolean strictlyServerHealthy = Constants.DEFAULT_STRICT_SERVER_HEALTHY;

        protected String graph = null;

        // the time zone, used to parse ZonedTime and ZonedDatetime
        protected String              timeZone          = null;
        protected String              schema            = null;
        protected Map<String, String> parameters        = new HashMap<>();
        protected int                 scanParallel      = Constants.DEFAULT_SCAN_PARALLEL;
        protected boolean             enableTls         = DEFAULT_ENABLE_TLS;
        protected String              tlsCa;
        protected String              tlsCert;
        protected String              tlsKey;
        protected boolean             tlsPeerNameVerify = DEFAULT_TLS_PEER_NAME_VERIFY;
        protected String              tlsPeerName;

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
            if (connectTimeoutMills <= 0
                    || connectTimeoutMills > Constants.DEFAULT_MAX_TIMEOUT_MS) {
                this.connectTimeoutMills = Constants.DEFAULT_MAX_TIMEOUT_MS;
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
            if (requestTimeoutMills <= 0
                    || requestTimeoutMills > Constants.DEFAULT_MAX_TIMEOUT_MS) {
                this.requestTimeoutMills = Constants.DEFAULT_MAX_TIMEOUT_MS;
            } else {
                this.requestTimeoutMills = requestTimeoutMills;
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
         * config if check the client when borrow from pool.
         *
         * @param testOnBorrow if ping when borrow client from pool
         * @return NebulaPool.Builder
         */
        public Builder withTestOnBorrow(boolean testOnBorrow) {
            this.testOnBorrow = testOnBorrow;
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
         * @param graph home graph name
         * @return NebulaPool.Builder
         */
        public Builder withGraph(String graph) {
            this.graph = graph;
            return this;
        }

        /**
         * config the initial ZonedId for NebulaClient in NebulaPool
         *
         * @param zoneId zone id
         * @return NebulaPool.Builder
         */
        public Builder withTimeZone(String zoneId) {
            this.timeZone = zoneId;
            return this;
        }

        /**
         * config the initial schema for NebulaClient in NebulaPool
         *
         * @param schema home schema path
         * @return NebulaPool.Builder
         */
        public Builder withSchema(String schema) {
            this.schema = schema;
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


        /**
         * config whether enable tls
         *
         * @param enableTls true if enable the tls
         * @return NebulaPool.Builder
         */
        public NebulaPool.Builder withEnableTls(boolean enableTls) {
            this.enableTls = enableTls;
            return this;
        }

        /**
         * config the ca certificate for TLS
         *
         * @param ca path to the trusted CA certificate file used to authenticate the server
         * @return NebulaPool.Builder
         */
        public NebulaPool.Builder withTlsCa(String ca) {
            this.tlsCa = ca;
            return this;
        }

        /**
         * config the TLS Cert options, necessary only if mTLS is enabled on Graph Server side
         *
         * @param cert certificate of client
         * @param key  private key of client certificate
         * @return NebulaPool.Builder
         */
        public NebulaPool.Builder withTlsCert(String cert, String key) {
            this.tlsCert = cert;
            this.tlsKey = key;
            return this;
        }


        /**
         * Peer name used to verify the CN or SAN,
         * hostname or IP will be used if empty or not specified
         *
         * @param tlsPeerName peer name
         * @return NebulaPool.Builder
         */
        public NebulaPool.Builder withTlsPeerName(String tlsPeerName) {
            this.tlsPeerName = tlsPeerName;
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
            if (enableTls && tlsCa == null) {
                throw new IllegalArgumentException("no CA certificate provide.");
            }
            if (enableTls && tlsPeerNameVerify && tlsPeerName == null) {
                throw new IllegalArgumentException(
                        "no peer name provide.");
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
