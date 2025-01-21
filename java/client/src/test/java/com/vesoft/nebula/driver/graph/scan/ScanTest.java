package com.vesoft.nebula.driver.graph.scan;

import com.vesoft.nebula.driver.graph.ErrorCode;
import com.vesoft.nebula.driver.graph.ServerConstant;
import com.vesoft.nebula.driver.graph.data.ResultSet;
import com.vesoft.nebula.driver.graph.exception.AuthFailedException;
import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import com.vesoft.nebula.driver.graph.net.NebulaClient;
import com.vesoft.nebula.driver.graph.util.MockGraph;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class ScanTest {
    String addresses = ServerConstant.address;
    String user      = ServerConstant.user;
    String passwd    = ServerConstant.passwd;
    String graphName = "nba";
    String nodeType  = "node_type_player";
    String edgeType  = "edge_type_follow";

    NebulaClient client;

    @Before
    public void setup() {
        try {
            client = NebulaClient.builder(addresses, user, passwd)
                    .withRequestTimeoutMills(3000)
                    .build();
            ResultSet resultSet = client.execute("return 1");
            Assert.assertEquals(resultSet.getErrorCode(), ErrorCode.SUCCESSFUL_COMPLETION);
        } catch (IOErrorException | AuthFailedException e) {
            Assert.fail(e.getMessage());
        }
        MockGraph.mockGraphData();
        MockGraph.mockSpecialGraphType();
    }

    @After
    public void tearDown() {
        if (client != null) {
            client.close();
        }
    }

    @Test
    public void scanNodeForReturnCols() {
        // test specific return columns
        List<String>           returnCols = Arrays.asList("id", "name");
        ScanNodeResultIterator iterator   = client.scanNode(graphName, nodeType, returnCols, 1);
        List<TableRow>         rows       = new ArrayList<>();
        List<String>           columns    = new ArrayList<>();
        while (iterator.hasNext()) {
            ScanNodeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }

        }
        assert (rows.size() == 3);
        assert (columns.size() == 2);

        // test null return columns
        rows.clear();
        columns.clear();
        returnCols = null;
        iterator = client.scanNode(graphName, nodeType, returnCols);
        while (iterator.hasNext()) {
            ScanNodeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 3);
        assert (columns.size() == 5);

        // test empty return columns
        rows.clear();
        columns.clear();
        returnCols = new ArrayList<>();
        iterator = client.scanNode(graphName, nodeType, returnCols);
        while (iterator.hasNext()) {
            ScanNodeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 3);
        assert (columns.size() == 1);

        // test batch size
        rows.clear();
        columns.clear();
        returnCols = new ArrayList<>();
        iterator = client.scanNode(graphName, nodeType, returnCols, 2, 1);
        while (iterator.hasNext()) {
            ScanNodeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 1);
        assert (columns.size() == 1);
    }

    @Test
    public void scanNodeForEscapeProp() {
        String graphName = "图";
        String nodeType  = "user";
        // test specific return columns
        List<String>           returnCols = Arrays.asList("id", "type", "姓名", "age\t");
        ScanNodeResultIterator iterator   = client.scanNode(graphName, nodeType, returnCols, 1);
        List<TableRow>         rows       = new ArrayList<>();
        List<String>           columns    = new ArrayList<>();
        while (iterator.hasNext()) {
            ScanNodeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }

        }
        assert (rows.size() == 2);
        assert (columns.size() == 4);

        // test null return columns
        rows.clear();
        columns.clear();
        returnCols = null;
        iterator = client.scanNode(graphName, nodeType, returnCols);
        while (iterator.hasNext()) {
            ScanNodeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 2);
        assert (columns.size() == 4);

        // test empty return columns
        rows.clear();
        columns.clear();
        returnCols = new ArrayList<>();
        iterator = client.scanNode(graphName, nodeType, returnCols);
        while (iterator.hasNext()) {
            ScanNodeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 2);
        assert (columns.size() == 1);
    }

    @Test
    public void scanEdgeForReturnCols() {
        // test specific return columns
        List<String>           returnCols = Arrays.asList("likeness");
        ScanEdgeResultIterator iterator   = client.scanEdge(graphName, edgeType, returnCols, 1);
        List<TableRow>         rows       = new ArrayList<>();
        List<String>           columns    = new ArrayList<>();
        while (iterator.hasNext()) {
            ScanEdgeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 2);
        assert (columns.size() == 3);

        // test null return columns
        rows.clear();
        columns.clear();
        returnCols = null;
        iterator = client.scanEdge(graphName, edgeType, returnCols);
        while (iterator.hasNext()) {
            ScanEdgeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 2);
        assert (columns.size() == 4);

        // test empty return columns
        rows.clear();
        columns.clear();
        returnCols = new ArrayList<>();
        iterator = client.scanEdge(graphName, edgeType, returnCols);
        while (iterator.hasNext()) {
            ScanEdgeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 2);
        assert (columns.size() == 2);

        // test batch size
        rows.clear();
        columns.clear();
        returnCols = new ArrayList<>();
        iterator = client.scanEdge(graphName, edgeType, returnCols, 2, 1);
        while (iterator.hasNext()) {
            ScanEdgeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 1);
        assert (columns.size() == 2);
    }


    @Test
    public void scanEdgeForEscapeProp() {
        String graphName = "图";
        String edgeType  = "edge";
        // test specific return columns
        List<String>           returnCols = Arrays.asList("likeness", "type\r");
        ScanEdgeResultIterator iterator   = client.scanEdge(graphName, edgeType, returnCols, 10);
        List<TableRow>         rows       = new ArrayList<>();
        List<String>           columns    = new ArrayList<>();
        while (iterator.hasNext()) {
            ScanEdgeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 2);
        assert (columns.size() == 4);

        // test null return columns
        rows.clear();
        columns.clear();
        returnCols = null;
        iterator = client.scanEdge(graphName, edgeType, returnCols);
        while (iterator.hasNext()) {
            ScanEdgeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 2);
        assert (columns.size() == 5);

        // test empty return columns
        rows.clear();
        columns.clear();
        returnCols = new ArrayList<>();
        iterator = client.scanEdge(graphName, edgeType, returnCols);
        while (iterator.hasNext()) {
            ScanEdgeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 2);
        assert (columns.size() == 2);
    }
}
