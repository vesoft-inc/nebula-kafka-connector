/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.source;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.CONNECT_POLLING_SLEEP_MS;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.NEBULA_PARTS_FOR_EACH_TASK;

import com.vesoft.nebula.client.graph.scan.TableRow;
import com.vesoft.nebula.connector.config.NebulaSourceConnectConfig;
import com.vesoft.nebula.connector.connection.NebulaGraphProvider;
import com.vesoft.nebula.connector.converter.NebulaRecordConverter;
import com.vesoft.nebula.connector.util.Version;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicBoolean;
import org.apache.kafka.common.config.ConfigException;
import org.apache.kafka.connect.source.SourceRecord;
import org.apache.kafka.connect.source.SourceTask;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class NebulaSourceTask extends SourceTask {
    private static final Logger LOG = LoggerFactory.getLogger(NebulaSourceTask.class);

    private final AtomicBoolean running = new AtomicBoolean(false);

    NebulaSourceConnectConfig config;
    List<NebulaReader> readers;
    private String connectorName = null;
    private String taskId = null;

    private List<Integer> parts = new ArrayList<>();

    private long pollingMs = 0;

    @Override
    public String version() {
        return Version.getVersion();
    }

    // TODO restore the latest offset
    @Override
    public void start(Map<String, String> props) {
        running.set(true);
        connectorName = props.get("name");
        taskId = props.get("task.id");
        pollingMs = Long.valueOf(props.get(CONNECT_POLLING_SLEEP_MS));
        LOG.info("Starting Nebula Source task {}-{}", connectorName, taskId);

        try {
            String partsStr = props.get(NEBULA_PARTS_FOR_EACH_TASK);
            for (String part : partsStr.split(",")) {
                parts.add(Integer.valueOf(part));
            }
            config = new NebulaSourceConnectConfig(props);
        } catch (ConfigException e) {
            throw new ConfigException("Couldn't start NebulaSourceTask due to configuration error",
                e);
        }
        // init the reader for each node type and edge type
        readers = new ArrayList<>();
        for (String nodeType : config.graphNodeTypes) {
            NebulaNodeReader nodeReader = new NebulaNodeReader(
                config,
                new NebulaGraphProvider(config),
                parts,
                nodeType);
            readers.add(nodeReader);
        }
        for (String edgeType : config.graphEdgeTypes) {
            NebulaEdgeReader edgeReader = new NebulaEdgeReader(
                config,
                new NebulaGraphProvider(config),
                parts,
                edgeType);
            readers.add(edgeReader);
        }
    }

    @Override
    public List<SourceRecord> poll() throws InterruptedException {
        List<SourceRecord> records = new ArrayList<>();
        String topicPrefix = config.topicPrefix + "." + config.graphName;
        if (!running.get()) {
            shutdown();
            return null;
        }

        for (NebulaReader reader : readers) {
            if (!reader.hasNext) {
                continue;
            }
            Map<String, String> nebulaSchema = reader.getSchema();
            List<TableRow> result = reader.next();
            String topicName = topicPrefix + "." + reader.typeName;
            List<SourceRecord> sourceRecords = NebulaRecordConverter
                .convertTableRows(reader.rowNames, result, topicName, nebulaSchema);
            records.addAll(sourceRecords);
        }
        if (records.isEmpty()) {
            LOG.info("no data found, sleep for {} ms", pollingMs);
            Thread.sleep(pollingMs);
        }
        return records;
    }

    private void shutdown() {
        for (NebulaReader reader : readers) {
            reader.close();
        }
        readers.clear();
    }

    @Override
    public void stop() {
        running.set(false);
    }
}
