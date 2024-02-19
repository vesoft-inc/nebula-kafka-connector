/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.util;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.net.NebulaClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class MockGraph {
    private static final Logger log = LoggerFactory.getLogger(MockGraph.class);

    static String address = "192.168.8.6:3820";
    static String user = "root";
    static String passwd = "nebula";

    public static void mockGraphData() {
        try {
            NebulaClient client = NebulaClient.builder(address, user, passwd).build();
            String createGraphType = "CREATE GRAPH TYPE graph_type_nba AS "
                    + "{(node_type_player LABEL player "
                    + "{id INT PRIMARY KEY, name STRING, score FLOAT, gender bool, rate DOUBLE}),"
                    + "(node_type_player)-[edge_type_follow LABEL follow "
                    + "{followness INT, likeness FLOAT64}]->(node_type_player)}";

            ResultSet resultSet = client.execute(createGraphType);
            if (!resultSet.isSucceeded()) {
                log.error("create graph type `graph_type_nba` failed.",
                        resultSet.getErrorMessage());
                System.exit(1);
            }

            String createGraph = "CREATE GRAPH nba graph_type_nba";
            resultSet = client.execute(createGraph);
            if (!resultSet.isSucceeded()) {
                log.error("create graph `nba` failed.", resultSet.getErrorMessage());
                System.exit(1);
            }

            String insertNode = "USE nba INSERT NODE node_type_player "
                    + "({id:1, name:\"Tim\", score: 87.0, gender: true, rate: 7.32}),"
                    + "({id:2, name:\"Jerry\", score: 95.0, gender: false, rate: 4.01}),"
                    + "({id:3, name:\"Kyle\", score: 100, gender: true, rate: 9.99})";
            resultSet = client.execute(insertNode);
            if (!resultSet.isSucceeded()) {
                log.error("insert node failed for node_type_player.", resultSet.getErrorMessage());
                System.exit(1);
            }

            String insertEdge = "USE nba INSERT EDGE edge_type_follow ({id:1})-[{followness:90, "
                    + "likeness: 66.8}]->({id:2}),({id:2})-[{followness:100, likeness: 93.35}]->"
                    + "({id:3})";
            resultSet = client.execute(insertEdge);
            if (!resultSet.isSucceeded()) {
                log.error("insert edge failed for edge_type_follow.", resultSet.getErrorMessage());
                System.exit(1);
            }
        } catch (Exception e) {
            log.error("mock graph failed", e);
            System.exit(1);
        }
    }
}
