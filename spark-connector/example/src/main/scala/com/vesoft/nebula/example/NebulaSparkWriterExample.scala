/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.example

import com.sun.org.slf4j.internal.LoggerFactory
import com.vesoft.nebula.connector.NebulaDataFrameWriter
import com.vesoft.nebula.spark.common.{NebulaConnectionConfig, WriteMode, WriteNebulaEdgeConfig, WriteNebulaVertexConfig}
import org.apache.spark.SparkConf
import org.apache.spark.sql.SparkSession
import org.apache.spark.storage.StorageLevel

object NebulaSparkWriterExample {
  private val LOG = LoggerFactory.getLogger(this.getClass)

  def main(args: Array[String]): Unit = {

    val spark = SparkSession
      .builder()
      .master("local")
      .getOrCreate()

    writeVertex(spark)
    writeEdge(spark)

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
    * for this example, your nebula tag schema should have property names: name, age, born
    * if your withVidAsProp is true, then tag schema also should have property name: id
    */
  def writeVertex(spark: SparkSession): Unit = {
    val df = spark.read.json("spark-connector/example/src/main/resources/vertex")
    df.show()

    val nebulaWriteVertexConfig: WriteNebulaVertexConfig = WriteNebulaVertexConfig
      .builder()
      .withGraphName("nba")
      .withNodeType("node_type_player")
      .withPrimaryKeyField("id")
      .withWriteMode(WriteMode.INSERT)
      .withBatchSize(1)
      .build()
    df.write.nebula(getNebulaConnectionConfig, nebulaWriteVertexConfig).writeVertices()
  }

  /**
    * for this example, your nebula edge schema should have property names: descr, timp
    * if your withSrcAsProperty is true, then edge schema also should have property name: src
    * if your withDstAsProperty is true, then edge schema also should have property name: dst
    */
  def writeEdge(spark: SparkSession): Unit = {
    val df = spark.read.json("spark-connector/example/src/main/resources/vertex")
    df.show()
    df.persist(StorageLevel.MEMORY_AND_DISK_SER)

    val nebulaWriteEdgeConfig: WriteNebulaEdgeConfig = WriteNebulaEdgeConfig
      .builder()
      .withGraphName("test")
      .withEdge("friend")
      .withSrcPkField("src")
      .withDstPkField("dst")
      .withSrcPkAsProperty(false)
      .withDstPkAsProperty(false)
      .withBatchSize(1000)
      .build()
    df.write.nebula(getNebulaConnectionConfig, nebulaWriteEdgeConfig).writeEdges()
  }

}
