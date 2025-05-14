
package com.vesoft.nebula.connector.sink;

import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import com.vesoft.nebula.connector.util.Version;
import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import com.vesoft.nebula.driver.graph.exception.NoValidSessionException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.apache.kafka.clients.consumer.OffsetAndMetadata;
import org.apache.kafka.common.TopicPartition;
import org.apache.kafka.connect.data.Field;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.sink.SinkRecord;
import org.apache.kafka.connect.sink.SinkTask;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class NebulaSinkTask extends SinkTask {
    private static final Logger log = LoggerFactory.getLogger(NebulaSinkTask.class);

    NebulaSinkConnectConfig config;
    NebulaWriter            writer;

    private String connectorName = null;
    private String taskId        = null;


    @Override
    public String version() {
        return Version.getVersion();
    }

    @Override
    public void start(Map<String, String> props) {
        config = new NebulaSinkConnectConfig(props);
        connectorName = config.connectorName;
        taskId = props.get("task.id");
        log.info("Starting Nebula Sink task {}-{}", connectorName, taskId);

        initWriter();
    }

    private void initWriter() {
        log.info("Initializing Nebula Writer");
        writer = new NebulaWriter(config);
        log.info("Nebula Writer initialized");
    }


    @Override
    public void put(Collection<SinkRecord> records) {
        if (records.isEmpty()) {
            return;
        }
        final SinkRecord first        = records.iterator().next();
        final int        recordsCount = records.size();
        List<String>     schemaNames  = new ArrayList<>();
        if (first.valueSchema() != null) {
            if (first.valueSchema().type() == Schema.Type.STRUCT) {
                final List<Field> valueSchemaFields = first.valueSchema().fields();
                for (Field field : valueSchemaFields) {
                    schemaNames.add(field.name());
                }
            }

        } else if (first.value() instanceof Map) {
            Map<String, Object> properties = (HashMap<String, Object>) first.value();
            schemaNames.addAll(properties.keySet());
        } else {
            log.warn("cannot get schema names of records, simple record:{}", first.value());
        }

        log.info("Received {} records. First record kafka coordinates:({}-{}-{}), record "
                         + "schema:{}. Writing them to nebula...",
                 recordsCount, first.topic(), first.kafkaPartition(), first.kafkaOffset(),
                 schemaNames);

        try {
            writer.write(records);
        } catch (Exception e) {
            log.error("failed to write {} records, reason: {}", records.size(), e.getMessage());
            throw new RuntimeException(e.getMessage());
        }
    }

    @Override
    public void flush(Map<TopicPartition, OffsetAndMetadata> currentOffsets) {
    }

    @Override
    public void stop() {
        log.info("Stopping Nebula Sink Task {}-{}", connectorName, taskId);
        try {
            if (writer != null) {
                writer.close();
            }
        } finally {
            //
        }
    }
}
