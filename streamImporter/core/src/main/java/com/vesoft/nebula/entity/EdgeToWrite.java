
package com.vesoft.nebula.entity;

import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;

public class EdgeToWrite implements Serializable {
    private long srcId;
    private long dstId;
    private long rank = 0L;
    private Map<String, String> properties = new HashMap<>();

    public EdgeToWrite(long srcId, long dstId, Map<String, String> properties) {
        this.srcId = srcId;
        this.dstId = dstId;
        this.properties = properties;
    }

    public long getSrcId() {
        return srcId;
    }

    public void setSrcId(long srcId) {
        this.srcId = srcId;
    }

    public long getDstId() {
        return dstId;
    }

    public void setDstId(long dstId) {
        this.dstId = dstId;
    }

    public long getRank() {
        return rank;
    }

    public void setRank(long rank) {
        this.rank = rank;
    }

    public Map<String, String> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, String> properties) {
        this.properties = properties;
    }

    @Override
    public String toString() {
        return "EdgeToWrite{" +
                "srcId=" + srcId +
                ", dstId=" + dstId +
                ", rank=" + rank +
                ", properties=" + properties +
                '}';
    }
}
