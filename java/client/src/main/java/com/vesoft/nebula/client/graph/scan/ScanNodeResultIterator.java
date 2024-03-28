/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.scan;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.data.ValueWrapper;
import com.vesoft.nebula.client.graph.net.Session;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import org.apache.commons.pool2.impl.GenericObjectPool;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ScanNodeResultIterator extends ScanResultIterator {
    private static final Logger LOGGER = LoggerFactory.getLogger(ScanNodeResultIterator.class);

    private static final String SCAN_NODE_TEMPLATE =
            "USE %s CALL cursor_node_scan(\"%s\",\"%s\",%s,%d,\"%s\", %d) return *";

    public ScanNodeResultIterator(GenericObjectPool<Session> pool, String graphName,
                                  String label, List<String> propNames, List<Integer> parts,
                                  int batchSize, ExecutorService threadPool, int retryTimes,
                                  int intervalTime, long timeoutMs) {
        super(pool, graphName, label, propNames, parts, batchSize, threadPool, retryTimes,
                intervalTime, timeoutMs);
    }


    public ScanNodeResult next() {
        if (!hasNext) {
            throw new NoSuchElementException("iterator has no more data");
        }
        final List<ResultSet> results =
                Collections.synchronizedList(new ArrayList<>(partCursor.size()));
        List<Exception> exceptions =
                Collections.synchronizedList(new ArrayList<>(partCursor.size()));
        CountDownLatch countDownLatch = new CountDownLatch(partCursor.size());
        for (Map.Entry<Integer, String> partCur : partCursor.entrySet()) {
            threadPool.submit(() -> {
                try {
                    ResultSet result = scan(SCAN_NODE_TEMPLATE, partCur);
                    // collect results and update the cursor
                    if (result.isSucceeded()) {
                        String cursor = getCursor(result);
                        partCursor.put(partCur.getKey(), cursor);
                        results.add(result);
                    } else {
                        LOGGER.error(String.format("Scan part %d of node %s failed for %s, "
                                        + "scan again in the next next()",
                                partCur.getKey(),
                                labelName,
                                result.getErrorMessage()));
                        exceptions.add(new Exception(String.format("part %d of %s scan error: ",
                                partCur.getKey(), labelName, result.getErrorMessage())));
                    }
                } catch (Exception e) {
                    LOGGER.error(String.format("Scan node error for %s", e.getMessage()), e);
                    exceptions.add(new Exception(String.format("part %d of %s scan failed: %s",
                            partCur.getKey(), labelName, e.getMessage()), e));
                } finally {
                    countDownLatch.countDown();
                }
            });
        }

        try {
            countDownLatch.await();
        } catch (InterruptedException interruptedException) {
            LOGGER.error("scan interrupted:", interruptedException);
            throw new RuntimeException("scan interrupted", interruptedException);
        }

        // As long as one part fails, the current iteration is considered as failed.
        if (!exceptions.isEmpty()) {
            List<String> exceptionMsg = new ArrayList<>();
            for (Exception e : exceptions) {
                exceptionMsg.add(e.getMessage());
            }
            throw new RuntimeException("scan node failed for current iterator: " + exceptionMsg);
        }

        hasNext = false;
        List<Integer> partKeyNeedToRemove = new ArrayList<>();
        for (Map.Entry<Integer, String> partCur : partCursor.entrySet()) {
            if (!"".equals(partCur.getValue())) {
                hasNext = true;
                break;
            } else {
                partKeyNeedToRemove.add(partCur.getKey());
            }
        }
        for (Integer part : partKeyNeedToRemove) {
            partCursor.remove(part);
        }
        return new ScanNodeResult(results, propNames);
    }

}
