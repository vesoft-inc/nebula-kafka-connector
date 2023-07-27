/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.entity

//USE nba INSERT NODE node_type_player ({id:1, name:"Tim", score: 87.0, gender: true, rate: 7.32}),({id:2, name:"Jerry", score: 95.0, gender: false, rate: 4.01}),({id:3, name:"Kyle", score: 100, gender: true, rate: 9.99})
case class Vertex(vertexID: String, values: Map[String, String]) {
  def getVertexString: String = s"({id:$vertexID,$getMapValues})"

  def getMapValues: String = values.map(kv => s"${kv._1}:${kv._2}").mkString(",")
}

// USE nba INSERT EDGE edge_type_follow ({id:1})-[{followness:90, likeness: 66.8}]->({id:2}),({id:2})-[{followness:100, likeness: 93.35}]->({id:3})
case class Edge(sourceID: String, targetID: String, values: Map[String, String]) {
  def getEdgeString: String = s"({id:$sourceID})-[{$getMapValues}]->({id:$targetID})"

  def getMapValues: String = values.map(kv => s"${kv._1}:${kv._2}").mkString(",")
}

class GQLTemplate {
  // USE graphName INSERT NODE/EDGE NODE_TYPE/EDGE_TYPE values
  val BATCH_INSERT_TEMPLATE = "USE %s INSERT %s %s %s"
}
