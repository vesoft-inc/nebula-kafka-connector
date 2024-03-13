/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.scan;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.data.ValueWrapper;
import com.vesoft.nebula.client.graph.net.Session;
import java.io.Serializable;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import org.apache.commons.pool2.impl.GenericObjectPool;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ScanResultIterator implements Serializable {
    private static final Logger LOGGER = LoggerFactory.getLogger(ScanResultIterator.class);

    protected static final String CURSOR_NAME = "__cursor__";
    protected boolean hasNext = true;

    protected final Map<Integer, String> partCursor = new HashMap<>();

    protected final String graphName;
    protected final String labelName;


    protected final GenericObjectPool<Session> pool;

    protected List<String> propNames;
    protected int batchSize;

    protected final ExecutorService threadPool;

    protected final int retryTimes;
    protected final int intervalTime;
    protected final long timeoutMs;

    protected ScanResultIterator(GenericObjectPool<Session> pool, String graphName,
                                 String labelName, List<String> propNames, List<Integer> parts,
                                 int batchSize, ExecutorService threadPool, int retryTimes,
                                 int intervalTime, long timeoutMs) {
        this.graphName = graphName;
        this.labelName = labelName;
        this.pool = pool;
        this.propNames = propNames;
        this.batchSize = batchSize;
        for (int part : parts) {
            partCursor.put(part, "");
        }
        this.threadPool = threadPool;
        this.retryTimes = retryTimes;
        this.intervalTime = intervalTime;
        this.timeoutMs = timeoutMs;
    }

    /**
     * if iter has more vertex data
     *
     * @return true if the scan cursor is not at end.
     */
    public boolean hasNext() {
        return hasNext;
    }

    protected String getPropertyListString() {
        StringBuilder properties = new StringBuilder();
        String propertyListPrefix = "list[";
        properties.append(propertyListPrefix);
        for (String column : propNames) {
            properties.append("\"");
            properties.append(column);
            properties.append("\"");
            properties.append(",");
        }
        if (properties.length() > propertyListPrefix.length()) {
            properties.deleteCharAt(properties.length() - 1);
        }
        String propertyListSuffix = "]";
        properties.append(propertyListSuffix);
        return properties.toString();
    }

    protected ResultSet scan(String scanTemplate, Map.Entry<Integer, String> partCur)
            throws Exception {
        // construct the scan producer
        String producer = String.format(scanTemplate, graphName, graphName, labelName,
                getPropertyListString(), partCur.getKey(), partCur.getValue(), batchSize);
        Session session = null;
        ResultSet result;
        try {
            session = pool.borrowObject(Long.MAX_VALUE);
            result = session.execute(producer, timeoutMs);
            int retry = retryTimes;
            while (retry > 0) {
                retry--;
                Thread.sleep(intervalTime);
                result = session.execute(producer, timeoutMs);
                if (result.isSucceeded()) {
                    break;
                }
            }
        } finally {
            pool.returnObject(session);
        }
        return result;
    }

    protected String getCursor(ResultSet resultSet) {
        Map<String, ValueWrapper> extraInfo = resultSet.getExtraInfo();
        if (!extraInfo.containsKey("cursor")) {
            throw new RuntimeException("result does not contain cursor in extra info.");
        }
        return extraInfo.get("cursor").asString();
    }
}
