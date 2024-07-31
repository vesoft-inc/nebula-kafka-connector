package com.vesoft.nebula.client.scan;

import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.net.NebulaClient;
import com.vesoft.nebula.client.util.MockGraph;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.net.UnknownHostException;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class SchemaParseTest {
    private final Logger log = LoggerFactory.getLogger(this.getClass());

    String address = "192.168.8.6:3820";
    String user = "root";
    String passwd = "nebula";


    //@Before
    public void setup() {
        MockGraph.mockGraphData();
    }


    @Test()
    public void testGetGraphType() {
        try {
            NebulaClient client = NebulaClient.builder(address, user, passwd).build();
            Method getGraphType = client.getClass().getDeclaredMethod("getGraphType",
                    String.class);
            getGraphType.setAccessible(true);
            String result = (String) getGraphType.invoke(client, "nba");
            assert result.equals("graph_type_nba");
        } catch (IOErrorException | AuthFailedException | InvocationTargetException
                 | IllegalAccessException | NoSuchMethodException e) {
            log.error("test getNodeProperties error", e);
            assert false;
        }
    }

    @Test()
    public void testGetGraphTypeWithWrongGraphName() {
        try {
            NebulaClient client = NebulaClient.builder(address, user, passwd).build();
            Method getGraphType = client.getClass().getDeclaredMethod("getGraphType",
                    String.class);
            getGraphType.setAccessible(true);
            getGraphType.invoke(client, "nba_not_exist");
        } catch (InvocationTargetException e) {
            if (e.getTargetException().getMessage()
                    .contains("graphName nba_not_exist does not exist")) {
                log.info("expected to get InvocationTargetException");
                assert true;
            } else {
                log.error("test getGraphType error.");
                assert false;
            }
        } catch (IOErrorException | NoSuchMethodException
                 | AuthFailedException | IllegalAccessException e) {
            log.error("test getGraphType error", e);
            assert false;
        }
    }

    @Test()
    public void testGetNodeProperties() {
        try {
            NebulaClient client = NebulaClient.builder(address, user, passwd).build();
            Method getNodeProperties = client.getClass().getDeclaredMethod("getNodeProperties",
                    String.class, String.class);
            getNodeProperties.setAccessible(true);
            List<String> result = (List<String>) getNodeProperties.invoke(client,
                    "nba",
                    "node_type_player");
            assert result.size() == 5;
            assert result.stream().distinct().count() == 5;
            assert result.contains("id");
        } catch (IOErrorException | AuthFailedException | NoSuchMethodException
                 | InvocationTargetException | IllegalAccessException e) {
            log.error("test getNodeProperties error", e);
            assert false;
        }
    }

    @Test()
    public void testGetEdgeProperties() {
        try {
            NebulaClient client = NebulaClient.builder(address, user, passwd).build();
            Method getNodeProperties = client.getClass().getDeclaredMethod("getEdgeProperties",
                    String.class, String.class);
            getNodeProperties.setAccessible(true);
            List<String> result = (List<String>) getNodeProperties.invoke(client,
                    "nba",
                    "edge_type_follow");
            assert result.size() == 2;
            assert result.stream().distinct().count() == 2;
            assert result.contains("followness");
            assert result.contains("likeness");
        } catch (IOErrorException | AuthFailedException | NoSuchMethodException
                 | InvocationTargetException | IllegalAccessException e) {
            log.error("test getEdgeProperties error", e);
            assert false;
        }
    }

    @Test
    public void testGetAllParts() {
        try {
            NebulaClient client = NebulaClient.builder(address, user, passwd).build();
            Method getNodeProperties = client.getClass().getDeclaredMethod("getAllParts");
            getNodeProperties.setAccessible(true);
            List<Integer> result = (List<Integer>) getNodeProperties.invoke(client);
            assert result.size() == 10;
            Collections.sort(result);
            List<Integer> expectList = Arrays.asList(1, 2, 3, 4, 5, 6, 7, 8, 9, 10);
            Collections.sort(expectList);
            assert result.equals(expectList);
        } catch (IOErrorException | AuthFailedException | NoSuchMethodException
                 | InvocationTargetException | IllegalAccessException e) {
            log.error("test getAllParts error", e);
            assert false;
        }
    }


}
