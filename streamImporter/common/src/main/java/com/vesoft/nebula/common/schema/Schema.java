
package com.vesoft.nebula.common.schema;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.apache.commons.text.similarity.JaccardSimilarity;

public class Schema {
    // properties for datasource node, property name and property datatype
    protected Map<String, String> properties = new HashMap<>();

    // properties for nebula node, property name and property datatype
    protected Map<String, String> nebulaProperties = new HashMap<>();

    public Map<String, String> getProperties() {
        return properties;
    }

    public Schema setProperties(Map<String, String> properties) {
        this.properties = properties;
        return this;
    }

    public Map<String, String> getNebulaProperties() {
        return nebulaProperties;
    }

    public void setNebulaProperties(Map<String, String> nebulaProperties) {
        this.nebulaProperties = nebulaProperties;
    }

    public Schema addProperty(String propName, String propDataType) {
        this.properties.put(propName, propDataType);
        return this;
    }

    public List<String> getPropertyFields() {
        return new ArrayList<String>(properties.keySet());
    }

    public List<String> getNebulaPropertyFields(){
        return new ArrayList<>(nebulaProperties.keySet());
    }
    public String getSchemaString() {
        StringBuilder propString = new StringBuilder();
        for (Map.Entry<String, String> props : properties.entrySet()) {
            propString.append(props.getKey());
            propString.append(" ");
            propString.append(props.getValue());
            propString.append(",");
        }
        String props = propString.substring(0, propString.length() - 1);
        return props;
    }


    public Map<String, String> getPropMapping() {
        List<String> sourcePropNames = new ArrayList<>(properties.keySet());
        List<String> nebulaPropNames = new ArrayList<>(nebulaProperties.keySet());

        Map<String,String> propMapping = new HashMap<>();

        List<String> sourcePropNamesCopy = new ArrayList<>(properties.keySet());
        List<String> nebulaPropNamesCopy = new ArrayList<>(nebulaProperties.keySet());
        // similar name mapping: compute the similarity of two property name
        for (String sourcePropName : sourcePropNamesCopy) {
            double similarity = 0;
            String similarPropName = null;
            for(String nebulaPropName:nebulaPropNamesCopy){
                double similarityValue = similarity(sourcePropName, nebulaPropName);
                if(similarityValue>similarity){
                    similarity = similarityValue;
                    similarPropName = nebulaPropName;
                }
            }

           if(similarity > 0.5){
               propMapping.put(sourcePropName, similarPropName);
               sourcePropNames.remove(sourcePropName);
               nebulaPropNames.remove(similarPropName);
           }
        }
        // sequential mapping
        for(int i=0; i<sourcePropNames.size();i++){
            propMapping.put(sourcePropNames.get(i), nebulaPropNames.get(i));
        }
        return propMapping;
    }

    /**
     * compute the JaccardSimilarity of two string values.
     *
     * @param a one string value
     * @param b another string value
     * @return similarity
     * */
    private double similarity(String a, String b){
        JaccardSimilarity similarity =  new JaccardSimilarity();
        return similarity.apply(a, b);
    }
}
