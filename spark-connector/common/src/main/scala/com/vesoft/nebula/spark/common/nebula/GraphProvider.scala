
package com.vesoft.nebula.spark.common.nebula

import com.vesoft.nebula.client.graph.data.{ResultSet, ValueWrapper}
import com.vesoft.nebula.client.graph.net.NebulaClient
import com.vesoft.nebula.client.graph.scan.{ScanEdgeResultIterator, ScanNodeResultIterator}
import org.slf4j.LoggerFactory


import java.util
import java.util.List
import scala.collection.mutable

/**
 * GraphProvider for Nebula Graph Service
 */
class GraphProvider(addresses: String,
                    user: String,
                    authOptions: java.util.HashMap[String, Object],
                    timeout: Int)
  extends AutoCloseable
    with Serializable {
  @transient private[this] lazy val LOG = LoggerFactory.getLogger(this.getClass)

  @transient val client: NebulaClient = NebulaClient
    .builder(addresses, user)
    .withAuthOptions(authOptions)
    .withRequestTimeoutMills(timeout * 1000)
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
    client.scanNode(graphName, nodeType, part, batchSize)
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
    client.scanEdge(graphName, edgeType, part, batchSize)
  }


  def scanEdge(graphName: String, edgeType: String, returnCols: util.List[String], part: Int, batchSize: Int): ScanEdgeResultIterator = {
    client.scanEdge(graphName, edgeType, returnCols, part, batchSize)
  }

  /**
   * get all part list for NebulaGraph
   */
  def getAllParts(graphName: String): List[Integer] = {
    val showPartitions: String = "CALL show_partitions() RETURN *"
    var resultSet: ResultSet = null
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
    val graphType = getGraphType(graphName)

    val descNodeType = s"DESCRIBE NODE TYPE $nodeType OF $graphType"
    val result = client.execute(descNodeType)
    if (!result.isSucceeded || result.isEmpty) {
      LOG.error(s"get 'describe' of $nodeType failed for ${result.getErrorMessage}")
      throw new IllegalArgumentException(s"node type $nodeType does not exist in $graphName.")
    }

    // for now, the pk is one property, composite pk is not support yet.
    var pk: String = null
    var pkDataType: String = null

    while (result.hasNext) {
      val record = result.next();
      schema += (record.get("property_name").asString() -> record.get("data_type").asString())
      if ("Y".equals(record.get("primary_key").asString())) {
        pk = record.get("property_name").asString()
        pkDataType = record.get("data_type").asString()
      }
    }

    if (pk == null) {
      LOG.error(s"node type $nodeType has no primary key.")
      throw new RuntimeException(s"node type $nodeType has no primary key")
    }
    val idDataType: VidType.Value = VidType.withName(pkDataType)

    NodeDesc(nodeType, pk, idDataType, schema.toMap)
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
    val graphType = getGraphType(graphName)

    val descEdgeType =
      s"call describe_graph_type('$graphType') filter type_name='$edgeType' return type_pattern next call describe_edge_type('$graphType', '$edgeType') return *"

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

    var srcNodeType: String = null
    var dstNodeType: String = null
    // regularly match two types of edge:()-[]->() or ()~[]~() to get the srcNodeType and dstNodeType.
    val edgeDirectionPattern = """\((.*?)\)-\[.*?\]->\((.*?)\)"""
    val edgeUnDirectionPattern = """\((.*?)\)~\[.*?\]~\((.*?)\)"""
    val regexWithEdgeDirection = edgeDirectionPattern.r
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

    val srcNodeIdDataType = getIdType(graphName, srcNodeType)
    val dstNodeIdDataType = getIdType(graphName, dstNodeType)

    val srcNodePkName = getNodeDesc(graphName, srcNodeType).nodePkName
    val dstNodePkName = getNodeDesc(graphName, dstNodeType).nodePkName


    EdgeDesc(edgeType, srcNodeType, srcNodePkName, srcNodeIdDataType, dstNodeType, dstNodePkName, dstNodeIdDataType, schema.toMap)
  }

  private def getGraphType(graphName: String): String = {
    var resultSet = client.execute(s"DESCRIBE GRAPH $graphName")
    val graphType = if (resultSet.isSucceeded && !resultSet.isEmpty) {
      resultSet.next().values().get(1).asString
    } else {
      throw new IllegalArgumentException(s"graphName $graphName does not exist.")
    }
    graphType
  }
}

object VidType extends Enumeration {
  type Type = Value

  val STRING = Value("STRING")
  val INT = Value("INT64")
}
