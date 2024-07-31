
package com.vesoft.nebula.common.connect

import com.vesoft.nebula.client.graph.data.ResultSet
import com.vesoft.nebula.client.graph.net.NebulaClient
import com.vesoft.nebula.common.configuration.IDType
import org.apache.log4j.Logger

import scala.collection.mutable

class GraphProvider(addresses: String,
                    user: String,
                    passwd: String,
                    connTimeout: Int,
                    requestTimeout: Int,
                    retryIntervalTime: Int)
    extends AutoCloseable
    with Serializable {
  private[this] lazy val LOG = Logger.getLogger(this.getClass)

  val graphClient = NebulaClient
    .builder(addresses, user, passwd)
    .setConnectTimeoutMills(connTimeout)
    .setRequestTimeoutMills(requestTimeout)
    .setMaxSessionSize(1)
    .setMinSessionSize(1)
    .setRetryTimes(3)
    .setIntervalTimeMills(retryIntervalTime)
    .setReconnect(true)
    .setBlockWhenExhausted(true)
    .setMaxWaitMills(1000)
    .setStrictlyServerHealthy(true)
    .build

  override def close(): Unit = {
    graphClient.close()
  }

  /**
    * execute query
    */
  def query(statement: String): ResultSet = {
    graphClient.execute(statement)
  }

  /**
    * get NebulaGraph server's bucket number
    * TODO(Anqi): implement the interface, for now MOCK the value to 65535
    */
  def getBucketNum(): Int = {
    65535
  }

  /**
    * get the graph id in NebulaGraph
    * TODO(Anqi): implement the interface, for now MOCK 0, depends on GQL.
    */
  def getGraphId(graphName: String): Int = {
    0
  }

  /**
    * get node schemas
    */
  def getNodeSchemas(graphName: String, nodetype: String): Map[String, String] = {
    val schema: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    // TODO query node schema
    schema.toMap
  }

  /**
    * get edge schemas
    */
  def getEdgeSchemas(graphName: String, edgetype: String): Map[String, String] = {
    val schema: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    // TODO query edge schema
    schema.toMap
  }

  /**
    * get Id data type
    * TODO(Anqi): implement the interface, for now MOCK data type to STRING
    * */
  def getNodeIdType(graphName: String, nodeType: String): IDType.Value = {
    val schema = getNodeSchemas(graphName, nodeType)
    val idType = schema("ID")
    idType.toUpperCase() match {
      case "STRING" => IDType.STRING
      case "INT"    => IDType.INT
      case _ =>
        throw new IllegalArgumentException(
          s"do not support the $idType ID data type for $nodeType, only support 'INT' and 'STRING'.")
    }
    IDType.STRING
  }

  /**
    * get Source Id and Target Id data type
    * TODO(Anqi): implement the interface, for now MOCK data type to STRING
    */
  def getEdgeNodesIdTypes(graphName: String, edgeType: String): (IDType.Value, IDType.Value) = {
    val schema         = getEdgeSchemas(graphName, edgeType)
    val sourceNodeType = "source node type"
    val targetNodeType = "target node type"
    val sourceIdType   = getNodeIdType(graphName, sourceNodeType)
    val targetIdType   = getNodeIdType(graphName, targetNodeType)
    (sourceIdType, targetIdType)
    // MOCK
    (IDType.STRING, IDType.STRING)
  }

  /**
    * get node type
    * TODO(Anqi): implement the interface, for now Mock as "Person". depends on GQL.
    *
    */
  def getNodeType(graphName: String, nodeType: String): String = {
    val schema = getNodeSchemas(graphName, nodeType)
    // MOCK
    "Person"
  }

  /**
    * get edge's source node type and target node type
    * TODO(Anqi): implement the interface, for now Mock as "Person" and "Player". depends on GQL.
    *
    */
  def getNodesType(graphName: String, edgeType: String): (String, String) = {
    val schema = getEdgeSchemas(graphName, edgeType)
    // MOCK
    ("Person", "Player")
  }

  def getInternalIds(nodeType: String, primaryKeys: List[String]): Map[String, String] = {
    val primaryKey2InternalId = Map[String, String]()
    primaryKey2InternalId
  }

  def generateInternalIds(nodeType: String, primaryKeys: List[String]): Map[String, String] = {
    val primaryKey2InternalId = Map[String, String]()
    primaryKey2InternalId
  }

  def getNodeTypeId(graphName:String, nodeName:String): Int = {
    0
  }

  def getEdgeTypeId(graphName:String, edgeName:String):Int={
    1
  }
}
