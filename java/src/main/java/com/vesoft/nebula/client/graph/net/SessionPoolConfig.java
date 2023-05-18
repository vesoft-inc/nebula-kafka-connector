/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.data.HostAddress;
import java.io.Serializable;
import java.util.List;

public class SessionPoolConfig implements Serializable {

    private final List<HostAddress> graphAddressList;

    private final String username;
    private final String password;

    // The min connections in pool for all addresses
    private int minSessionSize = 1;

    // The max connections in pool for all addresses
    private int maxSessionSize = 10;

    // Socket connection timeout, unit: millisecond. 0 means never timeout
    private int connTimeout = 0;

    // Socket request timeout, unit: millisecond. 0 means never timeout
    private int requestTimeout = 0;

    // The healthCheckTime for schedule check the health of session, unit: second.
    // 0 means never check the health of session.
    private int healthCheckTime = 600;

    // retry times for failed execute
    private int retryTimes = 3;

    // interval time for retry, unit ms
    private int intervalTime = 0;

    // whether reconnect when create session using a broken graphd server
    private boolean reconnect = false;

    // if block when session is exhausted, if false, throw exception.
    private boolean blockWhenExhausted = false;

    // the max wait time if blockWhenExhausted is true. if value is less than 0, always wait.
    // unit: millisecond
    private long maxWaitMills = -1;

    // if need all servers are strictly healthy.
    // if true, all addresses must be available, if false, at least one address is available.
    private boolean strictlyServerHealthy = false;


    public SessionPoolConfig(List<HostAddress> addresses,
                             String username,
                             String password) {
        this.graphAddressList = addresses;
        this.username = username;
        this.password = password;
    }

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
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

    public int getConnTimeout() {
        return connTimeout;
    }

    public SessionPoolConfig setConnTimeout(int connTimeout) {
        if (connTimeout < 0) {
            throw new IllegalArgumentException("connect timeout cannot be less than 0.");
        }
        this.connTimeout = connTimeout;
        return this;
    }

    public int getRequestTimeout() {
        return requestTimeout;
    }

    public SessionPoolConfig setRequestTimeout(int requestTimeout) {
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

    public boolean isReconnect() {
        return reconnect;
    }

    public SessionPoolConfig setReconnect(boolean reconnect) {
        this.reconnect = reconnect;
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
                + ", password='" + password + '\''
                + ", minSessionSize=" + minSessionSize
                + ", maxSessionSize=" + maxSessionSize
                + ", connTimeout=" + connTimeout
                + ", requestTimeout=" + requestTimeout
                + ", healthCheckTime=" + healthCheckTime
                + ", retryTimes=" + retryTimes
                + ", intervalTime=" + intervalTime
                + ", reconnect=" + reconnect
                + ", blockWhenExhausted=" + blockWhenExhausted
                + ", maxWaitMills=" + maxWaitMills
                + ", strictlyServerHealthy=" + strictlyServerHealthy
                + '}';
    }
}
