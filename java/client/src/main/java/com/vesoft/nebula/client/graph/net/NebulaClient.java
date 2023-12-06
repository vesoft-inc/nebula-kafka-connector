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
import com.vesoft.nebula.client.graph.scan.ScanEdgeResultIterator;
import com.vesoft.nebula.client.graph.scan.ScanNodeResultIterator;
import java.io.Serializable;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
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
        this.maxWaitMills = builder.maxWaitMills;
        this.scanParallel = builder.maxSessionSize;

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
                .setIdleEvictScheduleMills(builder.idleEvictScheduleMills)
                .setMinEvictableIdleTimeMillis(builder.minEvictableIdleTimeMillis)
                .setStrictlyServerHealthy(builder.strictlyServerHealthy);
        // pool config
        GenericObjectPoolConfig objConfig = new GenericObjectPoolConfig();
        objConfig.setMaxIdle(builder.maxSessionSize);
        objConfig.setMinIdle(builder.minSessionSize);
        objConfig.setMaxTotal(builder.maxSessionSize);
        objConfig.setBlockWhenExhausted(builder.blockWhenExhausted);
        objConfig.setMaxWaitMillis(builder.maxWaitMills);
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
    public ResultSet execute(String stmt) throws IOErrorException, NoValidSessionException {
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
                resultSet = session.execute(stmt);
                if (resultSet.isSucceeded() || "E_SEMANTIC_ERROR".equals(resultSet.getGqlStatus())
                        || "E_SYNTAX_ERROR".equals(resultSet.getGqlStatus())) {
                    return resultSet;
                }
                if ("E_SESSION_NOT_FOUND".equals(resultSet.getGqlStatus())
                        || "E_SESSION_INVALID".equals(resultSet.getGqlStatus())
                        || "E_SESSION_TIMEOUT".equals(resultSet.getGqlStatus())) {
                    isBadSession = true;
                }
                log.warn(String.format("execute error,  message: %s, retry: %d",
                        resultSet.getGqlStatus(), tryTimes));
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
    public ScanNodeResultIterator scanNode(String graphName,
                                           String nodeType) {
        Map<String, String> schema = null;
        try {
            schema = getNodeProperties(graphName, nodeType);
        } catch (IOErrorException | NoValidSessionException e) {
            log.error("get node schema failed.", e);
            throw new RuntimeException(e);
        }
        List<String> propertyList = new ArrayList<>(schema.keySet());
        // TODO get all parts
        List<Integer> parts = new ArrayList<>();
        return scanNode(graphName, nodeType, propertyList, parts, defaultBatchSize);
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
        Map<String, String> schema = null;
        try {
            schema = getNodeProperties(graphName, nodeType);
        } catch (IOErrorException | NoValidSessionException e) {
            log.error("get node schema failed.", e);
            throw new RuntimeException(e);
        }
        List<String> propertyList = new ArrayList<>(schema.keySet());
        return scanNode(graphName, nodeType, propertyList, part, batchSize);
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
        List<String> propertyList = new ArrayList<>();
        // TODO return primarykey by default, need a gql to get node type's primary key
        String primaryKey = null;
        // if returnProperties is empty, just return the value of primary key.
        // propertyList.add(primaryKey);

        // remove the primaryKey in parameter returnProperties
        // to keep the primary key in the first column of propertyList
        returnProperties.remove(primaryKey);
        propertyList.addAll(returnProperties);
        return scanNode(graphName, nodeType, propertyList,
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
     * @param parts            part list to scan
     * @param batchSize        the data size for one request for one part,
     *                         the ScanNodeResultIterator.next() will return at most batchSize
     *                         node records
     * @return ScanNodeResultIterator
     */
    private ScanNodeResultIterator scanNode(String graphName, String nodeType,
                                            List<String> returnProperties,
                                            List<Integer> parts, int batchSize) {
        initScanThreadPool();
        return new ScanNodeResultIterator(pool, graphName, nodeType, returnProperties,
                parts, batchSize, threadPool, retryTimes, intervalTime);
    }


    // not implemented.
    public ScanEdgeResultIterator scanEdge(String graphName, String edgeType) {
        // get all property of edge
        Map<String, String> schema = null;
        try {
            schema = getEdgeProperties(graphName, edgeType);
        } catch (IOErrorException | NoValidSessionException e) {
            log.error("get node schema failed.", e);
            throw new RuntimeException(e);
        }
        List<String> propertyList = new ArrayList<>(schema.keySet());
        // TODO get all parts
        List<Integer> parts = new ArrayList<>();
        return scanEdge(graphName, edgeType, propertyList, parts, defaultBatchSize);
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
        Map<String, String> schema = null;
        try {
            schema = getEdgeProperties(graphName, edgeType);
        } catch (IOErrorException | NoValidSessionException e) {
            log.error("get node schema failed.", e);
            throw new RuntimeException(e);
        }
        List<String> propertyList = new ArrayList<>(schema.keySet());
        return scanEdge(graphName, edgeType, propertyList, part, batchSize);
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
        return scanEdge(graphName, edgeType, returnProperties, Collections.singletonList(part),
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
    private ScanEdgeResultIterator scanEdge(String graphName,
                                            String edgeType,
                                            List<String> returnProperties,
                                            List<Integer> parts,
                                            int batchSize) {
        initScanThreadPool();
        return new ScanEdgeResultIterator(pool, graphName, edgeType, returnProperties,
                parts, batchSize, threadPool, retryTimes, intervalTime);
    }

    /**
     * get node type's properties
     *
     * @param graphName NebulaGraph name
     * @param nodeType  node type name
     * @return Map for property name and property data type
     */
    private Map<String, String> getNodeProperties(String graphName, String nodeType)
            throws IOErrorException, NoValidSessionException {
        Map<String, String> schema = new HashMap<>();
        ResultSet result = getGraphDesc(graphName);
        List<ResultSet.Record> records = result.getRows();

        for (ResultSet.Record record : records) {
            if (record.get("Field").asString().equalsIgnoreCase(nodeType)) {
                String propertyString = record.get("Properties").asString();
                String[] proeprties =
                        propertyString.substring(1, propertyString.length() - 1).split(",");
                for (String prop : proeprties) {
                    String[] nameAndType = prop.trim().split(" ");
                    schema.put(nameAndType[0], nameAndType[1]);
                }
                return schema;
            }
        }
        throw new IllegalArgumentException("node type " + nodeType + " does not exist!");
    }


    /**
     * get edge type's properties
     *
     * @param graphName NebulaGraph name
     * @param edgeType  edge type name
     * @return Map of property name and property data type
     */
    private Map<String, String> getEdgeProperties(String graphName, String edgeType)
            throws IOErrorException, NoValidSessionException {
        Map<String, String> schema = new HashMap<>();
        ResultSet result = getGraphDesc(graphName);
        List<ResultSet.Record> records = result.getRows();

        for (ResultSet.Record record : records) {
            if (record.get("Kind").asString().equals("Edge")) {
                String fullEdgeType = record.get("Field").asString();
                String regex = "\\((.*)\\)-\\[(.*)\\]->\\((.*)\\)";
                Pattern r = Pattern.compile(regex);
                Matcher m = r.matcher(fullEdgeType);
                if (m.find()) {
                    if (edgeType.equalsIgnoreCase(m.group(2))) {
                        String propertiesString = record.get("Properties").asString();
                        String[] properties = propertiesString.substring(1,
                                propertiesString.length() - 1).split(",");
                        for (String prop : properties) {
                            String[] nameAndType = prop.trim().split(" ");
                            schema.put(nameAndType[0], nameAndType[1]);
                        }
                        return schema;
                    }
                }
            }
        }
        throw new IllegalArgumentException("edgeType " + edgeType + " does not exist.");
    }


    /**
     * get the graph's schema info
     * TODO add `` for graph type and graph name
     *
     * @param graphName NebulaGraph name
     * @return ResultSet
     */
    private ResultSet getGraphDesc(String graphName) throws IOErrorException,
            NoValidSessionException {
        ResultSet resultSet = execute(String.format("DESCRIBE GRAPH %s", graphName));
        String graphType;
        if (resultSet.isSucceeded() && !resultSet.isEmpty()) {
            graphType = resultSet.getRows().get(0).values().get(1).asString();
        } else {
            throw new IllegalArgumentException("graphName " + graphName + " does not exist.");
        }

        String queryStatement = String.format("DESCRIBE GRAPH TYPE %s", graphType);
        resultSet = execute(queryStatement);
        if (!resultSet.isSucceeded()) {
            throw new RuntimeException("query error with " + queryStatement
                    + " for " + resultSet.getGqlStatus());
        }
        return resultSet;
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

        // the schedule time for test the idle session and evict it. if value is less than 0,
        // never evict the idle sessions.
        private long idleEvictScheduleMills = -1;

        // the min idle time for idle session
        private long minEvictableIdleTimeMillis = 1000L * 60L * 30L;

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
