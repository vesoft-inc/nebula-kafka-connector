
package com.vesoft.nebula.entity;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

public class NebulaEdgesRequest implements Serializable {
    private int graphId;
    private int partId;
    private int edgeTypeId;
    private List<EdgeToWrite> edges=new ArrayList<>();

    public int getGraphId() {
        return graphId;
    }

    public void setGraphId(int graphId) {
        this.graphId = graphId;
    }

    public int getPartId() {
        return partId;
    }

    public void setPartId(int partId) {
        this.partId = partId;
    }

    public int getEdgeTypeId() {
        return edgeTypeId;
    }

    public void setEdgeTypeId(int edgeTypeId) {
        this.edgeTypeId = edgeTypeId;
    }

    public List<EdgeToWrite> getEdges() {
        return edges;
    }

    public void setEdges(List<EdgeToWrite> edges) {
        this.edges = edges;
    }

    @Override
    public String toString() {
        return "NebulaEdgesRequest{" +
                "graphId=" + graphId +
                ", partId=" + partId +
                ", edgeTypeId=" + edgeTypeId +
                ", edges=" + edges +
                '}';
    }
}
