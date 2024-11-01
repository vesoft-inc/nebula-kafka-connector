package com.vesoft.nebula.driver.graph.net;

import com.vesoft.nebula.driver.graph.data.HostAddress;
import com.vesoft.nebula.driver.graph.exception.AuthFailedException;
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
    private static final Logger logger = LoggerFactory.getLogger(RoundRobinLoadBalancer.class);

    private static final int                       S_OK          = 0;
    private static final int                       S_BAD         = 1;
    private final        List<HostAddress>         addresses     = new ArrayList<>();
    private final        Map<HostAddress, Integer> serversStatus = new ConcurrentHashMap<>();
    private final        boolean                   strictlyServerHealthy;

    private final String              userName;
    private final Map<String, Object> authOptions;

    private final AtomicInteger            pos = new AtomicInteger(0);
    private       ScheduledExecutorService schedule;

    private boolean enableTls;
    private String  tlsCa;
    private String  tlsCert;
    private String  tlsKey;
    private boolean tlsPeerNameVerify;
    private String  tlsPeerName;

    public RoundRobinLoadBalancer(NebulaPool.Builder builder) {
        for (HostAddress addr : builder.address) {
            this.addresses.add(addr);
            this.serversStatus.put(addr, S_BAD);
        }
        this.strictlyServerHealthy = builder.strictlyServerHealthy;
        this.userName = builder.userName;
        this.authOptions = builder.authOptions;

        if (builder.healthCheckTimeMills > 0) {
            schedule = Executors.newScheduledThreadPool(1);
            schedule.scheduleAtFixedRate(this::scheduleTask,
                                         0,
                                         builder.healthCheckTimeMills,
                                         TimeUnit.MILLISECONDS);
        }
        enableTls = builder.enableTls;
        tlsCa = builder.tlsCa;
        tlsCert = builder.tlsCert;
        tlsKey = builder.tlsKey;
        tlsPeerNameVerify = builder.tlsPeerNameVerify;
        tlsPeerName = builder.tlsPeerName;
    }

    public void close() {
        if (schedule != null && !schedule.isShutdown()) {
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

    public void updateServersStatus() throws AuthFailedException {
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

    public boolean ping(HostAddress addr) throws AuthFailedException {
        try {
            NebulaClient client = NebulaClient
                    .builder(addr.toString(), userName)
                    .withAuthOptions(authOptions)
                    .withEnableTls(enableTls)
                    .withTlsCa(tlsCa)
                    .withTlsCert(tlsCert, tlsKey)
                    .withTlsPeerName(tlsPeerName)
                    .build();
            client.close();
            return true;
        } catch (AuthFailedException e) {
            logger.error("auth failed,", e);
            throw e;
        } catch (Exception e) {
            logger.error("ping failed, ", e);
            return false;
        }
    }

    public boolean isServersOK() throws AuthFailedException {
        this.updateServersStatus();
        int numServersWithOkStatus  = 0;
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
        try {
            updateServersStatus();
        } catch (AuthFailedException e) {
            logger.error("auth failed, ", e);
        }
    }
}
