package com.vesoft.nebula.client.scan;

import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.net.NebulaClient;
import com.vesoft.nebula.client.graph.scan.ScanNodeResult;
import com.vesoft.nebula.client.graph.scan.ScanNodeResultIterator;
import com.vesoft.nebula.client.graph.scan.TableRow;
import com.vesoft.nebula.client.util.MockGraph;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class ScanTest {
    String addresses = "127.0.0.1:9669";
    String user      = "root";
    String passwd    = "NebulaGraph01";
    String graphName = "nba";
    String nodeType  = "node_type_player";

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
    }

    @Test
    public void scanNodeForReturnCols() {
        // test specific return columns
        List<String>           returnCols = Arrays.asList("id", "name");
        ScanNodeResultIterator iterator   = client.scanNode(graphName, nodeType, returnCols);
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
        iterator = client.scanNode(graphName, nodeType, returnCols,2, 1);
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

}
