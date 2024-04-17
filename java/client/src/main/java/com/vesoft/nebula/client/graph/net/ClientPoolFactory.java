/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import java.io.Serializable;
import java.time.ZoneId;
import java.util.List;
import java.util.Map;
import org.apache.commons.pool2.BasePooledObjectFactory;
import org.apache.commons.pool2.PooledObject;
import org.apache.commons.pool2.impl.DefaultPooledObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ClientPoolFactory extends BasePooledObjectFactory<NebulaClient>
        implements Serializable {
    private final Logger logger = LoggerFactory.getLogger(this.getClass());

    private final LoadBalancer loadBalancer;
    private String userName;
    private Map<String, Object> authOptions;
    private long requestTimeout;
    private int retryTimes;
    private long intervalTimeMs;
    private int batchSize;
    private int scanParallel;
    private String workingGraph;
    private ZoneId zoneId;

    public ClientPoolFactory(
            LoadBalancer loadBalancer,
            String userName,
            Map<String, Object> authOptions,
            long requestTimeoutMs,
            int retryTimes,
            long intervalTimeMs,
            int batchSize,
            int scanParallel,
            String workingGraph,
            ZoneId zoneId) {
        this.loadBalancer = loadBalancer;
        this.userName = userName;
        this.authOptions = authOptions;
        this.requestTimeout = requestTimeoutMs;
        this.retryTimes = retryTimes;
        this.intervalTimeMs = intervalTimeMs;
        this.batchSize = batchSize;
        this.scanParallel = scanParallel;
        this.workingGraph = workingGraph;
        this.zoneId = zoneId;
    }


    @Override
    public NebulaClient create() throws IOErrorException, AuthFailedException {
        List<HostAddress> goodHosts = loadBalancer.getGoodAddresses();
        if (goodHosts.size() == 0) {
            throw new IOErrorException(IOErrorException.E_ALL_BROKEN, "No servers host is "
                    + "available, please check your servers is up and network between client and "
                    + "server is connected.");
        }

        GrpcConnection connection = new GrpcConnection();
        int tryConnect = goodHosts.size();
        while (tryConnect-- > 0) {
            try {
                connection.open(loadBalancer.getAddress(), requestTimeout);
                break;
            } catch (Exception e) {
                if (tryConnect == 0) {
                    throw e;
                } else {
                    logger.warn("connect failed, " + e.getMessage());
                }
            }
        }
        AuthResult authResult;
        try {
            authResult = connection.authenticate(userName, authOptions);
        } catch (AuthFailedException e) {
            logger.error(e.getMessage());
            throw e;
        }

        NebulaClient client = NebulaClient
                .builder(goodHosts.toString(), userName)
                .setAuthOptions(authOptions)
                .setRequestTimeoutMills(requestTimeout)
                .setRetryTimes(retryTimes)
                .setIntervalTimeMills(intervalTimeMs)
                .setScanParallel(scanParallel)
                .setWorkingGraph(workingGraph)
                .setTimeZone(zoneId)
                .build();
        return client;
    }

    @Override
    public PooledObject<NebulaClient> wrap(NebulaClient client) {
        return new DefaultPooledObject<>(client);
    }

    @Override
    public void destroyObject(PooledObject<NebulaClient> clientObject) throws Exception {
        NebulaClient client = clientObject.getObject();
        try {
            client.close();
        } catch (Exception e) {
            logger.warn("session release failed ", e);
        }
        super.destroyObject(clientObject);
    }

    @Override
    public boolean validateObject(PooledObject<NebulaClient> clientObject) {
        NebulaClient client = clientObject.getObject();
        return client.ping();
    }
}
