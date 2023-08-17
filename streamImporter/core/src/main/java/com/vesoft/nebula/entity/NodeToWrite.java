/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.entity;

import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;

public class NodeToWrite implements Serializable {
    private long internalId;
    private Map<String, String> properties = new HashMap<>();

    public NodeToWrite(long internalId, Map<String,String> properties){
        this.internalId = internalId;
        this.properties = properties;
    }

    public long getInternalId() {
        return internalId;
    }

    public void setInternalId(long internalId) {
        this.internalId = internalId;
    }

    public Map<String, String> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, String> properties) {
        this.properties = properties;
    }

    @Override
    public String toString() {
        return "NodeToWrite{" +
                "primaryKey='" + internalId + '\'' +
                ", properties=" + properties +
                '}';
    }
}
