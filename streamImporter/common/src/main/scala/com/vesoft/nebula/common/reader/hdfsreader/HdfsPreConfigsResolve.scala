
package com.vesoft.nebula.common.reader.hdfsreader

import com.typesafe.config.Config
import com.vesoft.nebula.common.configuration.ConfigConstant.{DEFAULT_CSV_HEADER, DEFAULT_CSV_SEPARATOR, DEFAULT_READ_PARTITION}
import com.vesoft.nebula.common.configuration.ConfigUtil.getOrElse
import com.vesoft.nebula.common.configuration.preConfig.PreConfigResolve
import com.vesoft.nebula.common.configuration.{DataPreProcessConfig, DataSourceConfigEntry, FileFormatCategory, SchemaConfig, SourceCategory}

class HdfsPreConfigsResolve {

  def getSourceConfigEntry(category: SourceCategory.Value,
                                    fileFormat: Option[FileFormatCategory.Value],
                                    sourceConfig: Config,
                                    preProcessConfig: List[DataPreProcessConfig],
                                    schemaConfigs: List[SchemaConfig]): DataSourceConfigEntry = {
    val path            = sourceConfig.getString("path")
    val fileFormatValue = if (fileFormat.isDefined) fileFormat.get else FileFormatCategory.CSV
    val header          = getOrElse(Option(sourceConfig), "header", DEFAULT_CSV_HEADER)
    val separator       = getOrElse(Option(sourceConfig), "separator", DEFAULT_CSV_SEPARATOR)
    val readPartition   = getOrElse(Option(sourceConfig), "readPartition", DEFAULT_READ_PARTITION)
    HdfsSourceConfigEntry(category,
                          readPartition,
                          fileFormatValue,
                          path,
                          separator,
                          header,
                          schemaConfigs,
                          preProcessConfig)
  }

}
