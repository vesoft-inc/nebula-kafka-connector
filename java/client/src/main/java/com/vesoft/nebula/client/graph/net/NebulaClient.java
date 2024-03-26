/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.google.common.net.InetAddresses;
import com.google.common.net.InternetDomainName;
import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.data.ValueWrapper;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import com.vesoft.nebula.client.graph.scan.ScanEdgeResultIterator;
import com.vesoft.nebula.client.graph.scan.ScanNodeResultIterator;
import java.io.Serializable;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.time.ZoneId;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
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
    private long timeoutMs;

    private ZoneId zoneId;

    // the default batch size for scan data, config for scan
    private final int defaultBatchSize = 10;

    // the default parallel for scan data, config for scan
    private int scanParallel;

    private ExecutorService threadPool = null;

    private NebulaClient(Builder builder) throws IOErrorException {
        if (hasInit.get()) {
            return;
        }
        this.retryTimes = builder.retryTimes;
        this.intervalTime = builder.intervalTime;
        this.maxWaitMills = builder.maxWaitMills < 0 ? Long.MAX_VALUE : builder.maxWaitMills;
        this.timeoutMs = builder.requestTimeout <= 0 ? Long.MAX_VALUE : builder.requestTimeout;
        this.scanParallel = builder.maxSessionSize;
        this.zoneId = builder.zoneId;

        this.loadBalancer = new RoundRobinLoadBalancer(builder.address, this.timeoutMs,
                builder.strictlyServerHealthy, builder.healthCheckTime);
        if (!loadBalancer.isServersOK()) {
            loadBalancer.close();
            log.error("servers status is not ok, please check the server status or network.");
            throw new IOErrorException(IOErrorException.E_SERVER_BAD, "Servers status is not ok.");
        }

        sessionPoolConfig = new SessionPoolConfig(builder.address, builder.userName,
                builder.authOptions)
                .setMaxSessionSize(builder.maxSessionSize)
                .setMinSessionSize(builder.minSessionSize)
                .setRequestTimeout(this.timeoutMs)
                .setRetryTimes(this.retryTimes)
                .setIntervalTime(this.intervalTime)
                .setHealthCheckTime(builder.healthCheckTime)
                .setBlockWhenExhausted(builder.blockWhenExhausted)
                .setMaxWaitMills(this.maxWaitMills)
                .setIdleEvictScheduleMills(builder.idleEvictScheduleMills)
                .setMinEvictableIdleTimeMillis(builder.minEvictableIdleTimeMillis)
                .setStrictlyServerHealthy(builder.strictlyServerHealthy);
        // pool config
        GenericObjectPoolConfig objConfig = new GenericObjectPoolConfig();
        objConfig.setMaxIdle(builder.maxSessionSize);
        objConfig.setMinIdle(builder.minSessionSize);
        objConfig.setMaxTotal(builder.maxSessionSize);
        objConfig.setBlockWhenExhausted(builder.blockWhenExhausted);
        objConfig.setMaxWaitMillis(this.maxWaitMills);
        objConfig.setTimeBetweenEvictionRunsMillis(builder.idleEvictScheduleMills);
        objConfig.setMinEvictableIdleTimeMillis(builder.minEvictableIdleTimeMillis);
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


    public ResultSet execute(String stmt) throws IOErrorException, NoValidSessionException {
        return execute(stmt, this.timeoutMs);
    }

    /**
     * execute the graph statement
     * There are four cases for execution result:
     * 1. execution succeed or execution failed with the syntax and semantics error,
     * then the session will be put back into the pool and the result will be returned.
     * 2. execution failed with SESSION error, then the session will be invalidated
     * and retried the execution.
     * 3. execution failed with other type error, then the session will be put
     * back into the pool and tried the execution.
     * 4. execution exception for IOException, then the session will be invalidated
     * and retried the execution.
     *
     * @param stmt graph statement
     * @return {@link ResultSet}
     */
    public ResultSet execute(String stmt, long timeoutMs)
            throws IOErrorException, NoValidSessionException {
        checkClosed();
        int tryTimes = 0;
        boolean isBadSession = false;
        ResultSet resultSet = null;

        // execute times will be (retryTimes + 1)
        while (tryTimes++ < retryTimes + 1) {
            if (tryTimes > 1) {
                sleep();
            }
            Session session = getSession();
            try {
                resultSet = session.execute(stmt, timeoutMs);
                if (resultSet.isSucceeded()
                        || resultSet.getErrorCode().isSemanticError()
                        || resultSet.getErrorCode().isSyntaxError()) {
                    return resultSet;
                }
                if (resultSet.getErrorCode().isSessionError()) {
                    isBadSession = true;
                }
                log.warn(String.format("execute error for times %s,  message: %s", tryTimes + 1,
                        resultSet.getErrorMessage()));
                if (tryTimes <= retryTimes) {
                    log.info("now retry the execute...");
                }
            } catch (IOErrorException e) {
                isBadSession = true;
                loadBalancer.updateServersStatus();
                // still has retry time.
                if (tryTimes <= retryTimes) {
                    log.warn(String.format("execute failed for IOErrorException, "
                            + "message: %s, tryTime: %d", e.getMessage(), tryTimes));
                } else {
                    // retry time is exhausted
                    throw e;
                }
            } finally {
                if (isBadSession) {
                    try {
                        pool.invalidateObject(session);
                    } catch (Exception e) {
                        log.warn("invalidate session failed", e);
                    }
                } else {
                    pool.returnObject(session);
                }
            }
        }
        return resultSet;
    }

    // not implemented.
    public ScanNodeResultIterator scanNode(String graphName, String nodeType) {
        List<Integer> parts = getAllParts();
        return scanNode(graphName, nodeType, new ArrayList<>(), true, parts, defaultBatchSize);
    }

    /**
     * scan the data of specific nodeType.
     * the result will contain all property of nodeType.
     *
     * @param graphName NebulaGraph name
     * @param nodeType  node type name
     * @param part      graph part id
     * @param batchSize the data size for one request for one part,
     *                  the ScanNodeResultIterator.next() will return at most batchSize
     *                  node records
     * @return ScanNodeResultIterator
     */
    public ScanNodeResultIterator scanNode(String graphName,
                                           String nodeType,
                                           int part,
                                           int batchSize) {
        return scanNode(graphName, nodeType, new ArrayList<>(), true,
                Collections.singletonList(part), batchSize);
    }

    /**
     * scan the data of specific nodeType
     * the result will contain primary key and specific return properties.
     *
     * @param graphName        NebulaGraph name
     * @param nodeType         node type name
     * @param returnProperties the property list to scan, if list is empty,
     *                         then the result just contain primary key.
     * @param part             graph part id
     * @param batchSize        the data size for one request for one part,
     *                         the ScanNodeResultIterator.next() will return at most batchSize
     *                         node records
     * @return ScanNodeResultIterator
     */
    public ScanNodeResultIterator scanNode(String graphName,
                                           String nodeType,
                                           List<String> returnProperties,
                                           int part,
                                           int batchSize) {
        return scanNode(graphName, nodeType, returnProperties,
                Collections.singletonList(part), batchSize);
    }

    /**
     * scan the data of specific nodeType
     * the result will contain primary key and specific return properties.
     *
     * @param graphName        NebulaGraph name
     * @param nodeType         node type name
     * @param returnProperties the property list to scan, if list is empty,
     *                         then the result just contain primary key.
     * @param parts            list of graph part id
     * @param batchSize        the data size for one request for one part,
     *                         the ScanNodeResultIterator.next() will return at most batchSize
     *                         node records
     * @return ScanNodeResultIterator
     */
    public ScanNodeResultIterator scanNode(String graphName,
                                           String nodeType,
                                           List<String> returnProperties,
                                           List<Integer> parts,
                                           int batchSize) {
        boolean allProperties = false;
        if (returnProperties == null) {
            allProperties = true;
        }
        return scanNode(graphName, nodeType, returnProperties, allProperties, parts, batchSize);
    }

    /**
     * scan the data of specific nodeType
     * the result will contain primary key and specific return properties.
     *
     * @param graphName        NebulaGraph name
     * @param nodeType         node type name
     * @param returnProperties the property list to scan, if list is empty,
     *                         then the result just contain primary key.
     * @param allProperties    if return all the properties of node type
     * @param parts            part list to scan
     * @param batchSize        the data size for one request for one part,
     *                         the ScanNodeResultIterator.next() will return at most batchSize
     *                         node records
     * @return ScanNodeResultIterator
     */
    private ScanNodeResultIterator scanNode(String graphName,
                                            String nodeType,
                                            List<String> returnProperties,
                                            boolean allProperties,
                                            List<Integer> parts,
                                            int batchSize) {
        initScanThreadPool();
        // get node type's all property names
        List<String> nodeProperties = null;
        try {
            nodeProperties = getNodeProperties(graphName, nodeType);
        } catch (IOErrorException | NoValidSessionException e) {
            log.error("get node schema failed.", e);
            throw new RuntimeException(e);
        }

        // construct the return property list for scan
        List<String> propertyList = new ArrayList<>();
        if (allProperties) {
            propertyList.addAll(nodeProperties);
        } else {
            String primaryKey = nodeProperties.get(0);
            // put the primary key always on the head of the propertyList for scan
            propertyList.add(primaryKey);
            for (String propName : returnProperties) {
                if (propName.trim().equals(primaryKey)) {
                    continue;
                }
                propertyList.add(propName);
            }
        }

        return new ScanNodeResultIterator(pool, graphName, nodeType, propertyList,
                parts, batchSize, threadPool, retryTimes, intervalTime, timeoutMs);
    }


    // not implemented.
    public ScanEdgeResultIterator scanEdge(String graphName, String edgeType) {
        List<Integer> parts = getAllParts();
        return scanEdge(graphName, edgeType, new ArrayList<>(), true, parts, defaultBatchSize);
    }


    /**
     * scan the data of specific edgeType
     * the result will contain src node's primary key, dst node's primary key, edge's all property.
     *
     * @param graphName NebulaGraph name
     * @param edgeType  edge type name
     * @param part      graph part id
     * @param batchSize the data size for one request for one part,
     *                  the ScanEdgeResultIterator.next() will return at most batchSize
     *                  edge records
     * @return ScanEdgeResultIterator
     */
    public ScanEdgeResultIterator scanEdge(String graphName,
                                           String edgeType,
                                           int part,
                                           int batchSize) {
        return scanEdge(graphName, edgeType, new ArrayList<>(), true,
                Collections.singletonList(part), batchSize);
    }


    /**
     * scan the data of specific edgeType
     *
     * @param graphName        NebulaGraph name
     * @param edgeType         edge type name
     * @param returnProperties the property list to scan, if list is empty, then the result will
     *                         just contain src node's primary key and dst node's primary key
     * @param part             graph part id
     * @param batchSize        the data size for one request for one part,
     *                         the ScanEdgeResultIterator.next() will return at most batchSize
     *                         edge records
     * @return ScanEdgeResultIterator
     */
    public ScanEdgeResultIterator scanEdge(String graphName,
                                           String edgeType,
                                           List<String> returnProperties,
                                           int part,
                                           int batchSize) {
        boolean allProperties = false;
        if (returnProperties == null) {
            allProperties = true;
        }
        return scanEdge(graphName, edgeType, returnProperties, allProperties,
                Collections.singletonList(part),
                batchSize);
    }


    /**
     * scan the data of specific edgeType
     *
     * @param graphName        NebulaGraph name
     * @param edgeType         edge type name
     * @param returnProperties the property list to scan, if list is empty, then the result will
     *                         just contain src node's primary key and dst node's primary key
     * @param parts            part list to scan
     * @param batchSize        the data size for one request for one part,
     *                         the ScanEdgeResultIterator.next() will return at most batchSize
     *                         edge records
     * @return ScanEdgeResultIterator
     */
    public ScanEdgeResultIterator scanEdge(String graphName,
                                           String edgeType,
                                           List<String> returnProperties,
                                           List<Integer> parts,
                                           int batchSize) {
        initScanThreadPool();
        boolean allProperties = false;
        if (returnProperties == null) {
            allProperties = true;
        }
        return scanEdge(graphName, edgeType, returnProperties, false, parts, batchSize);
    }


    /**
     * scan the data of specific edgeType
     *
     * @param graphName        NebulaGraph name
     * @param edgeType         edge type name
     * @param returnProperties the property list to scan, if list is empty, then the result will
     *                         just contain src node's primary key and dst node's primary key
     * @param allProperties    if return all the properties of edge type
     * @param parts            part list to scan
     * @param batchSize        the data size for one request for one part,
     *                         the ScanEdgeResultIterator.next() will return at most batchSize
     *                         edge records
     * @return ScanEdgeResultIterator
     */
    private ScanEdgeResultIterator scanEdge(String graphName,
                                            String edgeType,
                                            List<String> returnProperties,
                                            boolean allProperties,
                                            List<Integer> parts,
                                            int batchSize) {
        initScanThreadPool();
        // get edge type's all property names
        List<String> edgeProperties = null;
        try {
            edgeProperties = getEdgeProperties(graphName, edgeType);
        } catch (IOErrorException | NoValidSessionException e) {
            log.error("get node schema failed.", e);
            throw new RuntimeException(e);
        }

        // construct the return property list for scan
        List<String> propertyList = new ArrayList<>();
        if (allProperties) {
            propertyList.addAll(edgeProperties);
        } else {
            propertyList.addAll(returnProperties);
        }
        return new ScanEdgeResultIterator(pool, graphName, edgeType, propertyList,
                parts, batchSize, threadPool, retryTimes, intervalTime, timeoutMs);
    }


    /**
     * get all parts
     *
     * @return list of part id
     */
    private List<Integer> getAllParts() {
        String showPartitions = "CALL show_partitions() RETURN *";
        ResultSet resultSet;
        try {
            resultSet = execute(showPartitions);
        } catch (Exception e) {
            log.error("get all partitions error", e);
            throw new RuntimeException("get all partitions error", e);
        }
        if (!resultSet.isSucceeded() || resultSet.isEmpty()) {
            log.error("get all partitions failed for {}", resultSet.getErrorMessage());
            throw new RuntimeException(
                    "get all partitions failed for " + resultSet.getErrorMessage());
        }
        List<ValueWrapper> partitionsValue = resultSet.next().values().get(0).asList();
        List<Integer> partitions = new ArrayList<>();
        for (ValueWrapper part : partitionsValue) {
            partitions.add(part.asInt());
        }
        return partitions;
    }


    /**
     * get node type's properties
     *
     * @param graphName NebulaGraph name
     * @param nodeType  node type name
     * @return List of property name
     */
    private List<String> getNodeProperties(String graphName, String nodeType)
            throws IOErrorException, NoValidSessionException {
        String graphType = getGraphType(graphName);
        String descNodeType = String.format("DESCRIBE NODE TYPE %s OF %s", nodeType, graphType);
        ResultSet resultSet = execute(descNodeType);
        if (!resultSet.isSucceeded() || resultSet.isEmpty()) {
            log.error(String.format("get description of %s failed for %s", nodeType,
                    resultSet.getErrorMessage()));
            throw new IllegalArgumentException(String.format("node type %s does not exist in %s",
                    nodeType, graphName));
        }

        List<ValueWrapper> pks = resultSet.next().get("primary_keys").asList();
        if (pks.isEmpty()) {
            log.error("node type " + nodeType + " has no primary key.");
            throw new RuntimeException("node type " + nodeType + " has no primary key");
        }
        String pk = pks.get(0).asString();

        // define the property name list, and put the pk on the head of list.
        List<String> propertyNames = new ArrayList<>();
        propertyNames.add(pk);
        List<ValueWrapper> properties = resultSet.next().get("properties").asList();
        for (ValueWrapper property : properties) {
            String propertyName = property.asString().split(":")[0];
            if (pk.equals(propertyName)) {
                continue;
            }
            propertyNames.add(propertyName);
        }
        return propertyNames;
    }


    /**
     * get edge type's properties
     *
     * @param graphName NebulaGraph name
     * @param edgeType  edge type name
     * @return List of property name
     */
    private List<String> getEdgeProperties(String graphName, String edgeType)
            throws IOErrorException, NoValidSessionException {

        String graphType = getGraphType(graphName);

        String descEdgeType = String.format("DESCRIBE EDGE TYPE %s OF %s", edgeType, graphType);
        ResultSet resultSet = execute(descEdgeType);
        if (!resultSet.isSucceeded() || resultSet.isEmpty()) {
            log.error(String.format("get description of %s failed for %s", edgeType,
                    resultSet.getErrorMessage()));
            throw new IllegalArgumentException(String.format("edge type %s does not exist in %s",
                    edgeType, graphName));
        }

        List<ValueWrapper> properties = resultSet.next().get("properties").asList();
        List<String> propertyNames = new ArrayList<>();
        for (ValueWrapper property : properties) {
            String propertyName = property.asString().split(":")[0];
            propertyNames.add(propertyName);
        }
        return propertyNames;
    }


    /**
     * get the graph's graph type
     * TODO add `` for graph type and graph name
     *
     * @param graphName NebulaGraph name
     * @return String
     */
    private String getGraphType(String graphName) throws IOErrorException,
            NoValidSessionException {
        ResultSet resultSet = execute(String.format("DESCRIBE GRAPH %s", graphName));
        String graphType;
        if (resultSet.isSucceeded() && !resultSet.isEmpty()) {
            graphType = resultSet.next().values().get(1).asString();
        } else {
            throw new IllegalArgumentException("graphName " + graphName + " does not exist.");
        }
        return graphType;
    }

    /**
     * init the thread pool for scan.
     * The max thread size is the value of maxSessionSize for pool
     * When the number of parts to be scanned is greater than maxSessionSize, the maximum
     * concurrency will be the maximum number of sessions that can be executed concurrently.
     * When the number of parts is less than maxSessionSize, the upper limit of the thread pool is
     * maxSessionSize. Threads will only be created when a task is submitted, so in the pool,
     * Only parts number of threads will be created.
     */
    private void initScanThreadPool() {
        if (threadPool == null) {
            synchronized (this) {
                if (threadPool == null) {
                    threadPool = Executors.newFixedThreadPool(scanParallel);
                }
            }
        }
    }

    /**
     * close the client
     */
    public void close() {
        if (!isClosed.get() && isClosed.compareAndSet(false, true)) {
            loadBalancer.close();
            if (pool != null && !pool.isClosed()) {
                pool.close();
            }

            GrpcConnection.closeChannel();
            if (threadPool != null && !threadPool.isShutdown()) {
                threadPool.shutdown();
            }
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
            Session session = pool.borrowObject(maxWaitMills);
            if (zoneId != null) {
                ResultSet resultSet = session.execute(
                        "SESSION SET TIME ZONE \"" + zoneId + "\"", this.timeoutMs);
                if (!resultSet.isSucceeded()) {
                    log.error("failed to set time zone for {}", resultSet.getErrorMessage());
                    throw new RuntimeException("failed to set timezone for "
                            + resultSet.getErrorMessage());
                }
            }
            return session;
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

        private Map<String, Object> authOptions = new HashMap<>();
        // define default value

        // The max sessions in pool
        private int maxSessionSize = 10;

        // The min sessions in pool
        private int minSessionSize = 1;

        // socket timeout for request, unit: millisecond
        private int requestTimeout = 0;

        // retry times for failed execute
        private int retryTimes = 0;

        // interval time for retry, unit: millisecond
        private int intervalTime = 0;

        // The healthCheckTime for schedule check the health of session, unit: millisecond
        private int healthCheckTime = 300000;

        // if block when session is exhausted, if false, throw exception.
        private boolean blockWhenExhausted = false;

        // the max wait time if blockWhenExhausted is true. if value is less than 0, always wait.
        // unit: millisecond
        private long maxWaitMills = -1;

        // the schedule time for test the idle session and evict it. if value is less than 0,
        // never evict the idle sessions.
        private long idleEvictScheduleMills = -1;

        // the min idle time for idle session
        private long minEvictableIdleTimeMillis = 1000L * 60L * 30L;

        // if need all servers are strictly healthy.
        // if true, all addresses must be available, if false, at least one address is available.
        private boolean strictlyServerHealthy = false;

        // the time zone, used to parse ZonedTime and ZonedDatetime
        private ZoneId zoneId = null;


        public Builder(String address, String userName, String password)
                throws UnknownHostException {
            this.address = validateAddress(address);
            this.userName = userName;
            this.password = password;
        }

        public Builder(String address, String userName)
                throws UnknownHostException {
            this.address = validateAddress(address);
            this.userName = userName;
            this.password = null;
        }

        public Builder setAuthOptions(Map<String, Object> authOptions) {
            if (authOptions != null) {
                this.authOptions.putAll(authOptions);
            }
            return this;
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

        public Builder setIdleEvictScheduleMills(long idleEvictScheduleMills) {
            this.idleEvictScheduleMills = idleEvictScheduleMills;
            return this;
        }

        public Builder setMinEvictableIdleTimeMillis(long minEvictableIdleTimeMillis) {
            this.minEvictableIdleTimeMillis = minEvictableIdleTimeMillis;
            return this;
        }

        public Builder setStrictlyServerHealthy(boolean strictlyServerHealthy) {
            this.strictlyServerHealthy = strictlyServerHealthy;
            return this;
        }

        public Builder setTimeZone(ZoneId zoneId) {
            this.zoneId = zoneId;
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
            if (password != null) {
                this.authOptions.put("password", password);
            }
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
