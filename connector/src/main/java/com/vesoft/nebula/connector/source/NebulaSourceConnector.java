/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.source;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.NEBULA_PARTS_FOR_EACH_TASK;

import com.vesoft.nebula.connector.config.NebulaSourceConnectConfig;
import com.vesoft.nebula.connector.util.Version;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.connect.connector.Task;
import org.apache.kafka.connect.source.SourceConnector;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class NebulaSourceConnector extends SourceConnector {
    private static final Logger LOG = LoggerFactory.getLogger(NebulaSourceConnector.class);


    private Map<String, String> configProps;

    @Override
    public void start(Map<String, String> props) {
        this.configProps = props;
    }

    @Override
    public Class<? extends Task> taskClass() {
        return NebulaSourceTask.class;
    }

    // distribute different part for each task
    @Override
    public List<Map<String, String>> taskConfigs(int maxTasks) {
        LOG.info("Setting task configurations for {} workers", maxTasks);
        final List<Map<String, String>> configs = new ArrayList<>(maxTasks);
        // compute the nebula parts each kafka task should read
        int nebulaParts = 10;
        List<List<String>> taskParts = new ArrayList<>();
        for (int i = 0; i < maxTasks; i++) {
            taskParts.add(new ArrayList<>());
        }
        for (int i = 1; i <= nebulaParts; i++) {
            int index = i % maxTasks == 0 ? (maxTasks - 1) : (i % maxTasks - 1);
            taskParts.get(index).add(String.valueOf(i));
        }

        // distribute the configs for tasks
        for (int i = 0; i < maxTasks; i++) {
            Map<String, String> taskConfigProps = new HashMap<>(configProps);
            taskConfigProps.put(NEBULA_PARTS_FOR_EACH_TASK,
                    String.join(",", taskParts.get(i)));
            configs.add(taskConfigProps);
        }
        return configs;
    }

    @Override
    public void stop() {
    }

    @Override
    public ConfigDef config() {
        ConfigDef configDef = NebulaSourceConnectConfig.configDef();
        return configDef;
    }

    @Override
    public String version() {
        return Version.getVersion();
    }
}
