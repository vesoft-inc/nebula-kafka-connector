
package com.vesoft.nebula.common.configuration

import java.io.{File, InputStreamReader}
import java.nio.file.Files
import com.typesafe.config.{Config, ConfigFactory}
import com.vesoft.nebula.common.configuration.ConfigConstant.{
  DEFAULT_BATCH_SIZE,
  DEFAULT_WRITE_PARALLEL
}
import com.vesoft.nebula.common.configuration.ConfigUtil.{
  getConfigsOrNone,
  getOrElse,
  toFileFormatCategory,
  toSourceCategory
}
import org.apache.hadoop.conf.Configuration
import org.apache.hadoop.fs.{FSDataInputStream, FileSystem, Path}
import org.apache.log4j.Logger

import scala.collection.JavaConverters.asScalaBufferConverter
import scala.collection.mutable

/**
  * config file parser
  *
  * each data source should extends {@link ConfigsResolve} and implement the parse and getSourceConfigEntry function.
  * for parse function, subclass just need to do this:
  *
  * override def parse(configPath:String): Configs = {super.parse(configPath)}
  *
  * */
object ConfigsResolve {
  private[this] val LOG = Logger.getLogger(this.getClass)
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

    // parse nebula config
    val nebulaGraphConfigEntry = ConfigUtil.parseNebulaConfig(config)
    nebulaGraphConfigEntry.check()

    // parse MQ cluster
    val mqClusterConfigEntry = ConfigUtil.parseMQConfig(config)
    mqClusterConfigEntry.check()

    // parse error
    val errorConfigEntry = ConfigUtil.parseErrorConfig(config)

    // parse sources
    val sourceConfigEntrys = mutable.ListBuffer[DataSourceConfigEntry]()
    val sourceConfigs      = getConfigsOrNone(config, "sources")
    if (sourceConfigs.isDefined) {
      for (sourceConfig <- sourceConfigs.get.asScala) {
        val category   = toSourceCategory(sourceConfig.getString("type"))
        val fileFormat = toFileFormatCategory(category, sourceConfig)
        // parse source data pre-process
        val preProcessConfigEntrys = ConfigUtil.parsePreProcessConfigs(sourceConfig)

        // parse schema config, include nodetype and edgetype
        val schemaConfigEntrys = mutable.ListBuffer[SchemaConfig]()

        val nodetypeConfigEntrys = mutable.ListBuffer[SchemaConfig]()
        val nodetypeConfigs      = getConfigsOrNone(sourceConfig, "nodetypes")
        if (nodetypeConfigs.isDefined) {
          for (nodetypeConfig <- nodetypeConfigs.get.asScala) {
            val name         = nodetypeConfig.getString("name")
            val primaryKey   = nodetypeConfig.getString("primaryKey")
            val sourceFields = nodetypeConfig.getStringList("sourceFields").asScala
            val nebulaFields = nodetypeConfig.getStringList("nebulaFields").asScala
            val batchSize    = getOrElse(Option(nodetypeConfig), "batchSize", DEFAULT_BATCH_SIZE)
            val writeParallel =
              getOrElse(Option(nodetypeConfig), "writeParallel", DEFAULT_WRITE_PARALLEL)
            nodetypeConfigEntrys += NodeConfig(name,
                                               sourceFields.toList,
                                               nebulaFields.toList,
                                               batchSize,
                                               writeParallel,
                                               primaryKey)
          }
        }

        val edgetypeConfigEntrys = mutable.ListBuffer[SchemaConfig]()
        val edgetypeConfigs      = getConfigsOrNone(sourceConfig, "edgetypes")
        if (edgetypeConfigs.isDefined) {
          for (edgetypeConfig <- edgetypeConfigs.get.asScala) {
            val name          = edgetypeConfig.getString("name")
            val srcPrimaryKey = edgetypeConfig.getString("sourceKey")
            val dstPrimaryKey = edgetypeConfig.getString("targetKey")
            val sourceFields  = edgetypeConfig.getStringList("sourceFields").asScala
            val nebulaFields  = edgetypeConfig.getStringList("nebulaFields").asScala
            val batchSize     = getOrElse(Option(edgetypeConfig), "batchSize", DEFAULT_BATCH_SIZE)
            val writeParallel =
              getOrElse(Option(edgetypeConfig), "writeParallel", DEFAULT_WRITE_PARALLEL)
            edgetypeConfigEntrys += EdgeConfig(name,
                                               sourceFields.toList,
                                               nebulaFields.toList,
                                               batchSize,
                                               writeParallel,
                                               srcPrimaryKey,
                                               dstPrimaryKey)

          }
        }
        schemaConfigEntrys.appendAll(nodetypeConfigEntrys)
        schemaConfigEntrys.appendAll(edgetypeConfigEntrys)
        sourceConfigEntrys += ConfigUtil.getSourceConfigEntry(category,
                                                              fileFormat,
                                                              sourceConfig,
                                                              preProcessConfigEntrys,
                                                              schemaConfigEntrys.toList)
      }
    }
    Configs(nebulaGraphConfigEntry,
            mqClusterConfigEntry,
            errorConfigEntry,
            sourceConfigEntrys.toList)
  }

}
