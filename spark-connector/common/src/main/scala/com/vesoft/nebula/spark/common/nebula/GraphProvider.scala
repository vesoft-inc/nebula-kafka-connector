/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common.nebula

import com.vesoft.nebula.client.graph.data.ResultSet
import com.vesoft.nebula.client.graph.net.NebulaClient
import org.slf4j.LoggerFactory

import scala.collection.JavaConverters.asScalaBufferConverter
import scala.collection.mutable

/**
  * GraphProvider for Nebula Graph Service
  */
class GraphProvider(addresses: String, user: String, password: String, timeout: Int, retryTime: Int)
    extends AutoCloseable
    with Serializable {
  @transient private[this] lazy val LOG = LoggerFactory.getLogger(this.getClass)
  @transient val client: NebulaClient = NebulaClient
    .builder(addresses, user, password)
    .setConnectTimeoutMills(timeout * 1000)
    .setRequestTimeoutMills(timeout * 1000)
    .setRetryTimes(retryTime)
    .setMaxSessionSize(20)
    .setMinSessionSize(20)
    .build()

  /**
    * close Nebula client
    */
  override def close(): Unit = {
    client.close()
  }

  /**
    * execute the statement
    *
    * @param statement insert tag/edge statement
    * @return execute result
    */
  def submit(statement: String): ResultSet =
    client.execute(statement)

  def getIdType(graphName: String, nodeType: String): VidType.Value = {
    val schema = getTagSchema(graphName, nodeType)
    for (entry <- schema) {
      if (entry._1.equals("id")) {
        return VidType.withName(entry._2)
      }
    }
    throw new IllegalArgumentException(s"graphName $graphName does not have NodeType $nodeType.")
  }

  def getIdsType(graphName: String, edgeType: String): (VidType.Value, VidType.Value) = {
    val (sourceNodeType, targetNodeType) = getNodesType(graphName, edgeType)
    (getIdType(graphName, sourceNodeType), getIdType(graphName, targetNodeType))
  }

  def getNodesType(graphName: String, edgeType: String): (String, String) = {
    var resultSet = client.execute(s"DESCRIBE GRAPH $graphName")
    val graphType = if (resultSet.isSucceeded && !resultSet.isEmpty) {
      resultSet.getRows.get(0).values().get(1).asString
    } else {
      throw new IllegalArgumentException(s"graphName $graphName does not exist.")
    }

    resultSet = client.execute(s"DESCRIBE GRAPH TYPE $graphType")
    var sourceNodeType: String = null
    var targetNodeType: String = null

    if (resultSet.isSucceeded) {
      val records = resultSet.getRows
      for (record: ResultSet.Record <- records.asScala) {
        if (record.get("Field").asString.toUpperCase.contains(edgeType.toUpperCase)) {
          val propertyString = record.get("Field").asString()
          val pattern        = """\((.*?)\)-\[.*?\]->\((.*?)\)""".r
          propertyString match {
            case pattern(start, end) =>
              sourceNodeType = start
              targetNodeType = end

            case _ => throw new IllegalArgumentException(s"edge type pattern parse failed.")
          }
          return (sourceNodeType, targetNodeType)
        }
      }
    }
    throw new IllegalArgumentException(s"graphName $graphName does not have EdgeType $edgeType.")
  }

  /**
    * get tag's schema info
    *
    * @param graphName
    * @param nodeType
    * @return Map, property name -> data type {@link PropertyType}
    */
  def getTagSchema(graphName: String, nodeType: String): Map[String, String] = {
    val schema: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    val resultSet                               = getGraphDesc(graphName)

    val records             = resultSet.getRows
    var existLabel: Boolean = false
    for (record: ResultSet.Record <- records.asScala) {
      if (record.get("Field").asString.equalsIgnoreCase(nodeType)) {
        existLabel = true
        val propertyString = record.get("Properties").asString()
        val properties     = propertyString.substring(1, propertyString.length - 1).split(",")
        for (prop <- properties) {
          val nameAndType = prop.trim.split(" ")
          schema += (nameAndType(0) -> nameAndType(1))
        }
      }
    }
    if (!existLabel) {
      throw new IllegalArgumentException(s"graphName $graphName does not have nodeType $nodeType")
    }
    schema.toMap
  }

  /**
    * get edge's schema info
    *
    * @param graphName
    * @param edgeType
    * @return Map, property name -> data type {@link PropertyType}
    */
  def getEdgeSchema(graphName: String, edgeType: String): Map[String, String] = {
    val schema: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    val resultSet                               = getGraphDesc(graphName)

    val records             = resultSet.getRows
    var existLabel: Boolean = false
    for (record: ResultSet.Record <- records.asScala) {
      if (record.get("Kind").asString().equals("Edge")) {
        val fullEdgeType         = record.get("Field").asString
        val pattern              = """\((.*?)\)-\[(.*?)\]->\((.*?)\)""".r
        var edgeTypeName: String = null
        fullEdgeType match {
          case pattern(start, edgeType, end) =>
            edgeTypeName = edgeType
          case _ => throw new IllegalArgumentException(s"edge type pattern parse failed.")
        }
        if (edgeTypeName.equalsIgnoreCase(edgeType)) {
          existLabel = true
          val propertyString = record.get("Properties").asString()
          val properties     = propertyString.substring(1, propertyString.length - 1).split(",")
          for (prop <- properties) {
            val nameAndType = prop.trim.split(" ")
            schema += (nameAndType(0) -> nameAndType(1))
          }
        }
      }
    }
    if (!existLabel) {
      throw new IllegalArgumentException(s"graphName $graphName does not have edgeType $edgeType")
    }
    schema.toMap
  }

  private def getGraphDesc(graphName: String): ResultSet = {
    var resultSet = client.execute(s"DESCRIBE GRAPH $graphName")
    val graphType = if (resultSet.isSucceeded && !resultSet.isEmpty) {
      resultSet.getRows.get(0).values().get(1).asString
    } else {
      throw new IllegalArgumentException(s"graphName $graphName does not exist.")
    }

    val queryStatement = s"DESCRIBE GRAPH TYPE $graphType"
    resultSet = client.execute(queryStatement)
    if (!resultSet.isSucceeded) {
      throw new RuntimeException(
        s"query error with `$queryStatement` for ${resultSet.getGqlStatus}")
    }

    resultSet

  }
}

object VidType extends Enumeration {
  type Type = Value

  val STRING = Value("STRING")
  val INT    = Value("INT64")
}
