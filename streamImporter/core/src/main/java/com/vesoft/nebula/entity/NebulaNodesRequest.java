
package com.vesoft.nebula.entity;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

public class NebulaNodesRequest implements Serializable {
    private int graphId;
    private int partId;
    private int nodeTypeId;
    private List<NodeToWrite> nodes = new ArrayList<>();

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

    public int getNodeTypeId() {
        return nodeTypeId;
    }

    public void setNodeTypeId(int nodeTypeId) {
        this.nodeTypeId = nodeTypeId;
    }

    public List<NodeToWrite> getNodes() {
        return nodes;
    }

    public void setNodes(List<NodeToWrite> nodes) {
        this.nodes = nodes;
    }

    @Override
    public String toString() {
        return "NebulaNodesRequest{" +
                "graphId=" + graphId +
                ", partId=" + partId +
                ", nodeTypeId=" + nodeTypeId +
                ", nodes=" + nodes +
                '}';
    }
}
