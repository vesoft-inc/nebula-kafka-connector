package com.vesoft.nebula.driver.graph.net;

import com.vesoft.nebula.driver.graph.data.HostAddress;
import com.vesoft.nebula.driver.graph.data.ResultSet;
import com.vesoft.nebula.driver.graph.exception.AuthFailedException;
import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import java.io.Serializable;
import java.time.ZoneId;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import org.apache.commons.pool2.BasePooledObjectFactory;
import org.apache.commons.pool2.PooledObject;
import org.apache.commons.pool2.impl.DefaultPooledObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ClientPoolFactory extends BasePooledObjectFactory<NebulaClient>
        implements Serializable {
    private final Logger logger = LoggerFactory.getLogger(this.getClass());

    private final LoadBalancer       loadBalancer;
    private final NebulaPool.Builder builder;

    public ClientPoolFactory(
            LoadBalancer loadBalancer,
            NebulaPool.Builder builder) {
        this.loadBalancer = loadBalancer;
        this.builder = builder;
    }


    @Override
    public NebulaClient create() throws IOErrorException, AuthFailedException {
        List<HostAddress> goodHosts = loadBalancer.getGoodAddresses();
        if (goodHosts.size() == 0) {
            throw new IOErrorException(IOErrorException.E_ALL_BROKEN, "No servers host is "
                    + "available, please check your servers is up and network between client and "
                    + "server is connected.");
        }

        String addrStr = goodHosts.stream()
                .map(HostAddress::toString)
                .collect(Collectors.joining(","));

        NebulaClient client = NebulaClient
                .builder(addrStr, builder.userName)
                .withAuthOptions(builder.authOptions)
                .withConnectTimeoutMills(builder.connectTimeoutMills)
                .withRequestTimeoutMills(builder.requestTimeoutMills)
                .withScanParallel(builder.scanParallel)
                .withEnableTls(builder.enableTls)
                .withTlsCa(builder.tlsCa)
                .withTlsCert(builder.tlsCert, builder.tlsKey)
                .withTlsPeerName(builder.tlsPeerName)
                .build();

        // set the working graph and time zone
        StringBuilder sessionSetStatement = new StringBuilder();
        if (builder.workingGraph != null) {
            sessionSetStatement
                    .append("SESSION SET GRAPH `")
                    .append(builder.workingGraph)
                    .append("` ");
        }
        if (builder.timeZone != null) {
            sessionSetStatement
                    .append("SESSION SET TIME ZONE \"")
                    .append(builder.timeZone)
                    .append("\"");
        }
        if (builder.schemaName != null) {
            sessionSetStatement
                    .append("SESSION SET SCHEMA `")
                    .append(builder.schemaName)
                    .append("` ");
        }

        if (!sessionSetStatement.toString().isEmpty()) {
            ResultSet result = client.execute(sessionSetStatement.toString());
            if (!result.isSucceeded()) {
                throw new RuntimeException("SESSION SET failed for " + result.getErrorMessage());
            }
        }
        for (Map.Entry<String, String> parameter : builder.parameters.entrySet()) {
            String setParameterStatement = String.format("SESSION SET VALUE $%s=%s",
                                                         parameter.getKey(),
                                                         parameter.getValue());
            ResultSet setParamResult = client.execute(setParameterStatement);
            if (!setParamResult.isSucceeded()) {
                throw new RuntimeException("SESSION SET failed for "
                                                   + setParamResult.getErrorMessage());
            }
        }
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
