package org.ldbcouncil.snb.impls.workloads.nebula;

import com.vesoft.nebula.driver.graph.exception.IOErrorException;

import com.vesoft.nebula.driver.graph.net.NebulaPool;
import org.ldbcouncil.snb.impls.workloads.BaseDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.QueryStore;

import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;


public class NebulaDbConnectionState<TDbQueryStore extends QueryStore> extends BaseDbConnectionState<TDbQueryStore> {

    private       List<NewNebulaPool>         pools                         = new ArrayList<>();
    private       Map<Integer, AtomicInteger> requestTypeAndClientPoolIndex = new ConcurrentHashMap<>();
    private       AtomicInteger               poolClientIndex               = new AtomicInteger(0);
    private final String                      graphName;

    final String[] graphAddresses;
    final String   userName;
    final String   password;
    final long     connectTimeout;
    final long     requestTimeout;

    // enable the log print for query statement, query type, query execution latency, query response time
    final boolean enableQueryInfoLog;

    private final int graphServerSize;

    public String getGraphName() {
        return graphName;
    }

    public NebulaDbConnectionState(Map<String, String> properties, TDbQueryStore store) throws UnknownHostException, IOErrorException {
        super(properties, store);

        final String endpointURI = properties.get("endpoint");
        userName = properties.get("user");
        password = properties.get("password");

        if (properties.containsKey("requestTimeout")) {
            connectTimeout = Integer.parseInt(properties.get("requestTimeout")) * 1000L;
        } else {
            connectTimeout = Long.MAX_VALUE;
        }

        if (properties.containsKey("requestTimeout")) {
            requestTimeout = Integer.parseInt(properties.get("requestTimeout")) * 1000L;
        } else {
            requestTimeout = Long.MAX_VALUE;
        }

        if (properties.containsKey("enableQueryInfoLog")) {
            enableQueryInfoLog = Boolean.parseBoolean(properties.get("enableQueryInfoLog"));
        } else {
            enableQueryInfoLog = false;
        }

        int  maxClientSize      = Integer.parseInt(properties.get("maxSessionSize"));
        long maxSessionWaitTime = 0;
        if (properties.containsKey("maxSessionWaitTime")) {
            maxSessionWaitTime = Integer.parseInt(properties.get("maxSessionWaitTime")) * 1000L;
        } else {
            maxSessionWaitTime = Integer.MAX_VALUE;
        }

        graphName = properties.get("graphName");
        graphAddresses = endpointURI.split(",");
        graphServerSize = graphAddresses.length;
        try {
            for (String addr : graphAddresses) {
                NebulaPool pool = NebulaPool.builder(addr, userName, password)
                        .withMaxClientSize(maxClientSize)
                        .withMinClientSize(maxClientSize)
                        .withConnectTimeoutMills(connectTimeout)
                        .withRequestTimeoutMills(requestTimeout)
                        .withBlockWhenExhausted(true)
                        .withMaxWaitMills(maxSessionWaitTime)
                        .withHealthCheckTimeMills(-1)
                        .build();
                pools.add(new NewNebulaPool(pool, addr));
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        // we just balance the graphd for query operation types: LdbcQyery1 to LdbcQuery14
        // see org.ldbcouncil.snb.driver.Operation.type
        for (int i = 1; i < 15; i++) {
            requestTypeAndClientPoolIndex.put(i, new AtomicInteger(0));
        }
    }

    public NewNebulaPool getPool(int requestType) {
        if (requestType < 15) {
            int currentClientPoolIndexValue = requestTypeAndClientPoolIndex.get(requestType).get();
            while (!requestTypeAndClientPoolIndex.get(requestType).compareAndSet(currentClientPoolIndexValue, currentClientPoolIndexValue + 1)) {
                currentClientPoolIndexValue = requestTypeAndClientPoolIndex.get(requestType).get();
            }
            return pools.get(currentClientPoolIndexValue % graphServerSize);
        }
        // if operation is not query, then just balance the graphd through the requests.
        int curClientIndexV = poolClientIndex.get();
        while (!poolClientIndex.compareAndSet(curClientIndexV, curClientIndexV + 1)) {
            curClientIndexV = poolClientIndex.get();
        }
        return pools.get(curClientIndexV % pools.size());
    }

    public boolean isEnableQueryInfoLog() {
        return enableQueryInfoLog;
    }

    @Override
    public void close() {
        for (NewNebulaPool clientPool : pools) {
            clientPool.getPool().close();
        }
    }
}
