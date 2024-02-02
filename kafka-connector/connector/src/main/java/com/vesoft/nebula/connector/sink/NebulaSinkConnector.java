/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import com.vesoft.nebula.connector.util.Version;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.connect.connector.Task;
import org.apache.kafka.connect.sink.SinkConnector;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class NebulaSinkConnector extends SinkConnector {
    private static final Logger log = LoggerFactory.getLogger(NebulaSinkConnector.class);

    private Map<String, String> configProps;

    @Override
    public void start(Map<String, String> props) {
        configProps = props;
    }

    @Override
    public Class<? extends Task> taskClass() {
        return NebulaSinkTask.class;
    }

    @Override
    public List<Map<String, String>> taskConfigs(int maxTasks) {
        log.info("Setting task configurations for {} workers", maxTasks);
        final List<Map<String, String>> configs = new ArrayList<>(maxTasks);
        for (int i = 0; i < maxTasks; i++) {
            configs.add(configProps);
        }
        return configs;
    }

    @Override
    public void stop() {
    }

    @Override
    public ConfigDef config() {
        ConfigDef configDef = NebulaSinkConnectConfig.CONFIG_DEF;
        return configDef;
    }

    @Override
    public String version() {
        return Version.getVersion();
    }


}
