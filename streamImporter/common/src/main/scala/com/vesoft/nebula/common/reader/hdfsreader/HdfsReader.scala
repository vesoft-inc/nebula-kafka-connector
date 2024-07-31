
package com.vesoft.nebula.common.reader.hdfsreader

import com.vesoft.nebula.common.configuration.{
  DataSourceConfigEntry,
  FileFormatCategory,
  NebulaDataType,
  NonValueConfig
}
import com.vesoft.nebula.common.reader.DataSourceReader
import org.apache.parquet.format.TimeType
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.types.{
  BooleanType,
  ByteType,
  CharType,
  DataType,
  DateType,
  DoubleType,
  FloatType,
  IntegerType,
  LongType,
  ShortType,
  StringType,
  StructField,
  StructType
}
import org.joda.time.DateTimeFieldType

import scala.collection.mutable

class HdfsReader extends DataSourceReader {

  override def readData(spark: SparkSession, datasourceConfig: DataSourceConfigEntry,options: Map[String, String]): DataFrame = {
    val sourceConfig                         = datasourceConfig.asInstanceOf[HdfsSourceConfigEntry]
    val fileFormat                           = sourceConfig.fileFormat
    val options: mutable.Map[String, String] = new mutable.HashMap[String, String]()

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
