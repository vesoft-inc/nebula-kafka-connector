/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration

import com.vesoft.nebula.common.configuration.preConfig.PreConfigResolve
import com.vesoft.nebula.common.schema.{Edge, Node}
import org.scalatest.PrivateMethodTester
import org.scalatest.funsuite.AnyFunSuite

class PreConfigResolveSuite extends AnyFunSuite with PrivateMethodTester {
  test("test resolve pre config file and generate DDL from source schema file") {
    val configFilePath = "streamImporter/common/src/test/resources/pre-import.conf"
    val (ddl, configs) = PreConfigResolve.parse(configFilePath)

    val nebulaGraphConfigEntry = configs.nebulaGraphConfigEntry
    val mqClusterConfigEntry   = configs.mqClusterConfigEntry
    val errorConfigEntry       = configs.errorConfigEntry
    val sourceConfigEntrys     = configs.sourceConfigEntrys

    // assert some default configs
    assert(nebulaGraphConfigEntry.generateDDL)
    assert(mqClusterConfigEntry.server.equalsIgnoreCase("127.0.0.1:9092"))
    assert(mqClusterConfigEntry.topic.equalsIgnoreCase("nebula"))
    assert(errorConfigEntry.path.equalsIgnoreCase("file:///tmp/errors/"))
    assert(errorConfigEntry.maxRecords == Int.MaxValue)
    assert(sourceConfigEntrys.size == 2)

    assert(
      ddl.equals(
        "CREATE GRAPH TYPE graph_type_nba IF NOT EXISTS AS GRAPH TYPE { " +
          "(player(a) LABEL player(c string,g int))," +
          "(person(b) LABEL person(d int,e string))," +
          "(player)-[friend LABEL friend {c string,d int,e datetime,f int}]->(player) }"))
  }

  test("test resolve schema file, get list of node and edge schema") {
    val resolveSchema = PrivateMethod[(List[Node], List[Edge])](Symbol("resolveSchema"))
    val (nodes, edges) = PreConfigResolve invokePrivate resolveSchema(
      "streamImporter/common/src/test/resources/schema.csv")
    assert(nodes.size == 2)
    assert(edges.size == 1)

    for (node <- nodes) {
      node.getNodeTypeName match {
        case "player" =>
          assert(node.getVidField.equalsIgnoreCase("a"))
          assert(node.getVidDataType.equalsIgnoreCase("string"))
          assert(node.getProperties.size() == 2)
        case "person" =>
          assert(node.getVidField.equalsIgnoreCase("b"))
          assert(node.getVidDataType.equalsIgnoreCase("string"))
          assert(node.getProperties.size() == 2)
      }
    }

    val edge = edges.head
    assert(edge.getEdgeTypeName.equalsIgnoreCase("friend"))
    assert(edge.getSrcField.equalsIgnoreCase("a"))
    assert(edge.getDstField.equalsIgnoreCase("b"))
    assert(edge.getProperties.size() == 4)
  }

  test("test the check for schema format") {
    val checkSchema = PrivateMethod[String](Symbol("check"))
    val edgeSchema =
      "edge:friend:srckey:a:player, edge:friend:dstkey:b:player, c:string, d:int, e:datetime, f:int"
    assert((PreConfigResolve invokePrivate checkSchema(edgeSchema)).equals("edge"))

    val nodeSchema = "node:player:key:a:string, c:string, g:int"
    assert((PreConfigResolve invokePrivate checkSchema(nodeSchema)).equals("node"))
  }

  test("test regex for schema match") {
    val nodeRegexPattern    = "^node:[^:]+:key:[^:]+:(string|int)$"
    val edgeSrcRegexPattern = "^edge:[^:]+:srckey:[^:]+:[^:]+$"
    val edgeDstRegexPattern = "^edge:[^:]+:dstkey:[^:]+:[^:]+$"
    val propRegexPattern =
      "[^:]+:(string|int8|int16|int32|int64|int|date|time|datetime|duration|bool|float|double)"
    val nodeSchema    = "node:player:key:a:string"
    val edgeSrcSchema = "edge:friend:srckey:a:player"
    val edgeDstSchema = "edge:friend:dstkey:b:player"
    val propSchema    = "a:string"

    assert(nodeSchema.matches(nodeRegexPattern))
    assert(edgeSrcSchema.matches(edgeSrcRegexPattern))
    assert(edgeDstSchema.matches(edgeDstRegexPattern))
    assert(propSchema.matches(propRegexPattern))
    assert("a:int".matches(propRegexPattern))
    assert("a:bool".matches(propRegexPattern))
    assert("a:double".matches(propRegexPattern))
    assert("a:float".matches(propRegexPattern))
    assert("a:datetime".matches(propRegexPattern))
    assert("a:time".matches(propRegexPattern))
    assert("a:date".matches(propRegexPattern))
    assert("a:duration".matches(propRegexPattern))
  }

  test("test save configs into file"){
    val configFilePath = "streamImporter/common/src/test/resources/pre-import.conf"
    val (ddl, configs) = PreConfigResolve.parse(configFilePath)
    val targetPath = "streamImporter/common/src/test/resources/generated_import.conf"
    PreConfigResolve.save(configs, targetPath)
    ConfigsResolve.parse(targetPath)
  }

}
