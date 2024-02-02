/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.NODE_VALUES_TEMPLATE;
import com.vesoft.nebula.connector.exceptions.DataFormatException;
import com.vesoft.nebula.connector.util.NebulaUtils;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class NebulaNode {
    private String vid;
    private Map<String, Object> properties;

    public NebulaNode(String vid, Map<String, Object> properties) {
        this.vid = vid;
        this.properties = properties;
    }

    // TODO update schema data type class
    public String getNodeStatement(NebulaNodeSchema nodeSchema) throws DataFormatException {
        List<String> propEntryStringList = new ArrayList<>();
        properties.put("id", vid);
        for (Map.Entry<String, Object> entry : properties.entrySet()) {
            String propString = entry.getKey() + ":" + NebulaUtils.extractPropertyValue(nodeSchema.getNodeProperties(),
                    entry);
            propEntryStringList.add(propString);
        }
        String props = String.join(",", propEntryStringList);
        return String.format(NODE_VALUES_TEMPLATE, props);
    }

    public String getVid() {
        return vid;
    }

    public void setVid(String vid) {
        this.vid = vid;
    }

    public Map<String, Object> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, Object> properties) {
        this.properties = properties;
    }

    public String getNodeString() {
        StringBuilder sb = new StringBuilder();
        sb.append("{vid:");
        sb.append(vid);
        sb.append(",");
        for (Map.Entry<String, Object> kv : properties.entrySet()) {
            sb.append(kv.getKey());
            sb.append(":");
            sb.append(kv.getValue());
            sb.append(",");
        }
        sb.deleteCharAt(sb.length() - 1);
        sb.append("}");
        return sb.toString();
    }

    @Override
    public String toString() {
        return "NebulaNode{" +
                "vid='" + vid + '\'' +
                ", properties=" + properties +
                '}';
    }
}
