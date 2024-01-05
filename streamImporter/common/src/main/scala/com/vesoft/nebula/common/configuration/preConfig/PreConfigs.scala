/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration.preConfig

import com.typesafe.config.{Config, ConfigFactory}
import com.vesoft.nebula.common.configuration.ConfigConstant.{
  DEFAULT_BATCH_SIZE,
  DEFAULT_CONNECT_TIMEOUT,
  DEFAULT_ERROR_MAX_RECORDS,
  DEFAULT_ERROR_PATH,
  DEFAULT_GRAPH_ADDR,
  DEFAULT_GRAPH_NAME,
  DEFAULT_GRAPH_PASSWD,
  DEFAULT_GRAPH_TYPE,
  DEFAULT_GRAPH_USER,
  DEFAULT_REDPANDA_REPLIC,
  DEFAULT_REDPANDA_SERER,
  DEFAULT_REDPANDA_TOPIC,
  DEFAULT_REQUEST_TIMEOUT,
  DEFAULT_RETRY_INTERVAL_TIME,
  DEFAULT_WRITE_MODE,
  DEFAULT_WRITE_PARALLEL
}
import com.vesoft.nebula.common.configuration.ConfigUtil.{toFileFormatCategory, toSourceCategory}
import com.vesoft.nebula.common.configuration.{
  ConfigUtil,
  Configs,
  DataPreProcessConfig,
  DataSourceConfigEntry,
  EdgeConfig,
  ErrorConfigEntry,
  FileFormatCategory,
  MQClusterConfigEntry,
  NebulaGraphConfigEntry,
  NodeConfig,
  SchemaConfig,
  SourceCategory
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

/**
  * pre execution to generate the import config file
  *
  * */
object PreConfigEntry {
  private[this] val LOG = Logger.getLogger(this.getClass)

  private val nodeRegexPattern    = """"^:node:[^:]+:key:[^:]+:(string|int)$""".r
  private val edgeSrcRegexPattern = """^:edge:[^:]+:srckey:[^:]+:[^:]$""".r
  private val edgeDstRegexPattern = """^:edge:[^:]+:dstkey:[^:]+:[^:]$""".r
  private val propRegexPattern =
    """^:[^:]+:(string|int8|int16|int32|int64|int|date|time|datetime|duration|bool|float|double)""".r

  def parse(configPath: String): Configs = {
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

    val (graphType, nebulaGraphConfigEntry) = parseNebulaConfig(config)
    val mqClusterConfigEntry                = parseMQConfig(config)
    val errorConfigEntry                    = parseErrorConfig(config)
    val (ddlStatements, sourceConfigEntrys) =
      parseSourceConfigs(config, graphType, nebulaGraphConfigEntry.generateDDL)

    if (nebulaGraphConfigEntry.generateDDL) {
      val schemaDDL =
        s"CREATE GRAPH TYPE graph_type_nba IF NOT EXISTS AS { ${ddlStatements.mkString(",")} }"
      LOG.info("********************* schema DDL ********************* \n")
      LOG.info(schemaDDL)
    }

    Configs(nebulaGraphConfigEntry, mqClusterConfigEntry, errorConfigEntry, sourceConfigEntrys)
  }

  /**
    * resolve the schema in file to {@link Node} and {@link Edge}
    *
    * @param sourceSchemas schemas in schema file, one schema for one line.
    * @return schemas of node and schemas of edge
    *
    * eg: schema in file:
    *         :edge:friend:srckey:a:player,:edge:friend:dstkey:b:player,:c:string,:d:int,:e:datetime,:f:int
    *         :node:player:key:a:string,:c:string,:g:int
    *
    *   result is:
    *   Node{nodeType=`player`,vidType=`string`,vidField=`a`,properties={c->string,g->int}}
    *   Edge{edgeType= `friend`,srcField=`a`,srcNodeType=`player`,dstField=`b`,dstNodeType=`player`,properties={c->string,d->int,e->datetime,f->int}
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
    for (schema <- sourceSchemas) {
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
                node
                  .setNodeTypeName(elements(1))
                  .setVidField(elements(3))
                  .setVidDataType(elements(4))
              case propRegexPattern(_*) =>
                val elements = item.split(":")
                node.addProperty(elements(1), elements(2))
            }
          }
        case "edge" =>
          edge = new Edge
          for (item <- items) {
            item match {
              case edgeSrcRegexPattern(_*) =>
                val elements = item.split(":")
                edge
                  .setEdgeTypeName(elements(1))
                  .setSrcField(elements(3))
                  .setSrcNodeTypeName(elements(4))
              case edgeDstRegexPattern(_*) =>
                val elements = item.split(":")
                edge.setDstField(elements(3)).setDstNodeTypeName(elements(4))
              case propRegexPattern(_*) =>
                val elements = item.split(":")
                edge.addProperty(elements(1), elements(2))
            }
          }
      }
      if (node != null) nodeConfigSchemas.append(node)
      if (edge != null) edgeConfigSchemas.append(edge)
    }
    (nodeConfigSchemas.toList, edgeConfigSchemas.toList)
  }

  private[this] def parseNebulaConfig(config: Config): (String, NebulaGraphConfigEntry) = {
    // parse nebula config
    val nebulaConfig = ConfigUtil.getConfigOrNone(config, "nebula")

    val graphAddress = ConfigUtil.getOrElse(nebulaConfig, "graphAddr", DEFAULT_GRAPH_ADDR)
    val graphType    = ConfigUtil.getOrElse(nebulaConfig, "graphType", DEFAULT_GRAPH_TYPE)
    val graphName    = ConfigUtil.getOrElse(nebulaConfig, "graphName", DEFAULT_GRAPH_NAME)

    val user   = ConfigUtil.getOrElse(nebulaConfig, "user", DEFAULT_GRAPH_USER)
    val passwd = ConfigUtil.getOrElse(nebulaConfig, "passwd", DEFAULT_GRAPH_PASSWD)
    val connectTimeout =
      ConfigUtil.getOrElse(nebulaConfig, "connectTimeout", DEFAULT_CONNECT_TIMEOUT)
    val requestTimeout =
      ConfigUtil.getOrElse(nebulaConfig, "requestTimeout", DEFAULT_REQUEST_TIMEOUT)

    val retryInterval =
      ConfigUtil.getOrElse(nebulaConfig, "retryInterval", DEFAULT_RETRY_INTERVAL_TIME)
    val mode =
      ConfigUtil.toSinkCategory(ConfigUtil.getOrElse(nebulaConfig, "mode", DEFAULT_WRITE_MODE))
    val generateDDL = ConfigUtil.getOrElse(nebulaConfig, "generateDDL", false)
    val nebulaGraphConfigEntry =
      NebulaGraphConfigEntry(graphAddress,
                             graphName,
                             user,
                             passwd,
                             connectTimeout,
                             requestTimeout,
                             retryInterval,
                             mode,
                             generateDDL)
    nebulaGraphConfigEntry.check()
    LOG.info(nebulaGraphConfigEntry.toString)
    (graphType, nebulaGraphConfigEntry)
  }

  private[this] def parseMQConfig(config: Config): MQClusterConfigEntry = {
    // parse mq config
    val redPandaConfig       = ConfigUtil.getConfigOrNone(config, "mq")
    val mqServer             = ConfigUtil.getOrElse(redPandaConfig, "server", DEFAULT_REDPANDA_SERER)
    val topic                = ConfigUtil.getOrElse(redPandaConfig, "topic", DEFAULT_REDPANDA_TOPIC)
    val replic               = ConfigUtil.getOrElse(redPandaConfig, "replic", DEFAULT_REDPANDA_REPLIC)
    val mqClusterConfigEntry = MQClusterConfigEntry(mqServer, topic, replic)
    LOG.info(mqClusterConfigEntry.toString)
    mqClusterConfigEntry
  }

  private[this] def parseErrorConfig(config: Config): ErrorConfigEntry = {
    // parse error config
    val errorConfig      = ConfigUtil.getConfigOrNone(config, "error")
    val errorPath        = ConfigUtil.getOrElse(errorConfig, "path", DEFAULT_ERROR_PATH)
    val errorMaxRecords  = ConfigUtil.getOrElse(errorConfig, "maxRecords", DEFAULT_ERROR_MAX_RECORDS)
    val errorConfigEntry = ErrorConfigEntry(errorPath, errorMaxRecords)
    LOG.info(errorConfigEntry.toString)
    errorConfigEntry
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
            node.setNebulaProperties(getNodePropertyFields(graphType, node.getNodeTypeName).asJava)
          }
          val sourceFields2NebulaFields        = node.getPropMapping.asScala
          val sourceFields: ListBuffer[String] = new ListBuffer[String]()
          val nebulaFields: ListBuffer[String] = new ListBuffer[String]()
          for (kv <- sourceFields2NebulaFields) {
            sourceFields += kv._1
            nebulaFields += kv._2
          }
          val nodeConfig = NodeConfig(node.getNodeTypeName,
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
            edge.setNebulaProperties(getEdgePropertyFields(graphType, edge.getEdgeTypeName).asJava)
          }
          val sourceFields2NebulaFields        = edge.getPropMapping.asScala
          val sourceFields: ListBuffer[String] = new ListBuffer[String]()
          val nebulaFields: ListBuffer[String] = new ListBuffer[String]()
          for (kv <- sourceFields2NebulaFields) {
            sourceFields += kv._1
            nebulaFields += kv._2
          }
          val edgeConfig = EdgeConfig(edge.getEdgeTypeName,
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
    (ddlStatements.toList, sourceConfigEntrys.toList)
  }

  def getSourceConfigEntry(category: SourceCategory.Value,
                           fileFormat: Option[FileFormatCategory.Value] = None,
                           sourceConfig: Config,
                           preProcessConfig: List[DataPreProcessConfig],
                           schemaConfigs: List[SchemaConfig]): DataSourceConfigEntry = ???

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
      item match {
        case nodeRegexPattern(_*)    => schemaType = "node"
        case edgeSrcRegexPattern(_*) => schemaType = "edge"
        case edgeDstRegexPattern(_*) => schemaType = "edge"
        case propRegexPattern(_*)    => null
        case _                       => throw new IllegalArgumentException("schema format is invalid.")
      }
    }
    schemaType
  }

  /**
    * generate NebulaGraph DDL for node type and edge type from source schemas
   **/
  private[this] def generateDDLFromSourceSchema(graphType: String,
                                                graphName: String,
                                                nodeConfigSchemas: List[Node],
                                                edgeConfigSchemas: List[Edge]): String = {
    val nodeSchemas: ListBuffer[String] = new ListBuffer[String]()
    for (node <- nodeConfigSchemas) {
      nodeSchemas.append(node.getSchemaString)
    }
    val edgeSchemas: ListBuffer[String] = new ListBuffer[String]()
    for (edge <- edgeConfigSchemas) {
      edgeSchemas.append(edge.getSchemaString)
    }

    val ddl = s"CREATE GRAPH TYPE $graphType IF NOT EXISTS AS {${nodeSchemas.mkString(
      ",")},${edgeSchemas.mkString(",")}}; \nCREATE GRAPH IF NOT EXISTS $graphName TYPED $graphType"
    ddl
  }

  /**
    * query the graph DDL from NebulaGraph
    *
    * //TODO implementation
    * @param graphType graph type name
    * @param nebulaGraphConfig NebulaGraph server configs
    *
    * @return DDL
    * */
  private[this] def getGraphDDL(graphType: String,
                                nebulaGraphConfig: NebulaGraphConfigEntry): String = {
    // TODO remove the MOCK DDL
    "CREATE GRAPH TYPE graph_type IF NOT EXISTS AS {" +
      "(node_type LABEL player {id INT PRIMARY KEY, name STRING})," +
      "(node_type)-[edge_type LABEL follow {followness INT}]->(node_type)}"
  }

  /**
    * get schema from hdfs csv data file
    */
  private[this] def getCsvSchemaFromHdfs(): Unit = {}

  /**
    * get schema from hdfs json data file
    */
  private[this] def getJsonSchemaFromHdfs(): Unit = {}

  /**
    * get schema from oss csv data file
    */
  private[this] def getCsvSchemaFromOss(): Unit = {}

  /**
    * get schema from oss json data file
    */
  private[this] def getJsonSchemaFromOss(): Unit = {}

  /**
    * get schema from s3 csv data file
    */
  private[this] def getCsvSchemaFromS3(): Unit = {}

  /**
    * get schema from s3 json data file
    */
  private[this] def getJsonSchemaFromS3(): Unit = {}

  /**
    * get schema from Hive
    * */
  private[this] def getSchemaFromHive(): Unit = {}

  /**
    * get schema from jdbc
    * */
  private[this] def getSchemaFromJdbc(): Unit = {}

  /**
    * auto generate schema mappings for edge type
    *
    * @return srcId, dstId, sourceFields, nebulaFields
    *
    * */
  private[this] def generateNodeTypeSchemaMapping(
      schemaFile: String,
      ddl: String): (String, String, List[String], List[String]) = {
    null
  }

}
