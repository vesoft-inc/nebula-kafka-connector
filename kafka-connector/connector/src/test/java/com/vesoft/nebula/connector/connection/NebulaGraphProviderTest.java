
package com.vesoft.nebula.connector.connection;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_ADDRESS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_DATA_TYPE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_EDGE_TYPE_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_NODE_TYPE_NAME;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_PASSWD;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_GRAPH_USER;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_PRIMARY_KEY;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import com.vesoft.nebula.connector.sink.NebulaEdgeSchema;
import com.vesoft.nebula.connector.sink.NebulaNodeSchema;
import java.util.HashMap;
import java.util.Map;
import org.junit.Before;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class NebulaGraphProviderTest {
    private static final Logger log = LoggerFactory.getLogger(NebulaGraphProviderTest.class);

    String host       = "192.168.8.6:3820";
    String user       = "root";
    String passwd     = "Nebula123";
    String graphName  = "nba";
    String graphType  = "NODE";
    String nodeType   = "node_type_player";
    String primaryKey = "id";
    String edgeType   = "edge_type_follow";

    NebulaGraphProvider provider;

    @Before
    public void setup() {
        Map<String, Object> props = new HashMap<>();
        props.put(CONNECT_GRAPH_ADDRESS, host);
        props.put(CONNECT_GRAPH_USER, user);
        props.put(CONNECT_GRAPH_PASSWD, passwd);
        props.put(CONNECT_GRAPH_NAME, graphName);
        props.put(CONNECT_GRAPH_DATA_TYPE, graphType);
        props.put(CONNECT_GRAPH_NODE_TYPE_NAME, nodeType);
        props.put(CONNECT_GRAPH_EDGE_TYPE_NAME, edgeType);
        props.put(CONNECT_PRIMARY_KEY, primaryKey);
        NebulaSinkConnectConfig config = new NebulaSinkConnectConfig(props);
        provider = new NebulaGraphProvider(config);
        mockGraphSchema();
    }

    @Test
    public void testGetNodeSchema() {
        try {
            NebulaNodeSchema schema = provider.getNodeSchema(graphName, nodeType);
            assert (schema.getNodeTypeName().equalsIgnoreCase(nodeType));
            assert (schema.getNodePkName().equals("id"));
            assert (schema.getNodePkType().equalsIgnoreCase("INT64"));
            assert (schema.getNodeProperties().size() == 5);
        } catch (IOErrorException | NoValidSessionException e) {
            e.printStackTrace();
            assert false;
        }
    }


    @Test
    public void testGetEdgeSchema() {
        try {
            NebulaEdgeSchema schema = provider.getEdgeSchema(graphName, edgeType);
            assert (schema.getEdgeTypeName().equalsIgnoreCase(edgeType));
            assert (schema.getSourceNodeTypeName().equalsIgnoreCase(nodeType));
            assert (schema.getTargetNodeTypeName().equalsIgnoreCase(nodeType));
            assert (schema.getSourceNodePkName().equals("id"));
            assert (schema.getTargetNodePkName().equals("id"));
            assert (schema.getSourceNodePkType().equalsIgnoreCase("INT64"));
            assert (schema.getTargetNodePkType().equalsIgnoreCase("INT64"));
            assert (schema.getProperties().size() == 2);
        } catch (Exception e) {
            e.printStackTrace();
            assert false;
        }
    }


    private void mockGraphSchema() {
        String createSchema = "CREATE GRAPH TYPE IF NOT EXISTS graph_type_nba AS {"
                + "NODE TYPE node_type_player (LABEL player {id INT PRIMARY KEY, name STRING, "
                + "score FLOAT, gender bool, rate DOUBLE}),"
                + "EDGE TYPE edge_type_follow(node_type_player)-[LABEL follow "
                + "{followness INT, likeness FLOAT64}]->(node_type_player)}";
        String createGraph = "CREATE GRAPH IF NOT EXISTS nba graph_type_nba";
        try {
            ResultSet result = provider.execute(createSchema);
            if (!result.isSucceeded()) {
                log.error(">>>>> create graph schema failed, {}", result.getErrorMessage());
                assert false;
            }
            result = provider.execute(createGraph);
            if (!result.isSucceeded()) {
                log.error(">>>>> create graph failed, {}", result.getErrorMessage());
                assert false;
            }
        } catch (Exception e) {
            e.printStackTrace();
            assert false;
            System.exit(-1);
        }
    }
}
