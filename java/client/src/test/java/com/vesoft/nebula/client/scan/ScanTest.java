/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.scan;

import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import com.vesoft.nebula.client.graph.net.NebulaClient;
import com.vesoft.nebula.client.graph.scan.ScanNodeResult;
import com.vesoft.nebula.client.graph.scan.ScanNodeResultIterator;
import com.vesoft.nebula.client.graph.scan.TableRow;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class ScanTest {
    String addresses = "192.168.8.6:3820";
    String user = "root";
    String passwd = "nebula";
    String graphName = "nba";
    String nodeType = "node_type_player";

    NebulaClient client;

    @Before
    public void setup() {
        try {
            client = NebulaClient.builder(addresses, user, passwd)
                    .setConnectTimeoutMills(1000)
                    .setRequestTimeoutMills(3000)
                    .setRetryTimes(3)
                    .setMaxSessionSize(1)
                    .setMinSessionSize(0)
                    .setHealthCheckTimeMills(5000)
                    .setReconnect(true)
                    .build();
            ResultSet resultSet = client.execute("return 1");
            Assert.assertEquals(resultSet.getErrorCode(), ErrorCode.SUCCESSFUL_COMPLETION.code);
        } catch (UnknownHostException | IOErrorException | NoValidSessionException e) {
            Assert.fail(e.getMessage());
        }
        // MockGraph.mockGraphData();
    }

    @Test
    public void scanNodeForReturnCols() {
        // test specific return columns
        List<String> returnCols = Arrays.asList("id", "name");
        ScanNodeResultIterator iterator = client.scanNode(graphName, nodeType, returnCols, 3, 1);
        List<TableRow> rows = new ArrayList<>();
        List<String> columns = new ArrayList<>();
        while (iterator.hasNext()) {
            ScanNodeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }

        }
        assert (rows.size() == 1);
        assert (columns.size() == 2);

        // test null return columns
        rows.clear();
        columns.clear();
        returnCols = null;
        iterator = client.scanNode(graphName, nodeType, returnCols, 3, 1);
        while (iterator.hasNext()) {
            ScanNodeResult result = iterator.next();
            if (!result.isEmpty()) {
                rows.addAll(result.getTableRows());
                columns.addAll(result.getPropNames());
            }
        }
        assert (rows.size() == 1);
        assert (columns.size() == 5);

        // test empty return columns
        rows.clear();
        columns.clear();
        returnCols = new ArrayList<>();
        iterator = client.scanNode(graphName, nodeType, returnCols, 3, 1);
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
