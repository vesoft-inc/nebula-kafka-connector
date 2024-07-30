/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.writer

import com.vesoft.nebula.spark.common.nebula.VidType
import com.vesoft.nebula.spark.common.{NebulaEdge, NebulaEdges, NebulaNode, NebulaNodes}
import com.vesoft.nebula.spark.common.writer.NebulaExecutor
import org.apache.log4j.BasicConfigurator
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.catalyst.expressions.GenericInternalRow
import org.apache.spark.sql.types.{BooleanType, DataTypes, LongType, StringType, StructField, StructType}
import org.scalatest.BeforeAndAfterAll
import org.scalatest.funsuite.AnyFunSuite

import scala.collection.mutable.ListBuffer

class NebulaExecutorSuite extends AnyFunSuite with BeforeAndAfterAll {
  val graphName = "nba"

  BasicConfigurator.configure()
  var schema: StructType = _
  var row: InternalRow = _

  override def beforeAll(): Unit = {
    val fields = new ListBuffer[StructField]
    fields.append(DataTypes.createStructField("col1", StringType, false))
    fields.append(DataTypes.createStructField("col2", BooleanType, false))
    fields.append(DataTypes.createStructField("col3", LongType, false))
    schema = new StructType(fields.toArray)

    val values = new ListBuffer[Any]
    values.append("aaa")
    values.append(true)
    values.append(1L)
    row = new GenericInternalRow(values.toArray)
  }

  override def afterAll(): Unit = super.afterAll()

  test("test extraID") {
    // test string primary key
    var index: Int = 0
    val isPkStringType: Boolean = true
    val stringId = NebulaExecutor.extraID(schema, row, index, isPkStringType)
    assert("\"aaa\"".equals(stringId))

    // test int primary key
    index = 2
    val hashId = NebulaExecutor.extraID(schema, row, index, false)
    assert("1".equals(hashId))
  }

  test("test extraID with null primary key value") {
    val values = new ListBuffer[Any]
    values.append(null)
    values.append(true)
    values.append(1L)
    val rowTest: InternalRow = new GenericInternalRow(values.toArray)
    val index: Int = 0
    val isPkStringType: Boolean = true
    assert(NebulaExecutor.extraID(schema, rowTest, index, isPkStringType) == null)
  }


  test("test src pk & dst pk all as prop for assignEdgeValues") {
    val fieldTypeMap: Map[String, String] =
      Map("col1" -> "STRING", "col2" -> "STRING", "col3" -> "STRING")

    val prop = NebulaExecutor.assignEdgeValues(schema, row, 0, 1, true, true, fieldTypeMap)
    assert(prop.size == 3)
  }

  test("test src pk & dst all not as prop for assignEdgeValues") {
    val fieldTypeMap: Map[String, String] =
      Map("col1" -> "STRING", "col2" -> "STRING", "col3" -> "STRING")

    val prop =
      NebulaExecutor.assignEdgeValues(schema, row, 0, 1, false, false, fieldTypeMap)
    assert(prop.size == 1)
  }

