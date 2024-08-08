package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.util.ProcessUtil;
import java.util.concurrent.TimeUnit;
import org.junit.Assert;
import org.junit.Test;

public class NebulaClientTest {
    String addresses = "127.0.0.1:9669,127.0.0.1:9670,127.0.0.1:9671";
    String user      = "root";
    String passwd    = "NebulaGraph01";

    @Test()
    public void testBuildNebulaClient() {
        System.out.println("<==== testBuildNebulaClient ====>");
        try {
            NebulaClient client = NebulaClient.builder(addresses, user, passwd)
                    .withRequestTimeoutMills(3000)
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
            builder.withRequestTimeoutMills(-1).build();
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        }
    }

    @Test
    public void testReconnectWithMultiServices() {
        System.out.println("<==== testReconnectWithMultiServices ====>");
        Runtime runtime = Runtime.getRuntime();
        String  cmd     = "docker stop docker-compose-graphd0-1";

        try {
            NebulaClient.builder("127.0.0.1:9669", user, passwd).build();
        } catch (Exception e) {
            assert true;
        }

        try {
            Process p = runtime.exec(cmd);
            p.waitFor(10, TimeUnit.SECONDS);
            ProcessUtil.printProcessStatus(cmd, p);
            Thread.sleep(5000);
            NebulaClient client    = NebulaClient.builder(addresses, user, passwd).build();
            ResultSet    resultSet = client.execute("RETURN 1");
            assert resultSet.isSucceeded();
        } catch (Exception e) {
            Assert.fail(e.getMessage());
        } finally {
            try {
                cmd = "docker start docker-compose-graphd0-1";
                Process p = runtime.exec(cmd);
                p.waitFor(10, TimeUnit.SECONDS);
                ProcessUtil.printProcessStatus(cmd, p);
                Thread.sleep(10000);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
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
            Assert.fail("client execute should fail.");
        } catch (RuntimeException e) {
            Assert.assertEquals(e.getMessage(), "The NebulaClient already closed.");
        } catch (IOErrorException e) {
            Assert.fail(e.getMessage());
        }
    }

    @Test
    public void testCloseClient() {
        NebulaClient nebulaClient = null;
        try {
            nebulaClient = NebulaClient.builder(addresses, user, passwd)
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

        try {
            nebulaClient.execute("RETURN 1");
            Assert.fail("client execute should fail.");
        } catch (RuntimeException e) {
            assert true;
        } catch (IOErrorException e) {
            Assert.fail(e.getMessage());
        }
    }

    @Test
    public void testPingClient() {
        System.out.println("<==== testPingClient ====>");
        NebulaClient client = null;
        try {
            client = NebulaClient.builder(addresses, user, passwd)
                    .build();
        } catch (IOErrorException | AuthFailedException e) {
            Assert.fail(e.getMessage());
        }
        assert client.ping();
    }

    @Test
    public void testGetHost() {
        System.out.println("<==== testGetHost ====>");
        NebulaClient client = null;
        try {
            for (int i = 0; i < 10; i++) {
                client = NebulaClient.builder(addresses, user, passwd)
                        .build();
                System.out.println(">>>> host:" + client.getHost());
                assert client.getHost().equals("127.0.0.1:9669")
                        || client.getHost().equals("127.0.0.1:9670")
                        || client.getHost().equals("127.0.0.1:9671");
            }
        } catch (IOErrorException | AuthFailedException e) {
            Assert.fail(e.getMessage());
        }
        assert client.ping();
    }
}
