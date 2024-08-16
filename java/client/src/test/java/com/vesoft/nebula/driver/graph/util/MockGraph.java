package com.vesoft.nebula.driver.graph.util;

import com.vesoft.nebula.driver.graph.data.ResultSet;
import com.vesoft.nebula.driver.graph.net.NebulaClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class MockGraph {
    private static final Logger log = LoggerFactory.getLogger(MockGraph.class);

    static String address = "127.0.0.1:9669";
    static String user    = "root";
    static String passwd  = "NebulaGraph01";

    public static void mockGraphData() {
        try {
            NebulaClient client = NebulaClient.builder(address, user, passwd).build();
            String createGraphType = "CREATE GRAPH TYPE IF NOT EXISTS graph_type_nba AS {"
                    + "NODE node_type_player (LABEL player "
                    + "{id INT PRIMARY KEY, name STRING, score FLOAT, gender bool, rate DOUBLE}),"
                    + "EDGE edge_type_follow (node_type_player)-[LABEL follow "
                    + "{followness INT, likeness FLOAT64}]->(node_type_player)"
                    + "}";

            ResultSet resultSet = client.execute(createGraphType);
            if (!resultSet.isSucceeded()) {
                log.error("create graph type `graph_type_nba` failed:{}",
                          resultSet.getErrorMessage());
                System.exit(1);
            }

            String createGraph = "CREATE GRAPH IF NOT EXISTS nba TYPED graph_type_nba";
            resultSet = client.execute(createGraph);
            if (!resultSet.isSucceeded()) {
                log.error("create graph `nba` failed:{}", resultSet.getErrorMessage());
                System.exit(1);
            }

            String insertNode = "USE nba INSERT OR REPLACE "
                    + "(@node_type_player{id:1, name:\"Tim\", "
                    + "score: 87.0, gender: true, rate: 7.32}),"
                    + "(@node_type_player{id:2, name:\"Jerry\", "
                    + "score: 95.0, gender: false, rate: 4.01}),"
                    + "(@node_type_player{id:3, name:\"Kyle\", "
                    + "score: 100, gender: true, rate: 9.99})";
            resultSet = client.execute(insertNode);
            if (!resultSet.isSucceeded()) {
                log.error("insert node failed for node_type_player:{}",
                          resultSet.getErrorMessage());
                System.exit(1);
            }

            String insertEdge = "table t{id1,id2,followness,likeness} = \n"
                    + "(1,2,90,66.8),(2,3,100,93.35)\n"
                    + "USE nba\n"
                    + "for r IN t\n"
                    + "MATCH(src@node_type_player) where src.id=r.id1\n"
                    + "MATCH(dst@node_type_player) where dst.id=r.id2\n"
                    + "INSERT OR REPLACE (src)-[@edge_type_follow"
                    + "{followness:r.followness,likeness:r.likeness}]->(dst)";
            resultSet = client.execute(insertEdge);
            if (!resultSet.isSucceeded()) {
                log.error("insert edge failed for edge_type_follow:{}",
                          resultSet.getErrorMessage());
                System.exit(1);
            }
        } catch (Exception e) {
            log.error("mock graph failed", e);
            System.exit(1);
        }
    }
}
