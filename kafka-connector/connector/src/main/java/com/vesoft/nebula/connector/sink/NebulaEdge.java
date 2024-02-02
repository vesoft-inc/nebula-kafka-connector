/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.EDGE_VALUES_TEMPLATE;
import com.vesoft.nebula.connector.exceptions.DataFormatException;
import com.vesoft.nebula.connector.util.NebulaUtils;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class NebulaEdge {
    private String srcPk;
    private String dstPk;

    private Map<String, Object> properties;

    public NebulaEdge(String srcPk, String dstPk, Map<String, Object> properties) {
        this.srcPk = srcPk;
        this.dstPk = dstPk;
        this.properties = properties;
    }

    // schema: 属性名-> 属性类型
    public String getEdgeStatement(NebulaEdgeSchema edgeSchema) throws DataFormatException {
        List<String> propEntryStringList = new ArrayList<>();
        for (Map.Entry<String, Object> entry : properties.entrySet()) {
            String propString = entry.getKey() + ":" + NebulaUtils.extractPropertyValue(edgeSchema.getProperties(),
                    entry);
            propEntryStringList.add(propString);
        }
        String props = String.join(",", propEntryStringList);

        String srcId = NebulaUtils.extractIdValue(edgeSchema.getSourceNodeIdType(), srcPk);
        String dstId = NebulaUtils.extractIdValue(edgeSchema.getTargetNodeIdType(), dstPk);

        return String.format(EDGE_VALUES_TEMPLATE, srcId, props, dstId);
    }

    public String getSrcPk() {
        return srcPk;
    }

    public void setSrcPk(String srcPk) {
        this.srcPk = srcPk;
    }

    public String getDstPk() {
        return dstPk;
    }

    public void setDstPk(String dstPk) {
        this.dstPk = dstPk;
    }

    public Map<String, Object> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, Object> properties) {
        this.properties = properties;
    }

    public String getEdgeString() {
        StringBuilder sb = new StringBuilder();
        sb.append("{srcPk:");
        sb.append(srcPk);
        sb.append(",dstPk:");
        sb.append(dstPk);
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
        return "NebulaEdge{" +
                "srcPk='" + srcPk + '\'' +
                ", dstPk='" + dstPk + '\'' +
                ", properties=" + properties +
                '}';
    }
}
