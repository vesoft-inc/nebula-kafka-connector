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
    private long affectedForwardEdges;
    // the number of affected reverse edges
    private long affectedReverseEdges;

    public ExtraInfo() {
        this.cursor = null;
        affectedNodes = 0;
        affectedForwardEdges = 0;
        affectedReverseEdges = 0;
    }

    public void setCursor(String cursor) {
        this.cursor = cursor;
    }

    public void setAffectedNodes(long affectedNodes) {
        this.affectedNodes = affectedNodes;
    }

    public void setAffectedForwardEdges(long affectedForwardEdges) {
        this.affectedForwardEdges = affectedForwardEdges;
    }

    public void setAffectedReverseEdges(long affectedReverseEdges) {
        this.affectedReverseEdges = affectedReverseEdges;
    }

    public String getCursor() {
        return this.cursor;
    }

    public long getAffectedNodes() {
        return this.affectedNodes;
    }

    public long getAffectedForwardEdges() {
        return this.affectedForwardEdges;
    }

    public long getAffectedReverseEdges() {
        return this.affectedReverseEdges;
    }

    @Override
    public String toString() {
        if (cursor != null) {
            return "ExtraInfo{"
                    + "cursor='" + cursor + '\''
                    + ", affectedNodes=" + affectedNodes
                    + ", affectedForwardEdges=" + affectedForwardEdges
                    + ", affectedReverseEdges=" + affectedReverseEdges
                    + '}';
        } else {
            return "ExtraInfo{"
                    + ", affectedNodes=" + affectedNodes
                    + ", affectedForwardEdges=" + affectedForwardEdges
                    + ", affectedReverseEdges=" + affectedReverseEdges
                    + '}';
        }
    }
}
