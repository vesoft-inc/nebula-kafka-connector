/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration

class Category {}

object SourceCategory extends Enumeration {
  type Type = Value
  val CSV  = Value("CSV")
  val JSON = Value("JSON")
  val HIVE = Value("HIVE")
  val JDBC = Value("JDBC")
}

object FileSystemCategory extends Enumeration {
  type Type = Value
  val HDFS = Value("HDFS")
  val S3   = Value("S3")
  val OSS  = Value("OSS")
}
