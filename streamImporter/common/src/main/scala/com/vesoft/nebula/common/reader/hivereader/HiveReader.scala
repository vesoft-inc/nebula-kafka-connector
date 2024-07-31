
package com.vesoft.nebula.common.reader.hivereader

import com.vesoft.nebula.common.configuration.DataSourceConfigEntry
import com.vesoft.nebula.common.reader.DataSourceReader
import com.vesoft.nebula.common.reader.hdfsreader.HdfsSourceConfigEntry
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.types.StructType

class HiveReader extends DataSourceReader {

  override def readData(spark: SparkSession,
                        datasourceConfig: DataSourceConfigEntry,
                        options: Map[String, String]): DataFrame = {
    val sourceConfig = datasourceConfig.asInstanceOf[HiveSourceConfigEntry]
    spark.sql(sourceConfig.statement)
  }
}
