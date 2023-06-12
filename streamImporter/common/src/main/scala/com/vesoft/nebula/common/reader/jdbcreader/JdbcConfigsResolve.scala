package com.vesoft.nebula.common.reader.jdbcreader

import com.typesafe.config.Config
import com.vesoft.nebula.common.configuration.ConfigConstant.DEFAULT_READ_PARTITION
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

/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

object JdbcConfigsResolve {
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
    val readPartition = getOrElse(Option(sourceConfig), "readPartition", DEFAULT_READ_PARTITION)
    val url           = sourceConfig.getString("url")
    val driver        = sourceConfig.getString("driver")
    val user          = sourceConfig.getString("user")
    val passwd        = sourceConfig.getString("passwd")
    val statement =
      if (sourceConfig.hasPath("statement")) sourceConfig.getString("statement") else null
    val table = if (sourceConfig.hasPath("table")) sourceConfig.getString("table") else null
    val prepareQuery =
      if (sourceConfig.hasPath("prepareQuery")) Option(sourceConfig.getString("prepareQuery"))
      else None
    val partitionColumn =
      if (sourceConfig.hasPath("partitionColumn"))
        Option(sourceConfig.getString("partitionColumn"))
      else None
    val lowerBound =
      if (sourceConfig.hasPath("lowerBound")) Option(sourceConfig.getLong("lowerBound"))
      else None
    val upperBound =
      if (sourceConfig.hasPath("upperBound")) Option(sourceConfig.getLong("upperBound"))
      else None
    val fetchSize =
      if (sourceConfig.hasPath("fetchSize")) Option(sourceConfig.getLong("fetchSize")) else None
    val config = JdbcSourceConfigEntry(
      category,
      readPartition,
      statement,
      schemaConfigs,
      url,
      driver,
      user,
      passwd,
      table,
      prepareQuery,
      partitionColumn,
      lowerBound,
      upperBound,
      fetchSize,
      preProcessConfig
    )
    LOG.info(config.toString)
    config
  }

}
