/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import org.apache.spark.sql.catalyst.util.CaseInsensitiveMap
import org.apache.spark.sql.types.{LongType, StructField, StructType}
import org.scalatest.funsuite.AnyFunSuite

class NebulaUtilsSuite extends AnyFunSuite {
  test("makeGetters") {
    val schema = StructType(
      List(
        StructField("col1", LongType, nullable = false),
        StructField("col2", LongType, nullable = true)
      ))
    assert(NebulaUtils.makeGetters(schema).length == 2)
  }

  test("isNumic") {
    assert(NebulaUtils.isNumic("123"))
    assert(NebulaUtils.isNumic("-123"))
    assert(!NebulaUtils.isNumic(""))
    assert(!NebulaUtils.isNumic("-"))
    assert(!NebulaUtils.isNumic("1.0"))
    assert(!NebulaUtils.isNumic("a123"))
    assert(!NebulaUtils.isNumic("123b"))
  }

  test("escapeUtil") {
    assert(NebulaUtils.escapeUtil("123").equals("123"))
    // a\bc -> a\\bc
    assert(NebulaUtils.escapeUtil("a\bc").equals("a\\bc"))
    // a\tbc -> a\\tbc
    assert(NebulaUtils.escapeUtil("a\tbc").equals("a\\tbc"))
    // a\nbc -> a\\nbc
    assert(NebulaUtils.escapeUtil("a\nbc").equals("a\\nbc"))
    // a\"bc -> a\\"bc
    assert(NebulaUtils.escapeUtil("a\"bc").equals("a\\\"bc"))
    // a\'bc -> a\\'bc
    assert(NebulaUtils.escapeUtil("a\'bc").equals("a\\'bc"))
    // a\rbc -> a\\rbc
    assert(NebulaUtils.escapeUtil("a\rbc").equals("a\\rbc"))
    // a\bbc -> a\\bbc
    assert(NebulaUtils.escapeUtil("a\bbc").equals("a\\bbc"))
  }

  test("getSchema for node with all properties"){
    var parameters: Map[String, String] = Map()
    parameters += (NebulaOptions.TYPE -> "NODE")
    parameters += (NebulaOptions.OPERATE_TYPE -> "read")
    parameters += (NebulaOptions.GRAPH_ADDRESS -> "192.168.8.6:3820")
    parameters += (NebulaOptions.GRAPH_NAME -> "nba")
    parameters += (NebulaOptions.USER_NAME -> "root")
    parameters += (NebulaOptions.PASSWD -> "Nebula123")
    parameters += (NebulaOptions.LABEL -> "node_type_player")
    parameters += (NebulaOptions.RETURN_COLS -> "$null")
    val nebulaOptions = new NebulaOptions(CaseInsensitiveMap(parameters))
    val schema = NebulaUtils.getSchema(nebulaOptions)
    assert(schema.fields.length == 5)
    assert(schema.fields.map(field => field.name).contains("id"))
    assert(schema.fields.map(field => field.name).contains("name"))
    assert(schema.fields.map(field => field.name).contains("rate"))
    assert(schema.fields.map(field => field.name).contains("score"))
    assert(schema.fields.map(field => field.name).contains("gender"))
  }

  test("getSchema for node with no properties") {
    var parameters: Map[String, String] = Map()
    parameters += (NebulaOptions.TYPE -> "NODE")
    parameters += (NebulaOptions.OPERATE_TYPE -> "read")
    parameters += (NebulaOptions.GRAPH_ADDRESS -> "192.168.8.6:3820")
    parameters += (NebulaOptions.GRAPH_NAME -> "nba")
    parameters += (NebulaOptions.USER_NAME -> "root")
    parameters += (NebulaOptions.PASSWD -> "Nebula123")
    parameters += (NebulaOptions.LABEL -> "node_type_player")
    parameters += (NebulaOptions.RETURN_COLS -> "")
    val nebulaOptions = new NebulaOptions(CaseInsensitiveMap(parameters))
    val schema = NebulaUtils.getSchema(nebulaOptions)
    assert(schema.fields.length == 1)
    assert(schema.fields.map(field => field.name).contains("id"))
  }

