/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration

import org.apache.log4j.Logger

sealed trait DataSourceConfigEntry {

  /** one datasource may have multi node type or edge type configs */
  def schemaConfigs: List[SchemaConfigs]

  /** */
  def category: SourceCategory.Value
}

/**
  * file source base config, suitable for CSV,JSON source
  * */
sealed trait FileDataSourceConfigEntry extends DataSourceConfigEntry {

  /** file system the file data located  */
  def fileSystem: FileSystemCategory.Value

  /** file data path */
  def path: String
}

/**
  * file source config， support for csv,json
  */
case class FileBaseSourceConfigEntry(override val category: SourceCategory.Value,
                                     override val fileSystem: FileSystemCategory.Value,
                                     override val path: String,
                                     delimiter: Option[String] = None,
                                     header: Option[Boolean] = None,
                                     override val schemaConfigs: List[SchemaConfigs])
    extends FileDataSourceConfigEntry {

  override def toString: String =
    s"FileBaseSourceConfigEntry{category:$category, fileSystem:$fileSystem, path:$path, " +
      s"delimiter:$delimiter, header:$header, schemaConfigs:${schemaConfigs.map(_.toString)}"
}

/**
  * database server source config
  */
sealed trait DataBaseServerSourceConfigEntry extends DataSourceConfigEntry {
  def sentence: String
}

/**
  * hive source config
  *
  * @param category source category
  * @param sentence source query statement
  * @param schemaConfigs list schema configs for this source data
  */
case class HiveSourceConfigEntry(override val category: SourceCategory.Value,
                                 override val sentence: String,
                                 override val schemaConfigs: List[SchemaConfigs])
    extends DataBaseServerSourceConfigEntry {

  require(sentence != null && !sentence.isEmpty, "sentence for hive source cannot be null.")
  override def toString: String =
    s"HiveSourceConfigEntry{" +
      s"category:$category, sentence:$sentence, schemaConfigs:${schemaConfigs.map(_.toString)}"

}

/**
  * jdbc source config
  *
  * @param category source category
  * @param sentence source query statement
  * @param schemaConfigs list schema configs for this source data
  * @param url jdbc url
  * @param driver jdbc driver
  * @param user user name
  * @param passwd password
  * @param table database table name
  * @param partitionColumn column name can be partitioned, the column datatype must be numeric, date, timestamp
  * @param lowerBound used with {@link partitionColumn}
  * @param upperBound used with {@link partitionColumn}
  * @param numPartitions the max Spark partition number when reading source data
  * @param fetchSize the data amount for each request from the source database
  *
  *
  */
case class JdbcSourceConfigEntry(override val category: SourceCategory.Value,
                                 override val sentence: String,
                                 override val schemaConfigs: List[SchemaConfigs],
                                 url: String,
                                 driver: String,
                                 user: String,
                                 passwd: String,
                                 table: String,
                                 partitionColumn: Option[String] = None,
                                 lowerBound: Option[Long] = None,
                                 upperBound: Option[Long] = None,
                                 numPartitions: Option[Long] = None,
                                 fetchSize: Option[Long] = None)
    extends DataBaseServerSourceConfigEntry {
  private[this] val LOG = Logger.getLogger(this.getClass)

  require(url != null && !url.isEmpty, "url for jdbc source cannot be null.")
  require(driver != null && !driver.isEmpty, "driver for jdbc source cannot be null.")
  require(user != null && !user.isEmpty, "user for jdbc source cannot be null.")
  require(passwd != null && !passwd.isEmpty, "passwd for jdbc source cannot be null.")

  require(table == null || sentence == null, "sentence and table cannot be used at the same time.")
  require(table != null || sentence != null, "either sentence or table should be config.")

  if (sentence != null && partitionColumn.isDefined) {
    LOG.warn("sentence is configured, partitionColumn will be ignored.")
  }

  override def toString: String =
    s"JdbcSourceConfigEntry{category:$category, sentence:$sentence, url:$url, driver:$driver, " +
      s"user:$user, passwd:***, table:$table, " +
      s"partitionColumn:${if (partitionColumn.isDefined) partitionColumn.get}, " +
      s"lowerBound:${if (lowerBound.isDefined) lowerBound.get}," +
      s"upperBound:${if (upperBound.isDefined) upperBound.get}, " +
      s"numPartitions:${if (numPartitions.isDefined) numPartitions.get}, " +
      s"fetchSize:${if (fetchSize.isDefined) fetchSize.get}, " +
      s"schemaConfigs:${schemaConfigs.map(_.toString)}" +
      s" }"
}
