/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.example

import com.vesoft.nebula.connector.NebulaDataFrameReader
import com.vesoft.nebula.spark.common.{NebulaConnectionConfig, ReadNebulaConfig}
import org.apache.spark.sql.SparkSession

class NebulaSparkReaderExample {
  def main(args: Array[String]): Unit = {
    val spark = SparkSession
      .builder()
      .master("local")
      .getOrCreate()

    readNode(spark)
    readEdge(spark)

    spark.close()
  }

  private def getNebulaConnectionConfig: NebulaConnectionConfig = {
    NebulaConnectionConfig
      .builder()
      .withGraphAddress("192.168.8.6:3713")
      .withUser("root")
      .withPasswd("nebula")
      .build()
  }


  /**
   * for this example, you can config the read config to read node
   */
  private def readNode(spark: SparkSession): Unit = {
    val df = spark.read.json("spark-connector/example/src/main/resources/vertex")
    df.show()

    val nebulaNodeReadConfig: ReadNebulaConfig = ReadNebulaConfig
      .builder()
      .withGraphName("nba")
      .withTypeName("node_type_player")
      .withReturnCols(List("name"))
      .withBatchSize(10)
      .withPartitionNum(1)
      .build()
      spark.read.nebula(getNebulaConnectionConfig, nebulaNodeReadConfig).loadNode()
  }

  /**
   * for this example, you can config the read config to read edge
   */
  private def readEdge(spark: SparkSession): Unit = {
    val nebulaReadEdgeConfig: ReadNebulaConfig = ReadNebulaConfig
      .builder()
      .withGraphName("test")
      .withTypeName("edge_follow")
      .withReturnCols(List("degree"))
      .withBatchSize(1000)
      .withPartitionNum(1)
      .build()
    spark.read.nebula(getNebulaConnectionConfig, nebulaReadEdgeConfig).loadEdge()
  }
}
