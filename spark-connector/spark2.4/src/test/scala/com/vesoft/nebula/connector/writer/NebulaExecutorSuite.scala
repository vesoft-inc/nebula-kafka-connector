/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.writer

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

    val nebulaNodes = NebulaNodes(nodeType, nodes.toList)
    val nodeStatement =
      NebulaExecutor.toExecuteSentence(graphName, nodeType, nebulaNodes)

    val expectStatement =
      s"""USE $graphName INSERT NODE `$nodeType` ({`col_string`:\"Tom\",`col_fixed_string`:\"Tom\",`col_bool`:true,`col_int`:10,`col_int64`:100,`col_double`:1.0,`col_date`:date(\"2021-11-12\")}),({`col_string`:\"Bob\",`col_fixed_string`:\"Bob\",`col_bool`:false,`col_int`:20,`col_int64`:200,`col_double`:2.0,`col_date`:date(\"2021-05-01\")})""".stripMargin

    expectStatement.toCharArray.sorted.mkString("")
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
    edges.append(NebulaEdge("id", "\"vid1\"", "id", "\"vid2\"", props1))
    edges.append(NebulaEdge("id", "\"vid2\"", "id", "\"vid1\"", props2))

    val nebulaEdges: NebulaEdges = NebulaEdges(edgeType, edges.toList)
    val edgeStatement = NebulaExecutor.toExecuteSentence(graphName, edgeType, nebulaEdges)

    val expectStatement =
      s"""USE $graphName INSERT EDGE `$edgeType` ({`id`:\"vid1\"})-[{`col_string`:\"Tom\",`col_fixed_string`:\"Tom\",`col_bool`:true,`col_int`:10,`col_int64`:100,`col_double`:1.0,`col_date`:date(\"2021-11-12\")}]->({`id`:\"vid2\"}),({`id`:\"vid2\"})-[{`col_string`:\"Bob\",`col_fixed_string`:\"Bob\",`col_bool`:false,`col_int`:20,`col_int64`:200,`col_double`:2.0,`col_date`:date(\"2021-05-01\")}]->({`id`:\"vid1\"})""".stripMargin
    assert(expectStatement.toCharArray.sorted.mkString("").equals(edgeStatement.toCharArray.sorted.mkString("")))
  }
}
