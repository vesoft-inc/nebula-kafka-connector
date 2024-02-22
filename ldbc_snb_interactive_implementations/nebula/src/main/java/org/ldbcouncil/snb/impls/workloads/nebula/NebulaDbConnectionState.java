package org.ldbcouncil.snb.impls.workloads.nebula;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;

import com.vesoft.nebula.client.graph.net.NebulaClient;
import com.vesoft.nebula.client.graph.net.Session;
import org.ldbcouncil.snb.impls.workloads.BaseDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.QueryStore;

import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;


public class NebulaDbConnectionState<TDbQueryStore extends QueryStore> extends BaseDbConnectionState<TDbQueryStore> {
    private List<NebulaNewClient> clients = new ArrayList<>();
    private Map<Integer, Integer> requestTypeAndClientIndex = new HashMap<>();
    private AtomicInteger clientIndex = new AtomicInteger(0);
    private final String graphName;

    public String getGraphName() {
        return graphName;
    }

    public NebulaDbConnectionState(Map<String, String> properties, TDbQueryStore store) throws UnknownHostException, IOErrorException {
        super(properties, store);

        final String endpointURI = properties.get("endpoint");
        final String username = properties.get("user");
        final String password = properties.get("password");
        final int requestTimeout = Integer.parseInt(properties.get("requestTimeout"));
        final int maxSessionSize = Integer.parseInt(properties.get("maxSessionSize"));
        final int maxSessionWaitTime = Integer.parseInt(properties.get("maxSessionWaitTime"));
        final int retryTimes = Integer.parseInt(properties.get("retryTimes"));
        final int intervalTimeBetweenRetrys = Integer.parseInt(properties.get("intervalTimeBetweenRetrys"));
        graphName = properties.get("graphName");
        String[] graphAddresses = endpointURI.split(",");
        try {
            for (String addr : graphAddresses) {
                NebulaClient client = NebulaClient.builder(addr, username, password)
                        .setRequestTimeoutMills(requestTimeout * 1000)
                        .setMaxSessionSize(maxSessionSize)
                        .setMinSessionSize(maxSessionSize)
                        .setRetryTimes(retryTimes)
                        .setIntervalTimeMills(intervalTimeBetweenRetrys * 1000)
                        .setReconnect(true)
                        .setBlockWhenExhausted(true)
                        .setMaxWaitMills(maxSessionWaitTime * 1000L)
                        .setStrictlyServerHealthy(true)
                        .build();
                clients.add(new NebulaNewClient(client, addr));
            }
        } catch (UnknownHostException | IOErrorException e) {
            throw new RuntimeException(e);
        }
        // we just balance the graphd for query operation types: LdbcQyery1 to LdbcQuery14
        // see org.ldbcouncil.snb.driver.Operation.type
        for (int i = 1; i < 15; i++) {
            requestTypeAndClientIndex.put(i, 0);
        }
    }

    public NebulaNewClient getClient(int requestType) {
        if (requestType < 15) {
            int index = requestTypeAndClientIndex.get(requestType);
            int newIndex = (index + 1) % clients.size();
            requestTypeAndClientIndex.put(requestType, newIndex);
            return clients.get(index);
        }
        // if operation is not query, then just balance the graphd through the requests.
        return clients.get(clientIndex.getAndIncrement() % clients.size());
    }

    @Override
    public void close() {
        for (NebulaNewClient client : clients) {
            client.getClient().close();
        }
    }
}
