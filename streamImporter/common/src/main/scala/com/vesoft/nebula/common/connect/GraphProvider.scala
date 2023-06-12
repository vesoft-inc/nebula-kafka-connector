/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.connect

import com.vesoft.nebula.client.graph.data.ResultSet
import com.vesoft.nebula.client.graph.net.NebulaClient
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
    */
  def getBucketNum(): Int = {
    65535
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

}
