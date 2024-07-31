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

                while (resultSet.hasNext()) {
                    ResultSet.Record record = resultSet.next();
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
