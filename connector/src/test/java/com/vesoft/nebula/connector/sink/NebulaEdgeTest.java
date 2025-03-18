
package com.vesoft.nebula.connector.sink;

import static junit.framework.TestCase.assertEquals;

import com.vesoft.nebula.connector.config.NebulaSinkConnectConfig;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.junit.Before;
import org.junit.Test;

public class NebulaEdgeTest {
    private String           graphName        = "nba";
    private String           srcField         = "src";
    private String           dstField         = "dst";
    private List<String>     kafkaFieldNames  = Arrays.asList("dura", "ty", "de");
    private List<String>     nebulaFieldNames = Arrays.asList("duration", "type", "degree");
    private NebulaEdge       edge             = null;
    private NebulaEdgeSchema edgeSchema       = new NebulaEdgeSchema();

    @Before
    public void setup() {
        Map<String, String> props = new HashMap<>();
        props.put("duration", "10");
        props.put("type", "\"friend\"");
        props.put("degree", "5");
        edge = new NebulaEdge("1", "2", props);

        Map<String, String> schema = new HashMap<>();
        schema.put("duration", "INT64");
        schema.put("type", "STRING");
        schema.put("degree", "INT64");
        edgeSchema.setEdgeTypeName("friend");
        edgeSchema.setSourceNodeTypeName("person");
        edgeSchema.setSourceNodePkName("id");
        edgeSchema.setSourceNodePkType("INT64");
        edgeSchema.setTargetNodeTypeName("person");
        edgeSchema.setTargetNodePkName("id");
        edgeSchema.setTargetNodePkType("INT64");
        edgeSchema.setProperties(schema);
    }

    @Test
    public void testGetEdgeInsertStatement() {
        NebulaEdges nebulaEdges = new NebulaEdges(edgeSchema,
                                                  srcField,
                                                  dstField,
                                                  kafkaFieldNames,
                                                  nebulaFieldNames,
                                                  Arrays.asList(edge));
        String insertStatement = nebulaEdges
                .getInsertStatement(graphName, NebulaSinkConnectConfig.InsertMode.INSERT);
        String insertChars = insertStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        String expectStatement = "TABLE t{src,dst,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (src_node@person) WHERE src_node.id=r.src "
                + "OPTIONAL MATCH (dst_node@person) WHERE dst_node.id=r.dst \n"
                + "INSERT (src_node)-[@friend{duration:r.dura,type:r.ty,degree:r.de}]->(dst_node)";
        String expectChars = expectStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars, insertChars);

    }


    @Test
    public void testGetEdgeInsertIgnoreStatement() {
        NebulaEdges nebulaEdges = new NebulaEdges(edgeSchema,
                                                  srcField,
                                                  dstField,
                                                  kafkaFieldNames,
                                                  nebulaFieldNames,
                                                  Arrays.asList(edge));
        String insertStatement = nebulaEdges
                .getInsertStatement(graphName, NebulaSinkConnectConfig.InsertMode.INSERTIGNORE);
        String insertChars = insertStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        String expectStatement = "TABLE t{src,dst,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (src_node@person) WHERE src_node.id=r.src "
                + "OPTIONAL MATCH (dst_node@person) WHERE dst_node.id=r.dst \n"
                + "INSERT OR IGNORE (src_node)-[@friend{duration:r.dura,type:r.ty,degree:r.de}]"
                + "->(dst_node)";
        String expectChars = expectStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars, insertChars);

    }


    @Test
    public void testGetEdgeInsertReplaceStatement() {
        NebulaEdges nebulaEdges = new NebulaEdges(edgeSchema,
                                                  srcField,
                                                  dstField,
                                                  kafkaFieldNames,
                                                  nebulaFieldNames,
                                                  Arrays.asList(edge));
        String insertStatement = nebulaEdges
                .getInsertStatement(graphName, NebulaSinkConnectConfig.InsertMode.INSERTREPLACE);
        String insertChars = insertStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        String expectStatement = "TABLE t{src,dst,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (src_node@person) WHERE src_node.id=r.src "
                + "OPTIONAL MATCH (dst_node@person) WHERE dst_node.id=r.dst \n"
                + "INSERT OR REPLACE (src_node)-[@friend{duration:r.dura,type:r.ty,degree:r.de}]"
                + "->(dst_node)";
        String expectChars = expectStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars, insertChars);
    }


    @Test
    public void testGetEdgeUpdateStatement() {
        NebulaEdges nebulaEdges = new NebulaEdges(edgeSchema,
                                                  srcField,
                                                  dstField,
                                                  kafkaFieldNames,
                                                  nebulaFieldNames,
                                                  Arrays.asList(edge));
        String updateStatement = nebulaEdges
                .getUpdateStatement(graphName);
        String updateChars = updateStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        String expectStatement = "TABLE t{src,dst,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (src_node@person)-[e@friend]->(dst_node@person) "
                + "WHERE src_node.id=r.src AND dst_node.id=r.dst \n"
                + "SET e.duration=r.dura,e.type=r.ty,e.degree=r.de";
        String expectChars = expectStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars, updateChars);
    }


    @Test
    public void testGetEdgeDeleteStatement() {
        NebulaEdges nebulaEdges = new NebulaEdges(edgeSchema,
                                                  srcField,
                                                  dstField,
                                                  kafkaFieldNames,
                                                  nebulaFieldNames,
                                                  Arrays.asList(edge));
        String deleteStatement = nebulaEdges
                .getDeleteStatement(graphName);
        String deleteChars = deleteStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        String expectStatement = "TABLE t{src,dst,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (src_node@person)-[e@friend]->(dst_node@person) "
                + "WHERE src_node.id=r.src AND dst_node.id=r.dst \n"
                + "DELETE e";
        String expectChars = expectStatement
                .chars()
                .sorted()
                .collect(StringBuilder::new,
                         StringBuilder::appendCodePoint,
                         StringBuilder::append)
                .toString();
        assertEquals(expectChars, deleteChars);

    }
}
