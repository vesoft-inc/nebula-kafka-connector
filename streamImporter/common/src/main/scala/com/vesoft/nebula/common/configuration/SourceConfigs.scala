/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration

/**
  * data source config
  * any data source should extends either {@link FileDataSourceConfigEntry}
  * or {@link DataBaseServerSourceConfigEntry} or {@link StreamingSourceConfigEntry}
  *
  * */
sealed trait DataSourceConfigEntry {

  /** read partition number */
  def readParallel: Int

  /** one datasource may have multi node type or edge type configs */
  def schemaConfigs: List[SchemaConfig]

  /** source category */
  def category: SourceCategory.Value

  /** source data pre-process roles */
  def preProcessConfigs: List[DataPreProcessConfig]

  def check()
}

/**
  * file source base config, suitable for CSV,JSON source
  * */
trait FileDataSourceConfigEntry extends DataSourceConfigEntry {

  /** file format */
  def fileFormat: FileFormatCategory.Value

  /** file data path */
  def path: String
}

/**
  * database server source config
  */
trait DataBaseServerSourceConfigEntry extends DataSourceConfigEntry {
  def statement: String
}

/**
  * TODO add common config for streaming data source
  */
trait StreamingSourceConfigEntry extends DataSourceConfigEntry {}
