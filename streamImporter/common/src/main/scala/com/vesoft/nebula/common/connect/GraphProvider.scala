/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.connect

import com.vesoft.nebula.client.graph.data.ResultSet
import org.apache.log4j.Logger

class GraphProvider(addresses: String, connTimeout: Int, requestTimeout: Int)
    extends AutoCloseable
    with Serializable {
  private[this] lazy val LOG = Logger.getLogger(this.getClass)

  override def close(): Unit = {
    null
  }

  /**
    * execute query
    */
  def query(statement: String): ResultSet = {
    null
  }

  /**
    *
    */
  def getBucketNum(): Int = {
    65535
  }

}
