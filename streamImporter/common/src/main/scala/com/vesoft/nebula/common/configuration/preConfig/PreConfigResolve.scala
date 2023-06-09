/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration.preConfig

import com.typesafe.config.{Config, ConfigFactory}
import com.vesoft.nebula.common.configuration.ConfigConstant.{
  DEFAULT_BATCH_SIZE,
  DEFAULT_GRAPH_TYPE,
  DEFAULT_WRITE_PARALLEL
}
import com.vesoft.nebula.common.configuration.ConfigUtil.{
  getSourceConfigEntry,
  toFileFormatCategory,
  toSourceCategory
}
import com.vesoft.nebula.common.configuration.{
  ConfigUtil,
  Configs,
  DataSourceConfigEntry,
  EdgeConfig,
  NodeConfig,
  SchemaConfig
}
import com.vesoft.nebula.common.schema.{Edge, Node}
import org.apache.hadoop.conf.Configuration
import org.apache.hadoop.fs.{FSDataInputStream, FileSystem, Path}
import org.apache.log4j.Logger

import java.io.{BufferedReader, File, InputStreamReader}
import java.nio.file.Files
import scala.collection.JavaConverters.{
  asScalaBufferConverter,
  mapAsJavaMapConverter,
  mapAsScalaMapConverter
}
import scala.collection.mutable
import scala.collection.mutable.ListBuffer
import scala.io.{BufferedSource, Source}

object PreConfigResolve {
  private[this] val LOG = Logger.getLogger(this.getClass)

  private val nodeRegexPattern    = """^node:[^:]+:key:[^:]+:(string|int)$""".r
  private val edgeSrcRegexPattern = """^edge:[^:]+:srckey:[^:]+:[^:]+$""".r
  private val edgeDstRegexPattern = """^edge:[^:]+:dstkey:[^:]+:[^:]+$""".r
  private val propRegexPattern =
    """[^:]+:(string|int8|int16|int32|int64|int|date|time|datetime|duration|bool|float|double)""".r

  def parse(configPath: String): (String, Configs) = {
    var config: Config = null
    if (configPath.startsWith("hdfs://")) {
      val hadoopConfig: Configuration = new Configuration()
      val fs: FileSystem              = org.apache.hadoop.fs.FileSystem.get(hadoopConfig)
      val file: FSDataInputStream     = fs.open(new Path(configPath))
      val reader                      = new InputStreamReader(file)
      config = ConfigFactory.parseReader(reader)
    } else {
      if (!Files.exists(new File(configPath).toPath)) {
        throw new IllegalArgumentException(s"${configPath} not exist")
      }
      config = ConfigFactory.parseFile(new File(configPath))
    }

    val nebulaConfig = ConfigUtil.getConfigOrNone(config, "nebula")
    val graphType    = ConfigUtil.getOrElse(nebulaConfig, "graphType", DEFAULT_GRAPH_TYPE)

    val nebulaGraphConfigEntry = ConfigUtil.parseNebulaConfig(config)
    val mqClusterConfigEntry   = ConfigUtil.parseMQConfig(config)
    val errorConfigEntry       = ConfigUtil.parseErrorConfig(config)
    val (ddlStatements, sourceConfigEntrys) =
      parseSourceConfigs(config, graphType, nebulaGraphConfigEntry.generateDDL)

    var schemaDDL: String = null
    if (nebulaGraphConfigEntry.generateDDL) {
      schemaDDL =
        s"CREATE GRAPH TYPE graph_type_nba IF NOT EXISTS AS GRAPH TYPE { ${ddlStatements.mkString(",")} }"
      LOG.info("********************* schema DDL ********************* ")
      LOG.info(schemaDDL)
    }
    val configs =
      Configs(nebulaGraphConfigEntry, mqClusterConfigEntry, errorConfigEntry, sourceConfigEntrys)
    (schemaDDL, configs)
  }

