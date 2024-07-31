
package com.vesoft.nebula.common.reader.ossreader

import com.vesoft.nebula.common.configuration.{
  DataSourceConfigEntry,
  FileFormatCategory,
  NonValueConfig
}
import com.vesoft.nebula.common.reader.DataSourceReader
import org.apache.spark.sql.{DataFrame, SparkSession}

import scala.collection.mutable

class OssReader extends DataSourceReader {

  override def readData(spark: SparkSession,
                        datasourceConfig: DataSourceConfigEntry,
                        options: Map[String, String]): DataFrame = {
    val sourceConfig                         = datasourceConfig.asInstanceOf[OSSSourceConfigEntry]
    val fileFormat                           = sourceConfig.fileFormat

    val df = fileFormat match {
      case FileFormatCategory.CSV =>
        spark.read
          .options(options)
          .option("header", sourceConfig.header)
          .option("separator", sourceConfig.separator)
          .csv(sourceConfig.path)
      case FileFormatCategory.JSON =>
        spark.read.options(options).json(sourceConfig.path)
    }
    df
  }
}
