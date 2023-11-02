/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import com.vesoft.nebula.spark.common.nebula.{GraphProvider, VidType}
import org.scalatest.BeforeAndAfterAll
import org.scalatest.funsuite.AnyFunSuite

import java.util.concurrent.TimeUnit

class GraphProviderSuite extends AnyFunSuite with BeforeAndAfterAll {
  var graphProvider: GraphProvider = null

  override def beforeAll(): Unit = {
    val address = "192.168.8.6:3713"
    graphProvider = new GraphProvider(address, "root", "nebula", 3000, 1)

    val createSchema = "CREATE GRAPH TYPE graph_type_nba IF NOT EXISTS AS GRAPH TYPE {" +
      "(node_type_player(id) LABEL player {id INT, name STRING, score FLOAT, gender bool, rate DOUBLE})," +
      "(node_type_player)-[edge_type_follow LABEL follow {followness INT, likeness FLOAT64}]->(node_type_player)}"
    val resp = graphProvider.submit(createSchema)
    if (!resp.isSucceeded) {
      System.out.println("create graph type failed, " + resp.getGqlStatus)
      System.exit(1)
    }
    TimeUnit.SECONDS.sleep(5)
  }

  override def afterAll(): Unit = graphProvider.close()

  test("getIdType for node") {
    val idType = graphProvider.getIdType("nba", "node_type_player")
    assert(idType == VidType.INT)
  }

  test("getIdsType for edge") {
    val (sourceIdType, targetIdType) = graphProvider.getIdsType("nba", "edge_type_follow")
    assert(sourceIdType == VidType.INT)
    assert(targetIdType == VidType.INT)
  }

  test("getNodesType for edge") {
    val (sourceType, targetType) = graphProvider.getNodesType("nba", "edge_type_follow")
    assert(sourceType.equals("node_type_player"))
    assert(targetType.equals("node_type_player"))
  }

  test("getVertexSchema") {
    val schema = graphProvider.getTagSchema("nba", "node_type_player")
    assert(schema.size == 5)
    assert(schema.keySet.contains("id"))
    assert(schema.keySet.contains("name"))
    assert(schema.keySet.contains("score"))
    assert(schema.keySet.contains("gender"))
    assert(schema.keySet.contains("rate"))
  }

  test("getEdgeSchema") {
    val schema = graphProvider.getEdgeSchema("nba", "edge_type_follow")
    assert(schema.size == 2)
    assert(schema.keySet.contains("followness"))
    assert(schema.keySet.contains("likeness"))
  }

}
