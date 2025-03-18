
package com.vesoft.nebula.connector.sink;

import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.EDGE_VALUES_TEMPLATE;
import static com.vesoft.nebula.connector.config.NebulaConnectConfigName.PROPERTY_TEMPLATE;

import com.vesoft.nebula.connector.exceptions.DataFormatException;
import com.vesoft.nebula.connector.util.NebulaUtils;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class NebulaEdge {
    private String srcPk;
    private String dstPk;

    private Map<String, String> properties;

    public NebulaEdge(String srcPk, String dstPk,
                      Map<String, String> properties) {
        this.srcPk = srcPk;
        this.dstPk = dstPk;
        this.properties = properties;
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

    public Map<String, String> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, String> properties) {
        this.properties = properties;
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
