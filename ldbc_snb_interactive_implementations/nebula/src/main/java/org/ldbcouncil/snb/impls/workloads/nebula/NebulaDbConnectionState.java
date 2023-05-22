package org.ldbcouncil.snb.impls.workloads.nebula;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;

import com.vesoft.nebula.client.graph.net.NebulaClient;
import com.vesoft.nebula.client.graph.net.Session;
import com.vesoft.nebula.graph.GraphService;
import org.ldbcouncil.snb.impls.workloads.BaseDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.QueryStore;

import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;


public class NebulaDbConnectionState<TDbQueryStore extends QueryStore> extends BaseDbConnectionState<TDbQueryStore>
{
    private NebulaClient client = null;
    private final String graphName;

    public String getGraphName() {
        return graphName;
    }

    public NebulaDbConnectionState( Map<String, String> properties, TDbQueryStore store ) throws UnknownHostException, IOErrorException {
        super(properties, store);

        final String endpointURI = properties.get( "endpoint" );
        final String username = properties.get( "user" );
        final String password = properties.get( "password" );
        final int requestTimeout = Integer.parseInt(properties.get("requestTimeout"));
        graphName = properties.get("graphName");
        try {
            client = NebulaClient.builder(endpointURI, username, password)
                    .setRequestTimeoutMills(requestTimeout*1000)
                    .setMaxSessionSize(10)
                    .setMinSessionSize(1)
                    .setRetryTimes(3)
                    .setIntervalTimeMills(1000)
                    .setReconnect(true)
                    .setBlockWhenExhausted(true)
                    .setMaxWaitMills(1000)
                    .setStrictlyServerHealthy(true)
                    .build();
        } catch (UnknownHostException| IOErrorException e) {
            throw new RuntimeException(e);
        }
    }

    public NebulaClient getClient() {
        return client;
    }

    @Override
    public void close() {
        client.close();
    }
}
