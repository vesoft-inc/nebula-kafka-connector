/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.reader.hivereader

import com.typesafe.config.Config
import com.vesoft.nebula.common.configuration.ConfigConstant.DEFAULT_READ_PARTITION
import com.vesoft.nebula.common.configuration.ConfigUtil.getOrElse
import com.vesoft.nebula.common.configuration.{DataPreProcessConfig, DataSourceConfigEntry, FileFormatCategory, SchemaConfig, SourceCategory}

class HivePreConfigsResolve {

  def getSourceConfigEntry(category: SourceCategory.Value,
                                    fileFormat: Option[FileFormatCategory.Value],
                                    sourceConfig: Config,
                                    preProcessConfig: List[DataPreProcessConfig],
                                    schemaConfigs: List[SchemaConfig]): DataSourceConfigEntry = {
    val statement     = sourceConfig.getString("statement")
    val readPartition = getOrElse(Option(sourceConfig), "readPartition", DEFAULT_READ_PARTITION)
    val uris          = if (sourceConfig.hasPath("uris")) sourceConfig.getString("uris") else null
    HiveSourceConfigEntry(category, readPartition, statement, uris, schemaConfigs, preProcessConfig)
  }

}
