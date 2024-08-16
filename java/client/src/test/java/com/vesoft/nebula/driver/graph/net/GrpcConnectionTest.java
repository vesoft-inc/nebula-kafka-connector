package com.vesoft.nebula.driver.graph.net;

import com.google.common.base.Charsets;
import com.vesoft.nebula.driver.graph.data.HostAddress;
import com.vesoft.nebula.driver.graph.util.MockGraph;
import com.vesoft.nebula.proto.graph.ExecuteResponse;
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

    @Before
    public void setup() {
        MockGraph.mockGraphData();
    }


    @Test(timeout = 3000)
    public void testAll() {
        try {
            // Test open
            GrpcConnection connection = new GrpcConnection();
            connection.open(new HostAddress(host, port), 1000, 1000);

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
}
