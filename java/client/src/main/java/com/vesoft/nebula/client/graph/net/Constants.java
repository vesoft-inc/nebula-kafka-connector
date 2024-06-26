/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

public class Constants {

    static final int     DEFAULT_MAX_CLIENT_SIZE            = 10;
    static final int     DEFAULT_MIN_CLIENT_SIZE            = 1;
    static final long    DEFAULT_CONNECT_TIMEOUT            = 3600 * 1000; // 1 hour
    static final long    DEFAULT_REQUEST_TIMEOUT            = 3600 * 1000; // 1 hour
    static final long    DEFAULT_HEALTH_CHECK_TIME_MS       = 5 * 60 * 1000;
    static final boolean DEFAULT_BLOCK_WHEN_EXHAUSTED       = false;
    static final long    DEFAULT_MAX_WAIT_MS                = Long.MAX_VALUE;
    static final long    DEFAULT_IDLE_EVICT_SCHEDULE_MS     = -1;
    static final long    DEFAULT_MIN_EVICTABLE_IDLE_TIME_MS = 30 * 60 * 1000;
    static final boolean DEFAULT_STRICT_SERVER_HEALTHY      = false;
    static final int     DEFAULT_BATCH_SIZE                 = 1000;
    static final int     DEFAULT_SCAN_PARALLEL              = 10;

}
