package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class RoundRobinLoadBalancer implements LoadBalancer, Serializable {
    private static final Logger LOGGER = LoggerFactory.getLogger(RoundRobinLoadBalancer.class);
    private static final int S_OK = 0;
    private static final int S_BAD = 1;
    private final List<HostAddress> addresses = new ArrayList<>();
    private final Map<HostAddress, Integer> serversStatus = new ConcurrentHashMap<>();
    private final boolean strictlyServerHealthy;
    private final int connTimeout;
    private final int requestTimeout;
    private final AtomicInteger pos = new AtomicInteger(0);
    private final ScheduledExecutorService schedule = Executors.newScheduledThreadPool(1);

    public RoundRobinLoadBalancer(List<HostAddress> addresses, int connTimeout, int requestTimeout,
                                  boolean strictlyServerHealthy, int healthCheckTime) {
        this.connTimeout = connTimeout;
        this.requestTimeout = requestTimeout;
        for (HostAddress addr : addresses) {
            this.addresses.add(addr);
            this.serversStatus.put(addr, S_BAD);
        }
        this.strictlyServerHealthy = strictlyServerHealthy;

        schedule.scheduleAtFixedRate(this::scheduleTask, 0, healthCheckTime, TimeUnit.MILLISECONDS);
    }

    public void close() {
        if (!schedule.isShutdown()) {
            schedule.shutdownNow();
        }
    }

    @Override
    public HostAddress getAddress() {
        if (pos.get() == Integer.MAX_VALUE) {
            pos.set(0);
        }
        int tryCount = 0;
        int newPos;
        while (++tryCount <= addresses.size()) {
            newPos = (pos.getAndIncrement()) % addresses.size();
            HostAddress addr = addresses.get(newPos);
            if (serversStatus.get(addr) == S_OK) {
                return addr;
            }
        }
        return null;
    }

    public void updateServersStatus() {
        for (HostAddress hostAddress : addresses) {
            if (ping(hostAddress)) {
                serversStatus.put(hostAddress, S_OK);
            } else {
                serversStatus.put(hostAddress, S_BAD);
            }
        }
    }

    public List<HostAddress> getGoodAddresses() {
        List<HostAddress> goodHosts = new ArrayList<>();
        for (Map.Entry<HostAddress, Integer> server : serversStatus.entrySet()) {
            if (server.getValue() == S_OK) {
                goodHosts.add(server.getKey());
            }
        }
        return goodHosts;
    }

    public boolean ping(HostAddress addr) {
        try {
            Connection connection = new SyncConnection();
            connection.open(addr, this.connTimeout, this.requestTimeout);
            boolean pong = connection.ping();
            connection.close();
            return pong;
        } catch (IOErrorException e) {
            LOGGER.error("ping failed, ", e);
            return false;
        }
    }

    public boolean isServersOK() {
        this.updateServersStatus();
        int numServersWithOkStatus = 0;
        int numServersWithBadStatus = 0;
        for (HostAddress hostAddress : addresses) {
            if (serversStatus.get(hostAddress) == S_OK) {
                numServersWithOkStatus++;
            } else {
                numServersWithBadStatus++;
            }
        }
        return (strictlyServerHealthy && numServersWithBadStatus == 0)
                || (!strictlyServerHealthy && numServersWithOkStatus > 0);
    }

    private void scheduleTask() {
        updateServersStatus();
    }
}