  test("test toExecuteSentence for node") {
    val nodes: ListBuffer[NebulaNode] = new ListBuffer[NebulaNode]
    val nodeType = "person"

    val props1: Map[String, String] = Map(
      "col_string" -> "\"Tom\"",
      "col_fixed_string" -> "\"Tom\"",
      "col_bool" -> "true",
      "col_int" -> "10",
      "col_int64" -> "100",
      "col_double" -> "1.0",
      "col_date" -> "date(\"2021-11-12\")"
    )
    val props2: Map[String, String] =
      Map(
        "col_string" -> "\"Bob\"",
        "col_fixed_string" -> "\"Bob\"",
        "col_bool" -> "false",
        "col_int" -> "20",
        "col_int64" -> "200",
        "col_double" -> "2.0",
        "col_date" -> "date(\"2021-05-01\")"
      )
    nodes.append(NebulaNode(props1))
    nodes.append(NebulaNode(props2))
    val fieldTypeMap: Map[String, String] =
      Map("col_string" -> "STRING", "col_fixed_string" -> "STRING", "col_bool" -> "BOOL", "col_int" -> "INT32", "col_int64" -> "INT64", "col_double" -> "DOUBLE", "col_date" -> "DATE")
    val nebulaNodes = NebulaNodes(nodeType, nodes.toList, "col_string", fieldTypeMap)
    // test insert node
    var nodeStatement =
      NebulaExecutor.toExecuteSentence(graphName, nebulaNodes, "")
    var expectStatement =
      s"""
         |TABLE t {col_string,col_fixed_string,col_bool,col_int,col_int64,col_double,col_date} =
         |(\"Tom\",\"Tom\",true,10,100,1.0,date(\"2021-11-12\")),(\"Bob\",\"Bob\",false,20,200,2.0,date(\"2021-05-01\"))
         |USE `$graphName`
         |FOR r IN t
         |INSERT  (@`$nodeType`{`col_string`:CAST(r.col_string AS STRING),`col_fixed_string`:CAST(r.col_fixed_string AS STRING),`col_bool`:CAST(r.col_bool AS BOOL),`col_int`:CAST(r.col_int AS INT32),`col_int64`:CAST(r.col_int64 AS INT64),`col_double`:CAST(r.col_double AS DOUBLE),`col_date`:CAST(r.col_date AS DATE)})
         |""".stripMargin
    assert(expectStatement.toCharArray.sorted.mkString("").equals(nodeStatement.toCharArray.sorted.mkString("")))

    // test insert node with replace
    nodeStatement =
      NebulaExecutor.toExecuteSentence(graphName, nebulaNodes, "OR REPLACE")
    expectStatement =
      s"""
         |TABLE t {col_string,col_fixed_string,col_bool,col_int,col_int64,col_double,col_date} =
         |(\"Tom\",\"Tom\",true,10,100,1.0,date(\"2021-11-12\")),(\"Bob\",\"Bob\",false,20,200,2.0,date(\"2021-05-01\"))
         |USE `$graphName`
         |FOR r IN t
         |INSERT OR REPLACE (@`$nodeType`{`col_string`:CAST(r.col_string AS STRING),`col_fixed_string`:CAST(r.col_fixed_string AS STRING),`col_bool`:CAST(r.col_bool AS BOOL),`col_int`:CAST(r.col_int AS INT32),`col_int64`:CAST(r.col_int64 AS INT64),`col_double`:CAST(r.col_double AS DOUBLE),`col_date`:CAST(r.col_date AS DATE)})
         |""".stripMargin
    assert(expectStatement.toCharArray.sorted.mkString("").equals(nodeStatement.toCharArray.sorted.mkString("")))

    // test insert node with ignore
    nodeStatement =
      NebulaExecutor.toExecuteSentence(graphName, nebulaNodes, "OR IGNORE")
    expectStatement =
      s"""
         |TABLE t {col_string,col_fixed_string,col_bool,col_int,col_int64,col_double,col_date} =
         |(\"Tom\",\"Tom\",true,10,100,1.0,date(\"2021-11-12\")),(\"Bob\",\"Bob\",false,20,200,2.0,date(\"2021-05-01\"))
         |USE `$graphName`
         |FOR r IN t
         |INSERT OR IGNORE (@`$nodeType`{`col_string`:CAST(r.col_string AS STRING),`col_fixed_string`:CAST(r.col_fixed_string AS STRING),`col_bool`:CAST(r.col_bool AS BOOL),`col_int`:CAST(r.col_int AS INT32),`col_int64`:CAST(r.col_int64 AS INT64),`col_double`:CAST(r.col_double AS DOUBLE),`col_date`:CAST(r.col_date AS DATE)})
         |""".stripMargin
    assert(expectStatement.toCharArray.sorted.mkString("").equals(nodeStatement.toCharArray.sorted.mkString("")))

  }

  test("test toExecuteSentence for edge") {
    val edges: ListBuffer[NebulaEdge] = new ListBuffer[NebulaEdge]
    val edgeType = "friend"
    val props1: Map[String, String] = Map(
      "col_string" -> "\"Tom\"",
      "col_fixed_string" -> "\"Tom\"",
      "col_bool" -> "true",
      "col_int" -> "10",
      "col_int64" -> "100",
      "col_double" -> "1.0",
      "col_date" -> "date(\"2021-11-12\")"
    )
    val props2: Map[String, String] =
      Map(
        "col_string" -> "\"Bob\"",
        "col_fixed_string" -> "\"Bob\"",
        "col_bool" -> "false",
        "col_int" -> "20",
        "col_int64" -> "200",
        "col_double" -> "2.0",
        "col_date" -> "date(\"2021-05-01\")"
      )
    edges.append(NebulaEdge("\"vid1\"", "\"vid2\"", props1))
    edges.append(NebulaEdge("\"vid2\"", "\"vid1\"", props2))

    val fieldTypeMap: Map[String, String] =
      Map("col_string" -> "STRING", "col_fixed_string" -> "STRING", "col_bool" -> "BOOL", "col_int" -> "INT32", "col_int64" -> "INT64", "col_double" -> "DOUBLE", "col_date" -> "DATE")

    val nebulaEdges: NebulaEdges = NebulaEdges(edgeType, "person", "id", VidType.STRING, "id1", "person", "id", VidType.STRING, "id2", edges.toList, fieldTypeMap)
    val edgeStatement = NebulaExecutor.toExecuteSentence(graphName, nebulaEdges, "")

    val expectStatement =
      s"""
         |TABLE t {id1,id2,col_string,col_fixed_string,col_bool,col_int,col_int64,col_double,col_date} =
         |(\"vid1\",\"vid2\",\"Tom\",\"Tom\",true,10,100,1.0,date(\"2021-11-12\")),(\"vid2\",\"vid1\",\"Bob\",\"Bob\",false,20,200,2.0,date(\"2021-05-01\"))
         |USE `$graphName`
         |FOR r IN t
         |MATCH (src@`person`{`id`:CAST(r.id1 AS STRING)}),(dst@`person`{`id`:CAST(r.id2 AS STRING)})
         |INSERT  (src)-[e@`friend`{`col_string`:CAST(r.col_string AS STRING),`col_fixed_string`:CAST(r.col_fixed_string AS STRING),`col_bool`:CAST(r.col_bool AS BOOL),`col_int`:CAST(r.col_int AS INT32),`col_int64`:CAST(r.col_int64 AS INT64),`col_double`:CAST(r.col_double AS DOUBLE),`col_date`:CAST(r.col_date AS DATE)}]->(dst)
         |""".stripMargin

    assert(expectStatement.toCharArray.sorted.mkString("").equals(edgeStatement.toCharArray.sorted.mkString("")))
  }


