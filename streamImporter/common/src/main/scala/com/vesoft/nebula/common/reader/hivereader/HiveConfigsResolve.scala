
package com.vesoft.nebula.common.reader.hivereader

import com.typesafe.config.Config
import com.vesoft.nebula.common.configuration.ConfigConstant.DEFAULT_READ_PARTITION
import com.vesoft.nebula.common.configuration.ConfigUtil.getOrElse
import com.vesoft.nebula.common.configuration.{Configs, DataPreProcessConfig, DataSourceConfigEntry, FileFormatCategory, SchemaConfig, SourceCategory}
import org.apache.log4j.Logger

object HiveConfigsResolve {
  private[this] val LOG = Logger.getLogger(this.getClass)

  def parse(configPath: String): Configs = {
    null
  }

  def getSourceConfigEntry(category: SourceCategory.Value,
                           fileFormat: Option[FileFormatCategory.Value] = None,
                           sourceConfig: Config,
                           preProcessConfig: List[DataPreProcessConfig],
                           schemaConfigs: List[SchemaConfig]): DataSourceConfigEntry = {
    val statement     = sourceConfig.getString("statement")
    val readPartition = getOrElse(Option(sourceConfig), "readPartition", DEFAULT_READ_PARTITION)
    val uris          = if (sourceConfig.hasPath("uris")) sourceConfig.getString("uris") else null
    val config = HiveSourceConfigEntry(category, readPartition, statement, uris, schemaConfigs, preProcessConfig)
    LOG.info(config.toString)
    config
  }
}
