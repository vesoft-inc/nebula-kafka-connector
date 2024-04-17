package org.ldbcouncil.snb.impls.workloads.nebula;

import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;

import com.vesoft.nebula.client.graph.net.NebulaClient;
import org.ldbcouncil.snb.impls.workloads.BaseDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.QueryStore;

import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;


public class NebulaDbConnectionState<TDbQueryStore extends QueryStore> extends BaseDbConnectionState<TDbQueryStore> {

    private final ThreadLocal<List<NebulaNewClient>> threadLocal = new ThreadLocal<>();
    private Map<Integer, Integer> requestTypeAndClientIndex = new HashMap<>();
    private AtomicInteger clientIndex = new AtomicInteger(0);
    private final String graphName;

    final String[] graphAddresses;
    final String userName;
    final String password;
    final long requestTimeout;
    final int retryTimes;
    final long intervalTimeBetweenRetrys;

    private final int graphServerSize;

    public String getGraphName() {
        return graphName;
    }

    public NebulaDbConnectionState(Map<String, String> properties, TDbQueryStore store) throws UnknownHostException, IOErrorException {
        super(properties, store);

        final String endpointURI = properties.get("endpoint");
        userName = properties.get("user");
        password = properties.get("password");
        requestTimeout = Integer.parseInt(properties.get("requestTimeout")) * 1000L;
        retryTimes = Integer.parseInt(properties.get("retryTimes"));
        intervalTimeBetweenRetrys = Integer.parseInt(properties.get("intervalTimeBetweenRetrys")) * 1000L;
        graphName = properties.get("graphName");
        graphAddresses = endpointURI.split(",");
        graphServerSize = graphAddresses.length;
        // we just balance the graphd for query operation types: LdbcQyery1 to LdbcQuery14
        // see org.ldbcouncil.snb.driver.Operation.type
        for (int i = 1; i < 15; i++) {
            requestTypeAndClientIndex.put(i, 0);
        }
    }

    public NebulaNewClient getClient(int requestType) {
        if (threadLocal.get() == null) {
            List<NebulaNewClient> clients = new ArrayList<>();
            try {
                for (String addr : graphAddresses) {
                    NebulaClient client = NebulaClient.builder(addr, userName, password)
                            .setRequestTimeoutMills(requestTimeout * 1000L)
                            .setRetryTimes(retryTimes)
                            .setIntervalTimeMills(intervalTimeBetweenRetrys * 1000L)
                            .build();
                    clients.add(new NebulaNewClient(client, addr));
                }
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
            threadLocal.set(clients);
        }

        if (requestType < 15) {
            int currentClientIndex = requestTypeAndClientIndex.get(requestType);
            int nextClientIndex = (currentClientIndex + 1) % graphServerSize;
            requestTypeAndClientIndex.put(requestType, nextClientIndex);
            return threadLocal.get().get(currentClientIndex);
        }
        // if operation is not query, then just balance the graphd through the requests.
        return threadLocal.get().get(clientIndex.getAndIncrement() % graphServerSize);

    }

    @Override
    public void close() {
        if(threadLocal.get().isEmpty()){
            return;
        }

        for (NebulaNewClient client : threadLocal.get()) {
            client.getClient().close();
        }
    }
}