  test("test toDeleteSentence for node") {
    val nodes: ListBuffer[NebulaNode] = new ListBuffer[NebulaNode]
    val nodeType = "person"

    val props1: Map[String, String] = Map(
      "col_string" -> "\"Tom\"",
      "col_fixed_string" -> "\"Tom\"",
      "col_bool" -> "true",
      "col_int" -> "10",
      "col_int64" -> "100",
      "col_double" -> "1.0",
      "col_date" -> "date(\"2021-11-12\")"
    )
    val props2: Map[String, String] =
      Map(
        "col_string" -> "\"Bob\"",
        "col_fixed_string" -> "\"Bob\"",
        "col_bool" -> "false",
        "col_int" -> "20",
        "col_int64" -> "200",
        "col_double" -> "2.0",
        "col_date" -> "date(\"2021-05-01\")"
      )
    nodes.append(NebulaNode(props1))
    nodes.append(NebulaNode(props2))

    val fieldTypeMap: Map[String, String] =
      Map("col_string" -> "STRING", "col_fixed_string" -> "STRING", "col_bool" -> "BOOL", "col_int" -> "INT32", "col_int64" -> "INT64", "col_double" -> "DOUBLE", "col_date" -> "DATE")
    val nebulaNodes = NebulaNodes(nodeType, nodes.toList, "col_string", fieldTypeMap)
    val nodeStatement =
      NebulaExecutor.toDeleteSentence(graphName, nodeType, nebulaNodes)

    val expectStatement =
      s"""USE `$graphName` MATCH (a@$nodeType where a.col_string in [\"Bob\",\"Tom\"]) DETACH DELETE a""".stripMargin

    expectStatement.toCharArray.sorted.mkString("")
    assert(expectStatement.toCharArray.sorted.mkString("").equals(nodeStatement.toCharArray.sorted.mkString("")))
  }

  test("test toDeleteSentence for edge") {
    val edges: ListBuffer[NebulaEdge] = new ListBuffer[NebulaEdge]
    val edgeType = "friend"
    val props1: Map[String, String] = Map(
      "col_string" -> "\"Tom\"",
      "col_fixed_string" -> "\"Tom\"",
      "col_bool" -> "true",
      "col_int" -> "10",
      "col_int64" -> "100",
      "col_double" -> "1.0",
      "col_date" -> "date(\"2021-11-12\")"
    )
    edges.append(NebulaEdge("\"vid1\"", "\"vid2\"", props1))

    val fieldTypeMap: Map[String, String] =
      Map("col_string" -> "STRING", "col_fixed_string" -> "STRING", "col_bool" -> "BOOL", "col_int" -> "INT32", "col_int64" -> "INT64", "col_double" -> "DOUBLE", "col_date" -> "DATE")

    val nebulaEdges: NebulaEdges = NebulaEdges(edgeType, "person", "id", VidType.STRING, "col_string", "person", "id", VidType.STRING, "col_fixed_string", edges.toList, fieldTypeMap)

    val edgeStatement = NebulaExecutor.toDeleteSentence(graphName, edgeType, nebulaEdges)

    val expectStatement =
      s"""USE `$graphName` MATCH(a@person{id:\"vid1\"})-[e@friend]-(b@person{id:\"vid2\"}) DELETE e""".stripMargin
    assert(expectStatement.toCharArray.sorted.mkString("").equals(edgeStatement.toCharArray.sorted.mkString("")))
  }
}
