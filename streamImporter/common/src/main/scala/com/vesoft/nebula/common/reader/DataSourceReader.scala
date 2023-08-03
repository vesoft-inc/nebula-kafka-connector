/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.reader

import com.vesoft.nebula.common.configuration.{
  ConcatConfig,
  DataPreProcessConfig,
  DataSourceConfigEntry,
  FilterConfig,
  NebulaDataType,
  NonValueConfig,
  SeparatorConfig
}
import org.apache.parquet.format.TimeType
import org.apache.spark.sql.functions.{col, concat_ws, split}
import org.apache.spark.sql.types.{
  BooleanType,
  ByteType,
  DataType,
  DateType,
  DoubleType,
  FloatType,
  IntegerType,
  LongType,
  ShortType,
  StringType
}
import org.apache.spark.sql.{Column, DataFrame, SparkSession, functions}
import org.joda.time.DateTimeFieldType

import scala.collection.mutable
import scala.collection.mutable.ListBuffer

trait DataSourceReader {

  /**
    * get the source data's schema and convert it to NebulaGraph Schema
    */
  def readSchema(spark: SparkSession,
                 datasourceConfig: DataSourceConfigEntry): Map[String, String] = {

    val options: mutable.Map[String, String] = new mutable.HashMap[String, String]()
    var nonValue: NonValueConfig             = null
    datasourceConfig.preProcessConfigs.foreach {
      case config1: NonValueConfig => nonValue = config1
      case _                       =>
    }
    if (nonValue != null) {
      options += ("nullValue" -> nonValue.value)
    }
    options += ("inferSchema" -> "true")

    val fields = readData(spark, datasourceConfig, options.toMap).schema.fields
    val schema = fields.map(f => f.name -> convertDataType2NebulaType(f.dataType)).toMap
    schema
  }

  def getDataFrame(spark: SparkSession, dataSourceConfigEntry: DataSourceConfigEntry): DataFrame = {
    var df                = readData(spark, dataSourceConfigEntry, Map())
    val preProcessConfigs = dataSourceConfigEntry.preProcessConfigs
    for (preProcessConfig <- preProcessConfigs) {
      df = preProcessDF(df, preProcessConfig)
    }
    df
  }

  def readData(spark: SparkSession,
               datasourceConfig: DataSourceConfigEntry,
               options: Map[String, String]): DataFrame

  /**
    * pre-process data source
    * 1. concat multiple fields into one new field.
    * 2. separate one field into multiple new fields.
    * 3. filter
    * 4. convert specific value as NULL
    */
  private[this] def preProcessDF(df: DataFrame,
                                 preProcessConfig: DataPreProcessConfig): DataFrame = {
    var data = df
    preProcessConfig match {
      // concat multiple fields into one field
      case concatConfig: ConcatConfig =>
        val finalColNames: ListBuffer[Column] = new ListBuffer[Column]
        for (field <- df.schema.fieldNames.toList) {
          finalColNames.append(col(field))
        }
        finalColNames.append(
          concat_ws(concatConfig.sep, concatConfig.oldFields.map(c => col(c)): _*)
            .cast(StringType)
            .as(concatConfig.newFiled))
        data = data.select(finalColNames: _*)
      // separate one field into multiple fields
      case separatorConfig: SeparatorConfig =>
        val finalColNames: ListBuffer[Column] = new ListBuffer[Column]
        for (field <- df.schema.fieldNames.toList) {
          finalColNames.append(col(field))
        }
        val splitCols = functions.split(col(separatorConfig.oldField), separatorConfig.sep)
        for (i <- 0 to separatorConfig.newFields.size) {
          val newField    = separatorConfig.newFields(i)
          val splitColumn = splitCols.getItem(i)
          data = data.withColumn(newField, splitColumn)
        }
      // filter
      case filterConfig: FilterConfig =>
        for (filter <- filterConfig.conditions) {
          data = data.filter(filter)
        }
      case nonValueConfig: NonValueConfig => data = data
    }
    data
  }

  private[this] def convertDataType2NebulaType(dataType: DataType): String = {
    val nebulaType = dataType match {
      case StringType  => NebulaDataType.STRING
      case ByteType    => NebulaDataType.INT8
      case ShortType   => NebulaDataType.INT16
      case IntegerType => NebulaDataType.INT32
      case LongType    => NebulaDataType.ITN64
      case FloatType   => NebulaDataType.FLOAT
      case DoubleType  => NebulaDataType.DOUBLE
      case BooleanType => NebulaDataType.BOOL
      case DateType    => NebulaDataType.DATE
    }
    nebulaType.toString
  }
}
