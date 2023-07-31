/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.entity

case class Label(values:Map[String, String]){
  def getMapValues(schema:Map[String,String]):String = {
    values.map(kv => {
      val keyDataType = schema(kv._1)
      var value = kv._2
      keyDataType match {
        case "string" |"fixed_string" => value = s"`$value`"
        case "date" => value = s"date(\"$value\")"
        case "localtime" => value = s"localtime(\"$value\")"
        case "localdatetime" => value = s"localdatetime(\"$value\")"
        case "duration" => value = s"duration(\"$value\")"
        case _ => value = value
      }
      s"`${kv._1}`:$value"
    }).mkString(",")
  }
}

//USE nba INSERT NODE node_type_player ({id:1, name:"Tim", score: 87.0, gender: true, rate: 7.32}),({id:2, name:"Jerry", score: 95.0, gender: false, rate: 4.01}),({id:3, name:"Kyle", score: 100, gender: true, rate: 9.99})
case class Vertex(var vertexID: String, override val values: Map[String, String]) extends Label(values) {
  def getVertexString(schema:Map[String,String]): String = s"({id:$vertexID,${getMapValues(schema)})"
}

// USE nba INSERT EDGE edge_type_follow ({id:1})-[{followness:90, likeness: 66.8}]->({id:2}),({id:2})-[{followness:100, likeness: 93.35}]->({id:3})
case class Edge(var sourceID: String, var targetID: String, override val values: Map[String, String])extends Label(values){
  def getEdgeString(schema:Map[String,String]): String = s"({id:$sourceID})-[{${getMapValues(schema)}}]->({id:$targetID})"
}

object GQLTemplate {
  // USE graphName INSERT NODE/EDGE NODE_TYPE/EDGE_TYPE values
  val BATCH_INSERT_TEMPLATE = "USE %s INSERT %s `%s` %s"
}
