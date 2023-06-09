/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration

import com.typesafe.config.Config
import com.vesoft.nebula.common.configuration.ConfigConstant.{
  DEFAULT_CONNECT_TIMEOUT,
  DEFAULT_ERROR_MAX_RECORDS,
  DEFAULT_ERROR_PATH,
  DEFAULT_GRAPH_ADDR,
  DEFAULT_GRAPH_NAME,
  DEFAULT_GRAPH_PASSWD,
  DEFAULT_GRAPH_USER,
  DEFAULT_REDPANDA_SERER,
  DEFAULT_REDPANDA_TOPIC,
  DEFAULT_REQUEST_TIMEOUT,
  DEFAULT_RETRY_INTERVAL_TIME,
  DEFAULT_WRITE_MODE
}
import com.vesoft.nebula.common.reader.hdfsreader.HdfsConfigsResolve
import com.vesoft.nebula.common.reader.hivereader.HiveConfigsResolve
import com.vesoft.nebula.common.reader.jdbcreader.JdbcConfigsResolve
import com.vesoft.nebula.common.reader.ossreader.OSSConfigsResolve
import com.vesoft.nebula.common.reader.s3reader.S3ConfigsResolve
import org.apache.log4j.Logger

import scala.collection.JavaConverters.asScalaBufferConverter
import scala.collection.mutable

object ConfigConstant {

  // nebula default config
  val DEFAULT_GRAPH_ADDR          = "127.0.0.1:9669"
  val DEFAULT_GRAPH_TYPE          = "graph_type"
  val DEFAULT_GRAPH_NAME          = "graph"
  val DEFAULT_GRAPH_USER          = "root"
  val DEFAULT_GRAPH_PASSWD        = "nebula"
  val DEFAULT_CONNECT_TIMEOUT     = 0
  val DEFAULT_REQUEST_TIMEOUT     = 0
  val DEFAULT_RETRY_INTERVAL_TIME = 0

  // redpanda default config
  val DEFAULT_REDPANDA_SERER = "127.0.0.1:9092"
  val DEFAULT_REDPANDA_TOPIC = "nebula"

  // error default config
  val DEFAULT_ERROR_PATH             = "file:///tmp/errors/"
  val DEFAULT_ERROR_MAX_RECORDS: Int = Int.MaxValue

  // read default config
  val DEFAULT_READ_PARTITION = 32

  // csv default config
  val DEFAULT_CSV_HEADER    = false
  val DEFAULT_CSV_SEPARATOR = ","

  // write default config
  val DEFAULT_WRITE_MODE     = "IMPORT"
  val DEFAULT_BATCH_SIZE     = 2000
  val DEFAULT_WRITE_PARALLEL = 32

  // checkpoint default config
  val DEFAULT_CHECK_POINT_PATH = "file:///tmp/checkpoint"

}

object ConfigUtil {
  private[this] val LOG = Logger.getLogger(this.getClass)

  /**
    * Get the config list by the path.
    *
    * @param config The com.vesoft.exchange.common.config.
    * @param path   The path of the com.vesoft.exchange.common.config.
    * @return
    */
  def getConfigsOrNone(config: Config, path: String): Option[java.util.List[_ <: Config]] = {
    if (config.hasPath(path)) {
      Some(config.getConfigList(path))
    } else {
      None
    }
  }

  /**
    * Get the config by the path.
    *
    * @param config
    * @param path
    * @return
    */
  def getConfigOrNone(config: Config, path: String): Option[Config] = {
    if (config.hasPath(path)) {
      Some(config.getConfig(path))
    } else {
      None
    }
  }

  /**
    * Get the value from config by the path which is optional.
    * If the path not exist, return the default value.
    *
    * @param config
    * @param path
    * @param defaultValue
    * @tparam T
    * @return
    */
  def getOrElse[T](config: Option[Config], path: String, defaultValue: T): T = {
    if (config.isDefined && config.get.hasPath(path)) {
      config.get.getAnyRef(path).asInstanceOf[T]
    } else {
      defaultValue
    }
  }

  /**
    * get source category.
    *
    * @param category name
    * @return
    */
  def toSourceCategory(category: String): SourceCategory.Value = {
    category.trim.toUpperCase match {
      case "HDFS" => SourceCategory.HDFS
      case "S3"   => SourceCategory.S3
      case "OSS"  => SourceCategory.OSS
      case "HIVE" => SourceCategory.HIVE
      case "JDBC" => SourceCategory.JDBC
      case _      => throw new IllegalArgumentException(s"source ${category} not support")
    }
  }

  /**
    * get file format for hdfs/s3/oss file system
    *
    */
  def toFileFormatCategory(sourceCategory: SourceCategory.Value,
                           config: Config): Option[FileFormatCategory.Value] = {
    sourceCategory match {
      case SourceCategory.HDFS | SourceCategory.S3 | SourceCategory.OSS => {
        val fileFormat = config.getString("format")
        fileFormat.toUpperCase match {
          case "CSV"  => Option(FileFormatCategory.CSV)
          case "JSON" => Option(FileFormatCategory.JSON)
          case _      => throw new IllegalArgumentException(s"file format $fileFormat not support")
        }
      }
      case _ => None
    }
  }

  /**
    * get sink category.
    *
    * @param category name
    * @return
    */
  def toSinkCategory(category: String): SinkCategory.Value = {
    category.trim.toUpperCase match {
      case "IMPORT"   => SinkCategory.IMPORT
      case "BULKLOAD" => SinkCategory.BULKLOAD
      case _          => throw new IllegalArgumentException(s"sink mode ${category} not support")
    }
  }

