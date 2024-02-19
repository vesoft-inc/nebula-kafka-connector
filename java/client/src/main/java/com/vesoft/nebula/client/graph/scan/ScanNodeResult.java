/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.scan;

import com.vesoft.nebula.client.graph.data.ResultSet;
import com.vesoft.nebula.client.graph.data.ValueWrapper;
import com.vesoft.nebula.client.graph.data.Vertex;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class ScanNodeResult extends ScanResult {
    private List<String> propNames;

    public ScanNodeResult(List<ResultSet> results, List<String> propNames) {
        super(results);
        this.propNames = propNames;
    }

    /**
     * get node table row's column names
     *
     * @return list of row column names
     */
    public List<String> getPropNames() {
        return propNames;
    }

    protected void convertResultToRow() {
        if (isEmpty) {
            return;
        }
        if (tableRows.isEmpty()) {
            for (ResultSet resultSet : results) {
                List<ResultSet.Record> records = resultSet.getRows();
                for (ResultSet.Record record : records) {
                    List<ValueWrapper> values = record.values();
                    List<ValueWrapper> rowValues = new ArrayList<>();
                    Map<String, ValueWrapper> properties;
                    Vertex vertex = values.get(0).asNode();
                    properties = vertex.getProperties();
                    for (String propName : propNames) {
                        rowValues.add(properties.get(propName));
                    }
                    tableRows.add(new TableRow(rowValues));
                }
            }
        }
    }
}
