/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.proto.graph.Record;
import com.vesoft.nebula.proto.graph.Value;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;

public class NRecord {
    private Record record;

    public NRecord(Record record) {
        this.record = record;
    }

    /**
     * returns true if this record contains a mapping for the specified key. The key cannot be null.
     *
     * @param key key whose presence in this record is to be checked
     * @return true if this record contains a mapping for the specified key
     * @throws NullPointerException if the specified key is null
     */
    public boolean containsKey(String key) {
        if (key == null) {
            throw new NullPointerException("null map key");
        }
        return record.containsValues(key);
    }


    /**
     * get the Value of specified key in this record
     *
     * @param key key whose corresponding value in this record will be returned
     * @return The ValueWrapper of the specified key
     */
    public ValueWrapper getValue(String key) {
        if (!containsKey(key)) {
            return null;
        }
        return new ValueWrapper(record.getValuesOrDefault(key, null));
    }

    /**
     * Returns true if this record has no values
     *
     * @return true if this record contains no key-value mappings
     */
    public boolean isEmpty() {
        return record.getValuesCount() == 0;
    }

    /**
     * get the number of key-value mappings in this record.
     *
     * @return the number of key-value mappings in this record
     */
    public int size() {
        return record.getValuesCount();
    }

    /**
     * get the Map object for this record
     *
     * @return Map for this record
     */
    public Map<String, ValueWrapper> getValuesMap() {
        Map<String, ValueWrapper> values = new HashMap<>();
        for (Map.Entry<String, Value> entry : record.getValuesMap().entrySet()) {
            values.put(entry.getKey(), new ValueWrapper(entry.getValue()));
        }
        return values;
    }

    @Override
    public String toString() {
        return record.getValuesMap().toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        NRecord that = (NRecord) o;
        return record.getValuesMap().equals(that.record.getValuesMap());
    }

    @Override
    public int hashCode() {
        return Objects.hash(record.getValuesMap());
    }
}
