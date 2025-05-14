/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.source;

import com.vesoft.nebula.connector.config.NebulaSourceConnectConfig;
import com.vesoft.nebula.connector.connection.NebulaGraphProvider;
import com.vesoft.nebula.connector.sink.NebulaEdgeSchema;
import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import com.vesoft.nebula.driver.graph.exception.NoValidSessionException;
import com.vesoft.nebula.driver.graph.scan.ScanEdgeResult;
import com.vesoft.nebula.driver.graph.scan.ScanEdgeResultIterator;
import com.vesoft.nebula.driver.graph.scan.TableRow;
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
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        Map<String, String> props = schema.getProperties();

        for (String pk : schema.getSourcePkNameAndType().keySet()) {
            props.put("src_" + pk, schema.getSourcePkNameAndType().get(pk));
        }
        for (String pk : schema.getTargetPkNameAndType().keySet()) {
            props.put("dst_" + pk, schema.getTargetPkNameAndType().get(pk));
        }
        edgeSchema = props;
        return edgeSchema;
    }
}
