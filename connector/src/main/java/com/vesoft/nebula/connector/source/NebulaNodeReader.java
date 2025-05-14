/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.source;

import com.vesoft.nebula.connector.config.NebulaSourceConnectConfig;
import com.vesoft.nebula.connector.connection.NebulaGraphProvider;
import com.vesoft.nebula.connector.sink.NebulaNodeSchema;
import com.vesoft.nebula.driver.graph.data.ResultSet;
import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import com.vesoft.nebula.driver.graph.exception.NoValidSessionException;
import com.vesoft.nebula.driver.graph.scan.ScanNodeResult;
import com.vesoft.nebula.driver.graph.scan.ScanNodeResultIterator;
import com.vesoft.nebula.driver.graph.scan.TableRow;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class NebulaNodeReader extends NebulaReader {

    private ScanNodeResultIterator scanNodeResultIterator;

    private Map<String, String> nodeSchema = null;


    public NebulaNodeReader(NebulaSourceConnectConfig config,
                            NebulaGraphProvider graphProvider,
                            List<Integer> partsId,
                            String nodeType) {
        super(config, graphProvider, nodeType);
        scanNodeResultIterator = graphProvider.scanNode(
                config.graphName,
                nodeType,
                config.nebulaNodePropertyNames,
                partsId,
                config.batchSize);
        getSchema();
    }

    @Override
    public List<TableRow> next() {
        if (hasNext) {
            ScanNodeResult result = scanNodeResultIterator.next();
            if (rowNames == null) {
                rowNames = result.getPropNames();
            }
            List<TableRow> rows = result.getTableRows();
            return rows;
        } else {
            return new ArrayList<>();
        }
    }

    @Override
    public Map<String, String> getSchema() {
        if (nodeSchema != null) {
            return nodeSchema;
        }
        NebulaNodeSchema schema = null;
        try {
            schema = graphProvider.getNodeSchema(graphName, typeName);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        nodeSchema = schema.getNodeProperties();
        return nodeSchema;
    }
}
