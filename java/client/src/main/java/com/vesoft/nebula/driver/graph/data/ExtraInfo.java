package com.vesoft.nebula.driver.graph.data;

/**
 * the class maintains some additional information for execution result.
 */
public class ExtraInfo {

    // cursor for scan procedure
    private String cursor;
    // the number of affected nodes
    private long   affectedNodes;
    // the number of affected forward edges
    private long   affectedEdges;
    // the number of affected reverse edges

    private long totalServerTimeUs;

    private long buildTimeUs;

    private long optimizeTimeUs;

    private long serializeTimeUs;


    public ExtraInfo() {
        this.cursor = null;
        affectedNodes = 0;
        affectedEdges = 0;
    }

    public void setCursor(String cursor) {
        this.cursor = cursor;
    }

    public void setAffectedNodes(long affectedNodes) {
        this.affectedNodes = affectedNodes;
    }

    public void setAffectedEdges(long affectedEdges) {
        this.affectedEdges = affectedEdges;
    }

    public void setTotalServerTimeUs(long totalServerTimeUs) {
        this.totalServerTimeUs = totalServerTimeUs;
    }

    public void setBuildTimeUs(long buildTimeUs) {
        this.buildTimeUs = buildTimeUs;
    }

    public void setOptimizeTimeUs(long optimizeTimeUs) {
        this.optimizeTimeUs = optimizeTimeUs;
    }

    public void setSerializeTimeUs(long serializeTimeUs) {
        this.serializeTimeUs = serializeTimeUs;
    }

    public String getCursor() {
        return this.cursor;
    }

    public long getAffectedNodes() {
        return this.affectedNodes;
    }

    public long getAffectedEdges() {
        return this.affectedEdges;
    }

    public long getTotalServerTimeUs() {
        return totalServerTimeUs;
    }

    public long getBuildTimeUs() {
        return buildTimeUs;
    }

    public long getOptimizeTimeUs() {
        return optimizeTimeUs;
    }

    public long getSerializeTimeUs() {
        return serializeTimeUs;
    }

    @Override
    public String toString() {
        return "ExtraInfo{"
                + "cursor='" + cursor + '\''
                + ", affectedNodes=" + affectedNodes
                + ", affectedEdges=" + affectedEdges
                + ", totalServerTimeUs=" + totalServerTimeUs
                + ", buildTimeUs=" + buildTimeUs
                + ", optimizeTimeUs=" + optimizeTimeUs
                + ", serializeTimeUs=" + serializeTimeUs
                + '}';
    }
}
