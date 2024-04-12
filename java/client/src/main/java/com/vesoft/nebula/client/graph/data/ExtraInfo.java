/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

/**
 * the class maintains some additional information for execution result.
 */
public class ExtraInfo {

    // cursor for scan procedure
    private String cursor;
    // the number of affected nodes
    private long affectedNodes;
    // the number of affected forward edges
    private long affectedEdges;
    // the number of affected reverse edges

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

    public String getCursor() {
        return this.cursor;
    }

    public long getAffectedNodes() {
        return this.affectedNodes;
    }

    public long getAffectedEdges() {
        return this.affectedEdges;
    }

    @Override
    public String toString() {
        if (cursor != null) {
            return "ExtraInfo{"
                    + "cursor='" + cursor + '\''
                    + ", affectedNodes=" + affectedNodes
                    + ", affectedForwardEdges=" + affectedEdges
                    + '}';
        } else {
            return "ExtraInfo{"
                    + ", affectedNodes=" + affectedNodes
                    + ", affectedForwardEdges=" + affectedEdges
                    + '}';
        }
    }
}
