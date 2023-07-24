/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.core

import com.vesoft.nebula.common.configuration.{EdgeConfig, NodeConfig, SchemaConfig}
import org.apache.spark.sql.catalyst.expressions.GenericRowWithSchema
import org.apache.spark.sql.types.StructType
import org.apache.spark.sql.{DataFrame, Row, SparkSession}
import org.apache.spark.sql.types._
import org.apache.spark.util.LongAccumulator
import org.scalatest.PrivateMethodTester
import org.scalatest.funsuite.AnyFunSuite

import scala.collection.mutable.ListBuffer
import scala.reflect.runtime.{universe => ru}

class ProcessorSuite extends AnyFunSuite {
  val schemaConfigs: ListBuffer[SchemaConfig] = new ListBuffer[SchemaConfig]

  var processor: Processor = null
  val spark = SparkSession.builder().master("local").getOrCreate()
  val dirtyRecords: LongAccumulator = spark.sparkContext.longAccumulator(s"dirtyRecordSize")

  test("test filterDirtyData") {
    if (processor == null) {
      prepare()
    }

    val clearDataFrame = processor.filterDirtyData()
    assert(clearDataFrame.count() == 8)
  }

  test("test filterFun") {
    if (processor == null) {
      prepare()
    }

    val notNullFields = List("srcId", "dstId")

    val schema = StructType(
      List(
        StructField("srcId", StringType, nullable = true),
        StructField("dstId", StringType, nullable = true),
        StructField("name", StringType, nullable = true),
        StructField("age", StringType, nullable = true),
        StructField("gender", StringType, nullable = true),
        StructField("playerName", StringType, nullable = true),
        StructField("degree", StringType, nullable = true),
        StructField("time", StringType, nullable = true),
        StructField("type", StringType, nullable = true)
      ))

    val row1:Row = new GenericRowWithSchema(List("Tom", "player1", "Tom", "10", "男", "player1", "20", "2023-01-01", "friend").toArray, schema)
    assert(processor.filterFun(notNullFields, row1))

    val row2:Row = new GenericRowWithSchema(List(null, "player2", "Bob", "11", "男", "player2", "21", "2023-01-01", "friend").toArray, schema)
    assert(!processor.filterFun(notNullFields, row2))

    val row3:Row = new GenericRowWithSchema(List("Jina", null, "Jina", "12", "女", "player3", "22", "2023-01-01", "friend").toArray, schema)
    assert(!processor.filterFun(notNullFields, row3))
  }

  test("test getIdFields") {
    if (processor == null) {
      prepare()
    }
    val fields = processor.getIdFields
    assert(fields.size == 2)
    assert(fields.contains("srcId"))
    assert(fields.contains("dstId"))
  }

  test("test getNonEmptyPropertyName") {}

  private def prepare(): Unit = {
    mockSchemas()
    val data: DataFrame = mockDataFrame(spark)
    processor = new Processor(data, schemaConfigs.toList, dirtyRecords)
  }

  private def mockSchemas(): Unit = {
    schemaConfigs.appendAll(mockNodeSchemas())
    schemaConfigs.append(mockEdgeSchema())
  }

  private def mockNodeSchemas(): List[NodeConfig] = {
    val batchSize                           = 10
    val partition                           = 10
    val nodeConfigs: ListBuffer[NodeConfig] = new ListBuffer[NodeConfig]

    val personName         = "person"
    val personSourceFields = List("name", "age", "gender")
    val personNebulaFields = List("name", "age", "gender")
    val personVidField     = "srcId"
    val nodeConfig1 = NodeConfig(personName,
                                 personSourceFields,
                                 personNebulaFields,
                                 batchSize,
                                 partition,
                                 personVidField)

    val playerName         = "player"
    val playerSourceFields = List("playerName")
    val playerNebulaFields = List("playerName")
    val playerVidField     = "dstId"
    val nodeConfig2 = NodeConfig(playerName,
                                 playerSourceFields,
                                 playerNebulaFields,
                                 batchSize,
                                 partition,
                                 playerVidField)

    nodeConfigs.append(nodeConfig1)
    nodeConfigs.append(nodeConfig2)
    nodeConfigs.toList
  }

  private def mockEdgeSchema(): EdgeConfig = {
    val edgeName     = "friend"
    val sourceFields = List("degree", "time", "type")
    val nebulaFields = List("degree", "time", "type")
    val batchSize    = 10
    val partition    = 10
    val srcField     = "srcId"
    val dstField     = "dstId"
    EdgeConfig(edgeName, sourceFields, nebulaFields, batchSize, partition, srcField, dstField)
  }

  private def mockDataFrame(spark: SparkSession): DataFrame = {
    val schema = StructType(
      List(
        StructField("srcId", StringType, nullable = true),
        StructField("dstId", StringType, nullable = true),
        StructField("name", StringType, nullable = true),
        StructField("age", StringType, nullable = true),
        StructField("gender", StringType, nullable = true),
        StructField("playerName", StringType, nullable = true),
        StructField("degree", StringType, nullable = true),
        StructField("time", StringType, nullable = true),
        StructField("type", StringType, nullable = true)
      ))

    val rdd = spark.sparkContext.parallelize(
      Seq(
        Row("Tom", "player1", "Tom", "10", "男", "player1", "20", "2023-01-01", "friend"),
        Row(null, "player2", "Bob", "11", "男", "player2", "21", "2023-01-01", "friend"),
        Row("Jina", null, "Jina", "12", "女", "player3", "22", "2023-01-01", "friend"),
        Row("Nic", "player4", null, "13", "男", "player4", "23", "2023-01-01", "friend"),
        Row("Tim", "player5", "Tim", null, "男", "player5", "24", "2023-01-01", "friend"),
        Row("Tina", "player6", "Tina", "15", null, "player6", "25", "2023-01-01", "friend"),
        Row("Hope", "player7", "Hope", "16", "男", null, "26", "2023-01-01", "friend"),
        Row("Lura", "player8", "Hope", "17", "女", "player8", null, "2023-01-01", "friend"),
        Row("Nura", "player9", "Nura", "18", "男", "player9", "26", null, "friend"),
        Row("Dolla", "player10", "Dolla", "19", "女", "player10", "26", "2023-01-01", null),
        Row(null, null, null, null, null, null, null, null, null)
      ))
    spark.createDataFrame(rdd, schema)
  }

}
