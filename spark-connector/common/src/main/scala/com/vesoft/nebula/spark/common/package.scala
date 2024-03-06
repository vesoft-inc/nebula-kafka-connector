/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

case class NebulaNode(values: Map[String, String]) {
  def getNodeStr: String = s"({${getMapValues(values)}})"

  private[this] def getMapValues(values: Map[String, String]): String = {
    values
      .map(kv => s"`${kv._1}`:${kv._2}")
      .mkString(",")
  }
}

case class NebulaNodes(nodeType: String, values: List[NebulaNode]) {
  def getNodesStr = values.map(v => v.getNodeStr).mkString(",")
}

case class NebulaEdge(srcPkName:String, srcId: String, dstPkName:String, dstId: String, values: Map[String, String]) {
  def getEdgeStr: String = s"({`$srcPkName`:$srcId})-[{${getMapValues(values)}}]->({`$dstPkName`:$dstId})"

  private[this] def getMapValues(values: Map[String, String]): String = {
    values
      .map(kv => s"`${kv._1}`:${kv._2}")
      .mkString(",")
  }
}

case class NebulaEdges(edgeType: String, values: List[NebulaEdge]) {
  def getEdgesStr = values.map(e => e.getEdgeStr).mkString(",")
}
