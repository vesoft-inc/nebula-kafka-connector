
package com.vesoft.nebula.connector.writer

import com.vesoft.nebula.driver.graph.data.ResultSet
import com.vesoft.nebula.spark.common.NebulaOptions
import com.vesoft.nebula.spark.common.nebula.GraphProvider
import org.slf4j.LoggerFactory
import org.spark_project.guava.util.concurrent.RateLimiter

import java.util.concurrent.TimeUnit
import scala.collection.mutable.ListBuffer

class NebulaWriter(nebulaOptions: NebulaOptions) extends Serializable {
  private val LOG = LoggerFactory.getLogger(this.getClass)

  val failedExecs: ListBuffer[String] = new ListBuffer[String]

  val graphProvider = new GraphProvider(
    nebulaOptions.graphAddress,
    nebulaOptions.user,
    nebulaOptions.authOptions,
    nebulaOptions.timeout,
    nebulaOptions.schema,
    nebulaOptions.zonedDatetimeFormat,
    nebulaOptions.localDatetimeFormat,
    nebulaOptions.zonedTimeFormat,
    nebulaOptions.zonedTimeFormat
  )

  def submit(exec: String): ResultSet = {
    @transient val rateLimiter = RateLimiter.create(nebulaOptions.rateLimit)
    if (rateLimiter.tryAcquire(30, TimeUnit.MINUTES)) {
      val result = graphProvider.submit(exec)
      result
    } else {
      throw new RuntimeException("get rateLimiter timeout.")
    }
  }
}
