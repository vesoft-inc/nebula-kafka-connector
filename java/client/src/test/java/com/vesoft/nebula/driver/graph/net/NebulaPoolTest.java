/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.driver.graph.net;

import com.vesoft.nebula.driver.graph.ErrorCode;
import com.vesoft.nebula.driver.graph.data.ResultSet;
import com.vesoft.nebula.driver.graph.util.MockGraph;
import com.vesoft.nebula.driver.graph.util.ProcessUtil;
import java.time.ZoneId;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class NebulaPoolTest {
    String addresses = "127.0.0.1:9669";
    String user      = "root";
    String passwd    = "NebulaGraph01";

    @Before
    public void setup() {
        MockGraph.mockGraphData();
    }

    @Test
    public void testBuilder() {
        System.out.println("<==== testBuilder =====>");
        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withConnectTimeoutMills(1111)
                    .withRequestTimeoutMills(2222)
                    .withScanParallel(15)
                    .withHealthCheckTimeMills(3333)
                    .build();
            NebulaClient client = pool.getClient();
            Assert.assertEquals(1111L, client.getConnectTimeoutMills());
            Assert.assertEquals(2222L, client.getRequestTimeoutMills());
            Assert.assertEquals(15, client.getScanParallel());
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }

        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withConnectTimeoutMills(0)
                    .withRequestTimeoutMills(-1)
                    .build();
            NebulaClient client = pool.getClient();
            Assert.assertEquals(Integer.MAX_VALUE, client.getConnectTimeoutMills());
            Assert.assertEquals(Integer.MAX_VALUE, client.getRequestTimeoutMills());
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }
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
    public void testSessionSet() {
        System.out.println("<==== testSessionSet ====>");
        // test all session set configs
        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withGraph("nba")
                    .withSchema("/default_schema")
                    .withTimeZone(ZoneId.systemDefault().getId())
                    .build();
            Assert.assertEquals("nba",
                                pool
                                        .getClient()
                                        .execute("show current_session")
                                        .next()
                                        .get("home_graph_name")
                                        .asString());
            Assert.assertEquals("/default_schema",
                                pool
                                        .getClient()
                                        .execute("show current_session")
                                        .next()
                                        .get("home_schema_path")
                                        .asString());
            Assert.assertEquals("Asia/Shanghai",
                                pool
                                        .getClient()
                                        .execute("show current_session")
                                        .next()
                                        .get("time_zone")
                                        .asString());
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }

        // test one session set config

        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withGraph("nba")
                    .build();
            Assert.assertEquals("nba",
                                pool
                                        .getClient()
                                        .execute("show current_session")
                                        .next()
                                        .get("home_graph_name")
                                        .asString());
        } catch (Exception e) {
            assert e.getMessage().contains("SESSION SET home_graph_name=\"nba_not_exist\" failed");
        }

        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withSchema("/default_schema")
                    .build();
            Assert.assertEquals("/default_schema",
                                pool
                                        .getClient()
                                        .execute("show current_session")
                                        .next()
                                        .get("home_schema_path")
                                        .asString());

        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }

        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withTimeZone("UTC")
                    .build();
            Assert.assertEquals("UTC",
                                pool
                                        .getClient()
                                        .execute("show current_session")
                                        .next()
                                        .get("time_zone")
                                        .asString());

        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }

        // test wrong graph
        try {
            NebulaPool pool = NebulaPool.builder(addresses, user, passwd)
                    .withGraph("nba_not_exist")
                    .build();
            pool.getClient();
            Assert.fail("get client should fail.");
        } catch (Exception e) {
            assert e.getMessage().contains("SESSION SET home_graph_name=\"nba_not_exist\" failed");
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

    @Test
    public void testTls() {
        System.out.println("<==== testTls ====>");
        String address = "192.168.8.6:4820";


        String tlsCa   = "src/test/resources/tls/ca.pem";
        String tlsCert = "src/test/resources/tls/client/client.cert";
        String tlsKey  = "src/test/resources/tls/client/client-private.key";

        try {
            NebulaPool pool = NebulaPool
                    .builder(address, "root", "Nebula123")
                    .withEnableTls(true)
                    .withTlsCa(tlsCa)
                    .withTlsCert(tlsCert, tlsKey)
                    .withTlsPeerName("NICOLE")
                    .build();
            NebulaClient client = pool.getClient();
            client.execute("RETURN 1");
            assert true;
        } catch (Exception e) {
            e.printStackTrace();
            assert false;
        }
    }
}
