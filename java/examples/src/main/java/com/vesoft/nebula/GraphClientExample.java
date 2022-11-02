/* Copyright (c) 2022 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula;

import com.vesoft.nebula.client.graph.NebulaPoolConfig;
import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.data.ValueWrapper;
import com.vesoft.nebula.client.graph.net.NebulaPool;
import com.vesoft.nebula.client.graph.net.Session;
import java.io.UnsupportedEncodingException;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.TimeUnit;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class GraphClientExample {
    private static final Logger log = LoggerFactory.getLogger(GraphClientExample.class);
    static String ip = "127.0.0.1";
    static int port = 10669;

    public static void main(String[] args) {
        NebulaPool pool = new NebulaPool();
        Session session = null;
        try {
            // init the NebulaPool and get session
            NebulaPoolConfig nebulaPoolConfig = new NebulaPoolConfig();
            nebulaPoolConfig.setMaxConnSize(100);
            List<HostAddress> addresses = Arrays.asList(new HostAddress(ip, port));
            boolean initResult = pool.init(addresses, nebulaPoolConfig);
            if (!initResult) {
                log.error("pool init failed.");
                return;
            }
            session = pool.getSession("root", "nebula", false);

            // create schema and insert data
            String createSchema = "CREATE GRAPH TYPE graph_type IF NOT EXISTS AS GRAPH TYPE {"
                    + "(node_type(id) LABEL player {id INT, name STRING}),"
                    + "(node_type)-[edge_type LABEL follow {followness INT}]->(node_type)}";
            ResultSet resp = session.execute(createSchema);
            if (!resp.isSucceeded()) {
                log.error(String.format("Execute: `%s', failed: %s",
                        createSchema, resp.getGqlStatus()));
                System.exit(1);
            }
            TimeUnit.SECONDS.sleep(5);

            String createGraph = "CREATE GRAPH nba IF NOT EXISTS OF graph_type";
            resp = session.execute(createGraph);
            if (!resp.isSucceeded()) {
                log.error(String.format("Execute `%s`, failed: %s", createGraph,
                        resp.getGqlStatus()));
                System.exit(1);
            }
            TimeUnit.SECONDS.sleep(5);

            String insertVertexes = "USE nba INSERT NODE node_type ({id:1, name:\"Tim\"}),"
                    + "({id:2, name:\"Jerry\"}),({id:3, name:\"Kyle\"})";
            resp = session.execute(insertVertexes);
            if (!resp.isSucceeded()) {
                log.error(String.format("Execute: `%s', failed: %s",
                        insertVertexes, resp.getGqlStatus()));
                System.exit(1);
            }

            String insertEdges = "USE nba INSERT EDGE edge_type ({id:1})-[{followness:90}]->"
                    + "({id:2}),({id:2})-[{followness:100}]->({id:3})";
            resp = session.execute(insertEdges);
            if (!resp.isSucceeded()) {
                log.error(String.format("Execute: `%s', failed: %s",
                        insertEdges, resp.getGqlStatus()));
                System.exit(1);
            }


            // make query
            String query = "from nba match (m)-[:follow]->(p) return value "
                    + "{from nba match pa=(m)-[:follow]->(p) return pa} AS pas";
            resp = session.execute(query);
            if (!resp.isSucceeded()) {
                log.error(String.format("Execute: `%s', failed: %s",
                        query, resp.getGqlStatus()));
                System.exit(1);
            }
            printResult(resp);

        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        } finally {
            if (session != null) {
                session.release();
            }
            pool.close();
        }
    }


    private static void printResult(ResultSet resultSet) throws UnsupportedEncodingException {
        List<String> colNames = resultSet.keys();
        System.out.print("| ");
        for (String name : colNames) {
            System.out.printf("%s |", name);
        }
        System.out.println();
        for (int i = 0; i < resultSet.rowsSize(); i++) {
            System.out.print("| ");
            ResultSet.Record record = resultSet.rowValues(i);
            for (ValueWrapper value : record.values()) {
                if (value.isLong()) {
                    System.out.printf("%15s |", value.asLong());
                }
                if (value.isBoolean()) {
                    System.out.printf("%15s |", value.asBoolean());
                }
                if (value.isDouble()) {
                    System.out.printf("%15s |", value.asDouble());
                }
                if (value.isString()) {
                    System.out.printf("%15s |", value.asString());
                }
                if (value.isNode()) {
                    System.out.printf("%15s |", value.asNode());
                }
                if (value.isEdge()) {
                    System.out.printf("%15s |", value.asEdge());
                }
                if (value.isList()) {
                    System.out.printf("%15s |", value.asList());
                }
            }
            System.out.println();
        }
    }
}