  def mkStringForMap(map: Map[String, String], sep: String, sepForKV: String): String = {
    map.mkString(",").replaceAll(" -> ", sepForKV)
  }

  /**
    * parse NebulaGraph Config
    * */
  def parseNebulaConfig(config: Config): NebulaGraphConfigEntry = {
    // parse nebula config
    val nebulaConfig = ConfigUtil.getConfigOrNone(config, "nebula")

    val graphAddress = ConfigUtil.getOrElse(nebulaConfig, "graphAddr", DEFAULT_GRAPH_ADDR)
    val graphName    = ConfigUtil.getOrElse(nebulaConfig, "graphName", DEFAULT_GRAPH_NAME)
    val user         = ConfigUtil.getOrElse(nebulaConfig, "user", DEFAULT_GRAPH_USER)
    val passwd       = ConfigUtil.getOrElse(nebulaConfig, "passwd", DEFAULT_GRAPH_PASSWD)
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
    LOG.info(nebulaGraphConfigEntry.toString)
    nebulaGraphConfigEntry
  }

  /**
    * parse MQ config
    *
    * */
  def parseMQConfig(config: Config): MQClusterConfigEntry = {
    // parse mq config
    val redPandaConfig       = ConfigUtil.getConfigOrNone(config, "mq")
    val mqServer             = ConfigUtil.getOrElse(redPandaConfig, "server", DEFAULT_REDPANDA_SERER)
    val topic                = ConfigUtil.getOrElse(redPandaConfig, "topic", DEFAULT_REDPANDA_TOPIC)
    val mqClusterConfigEntry = MQClusterConfigEntry(mqServer, topic)
    LOG.info(mqClusterConfigEntry.toString)
    mqClusterConfigEntry
  }

  def parseErrorConfig(config: Config): ErrorConfigEntry = {
    // parse error config
    val errorConfig      = ConfigUtil.getConfigOrNone(config, "error")
    val errorPath        = ConfigUtil.getOrElse(errorConfig, "path", DEFAULT_ERROR_PATH)
    val errorMaxRecords  = ConfigUtil.getOrElse(errorConfig, "maxRecords", DEFAULT_ERROR_MAX_RECORDS)
    val errorConfigEntry = ErrorConfigEntry(errorPath, errorMaxRecords)
    LOG.info(errorConfigEntry.toString)
    errorConfigEntry
  }

  /**
    * parse the PreProcess config, use configured config or default config
    *
    */
  def parsePreProcessConfigs(sourceConfig: Config): List[DataPreProcessConfig] = {
    // parse source data pre-process
    val preProcessConfigEntrys = mutable.ListBuffer[DataPreProcessConfig]()
    if (!sourceConfig.hasPath("pre_processes")) {
      return preProcessConfigEntrys.toList
    }
    val config = sourceConfig.getConfig("pre_processes")
    // concat config
    val concatConfigs = getConfigsOrNone(config, "concat")
    if (concatConfigs.isDefined) {
      for (concatConfig <- concatConfigs.get.asScala) {
        val oldFields = concatConfig.getStringList("oldFields").asScala
        val newField  = concatConfig.getString("newField")
        val sep       = concatConfig.getString("sep")
        preProcessConfigEntrys += ConcatConfig(oldFields.toList, newField, sep)
      }
    }
    // separator config
    val separatorConfigs = getConfigsOrNone(config, "separate")
    if (separatorConfigs.isDefined) {
      for (separatorConfig <- separatorConfigs.get.asScala) {
        val oldField  = separatorConfig.getString("oldField")
        val newFields = separatorConfig.getStringList("newFields").asScala
        val sep       = separatorConfig.getString("sep")
        preProcessConfigEntrys += SeparatorConfig(oldField, newFields.toList, sep)
      }
    }
    // filter config
    if (config.hasPath("filter")) {
      val filterConditions = config.getStringList("filter").asScala
      preProcessConfigEntrys += FilterConfig(filterConditions.toList)
    }

    // none value config
    if (config.hasPath("nonValue")) {
      val noneValue = config.getString("nonValue")
      preProcessConfigEntrys += NonValueConfig(noneValue)
    }
    preProcessConfigEntrys.toList
  }

  /**
    * parse source config for multiple data sources
    *
    */
  def getSourceConfigEntry(category: SourceCategory.Value,
                           fileFormat: Option[FileFormatCategory.Value] = None,
                           sourceConfig: Config,
                           preProcessConfig: List[DataPreProcessConfig],
                           schemaConfigs: List[SchemaConfig]): DataSourceConfigEntry = {
    category match {
      case SourceCategory.HDFS =>
        HdfsConfigsResolve.getSourceConfigEntry(category,
                                                fileFormat,
                                                sourceConfig,
                                                preProcessConfig,
                                                schemaConfigs)
      case SourceCategory.S3 =>
        S3ConfigsResolve.getSourceConfigEntry(category,
                                              fileFormat,
                                              sourceConfig,
                                              preProcessConfig,
                                              schemaConfigs)
      case SourceCategory.OSS =>
        OSSConfigsResolve.getSourceConfigEntry(category,
                                               fileFormat,
                                               sourceConfig,
                                               preProcessConfig,
                                               schemaConfigs)
      case SourceCategory.HIVE =>
        HiveConfigsResolve.getSourceConfigEntry(category,
                                                fileFormat,
                                                sourceConfig,
                                                preProcessConfig,
                                                schemaConfigs)
      case SourceCategory.JDBC =>
        JdbcConfigsResolve.getSourceConfigEntry(category,
                                                fileFormat,
                                                sourceConfig,
                                                preProcessConfig,
                                                schemaConfigs)

    }
  }

}
