/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.reader.hdfsreader

import com.vesoft.nebula.common.configuration.{DataPreProcessConfig, FileDataSourceConfigEntry, FileFormatCategory, SchemaConfig, SourceCategory}

/**
  * HDFS file source config， support for csv,json
  * TODO support kerberos
  */
case class HdfsSourceConfigEntry(override val category: SourceCategory.Value,
                                 override val readParallel: Int,
                                 override val fileFormat: FileFormatCategory.Value,
                                 override val path: String,
                                 separator: String = null,
                                 header: Boolean = false,
                                 override val schemaConfigs: List[SchemaConfig],
                                 override val preProcessConfigs: List[DataPreProcessConfig])
    extends FileDataSourceConfigEntry {

  override def check(): Unit = {
    require(path != null && path.nonEmpty, "file path cannot be null")
  }

  override def toString: String =
    s"FileBaseSourceConfigEntry{category:$category, " +
      s"readParallel:$readParallel, " +
      s"fileFormat:$fileFormat, " +
      s"path:$path, " +
      s"separator:$separator, " +
      s"header:$header, " +
      s"schemaConfigs:${schemaConfigs.map(_.toString)}, " +
      s"preProcessConfigs:${preProcessConfigs.map(_.toString())}}"
}
