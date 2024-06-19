/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.source;

import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.exception.NoValidSessionException;
import com.vesoft.nebula.client.graph.scan.ScanEdgeResult;
import com.vesoft.nebula.client.graph.scan.ScanEdgeResultIterator;
import com.vesoft.nebula.client.graph.scan.TableRow;
import com.vesoft.nebula.connector.config.NebulaSourceConnectConfig;
import com.vesoft.nebula.connector.connection.NebulaGraphProvider;
import com.vesoft.nebula.connector.sink.NebulaEdgeSchema;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class NebulaEdgeReader extends NebulaReader {

    private ScanEdgeResultIterator scanEdgeResultIterator;

    private Map<String, String> edgeSchema = null;


    public NebulaEdgeReader(NebulaSourceConnectConfig config,
                            NebulaGraphProvider graphProvider,
                            List<Integer> partsId,
                            String edgeType) {
        super(config, graphProvider, edgeType);
        scanEdgeResultIterator = graphProvider.scanEdge(
                config.graphName,
                edgeType,
                config.nebulaEdgePropertyNames,
                partsId,
                config.batchSize);
        getSchema();
    }


    @Override
    public List<TableRow> next() {
        if (hasNext) {
            ScanEdgeResult result = scanEdgeResultIterator.next();
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
        if (edgeSchema != null) {
            return edgeSchema;
        }
        NebulaEdgeSchema schema = null;
        try {
            schema = graphProvider.getEdgeSchema(graphName, typeName);
        } catch (IOErrorException | NoValidSessionException e) {
            throw new RuntimeException(e);
        }
        Map<String, String> props = schema.getProperties();
        props.put(schema.getSourceNodePkName(), schema.getSourceNodePkType());
        props.put(schema.getTargetNodePkName(), schema.getTargetNodePkType());
        edgeSchema = props;
        return edgeSchema;
    }
}
