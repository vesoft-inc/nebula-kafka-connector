/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.sink;

import static junit.framework.TestCase.assertEquals;
import com.vesoft.nebula.connector.exceptions.DataFormatException;
import java.util.HashMap;
import java.util.Map;
import org.junit.Before;
import org.junit.Test;

public class NebulaEdgeTest {
    private NebulaEdge edge = null;
    private NebulaEdgeSchema edgeSchema = new NebulaEdgeSchema();

    @Before
    public void setup() {
        Map<String, Object> props = new HashMap<>();
        props.put("duration", 10);
        props.put("type", "friend");
        props.put("degree", 5);
        edge = new NebulaEdge("srcId","1","dstId", "2", props);

        Map<String, String> schema = new HashMap<>();
        schema.put("duration", "INT64");
        schema.put("type", "STRING");
        schema.put("degree", "INT64");
        edgeSchema.setSourceNodeTypeName("person");
        edgeSchema.setSourceNodePkType("INT64");
        edgeSchema.setTargetNodeTypeName("person");
        edgeSchema.setTargetNodePkType("INT64");
        edgeSchema.setProperties(schema);
    }

    @Test
    public void testGetEdgeStatement() {
        try {
            String statement =
                    edge.getEdgeStatement(edgeSchema)
                            .chars()
                            .sorted()
                            .collect(StringBuilder::new,
                                    StringBuilder::appendCodePoint,
                                    StringBuilder::append)
                            .toString();
            String expect = "({`id`:1})-[{`duration`:10,`type`:\"friend\",`degree`:5}]->({`id`:2})"
                    .chars()
                    .sorted()
                    .collect(StringBuilder::new,
                            StringBuilder::appendCodePoint,
                            StringBuilder::append)
                    .toString();
            assertEquals(expect, statement);
        } catch (DataFormatException e) {
            assert (false);
        }
    }

    @Test
    public void testGetEdgeString() {
        String edgeString = edge.getEdgeString()
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                        StringBuilder::appendCodePoint,
                        StringBuilder::append)
                .toString();
        String expect =
                "{`srcPk`:1,`dstPk`:2,`duration`:10,`degree`:5,`type`:friend}"
                        .chars()
                        .sorted()
                        .collect(StringBuilder::new,
                                StringBuilder::appendCodePoint,
                                StringBuilder::append)
                        .toString();
        assertEquals(expect, edgeString);
    }
}
