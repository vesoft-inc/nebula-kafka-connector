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

public class NebulaNodeTest {
    private NebulaNode node = null;
    private NebulaNodeSchema nodeSchema = new NebulaNodeSchema();

    @Before
    public void setup() {
        Map<String, Object> props = new HashMap<>();
        props.put("name", "Tom");
        props.put("age", 18);
        props.put("weight", 100);
        props.put("gender", "male");
        node = new NebulaNode(props);

        Map<String, String> schema = new HashMap<>();
        schema.put("id","STRING");
        schema.put("name", "STRING");
        schema.put("age", "INT64");
        schema.put("weight", "INT32");
        schema.put("gender", "STRING");
        nodeSchema.setNodeTypeName("player");
        nodeSchema.setNodePkType("STRING");
        nodeSchema.setNodeProperties(schema);
    }

    @Test
    public void testGetNodeStatement() {
        try {
            String statement =
                    node.getNodeStatement(nodeSchema)
                            .chars()
                            .sorted()
                            .collect(StringBuilder::new,
                                    StringBuilder::appendCodePoint,
                                    StringBuilder::append)
                            .toString();
            String expect = "({`id`:\"1\",`name`:\"Tom\",`age`:18,`weight`:100,`gender`:\"male\"})"
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
    public void testGetNodeString() {
        String nodeString = node.getNodeString()
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                        StringBuilder::appendCodePoint,
                        StringBuilder::append)
                .toString();
        String expect =
                "{`vid`:1,`name`:Tom,`age`:18,`gender`:male,`weight`:100}"
                        .chars()
                        .sorted()
                        .collect(StringBuilder::new,
                                StringBuilder::appendCodePoint,
                                StringBuilder::append)
                        .toString();
        assertEquals(expect, nodeString);
    }
}
