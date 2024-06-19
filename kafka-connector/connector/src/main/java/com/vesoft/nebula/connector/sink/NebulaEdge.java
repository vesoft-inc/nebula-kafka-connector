/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.EDGE_VALUES_TEMPLATE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.PROPERTY_TEMPLATE;

import com.vesoft.nebula.connector.exceptions.DataFormatException;
import com.vesoft.nebula.connector.util.NebulaUtils;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class NebulaEdge {
    private String srcPkName;
    private String srcPk;
    private String dstPkName;
    private String dstPk;

    private Map<String, Object> properties;

    public NebulaEdge(String srcPkName, String srcPk,
                      String dstPkName, String dstPk,
                      Map<String, Object> properties) {
        this.srcPkName = srcPkName;
        this.srcPk = srcPk;
        this.dstPkName = dstPkName;
        this.dstPk = dstPk;
        this.properties = properties;
    }

    // schema: 属性名-> 属性类型
    public String getEdgeStatement(NebulaEdgeSchema edgeSchema) throws DataFormatException {
        List<String> propEntryStringList = new ArrayList<>();
        for (Map.Entry<String, Object> entry : properties.entrySet()) {
            String propString = String.format(PROPERTY_TEMPLATE, entry.getKey(),
                    NebulaUtils.extractPropertyValue(edgeSchema.getProperties(), entry));
            propEntryStringList.add(propString);
        }
        String props = String.join(",", propEntryStringList);

        String srcPkValue = NebulaUtils.extractIdValue(edgeSchema.getSourceNodePkType(), srcPk);
        String dstPkValue = NebulaUtils.extractIdValue(edgeSchema.getTargetNodePkType(), dstPk);

        return String.format(EDGE_VALUES_TEMPLATE,
                srcPkName,
                srcPkValue,
                props,
                dstPkName,
                dstPkValue);
    }

    public String getSrcPkName() {
        return srcPkName;
    }

    public void setSrcPkName(String srcPkName) {
        this.srcPkName = srcPkName;
    }

    public String getSrcPk() {
        return srcPk;
    }

    public void setSrcPk(String srcPk) {
        this.srcPk = srcPk;
    }

    public String getDstPkName() {
        return dstPkName;
    }

    public void setDstPkName(String dstPkName) {
        this.dstPkName = dstPkName;
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
        sb.append("{");
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
        return "NebulaEdge{"
                + "srcPk='" + srcPk + '\''
                + ", dstPk='" + dstPk + '\''
                + ", properties=" + properties
                + '}';
    }
}
