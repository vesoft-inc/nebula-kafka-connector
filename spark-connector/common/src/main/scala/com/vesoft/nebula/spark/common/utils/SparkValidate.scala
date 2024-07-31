
package com.vesoft.nebula.spark.common.utils

import org.apache.spark.sql.SparkSession

object SparkValidate {
  def validate(supportedVersions: String*): Unit = {
    val sparkVersion = SparkSession.getActiveSession.map(_.version).getOrElse("UNKNOWN")
    if (sparkVersion != "UNKNOWN" && !supportedVersions.exists(sparkVersion.matches)) {
      throw new RuntimeException(
        s"""Your current spark version ${sparkVersion} is not supported by the current NebulaGraph Spark Connector.
           | """.stripMargin)
    }
  }
}
