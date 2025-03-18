
package com.vesoft.nebula.connector.sink;

import static junit.framework.TestCase.assertEquals;

import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import com.vesoft.nebula.connector.exceptions.DataFormatException;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import org.junit.Before;
import org.junit.Test;

public class NebulaNodeTest {
    private String           graphName  = "nba";
    private NebulaNode       node       = null;
    private NebulaNodeSchema nodeSchema = new NebulaNodeSchema();

    @Before
    public void setup() {
        Map<String, String> props = new HashMap<>();
        props.put("id", "\"1\"");
        props.put("name", "\"Tom\"");
        props.put("age", "18");
        props.put("weight", "100");
        props.put("gender", "\"male\"");
        node = new NebulaNode(props);

        Map<String, String> schema = new HashMap<>();
        schema.put("id", "STRING");
        schema.put("name", "STRING");
        schema.put("age", "INT64");
        schema.put("weight", "INT32");
        schema.put("gender", "STRING");
        nodeSchema.setNodeTypeName("player");
        nodeSchema.setNodePkName("id");
        nodeSchema.setNodePkType("STRING");
        nodeSchema.setNodeProperties(schema);
    }

    @Test
    public void testGetNodeInsertStatement() {
        NebulaNodes nodes = new NebulaNodes(nodeSchema, Arrays.asList(node));
        String statement =
                nodes.getInsertStatement(graphName, NebulaSinkConnectConfig.InsertMode.INSERT);
        String expectStatement = "TABLE t{id,name,age,weight,gender} = \n"
                + "(\"1\",\"Tom\",18,100,\"male\") \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "INSERT (@player{id:r.id,name:r.name,age:r.age,weight:r.weight,gender:r.gender})";
        String expectChars = expectStatement.chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars,
                     statement
                             .chars()
                             .sorted()
                             .collect(StringBuilder::new,
                                      StringBuilder::appendCodePoint,
                                      StringBuilder::append)
                             .toString());
    }

    @Test
    public void testGetNodeInsertIgnoreStatement() {
        NebulaNodes nodes = new NebulaNodes(nodeSchema, Arrays.asList(node));
        String statement =
                nodes.getInsertStatement(graphName,
                                         NebulaSinkConnectConfig.InsertMode.INSERTIGNORE);
        String expectStatement = "TABLE t{id,name,age,weight,gender} = \n"
                + "(\"1\",\"Tom\",18,100,\"male\") \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "INSERT OR IGNORE "
                + "(@player{id:r.id,name:r.name,age:r.age,weight:r.weight,gender:r.gender})";
        String expectChars = expectStatement.chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars,
                     statement
                             .chars()
                             .sorted()
                             .collect(StringBuilder::new,
                                      StringBuilder::appendCodePoint,
                                      StringBuilder::append)
                             .toString());
    }


    @Test
    public void testGetNodeInsertReplaceStatement() {
        NebulaNodes nodes = new NebulaNodes(nodeSchema, Arrays.asList(node));
        String statement =
                nodes.getInsertStatement(graphName,
                                         NebulaSinkConnectConfig.InsertMode.INSERTREPLACE);
        String expectStatement = "TABLE t{id,name,age,weight,gender} = \n"
                + "(\"1\",\"Tom\",18,100,\"male\") \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "INSERT OR REPLACE "
                + "(@player{id:r.id,name:r.name,age:r.age,weight:r.weight,gender:r.gender})";
        String expectChars = expectStatement.chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars,
                     statement
                             .chars()
                             .sorted()
                             .collect(StringBuilder::new,
                                      StringBuilder::appendCodePoint,
                                      StringBuilder::append)
                             .toString());
    }

    @Test
    public void testGetNodeUpdateStatement() {
        NebulaNodes nodes = new NebulaNodes(nodeSchema, Arrays.asList(node));
        String statement =
                nodes.getUpdateStatement(graphName);
        String expectStatement = "TABLE t{id,name,age,weight,gender} = \n"
                + "(\"1\",\"Tom\",18,100,\"male\") \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (v@player) WHERE v.id=r.id \n"
                + "SET v.name=r.name,v.age=r.age,v.weight=r.weight,v.gender=r.gender";
        String expectChars = expectStatement.chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars,
                     statement
                             .chars()
                             .sorted()
                             .collect(StringBuilder::new,
                                      StringBuilder::appendCodePoint,
                                      StringBuilder::append)
                             .toString());
    }


    @Test
    public void testGetNodeDeleteStatement() {
        NebulaNodes nodes = new NebulaNodes(nodeSchema, Arrays.asList(node));
        String statement =
                nodes.getDeleteStatement(graphName, NebulaSinkConnectConfig.InsertMode.DELETE);
        String expectStatement = "TABLE t{id,name,age,weight,gender} = \n"
                + "(\"1\",\"Tom\",18,100,\"male\") \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (v@player) WHERE v.id=r.id \n"
                + "DELETE v";
        String expectChars = expectStatement.chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars,
                     statement
                             .chars()
                             .sorted()
                             .collect(StringBuilder::new,
                                      StringBuilder::appendCodePoint,
                                      StringBuilder::append)
                             .toString());
    }


    @Test
    public void testGetNodeDetachDeleteStatement() {
        NebulaNodes nodes = new NebulaNodes(nodeSchema, Arrays.asList(node));
        String statement =
                nodes.getDeleteStatement(graphName,
                                         NebulaSinkConnectConfig.InsertMode.DETACHDELETE);
        String expectStatement = "TABLE t{id,name,age,weight,gender} = \n"
                + "(\"1\",\"Tom\",18,100,\"male\") \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (v@player) WHERE v.id=r.id \n"
                + "DETACH DELETE v";
        String expectChars = expectStatement.chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars,
                     statement
                             .chars()
                             .sorted()
                             .collect(StringBuilder::new,
                                      StringBuilder::appendCodePoint,
                                      StringBuilder::append)
                             .toString());
    }



}
