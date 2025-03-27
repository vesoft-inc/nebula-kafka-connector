
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
    private List<String>     srcFields        = Arrays.asList("src");
    private List<String>     dstFields        = Arrays.asList("dst");
    private List<String>     kafkaFieldNames  = Arrays.asList("dura", "ty", "de");
    private List<String>     nebulaFieldNames = Arrays.asList("duration", "type", "degree");
    private NebulaEdge       edge             = null;
    private NebulaEdgeSchema edgeSchema       = new NebulaEdgeSchema();

    @Before
    public void setup() {
        Map<String, String> props  = new HashMap<>();
        Map<String, String> srcPks = new HashMap<>();
        srcPks.put("id", "1");
        Map<String, String> dstPks = new HashMap<>();
        dstPks.put("id", "2");
        props.put("duration", "10");
        props.put("type", "\"friend\"");
        props.put("degree", "5");
        edge = new NebulaEdge(srcPks, dstPks, props);

        Map<String, String> schema = new HashMap<>();
        schema.put("duration", "INT64");
        schema.put("type", "STRING");
        schema.put("degree", "INT64");
        edgeSchema.setEdgeTypeName("friend");
        edgeSchema.setSourceNodeTypeName("person");
        Map<String, String> srcPksSchema = new HashMap<>();
        srcPksSchema.put("id", "INT64");
        Map<String, String> dstPksSchema = new HashMap<>();
        dstPksSchema.put("id", "INT64");
        edgeSchema.setSourcePkNameAndType(srcPksSchema);
        edgeSchema.setTargetPkNameAndType(dstPksSchema);
        edgeSchema.setTargetNodeTypeName("person");
        edgeSchema.setProperties(schema);
    }

    @Test
    public void testGetEdgeInsertStatement() {
        NebulaEdges nebulaEdges = new NebulaEdges(edgeSchema,
                                                  srcFields,
                                                  dstFields,
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
        String expectStatement = "TABLE t{src_id,dst_id,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (n_src@person) WHERE n_src.id=r.src_id "
                + "OPTIONAL MATCH (n_dst@person) WHERE n_dst.id=r.dst_id \n"
                + "INSERT (n_src)-[@friend{duration:r.dura,type:r.ty,degree:r.de}]->(n_dst)";
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
                                                  srcFields,
                                                  dstFields,
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
        String expectStatement = "TABLE t{src_id,dst_id,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (n_src@person) WHERE n_src.id=r.src_id "
                + "OPTIONAL MATCH (n_dst@person) WHERE n_dst.id=r.dst_id \n"
                + "INSERT OR IGNORE (n_src)-[@friend{duration:r.dura,type:r.ty,degree:r.de}]"
                + "->(n_dst)";
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
                                                  srcFields,
                                                  dstFields,
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
        String expectStatement = "TABLE t{src_id,dst_id,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (n_src@person) WHERE n_src.id=r.src_id "
                + "OPTIONAL MATCH (n_dst@person) WHERE n_dst.id=r.dst_id \n"
                + "INSERT OR REPLACE (n_src)-[@friend{duration:r.dura,type:r.ty,degree:r.de}]"
                + "->(n_dst)";
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
                                                  srcFields,
                                                  dstFields,
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
        String expectStatement = "TABLE t{src_id,dst_id,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (n_src@person)-[n_e@friend]->(n_dst@person) "
                + "WHERE n_src.id=r.src_id AND n_dst.id=r.dst_id \n"
                + "SET n_e.duration=r.dura,n_e.type=r.ty,n_e.degree=r.de";
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
                                                  srcFields,
                                                  dstFields,
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
        String expectStatement = "TABLE t{src_id,dst_id,dura,ty,de} = \n"
                + "(1,2,10,\"friend\",5) \n"
                + "USE nba \n"
                + "FOR r IN t \n"
                + "OPTIONAL MATCH (n_src@person)-[n_e@friend]->(n_dst@person) "
                + "WHERE n_src.id=r.src_id AND n_dst.id=r.dst_id \n"
                + "DELETE n_e";
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
