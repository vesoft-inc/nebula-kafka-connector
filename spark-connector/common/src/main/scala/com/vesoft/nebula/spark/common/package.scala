/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import com.vesoft.nebula.spark.common.nebula.VidType

import scala.collection.mutable.ListBuffer

case class NebulaNode(values: Map[String, String])

case class NebulaNodes(nodeType: String, values: List[NebulaNode], pkName: String, fieldTypeMap: Map[String, String]) {
  private val propNames = values.iterator.next().values.keySet.toSeq

  def getNodesStr = values.map(node => {
    s"(${propNames.map(prop => s"${node.values(prop)}").mkString(",")})"
  }).mkString(",")

  def propNamesStr = propNames.mkString(",")

  def propNamesWithTableStr = propNames.map(prop => s"`$prop`:CAST(r.$prop AS ${fieldTypeMap(prop)})").mkString(",")
}

case class NebulaEdge(srcId: String, dstId: String, values: Map[String, String])

case class NebulaEdges(edgeType: String,
                       srcType: String,
                       srcPkName: String,
                       srcPkDataType: VidType.Value,
                       dfSrcField: String,
                       dstType: String,
                       dstPkName: String,
                       dstPkDataType:VidType.Value,
                       dfDstField: String,
                       values: List[NebulaEdge],
                       fieldTypeMap: Map[String, String]) {
  private val propNames = values.iterator.next().values.keySet.toSeq

  def getEdgesStr = {
    values.map(edge => {
      s"(${edge.srcId},${edge.dstId}," +
        s"${
          propNames
            .filterNot(p => p.equals(dfSrcField) || p.equals(dfDstField))
            .map(prop => s"${edge.values(prop)}")
            .mkString(",")
        })"
    }).mkString(",")
  }

  def getSrcPkStr: String = s"{`${srcPkName}`:CAST(r.${dfSrcField} AS ${srcPkDataType.toString})}"

  def getDstPkStr: String = s"{`${dstPkName}`:CAST(r.${dfDstField} AS ${dstPkDataType.toString})}"


  def tableHeaders: String = {
    val propNamesWithoutPks = propNames.filterNot(p => p.equals(dfSrcField) || p.equals(dfDstField))
    val pksAndPropNames = Seq(dfSrcField, dfDstField) ++ propNamesWithoutPks
    pksAndPropNames.mkString(",")
  }

  def propNamesWithTableStr: String = propNames.map(prop => s"`$prop`:CAST(r.$prop AS ${fieldTypeMap(prop)})").mkString(",")
}
