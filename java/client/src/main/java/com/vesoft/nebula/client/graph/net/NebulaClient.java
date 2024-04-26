/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_BATCH_SIZE;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_CONNECT_TIMEOUT;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_REQUEST_TIMEOUT;
import static com.vesoft.nebula.client.graph.net.Constants.DEFAULT_SCAN_PARALLEL;

import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.data.ValueWrapper;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.scan.ScanEdgeResultIterator;
import com.vesoft.nebula.client.graph.scan.ScanNodeResultIterator;
import com.vesoft.nebula.client.graph.utils.AddressUtil;
import java.io.Serializable;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.stream.Collectors;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 *
 */
public class NebulaClient implements Serializable {
    private final Logger logger = LoggerFactory.getLogger(this.getClass());
    private List<HostAddress> servers;
    private long connectTimeout;
    private long requestTimeout;
    private final String userName;
    private final Map<String, Object> authOptions;
    private GrpcConnection connection;
    private long sessionId;

    private int scanParallel;

    private boolean isClosed = false;

    public static Builder builder(String addresses, String userName, String password) {
        return new Builder(addresses, userName, password);
    }

    public static Builder builder(String addresses, String userName) {
        return new Builder(addresses, userName, null);
    }

    private NebulaClient(NebulaClient.Builder builder) throws AuthFailedException {
        this.servers = builder.address;
        this.userName = builder.userName;
        this.authOptions = builder.authOptions;
        this.connectTimeout = builder.connectTimeoutMills;
        this.requestTimeout = builder.requestTimeoutMills;
        this.scanParallel = builder.scanParallel;
        initClient();
    }


    public ResultSet execute(String gql) throws IOErrorException {
        return execute(gql, requestTimeout);
    }

    public synchronized ResultSet execute(String gql, long requestTimeout) throws IOErrorException {
        checkClosed();
        return new ResultSet(connection.execute(sessionId, gql, requestTimeout));
    }

    /**
     * get the SessionId of the Client
     *
     * @return session id
     */
    public long getSessionId() {
        return sessionId;
    }

    /**
     * get the NebulaGraph host address that this client is connecting to.
     *
     * @return NebulaGraph server host.
     */
    public String getHost() {
        return connection.getServerAddress().toString();
    }

    /**
     * ping the NebulaGraph server
     *
     * @return true when client can ping server succeed.
     */
    public synchronized boolean ping() {
        checkClosed();
        try {
            return connection.ping(sessionId);
        } catch (IOErrorException e) {
            logger.error("ping error for host {}", getHost(), e);
            return false;
        }
    }

    /**
     * release and close the connection with NebulaGraph server.
     */
    public synchronized void close() {
        if (!isClosed) {
            isClosed = true;
            if (connection != null) {
                try {
                    connection.execute(sessionId, "SESSION CLOSE");
                } catch (Exception e) {
                    logger.warn("signout failed,", e);
                }
            }
            connection = null;
        }
    }

    /**
     * get the Client status
     *
     * @reuturn true if client is closed
     */
    public boolean isClosed() {
        return isClosed;
    }

    /**
     * check if the client already closed
     *
     * @throws RuntimeException if the client is closed.
     */
    private void checkClosed() {
        if (isClosed) {
            throw new RuntimeException("The NebulaClient already closed.");
        }
    }

    private void initClient() throws AuthFailedException {
        // create connection with NebulaGraph Server
        connection = new GrpcConnection();
        int tryConnectTimes = servers.size();
        while (tryConnectTimes-- > 0) {
            try {
                // TODO polling the address
                Collections.shuffle(servers);
                connection.open(servers.get(tryConnectTimes), connectTimeout, requestTimeout);
                break;
            } catch (Exception e) {
                if (tryConnectTimes == 0) {
                    logger.error("create NebulaClient failed.", e);
                    throw e;
                }
            }
        }

        // auth for user
        AuthResult authResult;
        try {
            authResult = connection.authenticate(userName, authOptions);
        } catch (AuthFailedException e) {
            logger.error("create NebulaClient failed.", e);
            throw e;
        }
        sessionId = authResult.getSessionId();
    }


