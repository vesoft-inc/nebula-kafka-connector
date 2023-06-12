/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.reader.hdfsreader

import com.typesafe.config.Config
import com.vesoft.nebula.common.configuration.ConfigConstant.{DEFAULT_CSV_HEADER, DEFAULT_CSV_SEPARATOR, DEFAULT_READ_PARTITION}
import com.vesoft.nebula.common.configuration.ConfigUtil.getOrElse
import com.vesoft.nebula.common.configuration.{Configs, ConfigsResolve, DataPreProcessConfig, DataSourceConfigEntry, FileFormatCategory, SchemaConfig, SourceCategory}
import org.apache.log4j.Logger

object HdfsConfigsResolve {
  private[this] val LOG = Logger.getLogger(this.getClass)
  def parse(configPath: String): Configs = {
    // super.parse(configPath)
    null
  }
  def getSourceConfigEntry(category: SourceCategory.Value,
                           fileFormat: Option[FileFormatCategory.Value] = None,
                           sourceConfig: Config,
                           preProcessConfig: List[DataPreProcessConfig],
                           schemaConfigs: List[SchemaConfig]): DataSourceConfigEntry = {

    val path            = sourceConfig.getString("path")
    val fileFormatValue = if (fileFormat.isDefined) fileFormat.get else FileFormatCategory.CSV
    val header          = getOrElse(Option(sourceConfig), "header", DEFAULT_CSV_HEADER)
    val separator       = getOrElse(Option(sourceConfig), "separator", DEFAULT_CSV_SEPARATOR)
    val readPartition   = getOrElse(Option(sourceConfig), "readParallel", DEFAULT_READ_PARTITION)
    val config = HdfsSourceConfigEntry(category,
                          readPartition,
                          fileFormatValue,
                          path,
                          separator,
                          header,
                          schemaConfigs,
                          preProcessConfig)
    LOG.info(config.toString)
    config
  }
}
