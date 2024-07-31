package com.vesoft.nebula.client.graph.data;

import java.io.Serializable;

public abstract class BaseDataObject implements Serializable {
    private String decodeType = "utf-8";

    public String getDecodeType() {
        return decodeType;
    }

    public BaseDataObject setDecodeType(String decodeType) {
        this.decodeType = decodeType;
        return this;
    }

}
