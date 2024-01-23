/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.net;

import com.google.protobuf.ByteString;
import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.proto.ExecuteResponse;
import java.nio.charset.Charset;
import org.junit.Assert;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class TestGrpcConnection {
    private static final Logger LOGGER = LoggerFactory.getLogger(GrpcConnection.class);

    @Test(timeout = 3000)
    public void testAll() {
        try {
            // Test open
            GrpcConnection connection = new GrpcConnection();
            connection.open(new HostAddress("127.0.0.1", 20584), 1000, 1000);

            // Test authenticate
            AuthResult authResult = connection.authenticate("root", "nebula");
            LOGGER.info(authResult.toString());

            Assert.assertNotEquals(0, authResult.getSessionId());

            // Test execute
            ExecuteResponse resp = connection.execute(authResult.getSessionId(), "SHOW GRAPHS");

            LOGGER.info(resp.toString());

            resp = connection.execute(authResult.getSessionId(), "USE ldbc\n"
                    + " MATCH (v:Person WHERE v.id IN LIST[2, 3, 4, 888])-[e:KNOWS]->{0,1}(b)\n"
                    + " RETURN v.id as src, b.id as dst");

            LOGGER.info(resp.toString());

            resp = connection.execute(authResult.getSessionId(), "RETURN RECORD {id: 1}.id");

            LOGGER.info(resp.toString());

            // Test sign out
            connection.signout(authResult.getSessionId());

        } catch (Exception e) {
            e.printStackTrace();
            assert (false);
        }
    }
}
