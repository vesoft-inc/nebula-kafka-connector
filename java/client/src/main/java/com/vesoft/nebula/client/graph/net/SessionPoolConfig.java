/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.data.HostAddress;
import java.io.Serializable;
import java.util.List;
import java.util.Map;

public class SessionPoolConfig implements Serializable {

    private final List<HostAddress> graphAddressList;

    private final String username;
    private final Map<String,Object> authOptions;

    // The min connections in pool for all addresses
    private int minSessionSize = 1;

    // The max connections in pool for all addresses
    private int maxSessionSize = 10;

    // Socket request timeout, unit: millisecond. 0 means never timeout
    private long requestTimeout = 0;

    // The healthCheckTime for schedule check the health of session, unit: millisecond.
    // 0 means never check the health of session.
    private int healthCheckTime = 300000;

    // retry times for failed execute
    private int retryTimes = 3;

    // interval time for retry, unit ms
    private int intervalTime = 0;

    // if block when session is exhausted, if false, throw exception.
    private boolean blockWhenExhausted = false;

    // the max wait time if blockWhenExhausted is true. if value is less than 0, always wait.
    // unit: millisecond
    private long maxWaitMills = -1;

    // the schedule time for test the idle session and evict it. if value is less than 0,
    // never evict the idle sessions.
    private long idleEvictScheduleMills = -1;

    // the min idle time for idle session, default is 30 minutes.
    private long minEvictableIdleTimeMillis = 1000L * 60L * 30L;

    // if need all servers are strictly healthy.
    // if true, all addresses must be available, if false, at least one address is available.
    private boolean strictlyServerHealthy = false;


    public SessionPoolConfig(List<HostAddress> addresses,
                             String username,
                             Map<String,Object> authOptions) {
        this.graphAddressList = addresses;
        this.username = username;
        this.authOptions = authOptions;
    }

    public String getUsername() {
        return username;
    }

    public Map<String, Object> getAuthOptions() {
        return authOptions;
    }

    public List<HostAddress> getGraphAddressList() {
        return graphAddressList;
    }

    public int getMinSessionSize() {
        return minSessionSize;
    }

    public SessionPoolConfig setMinSessionSize(int minSessionSize) {
        if (minSessionSize < 0) {
            throw new IllegalArgumentException("minSessionSize cannot be less than 0.");
        }
        this.minSessionSize = minSessionSize;
        return this;
    }

    public int getMaxSessionSize() {
        return maxSessionSize;
    }

    public SessionPoolConfig setMaxSessionSize(int maxSessionSize) {
        if (maxSessionSize < 1) {
            throw new IllegalArgumentException("maxSessionSize cannot be less than 1.");
        }
        this.maxSessionSize = maxSessionSize;
        return this;
    }

    public long getRequestTimeout() {
        return requestTimeout;
    }

    public SessionPoolConfig setRequestTimeout(long requestTimeout) {
        if (requestTimeout < 0) {
            throw new IllegalArgumentException("request timeout cannot be less than 0.");
        }
        this.requestTimeout = requestTimeout;
        return this;
    }

    public int getHealthCheckTime() {
        return healthCheckTime;
    }

    public SessionPoolConfig setHealthCheckTime(int healthCheckTime) {
        if (healthCheckTime < 0) {
            throw new IllegalArgumentException("cleanTime cannot be less than 0.");
        }
        this.healthCheckTime = healthCheckTime;
        return this;
    }

    public int getRetryTimes() {
        return retryTimes;
    }

    public SessionPoolConfig setRetryTimes(int retryTimes) {
        if (retryTimes < 0) {
            throw new IllegalArgumentException("retryTimes cannot be less than 0.");
        }
        this.retryTimes = retryTimes;
        return this;
    }

    public int getIntervalTime() {
        return intervalTime;
    }

    public SessionPoolConfig setIntervalTime(int intervalTime) {
        if (intervalTime < 0) {
            throw new IllegalArgumentException("intervalTime cannot be less than 0.");
        }
        this.intervalTime = intervalTime;
        return this;
    }

    public boolean isBlockWhenExhausted() {
        return blockWhenExhausted;
    }

    public SessionPoolConfig setBlockWhenExhausted(boolean blockWhenExhausted) {
        this.blockWhenExhausted = blockWhenExhausted;
        return this;
    }

    public long getMaxWaitMills() {
        return maxWaitMills;
    }

    public SessionPoolConfig setMaxWaitMills(long maxWaitMills) {
        this.maxWaitMills = maxWaitMills;
        return this;
    }

    public SessionPoolConfig setIdleEvictScheduleMills(long idleEvictScheduleMills) {
        this.idleEvictScheduleMills = idleEvictScheduleMills;
        return this;
    }

    public SessionPoolConfig setMinEvictableIdleTimeMillis(long minEvictableIdleTimeMillis) {
        this.minEvictableIdleTimeMillis = minEvictableIdleTimeMillis;
        return this;
    }

    public boolean isStrictlyServerHealthy() {
        return strictlyServerHealthy;
    }

    public SessionPoolConfig setStrictlyServerHealthy(boolean strictlyServerHealthy) {
        this.strictlyServerHealthy = strictlyServerHealthy;
        return this;
    }

    @Override
    public String toString() {
        return "SessionPoolConfig{"
                + "graphAddressList=" + graphAddressList
                + ", username='" + username + '\''
                + ", authOptions='" + authOptions + '\''
                + ", minSessionSize=" + minSessionSize
                + ", maxSessionSize=" + maxSessionSize
                + ", requestTimeout=" + requestTimeout
                + ", healthCheckTime=" + healthCheckTime
                + ", retryTimes=" + retryTimes
                + ", intervalTime=" + intervalTime
                + ", blockWhenExhausted=" + blockWhenExhausted
                + ", maxWaitMills=" + maxWaitMills
                + ", idleEvictScheduleMills=" + idleEvictScheduleMills
                + ", minEvictableIdleTimeMillis=" + minEvictableIdleTimeMillis
                + ", strictlyServerHealthy=" + strictlyServerHealthy
                + '}';
    }
}
