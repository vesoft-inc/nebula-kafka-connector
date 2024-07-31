
package com.vesoft.nebula.common.reader.ossreader

import com.typesafe.config.Config
import com.vesoft.nebula.common.configuration.ConfigConstant.{
  DEFAULT_CSV_HEADER,
  DEFAULT_CSV_SEPARATOR,
  DEFAULT_READ_PARTITION
}
import com.vesoft.nebula.common.configuration.ConfigUtil.getOrElse
import com.vesoft.nebula.common.configuration.{
  Configs,
  ConfigsResolve,
  DataPreProcessConfig,
  DataSourceConfigEntry,
  FileFormatCategory,
  SchemaConfig,
  SourceCategory
}
import org.apache.log4j.Logger

object OSSConfigsResolve {
  private[this] val LOG = Logger.getLogger(this.getClass)
  def parse(configPath: String): Configs = {
    //super.parse(configPath)
    null
  }

  def getSourceConfigEntry(category: SourceCategory.Value,
                           fileFormat: Option[FileFormatCategory.Value] = None,
                           sourceConfig: Config,
                           preProcessConfig: List[DataPreProcessConfig],
                           schemaConfigs: List[SchemaConfig]): DataSourceConfigEntry = {
    val path            = sourceConfig.getString("path")
    val readPartition   = getOrElse(Option(sourceConfig), "readPartition", DEFAULT_READ_PARTITION)
    val fileFormatValue = if (fileFormat.isDefined) fileFormat.get else FileFormatCategory.CSV
    val header          = getOrElse(Option(sourceConfig), "header", DEFAULT_CSV_HEADER)
    val separator       = getOrElse(Option(sourceConfig), "separator", DEFAULT_CSV_SEPARATOR)
    val endpoint        = sourceConfig.getString("ossEndpoint")
    val ossAccessKey    = sourceConfig.getString("ossAccessKey")
    val ossSecretKey    = sourceConfig.getString("ossSecretKey")
    val config = OSSSourceConfigEntry(category,
                                      readPartition,
                                      fileFormatValue,
                                      path,
                                      endpoint,
                                      ossAccessKey,
                                      ossSecretKey,
                                      separator,
                                      header,
                                      schemaConfigs,
                                      preProcessConfig)
    LOG.info(config.toString)
    config
  }
}
