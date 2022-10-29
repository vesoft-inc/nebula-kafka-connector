/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.client.graph.NebulaPoolConfig;
import com.vesoft.nebula.client.graph.net.NebulaPool;
import com.vesoft.nebula.client.graph.net.Session;
import java.util.Arrays;
import java.util.concurrent.TimeUnit;
import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class TestDataFromServer {
    private final NebulaPool pool = new NebulaPool();
    private Session session = null;

    @Before
    public void setUp() throws Exception {
        NebulaPoolConfig nebulaPoolConfig = new NebulaPoolConfig();
        nebulaPoolConfig.setMaxConnSize(1);
        Assert.assertTrue(pool.init(Arrays.asList(new HostAddress(
                        "127.0.0.1", 10669)),
                nebulaPoolConfig));
        session = pool.getSession("root", "nebula", true);
        ResultSet resp = session.execute(
                "CREATE GRAPH TYPE graph_type IF NOT EXISTS AS GRAPH TYPE {"
                        + "(node_type(id) LABEL player {id INT, name STRING}),"
                        + "(node_type)-[edge_type LABEL follow {followness INT}]->(node_type)}");
        Assert.assertTrue(resp.getGqlStatus(), resp.isSucceeded());
        TimeUnit.SECONDS.sleep(10);

        resp = session.execute("CREATE GRAPH nba IF NOT EXISTS OF graph_type");
        Assert.assertTrue(resp.getGqlStatus(), resp.isSucceeded());
        TimeUnit.SECONDS.sleep(10);

        resp = session.execute("USE nba INSERT NODE node_type ({id:1, name:\"Tim\"}),"
                + "({id:2, name:\"Jerry\"}),({id:3, name:\"Kyle\"})");
        Assert.assertTrue(resp.getGqlStatus(), resp.isSucceeded());
        TimeUnit.SECONDS.sleep(10);

        resp = session.execute("USE nba INSERT EDGE edge_type ({id:1})-[{followness:90}]->"
                + "({id:2}),({id:2})-[{followness:100}]->({id:3})");
        Assert.assertTrue(resp.getGqlStatus(), resp.isSucceeded());
        TimeUnit.SECONDS.sleep(10);
    }

    @After
    public void tearDown() {
        if (session != null) {
            session.release();
        }
        pool.close();
    }

    @Test
    public void testPropData() {
        try {
            ResultSet result = session.execute("FROM nba MATCH (v:player) RETURN v");
            assert (result.isSucceeded());
            Assert.assertEquals(3, result.rowsSize());
            result.toString();

            result = session.execute(
                    "FROM nba MATCH (m)-[:follow]->(p) RETURN m.id, p.id");
            Assert.assertTrue(result.isSucceeded());
            Assert.assertEquals(2, result.rowsSize());
            result.toString();

            result = session.execute("from nba match (m)-[:follow]->(p) return value "
                    + "{from nba match pa=(m)-[:follow]->(p) return pa} AS pas");
            Assert.assertTrue(result.getGqlStatus(), result.isSucceeded());
            Assert.assertEquals(2, result.rowsSize());
            result.toString();
        } catch (Exception e) {
            e.printStackTrace();
            assert false;
        } finally {
            session.release();
        }

    }

    @Test
    public void testErrorResult() {
        try {
            ResultSet result = session.execute(
                    "FETCH PROP ON no_exist_tag \"nobody\" yield vertex as vertices_");
            Assert.assertTrue(result.toString().contains("ExecutionResponse"));
            Assert.assertTrue(result.getGqlStatus().startsWith("ERROR"));
        } catch (Exception e) {
            e.printStackTrace();
            assert false;
        }
    }
}
