
package com.vesoft.nebula.spark.common

import com.vesoft.nebula.spark.common.nebula.{GraphProvider}
import org.scalatest.BeforeAndAfterAll
import org.scalatest.funsuite.AnyFunSuite

import java.util
import java.util.concurrent.TimeUnit
import scala.collection.JavaConversions.asJavaCollection
import scala.collection.JavaConverters.asScalaBufferConverter

class GraphProviderSuite extends AnyFunSuite with BeforeAndAfterAll {
  var graphProvider: GraphProvider = null

  override def beforeAll(): Unit = {
    val address = "192.168.8.6:3820"
    val authOptions = new util.HashMap[String, Object]()
    authOptions.put("password", "nebula")
    graphProvider = new GraphProvider(address, "root", authOptions, 3000)

    val createSchema = "CREATE GRAPH TYPE graph_type_nba AS {" +
      "(node_type_player LABEL player {id INT PRIMARY KEY, name STRING, score FLOAT, gender bool, rate DOUBLE})," +
      "(node_type_player)-[edge_type_follow LABEL follow {followness INT, likeness FLOAT64}]->(node_type_player)}"
    val resp = graphProvider.submit(createSchema)
    if (!resp.isSucceeded) {
      System.out.println("create graph type failed, " + resp.getErrorMessage)
      System.exit(1)
    }
    TimeUnit.SECONDS.sleep(5)
  }

  override def afterAll(): Unit = graphProvider.close()


  test("getNodeDesc") {
    val nodeDesc = graphProvider.getNodeDesc("nba", "node_type_player")
    assert(nodeDesc.nodeTypeName.equals("node_type_player"))
    assert(nodeDesc.properties.size == 5)
    assert(nodeDesc.properties.keySet.contains("id"))
    assert(nodeDesc.properties.keySet.contains("name"))
    assert(nodeDesc.properties.keySet.contains("score"))
    assert(nodeDesc.properties.keySet.contains("gender"))
    assert(nodeDesc.properties.keySet.contains("rate"))
  }

  test("getEdgeDesc") {
    val edgeDesc = graphProvider.getEdgeDesc("nba", "edge_type_follow")
    assert(edgeDesc.edgeTypeName.equals("edge_type_follow"))
    assert(edgeDesc.srcNodeTypeName.equals("node_type_player"))
    assert(edgeDesc.dstNodeTypeName.equals("node_type_player"))
    assert(edgeDesc.srcNodePkDataTypeMap("id").equals("INT64"))
    assert(edgeDesc.dstNodePkDataTypeMap("id").equals("INT64"))
    assert(edgeDesc.properties.size == 2)
    assert(edgeDesc.properties.keySet.contains("followness"))
    assert(edgeDesc.properties.keySet.contains("likeness"))
  }

  test("getAllParts") {
    val parts: List[Integer] = graphProvider.getAllParts.asScala.toList
    assert(parts.size() == 10)
    val expectParts = List(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    assert(parts.containsAll(expectParts))
  }

}