  test("getSchema for edge with all properties") {
    var parameters: Map[String, String] = Map()
    parameters += (NebulaOptions.TYPE -> "EDGE")
    parameters += (NebulaOptions.OPERATE_TYPE -> "read")
    parameters += (NebulaOptions.GRAPH_ADDRESS -> "192.168.8.6:3820")
    parameters += (NebulaOptions.GRAPH_NAME -> "nba")
    parameters += (NebulaOptions.USER_NAME -> "root")
    parameters += (NebulaOptions.PASSWD -> "Nebula123")
    parameters += (NebulaOptions.LABEL -> "edge_type_follow")
    parameters += (NebulaOptions.RETURN_COLS -> "$null")
    val nebulaOptions = new NebulaOptions(CaseInsensitiveMap(parameters))
    val schema = NebulaUtils.getSchema(nebulaOptions)
    assert(schema.fields.length == 4)
    assert(schema.fields.map(field => field.name).contains("id"))
    assert(schema.fields.map(field => field.name).contains("id"))
    assert(schema.fields.map(field => field.name).contains("followness"))
    assert(schema.fields.map(field => field.name).contains("likeness"))
  }

  test("getSchema for edge with no properties") {
    var parameters: Map[String, String] = Map()
    parameters += (NebulaOptions.TYPE -> "EDGE")
    parameters += (NebulaOptions.OPERATE_TYPE -> "read")
    parameters += (NebulaOptions.GRAPH_ADDRESS -> "192.168.8.6:3820")
    parameters += (NebulaOptions.GRAPH_NAME -> "nba")
    parameters += (NebulaOptions.USER_NAME -> "root")
    parameters += (NebulaOptions.PASSWD -> "Nebula123")
    parameters += (NebulaOptions.LABEL -> "edge_type_follow")
    parameters += (NebulaOptions.RETURN_COLS -> "")
    val nebulaOptions = new NebulaOptions(CaseInsensitiveMap(parameters))
    val schema = NebulaUtils.getSchema(nebulaOptions)
    assert(schema.fields.length == 2)
    assert(schema.fields.map(field => field.name).contains("id"))
    assert(schema.fields.map(field => field.name).contains("id"))
  }

  test("getGqlSchema") {
    var parameters: Map[String, String] = Map()
    parameters += (NebulaOptions.TYPE -> "GQL")
    parameters += (NebulaOptions.OPERATE_TYPE -> "read")
    parameters += (NebulaOptions.GRAPH_ADDRESS -> "192.168.8.6:3820")
    parameters += (NebulaOptions.GRAPH_NAME -> "nba")
    parameters += (NebulaOptions.USER_NAME -> "root")
    parameters += (NebulaOptions.PASSWD -> "Nebula123")

    parameters += (NebulaOptions.GQL -> "use nba match(v)-[e:follow]-(v1) return v.id, v1.id, e.followness as followness, e.likeness as likeness")
    var nebulaOptions = new NebulaOptions(CaseInsensitiveMap(parameters))
    var schema = NebulaUtils.getSchemaForGql(nebulaOptions)
    assert(schema.fields.length == 4)
    assert(schema.fields.map(field=>field.name).contains("v.id"))
    assert(schema.fields.map(field=>field.name).contains("v1.id"))
    assert(schema.fields.map(field=>field.name).contains("followness"))
    assert(schema.fields.map(field=>field.name).contains("likeness"))

    parameters += (NebulaOptions.GQL -> "use nba match(v)-[e:follow]-(v1) return v.id, v1.id, e.followness as followness, e.likeness as likeness limit 1")
    nebulaOptions = new NebulaOptions(CaseInsensitiveMap(parameters))
    schema = NebulaUtils.getSchemaForGql(nebulaOptions)
    assert(schema.fields.length == 4)
    assert(schema.fields.map(field => field.name).contains("followness"))
    assert(schema.fields.map(field => field.name).contains("likeness"))

    parameters += (NebulaOptions.GQL -> "use nba match(v)-[e:follow]-(v1) return e limit 1")
    nebulaOptions = new NebulaOptions(CaseInsensitiveMap(parameters))
    schema = NebulaUtils.getSchemaForGql(nebulaOptions)
    assert(schema.fields.length == 1)
    assert(schema.fields.map(field => field.name).contains("e"))


    parameters += (NebulaOptions.GQL -> "use nba match p=(v)-[e:follow]-(v1) return p limit 1")
    nebulaOptions = new NebulaOptions(CaseInsensitiveMap(parameters))
    schema = NebulaUtils.getSchemaForGql(nebulaOptions)
    assert(schema.fields.length == 1)
    assert(schema.fields.map(field => field.name).contains("p"))
  }

}