    public ScanNodeResultIterator scanNode(String graphName, String nodeType) {
        List<Integer> parts = getAllParts();
        return scanNode(graphName, nodeType, new ArrayList<>(), true, parts, DEFAULT_BATCH_SIZE);
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
        // get node type's all property names
        List<String> nodeProperties = null;
        try {
            nodeProperties = getNodeProperties(graphName, nodeType);
        } catch (IOErrorException e) {
            logger.error("get node schema failed.", e);
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

        return new ScanNodeResultIterator(graphName, nodeType, propertyList,
                parts, batchSize, scanParallel, servers, userName, authOptions, requestTimeout);
    }


    public ScanEdgeResultIterator scanEdge(String graphName, String edgeType) {
        List<Integer> parts = getAllParts();
        return scanEdge(graphName, edgeType, new ArrayList<>(), true, parts, DEFAULT_BATCH_SIZE);
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
        return scanEdge(
                graphName,
                edgeType,
                returnProperties,
                allProperties,
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
        // get edge type's all property names
        List<String> edgeProperties = null;
        try {
            edgeProperties = getEdgeProperties(graphName, edgeType);
        } catch (IOErrorException e) {
            logger.error("get node schema failed.", e);
            throw new RuntimeException(e);
        }

        // construct the return property list for scan
        List<String> propertyList = new ArrayList<>();
        if (allProperties) {
            propertyList.addAll(edgeProperties);
        } else {
            propertyList.addAll(returnProperties);
        }
        return new ScanEdgeResultIterator(graphName, edgeType, propertyList,
                parts, batchSize, scanParallel, servers, userName, authOptions, requestTimeout);
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
            logger.error("get all partitions error", e);
            throw new RuntimeException("get all partitions error", e);
        }
        if (!resultSet.isSucceeded() || resultSet.isEmpty()) {
            logger.error("get all partitions failed for {}", resultSet.getErrorMessage());
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
            throws IOErrorException {
        String graphType = getGraphType(graphName);
        String descNodeType = String.format("DESCRIBE NODE TYPE %s OF %s", nodeType, graphType);
        ResultSet resultSet = execute(descNodeType);
        if (!resultSet.isSucceeded() || resultSet.isEmpty()) {
            logger.error(String.format("get description of %s failed for %s", nodeType,
                    resultSet.getErrorMessage()));
            throw new IllegalArgumentException(String.format("node type %s does not exist in %s",
                    nodeType, graphName));
        }

        List<String> pks = new ArrayList<>();
        List<String> propNames = new ArrayList<>();
        while (resultSet.hasNext()) {
            ResultSet.Record record = resultSet.next();
            propNames.add(record.get("property_name").asString());
            if ("Y".equals(record.get("primary_key").asString())) {
                pks.add(record.get("property_name").asString());
            }
        }

        if (pks.isEmpty()) {
            logger.error("node type " + nodeType + " has no primary key.");
            throw new RuntimeException("node type " + nodeType + " has no primary key");
        }

        // define the property name list, and put the pk on the head of list.
        List<String> propertyNames = new ArrayList<>(pks);
        for (String property : propNames) {
            if (pks.contains(property)) {
                continue;
            }
            propertyNames.add(property);
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
            throws IOErrorException {

        String graphType = getGraphType(graphName);

        String descEdgeType = String.format(
                "CALL describe_graph_type(\"%s\") FILTER type_name=\"%s\" return properties",
                graphType, edgeType);
        ResultSet resultSet = execute(descEdgeType);
        if (!resultSet.isSucceeded() || resultSet.isEmpty()) {
            logger.error(String.format("get description of %s failed for %s", edgeType,
                    resultSet.getErrorMessage()));
            throw new IllegalArgumentException(String.format("edge type %s does not exist in %s",
                    edgeType, graphName));
        }

        List<ValueWrapper> properties = resultSet.next().get("properties").asList();
        return properties.stream().map(ValueWrapper::asString).collect(Collectors.toList());
    }


    /**
     * get the graph's graph type
     * TODO add `` for graph type and graph name
     *
     * @param graphName NebulaGraph name
     * @return String
     */
    private String getGraphType(String graphName) throws IOErrorException {
        ResultSet resultSet = execute(String.format("DESCRIBE GRAPH %s", graphName));
        String graphType;
        if (resultSet.isSucceeded() && !resultSet.isEmpty()) {
            graphType = resultSet.next().values().get(1).asString();
        } else {
            throw new IllegalArgumentException("graphName " + graphName + " does not exist.");
        }
        return graphType;
    }

    public static class Builder {
        private final List<HostAddress> address;
        private final String userName;
        private final String password;
        private Map<String, Object> authOptions = new HashMap<>();
        private long connectTimeoutMills = DEFAULT_CONNECT_TIMEOUT;
        private long requestTimeoutMills = DEFAULT_REQUEST_TIMEOUT;
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

        public Builder withConnectTimeoutMills(long connectTimeoutMills) {
            this.connectTimeoutMills =
                    connectTimeoutMills <= 0 ? Long.MAX_VALUE : connectTimeoutMills;
            return this;
        }

        public Builder withRequestTimeoutMills(long requestTimeoutMills) {
            this.requestTimeoutMills =
                    requestTimeoutMills <= 0 ? Long.MAX_VALUE : requestTimeoutMills;
            return this;
        }

        public Builder withScanParallel(int parallel) {
            this.scanParallel = parallel;
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

        public NebulaClient build() throws AuthFailedException, IOErrorException {
            check();
            if (password != null) {
                authOptions.put("password", password);
            }
            return new NebulaClient(this);
        }

    }
}
