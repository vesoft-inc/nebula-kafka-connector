/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.source;

import com.vesoft.nebula.client.graph.scan.TableRow;
import com.vesoft.nebula.connector.config.NebulaSourceConnectConfig;
import com.vesoft.nebula.connector.connection.NebulaGraphProvider;
import java.util.List;
import java.util.Map;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * NebulaReader is a reader to read node or edge data from NebulaGraph.
 * TODO to support checkpoint while reading, the reader should maintain the cursor.
 * TODO maybe maintain the element_id of each data after scan function is updated.
 */
public abstract class NebulaReader {
    private static final Logger LOG = LoggerFactory.getLogger(NebulaReader.class);

    protected NebulaSourceConnectConfig config;
    protected NebulaGraphProvider graphProvider;

    protected String graphName;
    protected String typeName;
    protected int batchSize;

    protected String cursor = "";
    protected boolean hasNext;

    protected List<String> rowNames = null;


    public NebulaReader(NebulaSourceConnectConfig config,
                        NebulaGraphProvider graphProvider,
                        String typeName) {
        this.config = config;
        this.graphProvider = graphProvider;
        this.graphName = config.graphName;
        this.typeName = typeName;
        this.batchSize = config.batchSize;
    }

    /**
     * get the next batch data of node type or edge type.
     */
    public abstract List<TableRow> next();

    /**
     * get the schema of NodeType or EdgeType.
     * The schema includes not only property information (property name and property data type),
     * but also includes source node primary key information(pk name and pk data type)
     * and target node primary key information(pk name and pk data type) of the edge.
     *
     * @return map for column name and column data type
     */
    public abstract Map<String, String> getSchema();

    /**
     * close the connections.
     */
    public void close() {
        graphProvider.close();
    }

}
