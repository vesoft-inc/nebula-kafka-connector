/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.util.MockGraph;
import com.vesoft.nebula.client.util.ProcessUtil;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class NebulaPoolTest {
    String addresses = "127.0.0.1:9669,127.0.0.1:9670,127.0.0.1:9671";
    String user      = "root";
    String passwd    = "NebulaGraph01";

    @Before
    public void setup() {
        MockGraph.mockGraphData();
    }

    @Test
    public void testNebulaPool() {
        System.out.println("<==== testNebulaPool ====>");
        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withMaxClientSize(10)
                    .withMinClientSize(1)
                    .build();
            NebulaClient client = pool.getClient();
            client.execute("RETURN 1");
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }
    }

    @Test
    public void testWorkingGraph() {
        System.out.println("<==== testWorkingGraph ====>");
        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withWorkingGraph("nba")
                    .build();
            NebulaClient client = pool.getClient();
            client.execute("RETURN 1");
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }

        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withWorkingGraph("nba_not_exist")
                    .build();
            pool.getClient();
            Assert.fail("get client should fail.");
        } catch (Exception e) {
            assert e.getMessage().contains("SESSION SET failed");
        }
    }


    @Test
    public void testStrictlyServerHealthy() {
        System.out.println("<==== testStrictlyServerHealthy ====>");
        // stop one graphd server
        Runtime runtime = Runtime.getRuntime();
        try {
            String  cmd = "docker stop docker-compose-graphd0-1";
            Process p   = runtime.exec(cmd);
            p.waitFor(10, TimeUnit.SECONDS);
            ProcessUtil.printProcessStatus(cmd, p);
            Thread.sleep(5000);
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withStrictlyServerHealthy(true)
                    .build();
        } catch (Exception e) {
            assert e.getMessage().contains("Servers status is not ok");
        }

        // start the graphd server
        try {
            String  cmd = "docker start docker-compose-graphd0-1";
            Process p   = runtime.exec(cmd);
            p.waitFor(10, TimeUnit.SECONDS);
            ProcessUtil.printProcessStatus(cmd, p);
            Thread.sleep(10000);
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withStrictlyServerHealthy(true)
                    .build();
            NebulaClient client    = pool.getClient();
            ResultSet    resultSet = client.execute("return 1");
            Assert.assertEquals(resultSet.getErrorCode(), ErrorCode.SUCCESSFUL_COMPLETION);
            pool.returnClient(client);
            pool.close();
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }
    }

    @Test
    public void testMultiThreads() {
        System.out.println("<==== testMultiThreads ====>");
        NebulaPool pool = null;
        try {
            pool = NebulaPool.builder(addresses, user, passwd).build();
            ExecutorService executorService = Executors.newFixedThreadPool(10);
            AtomicInteger   failedCount     = new AtomicInteger(0);
            CountDownLatch  countDownLatch  = new CountDownLatch(10);
            for (int i = 0; i < 10; i++) {
                NebulaPool finalPool = pool;
                executorService.submit(() -> {
                    try {
                        NebulaClient client = finalPool.getClient();
                        client.execute("SHOW GRAPHS");
                    } catch (Exception e) {
                        failedCount.incrementAndGet();
                    } finally {
                        countDownLatch.countDown();
                    }
                });
            }
            countDownLatch.await();
            executorService.shutdownNow();
        } catch (Exception e) {
            e.printStackTrace();
            Assert.fail(e.getMessage());
        } finally {
            if (pool != null) {
                pool.close();
            }
        }
    }
}