  /**
    * resolve the schema in file to {@link Node} and {@link Edge}
    *
    * @param sourceSchemas schemas in schema file, one schema for one line.
    * @return schemas of node and schemas of edge
    *
    *         eg: schema in file:
    *         :edge:friend:srckey:a:player,:edge:friend:dstkey:b:player,:c:string,:d:int,:e:datetime,:f:int
    *         :node:player:key:a:string,:c:string,:g:int
    *
    *         result is:
    *         Node{nodeType=`player`,vidType=`string`,vidField=`a`,properties={c->string,g->int}}
    *         Edge{edgeType= `friend`,srcField=`a`,srcNodeType=`player`,dstField=`b`,dstNodeType=`player`,properties={c->string,d->int,e->datetime,f->int}
    * */
  private[this] def resolveSchema(schemaFilePath: String): (List[Node], List[Edge]) = {
    val sourceSchemas: ListBuffer[String] = mutable.ListBuffer[String]()
    if (schemaFilePath == null) {
      //TODO schema file可以为空，即基于源数据 自动生成schema。
    } else {
      // get schema information from schema file
      if (schemaFilePath.startsWith("hdfs://")) {
        val hadoopConfig: Configuration = new Configuration()
        val fs: FileSystem              = org.apache.hadoop.fs.FileSystem.get(hadoopConfig)
        val file: FSDataInputStream     = fs.open(new Path(schemaFilePath))
        val reader                      = new BufferedReader(new InputStreamReader(file))
        var line: String                = reader.readLine()
        while (line != null) {
          sourceSchemas += reader.readLine()
          line = reader.readLine()
        }
      } else {
        if (!Files.exists(new File(schemaFilePath).toPath)) {
          throw new IllegalArgumentException(s"$schemaFilePath not exist")
        }
        var bufferedSource: BufferedSource = null
        try {
          bufferedSource = Source.fromFile(schemaFilePath)
          val lines = bufferedSource.getLines().toList
          sourceSchemas.appendAll(lines)
        } finally {
          if (bufferedSource != null) {
            bufferedSource.close()
          }
        }
      }
    }

    val nodeConfigSchemas: ListBuffer[Node] = new ListBuffer[Node]()
    val edgeConfigSchemas: ListBuffer[Edge] = new ListBuffer[Edge]()
    for (schema <- sourceSchemas if !schema.equals("")) {
      val schemaType = check(schema)
      val items      = schema.split(",")
      var node: Node = null
      var edge: Edge = null
      schemaType match {
        case "node" =>
          node = new Node
          for (item <- items) {
            item match {
              case nodeRegexPattern(_*) =>
                val elements = item.split(":")
                node.setNodeType(elements(1)).setVidField(elements(3)).setVidType(elements(4))
              case propRegexPattern(_*) =>
                val elements = item.split(":")
                node.addProperty(elements(0), elements(1))
            }
          }
        case "edge" =>
          edge = new Edge
          for (item <- items) {
            item match {
              case edgeSrcRegexPattern(_*) =>
                val elements = item.split(":")
                edge.setEdgeType(elements(1)).setSrcField(elements(3)).setSrcNodeType(elements(4))
              case edgeDstRegexPattern(_*) =>
                val elements = item.split(":")
                edge.setDstField(elements(3)).setDstNodeType(elements(4))
              case propRegexPattern(_*) =>
                val elements = item.split(":")
                edge.addProperty(elements(0), elements(1))
            }
          }
      }
      if (node != null) nodeConfigSchemas.append(node)
      if (edge != null) edgeConfigSchemas.append(edge)
    }
    (nodeConfigSchemas.toList, edgeConfigSchemas.toList)
  }

