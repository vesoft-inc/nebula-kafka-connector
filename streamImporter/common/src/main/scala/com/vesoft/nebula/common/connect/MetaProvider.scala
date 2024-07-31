
package com.vesoft.nebula.common.connect

import org.apache.log4j.Logger

abstract class MetaProvider(addresses: String, connTimeout: Int, requestTimeout: Int)
    extends AutoCloseable
    with Serializable {
  private[this] lazy val LOG = Logger.getLogger(this.getClass)

  override def close(): Unit = {
    null
  }

  def getVidType(graphName: String): Unit = {}
  def getNodeSchema(graphName: String, node: String): Map[String, String] = {
    null
  }

  def getEdgeSchema(graphName: String, edge: String): Map[String, String] = {
    null
  }

}
