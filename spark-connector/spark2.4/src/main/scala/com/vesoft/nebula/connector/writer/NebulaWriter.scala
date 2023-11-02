/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.writer

import com.vesoft.nebula.spark.common.{DataTypeEnum, NebulaOptions}
import com.vesoft.nebula.spark.common.nebula.GraphProvider
import org.slf4j.LoggerFactory
import org.spark_project.guava.util.concurrent.RateLimiter

import java.util.concurrent.TimeUnit
import scala.collection.mutable.ListBuffer

class NebulaWriter(nebulaOptions: NebulaOptions) extends Serializable {
  private val LOG = LoggerFactory.getLogger(this.getClass)

  val failedExecs: ListBuffer[String] = new ListBuffer[String]

  val graphProvider = new GraphProvider(
    nebulaOptions.getGraphAddress,
    nebulaOptions.user,
    nebulaOptions.passwd,
    nebulaOptions.timeout,
    nebulaOptions.executionRetry
  )

  def submit(exec: String): Unit = {
    @transient val rateLimiter = RateLimiter.create(nebulaOptions.rateLimit)
    if (rateLimiter.tryAcquire(Long.MaxValue, TimeUnit.SECONDS)) {
      val result = graphProvider.submit(exec)
      if (!result.isSucceeded) {
        failedExecs.append(exec)
        LOG.error(s"failed to write ${exec}\n error message: ${result.getGqlStatus}")
      } else {
        LOG.info(s"batch write succeed.")
        LOG.debug(s"batch write succeed: ${exec}")
      }
    }
  }
}
