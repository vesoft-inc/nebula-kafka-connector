/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import java.net.UnknownHostException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import org.junit.Assert;
import org.junit.Test;

public class NebulaClientTest {
    String addresses = "192.168.8.6:3820";
    String user = "root";
    String passwd = "Nebula123";

    @Test()
    public void testBuildNebulaClient() {
        System.out.println("<==== testBuildNebulaClient ====>");
        try {
            NebulaClient client = NebulaClient.builder(addresses, user, passwd)
                    .setRequestTimeoutMills(3000)
                    .setRetryTimes(3)
                    .build();
            ResultSet resultSet = client.execute("return 1");
            Assert.assertEquals(resultSet.getErrorCode(), ErrorCode.SUCCESSFUL_COMPLETION);
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }
    }

    @Test()
    public void testBuildNebulaClientFailed() {
        System.out.println("<==== testBuildNebulaClientFailed ====>");
        // config illegal, illegal address
        String illegalAddress = "127.0.0.1:100000";
        try {
            NebulaClient client = NebulaClient.builder(illegalAddress, user, passwd).build();
            client.close();
        } catch (AuthFailedException | IOErrorException e) {
            Assert.fail(e.getMessage());
        } catch (IllegalArgumentException e) {
            assert true;
        }

        // config illegal, host name is not exist
        String illegalHostName = "hostname:9669";
        try {
            NebulaClient client = NebulaClient.builder(illegalHostName, user, passwd).build();
            client.close();
        } catch (RuntimeException e) {
            if (e.getMessage().contains("UnknownHostException")) {
                Assert.assertTrue("expect reach here:" + e.getMessage(), true);
            }
        } catch (IOErrorException | AuthFailedException e) {
            Assert.fail(e.getMessage());
        }


        // common builder
        NebulaClient.Builder builder = null;
        try {
            builder = NebulaClient.builder(addresses, user, passwd);
        } catch (RuntimeException e) {
            Assert.fail(e.getMessage());
        }

        // config illegal, illegal requestTimeout
        try {
            builder.setRequestTimeoutMills(-1).build();
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }

        // config retry times
        try {
            builder.setRetryTimes(-1).build();
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }

        // config interval time
        try {
            builder.setIntervalTimeMills(-1).build();
        } catch (IOErrorException | AuthFailedException e) {
            Assert.fail(e.getMessage());
        }

    }

    @Test()
    public void testStrictlyServerHealthy() {
        System.out.println("<==== testStrictlyServerHealthy ====>");
        // TODO stop one graphd server
        String address = "127.0.0.1:9669,127.0.0.1:9670,127.0.0.1:9671";
        try {
            NebulaClient client = NebulaClient.builder(address, user, passwd)
                    .build();
            ResultSet resultSet = client.execute("return 1");
            Assert.assertEquals(resultSet.getErrorCode(), ErrorCode.SUCCESSFUL_COMPLETION);
            client.close();
        } catch (AuthFailedException | IOErrorException e) {
            Assert.fail(e.getMessage());
        }

        try {
            NebulaClient.builder(address, user, passwd).build();
            Assert.fail("client build should failed, strictlyServerHealthy is true, "
                    + "graphd servers are not all ok.");
        } catch (AuthFailedException e) {
            Assert.fail(e.getMessage());
        } catch (IOErrorException e) {
            Assert.assertTrue("expect here:" + e.getMessage(),
                    e.getMessage().contains("Servers status is not ok"));
        }

        // TODO start the graphd server
    }

    @Test
    public void testCloseNebulaClient() {
        System.out.println("<==== testCloseNebulaClient ====>");
        NebulaClient client = null;
        try {
            client = NebulaClient.builder(addresses, user, passwd)
                    .build();
        } catch (IOErrorException | AuthFailedException e) {
            Assert.fail(e.getMessage());
        }
        client.close();
        // test the idempotence of close()
        client.close();

        // test the check for execute using closed client
        try {
            client.execute("RETURN 1");
        } catch (RuntimeException e) {
            Assert.assertEquals(e.getMessage(), "NebulaClient has closed. Couldn't use again.");
        } catch (IOErrorException e) {
            Assert.fail(e.getMessage());
        }
    }

    @Test
    public void testExecuteRetryForExecutionError() {
        System.out.println("<==== testExecuteRetryForExecutionError ====>");
        // retry 50 times at most, wait 5 s between retry.
        NebulaClient client = null;
        try {
            client = NebulaClient.builder(addresses, user, passwd)
                    .setRetryTimes(50)
                    .setIntervalTimeMills(5000)
                    .build();
        } catch (IOErrorException | AuthFailedException e) {
            Assert.fail(e.getMessage());
        }

        // TODO stop one storaged server
        ExecutorService executorService = Executors.newFixedThreadPool(10);
        NebulaClient finalClient = client;
        AtomicInteger succeedCount = new AtomicInteger(0);
        for (int i = 0; i < 10; i++) {
            executorService.submit(() -> {
                try {
                    ResultSet resultSet = finalClient.execute("MATCH(v) RETURN v limit 10");
                    if (resultSet.isSucceeded()) {
                        succeedCount.incrementAndGet();
                    }
                } catch (Exception e) {
                    Assert.fail(e.getMessage());
                }
            });
        }

        // TODO start the storaged server
        executorService.shutdown();
        try {
            executorService.awaitTermination(Long.MAX_VALUE, TimeUnit.SECONDS);
        } catch (InterruptedException e) {
            e.printStackTrace();
            assert false;
        }
        client.close();
        Assert.assertEquals(10, succeedCount.get());
    }

