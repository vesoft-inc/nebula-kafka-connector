/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.reader.s3reader

import com.vesoft.nebula.common.configuration.{
  DataSourceConfigEntry,
  FileFormatCategory,
  NonValueConfig
}
import com.vesoft.nebula.common.reader.DataSourceReader
import com.vesoft.nebula.common.reader.ossreader.OSSSourceConfigEntry
import org.apache.spark.sql.{DataFrame, SparkSession}

import scala.collection.mutable

class S3Reader extends DataSourceReader {

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
