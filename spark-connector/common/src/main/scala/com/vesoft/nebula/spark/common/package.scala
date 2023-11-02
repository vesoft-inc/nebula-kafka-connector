/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

case class NebulaVertex(id: String, values: Map[String, String]) {
  def getVertexStr: String = s"({id:$id,${getMapValues(values)}})"

  private[this] def getMapValues(values: Map[String, String]): String = {
    values
      .map(kv => s"`${kv._1}`:${kv._2}")
      .mkString(",")
  }
}

case class NebulaVertices(nodeType: String, values: List[NebulaVertex]) {
  def getVerticesStr = values.map(v => v.getVertexStr).mkString(",")
}

case class NebulaEdge(srcId: String, dstId: String, values: Map[String, String]) {
  def getEdgeStr: String = s"({id:$srcId})-[{${getMapValues(values)}]->({id:$dstId})"

  private[this] def getMapValues(values: Map[String, String]): String = {
    values
      .map(kv => s"`${kv._1}`:${kv._2}")
      .mkString(",")
  }
}

case class NebulaEdges(edgeType: String, values: List[NebulaEdge]) {
  def getEdgesStr = values.map(e => e.getEdgeStr).mkString(",")
}