    @Test
    public void testExecuteRetryForSessionError() {
        // TODO set the session's invalidate time is 3s

        System.out.println("<==== testExecuteRetryForSessionError ====>");
        // retry 50 times at most, wait 2s between retry.
        NebulaClient client = null;
        try {
            client = NebulaClient.builder(addresses, user, passwd)
                    .setRetryTimes(50)
                    .setIntervalTimeMills(2000)
                    .build();
        } catch (IOErrorException | AuthFailedException e) {
            Assert.fail(e.getMessage());
        }
        try {
            Thread.sleep(3000);
            ResultSet resultSet = client.execute("RETURN 1");
            assert (resultSet.isSucceeded());
        } catch (IOErrorException | InterruptedException e) {
            Assert.fail(e.getMessage());
        }
    }

    @Test
    public void testExecuteRetryForConnectionError() {
        System.out.println("<==== testExecuteRetryForConnectionError ====>");
        // retry 50 times at most, wait 2s between retry.
        NebulaClient client = null;
        try {
            client = NebulaClient.builder(addresses, user, passwd)
                    .setRetryTimes(50)
                    .setIntervalTimeMills(2000)
                    .build();
        } catch (IOErrorException | AuthFailedException e) {
            Assert.fail(e.getMessage());
        }

        // TODO stop one graphd server
        ExecutorService executorService = Executors.newFixedThreadPool(10);
        NebulaClient finalClient = client;
        AtomicInteger succeedCount = new AtomicInteger(0);
        for (int i = 0; i < 10; i++) {
            executorService.submit(() -> {
                try {
                    ResultSet resultSet = finalClient.execute("RETURN 1");
                    if (resultSet.isSucceeded()) {
                        succeedCount.incrementAndGet();
                    }
                } catch (Exception e) {
                    Assert.fail(e.getMessage());
                }
            });
        }

        // TODO start the graphd server
        executorService.shutdown();
        try {
            executorService.awaitTermination(Long.MAX_VALUE, TimeUnit.SECONDS);
        } catch (InterruptedException e) {
            e.printStackTrace();
            assert false;
        }
        client.close();
        Assert.assertEquals(10, succeedCount.get());
    }

    @Test
    public void testIntervalTimeForExecuteRetry() {
        System.out.println("<==== testIntervalTimeForExecuteRetry ====>");
        // retry 50 times at most, wait 2s between retry.
        NebulaClient client = null;
        try {
            client = NebulaClient.builder(addresses, user, passwd)
                    .setRetryTimes(50)
                    .setIntervalTimeMills(2000)
                    .build();
        } catch (AuthFailedException | IOErrorException e) {
            Assert.fail(e.getMessage());
        }
        // TODO stop all graphd servers
        long beginTime = System.currentTimeMillis();
        ExecutorService executorService = Executors.newFixedThreadPool(1);
        NebulaClient finalClient = client;
        executorService.submit(() -> {
            try {
                ResultSet resultSet = finalClient.execute("RETURN 1");
                Assert.assertTrue(resultSet.isSucceeded());
            } catch (IOErrorException e) {
                Assert.fail(e.getMessage());
            }
            Assert.assertTrue((System.currentTimeMillis() - beginTime) >= 2000);
        });
        // TODO start all graphd servers
    }

    @Test
    public void testBadSessionRelease() {
        NebulaClient.Builder clientBuilder = null;
        NebulaClient nebulaClient = null;
        try {
            nebulaClient = NebulaClient.builder(addresses, user, passwd)
                    .setRetryTimes(50)
                    .setIntervalTimeMills(2000)
                    .build();
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }
        // TODO close one graphd server
        try {
            Thread.sleep(10000);
        } catch (InterruptedException e) {
            Assert.fail(e.getMessage());
        }

        nebulaClient.close();

        String goodAddresss = "127.0.0.1:9670";
        NebulaClient client = null;
        try {
            client = NebulaClient.builder(goodAddresss, user, passwd).build();
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }

        try {
            // TODO depend on server's `show sessions` statement.
            ResultSet resultSet = client.execute("show sessions");
            assert (resultSet.isSucceeded());
            assert (resultSet.rowSize() == 1);
        } catch (IOErrorException e) {
            Assert.fail(e.getMessage());
        }
    }

    @Test
    public void testCloseClient() {
        NebulaClient.Builder clientBuilder = null;
        NebulaClient nebulaClient = null;
        try {
            nebulaClient = NebulaClient.builder(addresses, user, passwd)
                    .setRetryTimes(50)
                    .setIntervalTimeMills(2000)
                    .build();
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }

        // execution lasts 20 seconds
        long start = System.currentTimeMillis();
        while ((System.currentTimeMillis() - start) < 20000) {
            try {
                nebulaClient.execute("RETURN 1");
            } catch (IOErrorException e) {
                Assert.fail(e.getMessage());
            }
        }
        nebulaClient.close();
    }
}
