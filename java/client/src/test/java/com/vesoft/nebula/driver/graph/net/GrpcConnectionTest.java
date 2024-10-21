package com.vesoft.nebula.driver.graph.net;

import com.google.common.base.Charsets;
import com.vesoft.nebula.driver.graph.data.HostAddress;
import com.vesoft.nebula.driver.graph.util.MockGraph;
import com.vesoft.nebula.proto.graph.ExecuteResponse;
import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class GrpcConnectionTest {
    private static final Logger LOGGER   = LoggerFactory.getLogger(GrpcConnection.class);
    private final        String host     = "127.0.0.1";
    private final        int    port     = 9669;
    private final        String user     = "root";
    private final        String password = "NebulaGraph01";

    private NebulaClient.Builder builder;

    @Before
    public void setup() {
        MockGraph.mockGraphData();
        builder = NebulaClient
                .builder(host + ":" + port, "root", "Nebula123")
                .withConnectTimeoutMills(1000)
                .withRequestTimeoutMills(1000)
                .withEnableTls(false);
    }


    @Test(timeout = 3000)
    public void testAll() {
        try {
            // Test open
            GrpcConnection connection = new GrpcConnection();
            NebulaClient.Builder builder = NebulaClient
                    .builder(host + port, "root", "Nebula123")
                    .withConnectTimeoutMills(1000)
                    .withRequestTimeoutMills(1000)
                    .withEnableTls(false);
            connection.open(new HostAddress(host, port), builder);

            // Test authenticate
            Map<String, Object> authInfo = new HashMap<>();
            authInfo.put("password", password);
            AuthResult authResult = connection.authenticate(user, authInfo);
            LOGGER.info(authResult.toString());

            Assert.assertNotEquals(0, authResult.getSessionId());

            // Test ping
            assert connection.ping(authResult.getSessionId(), 1000);

            // Test execute
            ExecuteResponse resp = connection.execute(authResult.getSessionId(), "SHOW GRAPHS");
            LOGGER.info(resp.toString());
            assert resp.getStatus().getCode().toString(Charsets.UTF_8).equals("00000");


            resp = connection.execute(authResult.getSessionId(), "USE nba\n"
                    + " MATCH (v:player WHERE v.id IN LIST[2, 3, 4, 888])-[e:follow]->{0,1}(b)\n"
                    + " RETURN v.id as src, b.id as dst");
            LOGGER.info(resp.toString());
            assert resp.getStatus().getCode().toString(Charsets.UTF_8).equals("00000");


            resp = connection.execute(authResult.getSessionId(), "RETURN RECORD {id: 1}.id");
            LOGGER.info(resp.toString());
            assert resp.getStatus().getCode().toString(Charsets.UTF_8).equals("00000");


            // Test sign out
            resp = connection.execute(authResult.getSessionId(), "SESSION CLOSE");
            LOGGER.info(resp.toString());
            assert resp.getStatus().getCode().toString(Charsets.UTF_8).equals("00000");
        } catch (Exception e) {
            e.printStackTrace();
            assert (false);
        }
    }

    @Test
    public void testNotReuseChannel() {
        Map<String, Object> authInfo = new HashMap<>();
        authInfo.put("password", password);
        GrpcConnection connection1 = null;
        GrpcConnection connection2 = null;

        try {
            connection1 = new GrpcConnection();
            connection1.open(new HostAddress(host, port), builder);
            AuthResult authResult = connection1.authenticate(user, authInfo);

            connection2 = new GrpcConnection();
            connection2.open(new HostAddress(host, port), builder);
            authResult = connection2.authenticate(user, authInfo);

            Thread.sleep(3000);
            // close connection1 and test whether connection2 is ok
            connection1.close();
            assert connection2.ping(authResult.getSessionId(), 1000);
        } catch (Exception e) {
            e.printStackTrace();
            assert false;
        } finally {
            if (connection1 != null) {
                connection1.close();
            }
            if (connection2 != null) {
                connection2.close();
            }
        }
    }

    @Test
    public void testTlsConnection() {
        String              host     = "192.168.8.6";
        int                 port     = 4820;
        Map<String, Object> authInfo = new HashMap<>();
        authInfo.put("password", "Nebula123");

        String tlsCa   = "src/test/resources/tls/ca.pem";
        String tlsCert = "src/test/resources/tls/client/client.cert";
        String tlsKey  = "src/test/resources/tls/client/client-private.key";

        NebulaClient.Builder tlsBuilder = NebulaClient
                .builder(host + ":" + port, "root", "Nebula123")
                .withConnectTimeoutMills(1000)
                .withRequestTimeoutMills(1000)
                .withEnableTls(true)
                .withTlsCa(tlsCa)
                .withTlsCert(tlsCert, tlsKey)
                .withTlsPeerName("NICOLE");

        try {
            GrpcConnection connection = new GrpcConnection();
            connection.open(new HostAddress(host, port), tlsBuilder);
            AuthResult authResult = connection.authenticate(user, authInfo);
            assert true;
        } catch (Exception e) {
            e.printStackTrace();
            assert (false);
        }
    }
}
