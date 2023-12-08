/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common.nebula

import com.vesoft.nebula.client.graph.data.{ResultSet, ValueWrapper}
import com.vesoft.nebula.client.graph.net.NebulaClient
import org.slf4j.LoggerFactory

import scala.collection.JavaConverters.{asScalaBufferConverter, mapAsScalaMapConverter}
import scala.collection.{breakOut, mutable}

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
    * @param statement insert node/edge statement
    * @return execute result
    */
  def submit(statement: String): ResultSet =
    client.execute(statement)

  def getIdType(graphName: String, nodeType: String): VidType.Value = {
    val nodeDesc = getNodeDesc(graphName, nodeType)
    nodeDesc.nodePkDataType
  }

  /**
    * get node's schema info
    *
    * @param graphName
    * @param nodeType
    * @return {@link NodeDesc}
    */
  def getNodeDesc(graphName: String, nodeType: String): NodeDesc = {
    val schema: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    val graphType                               = getGraphType(graphName)

    val descNodeType = s"DESCRIBE NODE TYPE $nodeType OF $graphType"
    val result       = client.execute(descNodeType)
    if (!result.isSucceeded || result.isEmpty) {
      LOG.error(s"get 'describe' of $nodeType failed for ${result.getGqlStatus}")
      throw new IllegalArgumentException(s"node type $nodeType does not exist in $graphName.")
    }
    val properties = result.getRows.get(0).get("properties").asMap()
    for ((k, v) <- properties.asScala) schema += (k -> v.asString())

    // for now, the pk is one property, composite pk is not support yet.
    val pks: List[ValueWrapper] = result.getRows.get(0).get("primary_keys").asList().asScala.toList
    if (pks.isEmpty) {
      LOG.error(s"node type $nodeType has no primary key.")
      throw new RuntimeException(s"node type $nodeType has no primary key")
    }
    val pk: String = pks.head.asString()

    var flag                      = true
    var idDataType: VidType.Value = null
    for (entry <- schema if flag) {
      if (entry._1.equals(pk)) {
        idDataType = VidType.withName(entry._2)
        flag = false
      }
    }
    if (idDataType == null) {
      throw new RuntimeException(s"can not get the pk $pk for $nodeType")
    }
    NodeDesc(nodeType, idDataType, schema.toMap)
  }

  /**
    * get edge description info
    *
    * @param graphName
    * @param edgeType
    * @return {@link EdgeDesc}
    */
  def getEdgeDesc(graphName: String, edgeType: String): EdgeDesc = {
    val schema: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    val graphType                               = getGraphType(graphName)

    val descEdgeType = s"DESCRIBE EDGE TYPE $edgeType OF $graphType"
    val result       = client.execute(descEdgeType)
    if (!result.isSucceeded || result.isEmpty) {
      LOG.error(s"get 'describe' of $edgeType failed for ${result.getGqlStatus}")
      throw new IllegalArgumentException(s"edge type $edgeType does not exist in $graphName.")
    }

    val types = result.getRows.get(0).get("types").asList()
    // the 'Types' in result for DESCRIBE EDGE TYPE should contain 'src node type', 'edge type', 'dst node type'
    if (types.size() < 3) {
      LOG.error(s"types size is less than 3 for edge type $edgeType")
      throw new RuntimeException(
        s"edge type $edgeType has unexpected 'Types', the types size is less than 3.")
    }
    val srcNodeType = types.get(0).asString()
    val dstNodeType = types.get(2).asString()

    val srcNodeIdDataType = getIdType(graphName, srcNodeType)
    val dstNodeIdDataType = getIdType(graphName, dstNodeType)

    val properties = result.getRows.get(0).get("properties").asMap()
    for ((k, v) <- properties.asScala) schema += (k -> v.asString())

    EdgeDesc(edgeType, srcNodeType, srcNodeIdDataType, dstNodeType, dstNodeIdDataType, schema.toMap)
  }

  private def getGraphType(graphName: String): String = {
    var resultSet = client.execute(s"DESCRIBE GRAPH $graphName")
    val graphType = if (resultSet.isSucceeded && !resultSet.isEmpty) {
      resultSet.getRows.get(0).values().get(1).asString
    } else {
      throw new IllegalArgumentException(s"graphName $graphName does not exist.")
    }
    graphType
  }
}

object VidType extends Enumeration {
  type Type = Value

  val STRING = Value("STRING")
  val INT    = Value("INT64")
}