  private[this] def parseSourceConfigs(
      config: Config,
      graphType: String,
      generateDDL: Boolean): (List[String], List[DataSourceConfigEntry]) = {
    val sourceConfigEntrys = mutable.ListBuffer[DataSourceConfigEntry]()
    val sourceConfigs      = ConfigUtil.getConfigsOrNone(config, "sources")

    val ddlStatements: ListBuffer[String] = new ListBuffer[String]()
    if (sourceConfigs.isDefined) {
      for (sourceConfig <- sourceConfigs.get.asScala) {
        // source config parse
        val category         = toSourceCategory(sourceConfig.getString("type"))
        val fileFormat       = toFileFormatCategory(category, sourceConfig)
        val preProcessConfig = ConfigUtil.parsePreProcessConfigs(sourceConfig)

        // schema parse
        val schemaFile = sourceConfig.getString("schemaFile")

        // 1. 从 schemaFile中解析 源数据node，edge 的schema
        val (nodeConfigSchemas, edgeConfigSchemas) = resolveSchema(schemaFile)

        val schemaConfigs: ListBuffer[SchemaConfig] = new ListBuffer[SchemaConfig]()
        // 2. 对每个node schema 获取对应的nebula 属性字段集合
        val nodeConfigs: ListBuffer[NodeConfig] = new ListBuffer[NodeConfig]()
        for (node <- nodeConfigSchemas) {
          if (generateDDL) {
            node.setNebulaProperties(node.getProperties)
            ddlStatements += node.getSchemaString
          } else {
            node.setNebulaProperties(getNodePropertyFields(graphType, node.getNodeType).asJava)
          }
          val sourceFields2NebulaFields        = node.getPropMapping.asScala
          val sourceFields: ListBuffer[String] = new ListBuffer[String]()
          val nebulaFields: ListBuffer[String] = new ListBuffer[String]()
          for (kv <- sourceFields2NebulaFields) {
            sourceFields += kv._1
            nebulaFields += kv._2
          }
          val nodeConfig = NodeConfig(node.getNodeType,
                                      sourceFields.toList,
                                      nebulaFields.toList,
                                      DEFAULT_BATCH_SIZE,
                                      DEFAULT_WRITE_PARALLEL,
                                      node.getVidField)
          schemaConfigs += nodeConfig
        }

        val edgeConfigs: ListBuffer[EdgeConfig] = new ListBuffer[EdgeConfig]()
        for (edge <- edgeConfigSchemas) {
          if (generateDDL) {
            edge.setNebulaProperties(edge.getProperties)
            ddlStatements += edge.getSchemaString
          } else {
            edge.setNebulaProperties(getEdgePropertyFields(graphType, edge.getEdgeType).asJava)
          }
          val sourceFields2NebulaFields        = edge.getPropMapping.asScala
          val sourceFields: ListBuffer[String] = new ListBuffer[String]()
          val nebulaFields: ListBuffer[String] = new ListBuffer[String]()
          for (kv <- sourceFields2NebulaFields) {
            sourceFields += kv._1
            nebulaFields += kv._2
          }
          val edgeConfig = EdgeConfig(edge.getEdgeType,
                                      sourceFields.toList,
                                      nebulaFields.toList,
                                      DEFAULT_BATCH_SIZE,
                                      DEFAULT_WRITE_PARALLEL,
                                      edge.getSrcField,
                                      edge.getDstField)
          schemaConfigs += edgeConfig
        }
        sourceConfigEntrys += getSourceConfigEntry(category,
                                                   fileFormat,
                                                   sourceConfig,
                                                   preProcessConfig,
                                                   schemaConfigs.toList)
      }
    }
    (ddlStatements.toSet.toList, sourceConfigEntrys.toList)
  }

  /**
    * TODO depend core statement to show nodeType's schema
    *
    * */
  private[this] def getNodePropertyFields(graphType: String,
                                          nodeType: String): Map[String, String] = {
    Map()
  }

  /**
    * TODO depend core statement to show nodeType's schema
    *
    * */
  private[this] def getEdgePropertyFields(graphType: String,
                                          edgeType: String): Map[String, String] = {
    Map()
  }

  /**
    * check if the schema in file valid and resolve the schema type
    *
    * */
  private[this] def check(schema: String): String = {
    val items              = schema.split(",")
    var schemaType: String = null
    for (item <- items) {
      item.trim match {
        case nodeRegexPattern(_*) =>
          schemaType = "node"
          return schemaType
        case edgeSrcRegexPattern(_*) =>
          schemaType = "edge"
          return schemaType
        case edgeDstRegexPattern(_*) =>
          schemaType = "edge"
          return schemaType
        case propRegexPattern(_*) => null
        case _                    => throw new IllegalArgumentException("schema format is invalid.")
      }
    }
    schemaType
  }
}
