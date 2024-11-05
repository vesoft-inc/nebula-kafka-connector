
package com.vesoft.nebula.spark.common.nebula

import com.vesoft.nebula.driver.graph.data.{ResultSet, ValueWrapper}
import com.vesoft.nebula.driver.graph.net.NebulaClient
import com.vesoft.nebula.driver.graph.scan.{ScanEdgeResultIterator, ScanNodeResultIterator}
import com.vesoft.nebula.spark.common.NebulaUtils
import org.slf4j.LoggerFactory

import java.util
import java.util.List
import scala.collection.mutable
import scala.collection.mutable.ListBuffer

/**
 * GraphProvider for Nebula Graph Service
 */
class GraphProvider(addresses: String,
                    user: String,
                    authOptions: java.util.HashMap[String, Object],
                    timeout: Int,
                    schema: String,
                    zonedDatetimeFormat: String,
                    localDatetimeFormat: String,
                    zonedTimeFormat: String,
                    localTimeFormat: String)
  extends AutoCloseable
    with Serializable {
  @transient private[this] lazy val LOG = LoggerFactory.getLogger(this.getClass)

  @transient val client: NebulaClient = NebulaClient
    .builder(addresses, user)
    .withAuthOptions(authOptions)
    .withRequestTimeoutMills(timeout * 1000)
    .build()
  if (schema != null) {
    val res = client.execute("SESSION SET SCHEMA `" + schema + "`")
    if (!res.isSucceeded) {
      throw new IllegalArgumentException(s"SESSION SET SCHEMA failed, ${res.getErrorMessage}")
    }
  }

  if (zonedDatetimeFormat != null) {
    val res = client.execute("SESSION SET zoned_datetime_format=\"" + zonedDatetimeFormat + "\"")
    if (!res.isSucceeded) {
      throw new IllegalArgumentException(s"SESSION SET zoneddatetime format failed, ${res.getErrorMessage}")
    }
  }

  if (localDatetimeFormat != null) {
    val res = client.execute("SESSION SET local_datetime_format=\"" + localDatetimeFormat + "\"")
    if (!res.isSucceeded) {
      throw new IllegalArgumentException(s"SESSION SET localdatetime format failed, ${res.getErrorMessage}")
    }
  }

  if (zonedTimeFormat != null) {
    val res = client.execute("SESSION SET zoned_time_format=\"" + zonedTimeFormat + "\"")
    if (!res.isSucceeded) {
      throw new IllegalArgumentException(s"SESSION SET zonedtime format failed, ${res.getErrorMessage}")
    }
  }

  if (localTimeFormat != null) {
    val res = client.execute("SESSION SET local_time_format=\"" + localTimeFormat + "\"")
    if (!res.isSucceeded) {
      throw new IllegalArgumentException(s"SESSION SET localtime format failed, ${res.getErrorMessage}")
    }
  }

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

  /**
   * scan node type
   *
   * @param graphName graph name
   * @param nodeType  node type name
   * @param part      NebulaGraph partition id
   * @param batchSize batchSize for each scan request
   * @return {@link ScanNodeResultIterator}
   */
  def scanNode(graphName: String, nodeType: String, part: Int, batchSize: Int): ScanNodeResultIterator = {
    client.scanNode(graphName, nodeType, null, part, batchSize)
  }


  def scanNode(graphName: String, nodeType: String, returnCols: util.List[String], part: Int, batchSize: Int): ScanNodeResultIterator = {
    client.scanNode(graphName, nodeType, returnCols, part, batchSize)
  }

  /**
   * scan edge type
   *
   * @param graphName graph name
   * @param edgeType  edge type name
   * @param part      NebulaGraph partition id
   * @param batchSize batchSize for each scan request
   * @return {@link ScanEdgeResultIterator}
   */
  def scanEdge(graphName: String, edgeType: String, part: Int, batchSize: Int): ScanEdgeResultIterator = {
    client.scanEdge(graphName, edgeType, null, part, batchSize)
  }


  def scanEdge(graphName: String, edgeType: String, returnCols: util.List[String], part: Int, batchSize: Int): ScanEdgeResultIterator = {
    client.scanEdge(graphName, edgeType, returnCols, part, batchSize)
  }

  /**
   * get all part list for NebulaGraph
   */
  def getAllParts: List[Integer] = {
    val showPartitions: String    = "CALL show_partitions() RETURN *"
    var resultSet     : ResultSet = null
    try resultSet = client.execute(showPartitions)
    catch {
      case e: Exception =>
        LOG.error("get all partitions error", e)
        throw new RuntimeException("get all partitions error", e)
    }
    if (!resultSet.isSucceeded || resultSet.isEmpty) {
      LOG.error("get all partitions failed for {}", resultSet.getErrorMessage)
      throw new RuntimeException("get all partitions failed for " + resultSet.getErrorMessage)
    }
    val partitions: util.List[Integer] = new util.ArrayList[Integer]
    while (resultSet.hasNext) {
      partitions.add(resultSet.next().get("partition_id").asInt())
    }
    partitions
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

    val escapedNodeType = NebulaUtils.escapeUtil(nodeType)
    val descNodeType    = s"DESCRIBE NODE TYPE `$escapedNodeType` OF `$graphType`"
    val result          = client.execute(descNodeType)
    if (!result.isSucceeded || result.isEmpty) {
      LOG.error(s"get 'describe' of $nodeType failed for ${result.getErrorMessage}")
      throw new IllegalArgumentException(s"node type $escapedNodeType does not exist in $graphName.")
    }

    val pkNames = new ListBuffer[String]

    while (result.hasNext) {
      val record = result.next();
      schema += (record.get("property_name").asString() -> record.get("data_type").asString())
      if (!record.get("primary_key").isNull && record.get("primary_key").asString().equals("Y")) {
        pkNames.append(record.get("property_name").asString())
      }
    }

    if (pkNames.isEmpty) {
      LOG.error(s"node type $nodeType has no primary key.")
      throw new RuntimeException(s"node type $nodeType has no primary key")
    }
    NodeDesc(nodeType, pkNames.toList, schema.toMap)
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

    val escapedEdgeType = NebulaUtils.escapeUtil(edgeType)
    val descEdgeType    =
      s"call describe_graph_type('$graphType') filter type_name='$escapedEdgeType' return type_pattern next call describe_edge_type('$graphType', '$escapedEdgeType') return *"

    val result = client.execute(descEdgeType)
    if (!result.isSucceeded || result.isEmpty) {
      LOG.error(s"get 'describe' of $edgeType failed for ${result.getErrorMessage}")
      throw new IllegalArgumentException(s"edge type $edgeType does not exist in $graphName.")
    }

    var edgeTypePattern: String = null
    while (result.hasNext) {
      val record = result.next()
      if (edgeTypePattern == null) {
        edgeTypePattern = record.get("type_pattern").asString()
      }
      schema += (record.get("property_name").asString() -> record.get("data_type").asString())
    }

    var srcNodeType: String       = null
    var dstNodeType: String       = null
    // regularly match two types of edge:()-[]->() or ()~[]~() to get the srcNodeType and dstNodeType.
    val edgeDirectionPattern      = """\((.*?)\)-\[.*?\]->\((.*?)\)"""
    val edgeUnDirectionPattern    = """\((.*?)\)~\[.*?\]~\((.*?)\)"""
    val regexWithEdgeDirection    = edgeDirectionPattern.r
    val regexWithoutEdgeDirection = edgeUnDirectionPattern.r
    if (edgeTypePattern.matches(edgeDirectionPattern)) {
      edgeTypePattern match {
        case regexWithEdgeDirection(start, end) =>
          srcNodeType = start
          dstNodeType = end
      }
    } else if (edgeTypePattern.matches(edgeUnDirectionPattern)) {
      edgeTypePattern match {
        case regexWithoutEdgeDirection(start, end) =>
          srcNodeType = start
          dstNodeType = end
      }
    } else {
      throw new RuntimeException("can not parse the edge type pattern.")
    }

    val srcNodeDesc       = getNodeDesc(graphName, srcNodeType)
    val dstNodeDesc       = getNodeDesc(graphName, dstNodeType)
    val srcNodePkDataType = srcNodeDesc.properties.filterKeys(srcNodeDesc.nodePkNames.contains)
    val dstNodeIdDataType = dstNodeDesc.properties.filterKeys(dstNodeDesc.nodePkNames.contains)

    EdgeDesc(edgeType,
             srcNodeType,
             srcNodeDesc.nodePkNames,
             srcNodePkDataType,
             dstNodeType,
             dstNodeDesc.nodePkNames,
             dstNodeIdDataType,
             schema.toMap)
  }

  private def getGraphType(graphName: String): String = {
    val escapedGraphName = NebulaUtils.escapeUtil(graphName)
    val resultSet        = client.execute(s"DESCRIBE GRAPH `$escapedGraphName`")
    val graphType        = if (resultSet.isSucceeded && !resultSet.isEmpty) {
      resultSet.next().values().get(1).asString
    } else {
      throw new IllegalArgumentException(s"graphName $graphName does not exist.")
    }
    graphType
  }
}

