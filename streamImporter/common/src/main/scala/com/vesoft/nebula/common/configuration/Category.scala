
package com.vesoft.nebula.common.configuration

/**
  * data source supported by import tool
  */
object SourceCategory extends Enumeration {
  type Type = Value
  val HDFS = Value("HDFS")
  val S3   = Value("S3")
  val OSS  = Value("OSS")
  val HIVE = Value("HIVE")
  val JDBC = Value("JDBC")
}

/**
  * file system supported by import tool
  */
object FileFormatCategory extends Enumeration {
  type Type = Value
  val CSV  = Value("CSV")
  val JSON = Value("JSON")
}

/**
  * two different data importing way
  * IMPORT: data will be construct into the gql statement and will be processed by Graphd server
  * BULKLOAD: data will be construct into struct and will be processed directly by Storaged server
  */
object SinkCategory extends Enumeration {
  type Type = Value
  val IMPORT   = Value("IMPORT")
  val BULKLOAD = Value("BULKLOAD")
}

object NebulaDataType extends Enumeration {
  type Type = Value
  val STRING   = Value("string")
  val INT8     = Value("int8")
  val INT16    = Value("int16")
  val INT32    = Value("int32")
  val ITN64    = Value("int64")
  val BOOL     = Value("bool")
  val FLOAT    = Value("float")
  val DOUBLE   = Value("double")
  val DATE     = Value("date")
  val TIME     = Value("localTime")
  val DATETIME = Value("localDatetime")
  val DURATION = Value("duration")

}

object IDType extends Enumeration {
  type Type = Value
  val STRING = Value("STRING")
  val INT    = Value("INT")
}
