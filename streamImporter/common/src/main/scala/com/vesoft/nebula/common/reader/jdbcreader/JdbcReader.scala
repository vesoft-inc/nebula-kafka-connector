package com.vesoft.nebula.common.reader.jdbcreader

import com.vesoft.nebula.common.configuration.DataSourceConfigEntry
import com.vesoft.nebula.common.reader.DataSourceReader
import org.apache.spark.sql.execution.datasources.jdbc.JDBCOptions.{JDBC_BATCH_FETCH_SIZE, JDBC_DRIVER_CLASS, JDBC_LOWER_BOUND, JDBC_NUM_PARTITIONS, JDBC_PARTITION_COLUMN, JDBC_QUERY_STRING, JDBC_TABLE_NAME, JDBC_UPPER_BOUND, JDBC_URL}
import org.apache.spark.sql.{DataFrame, SparkSession}


class JdbcReader extends DataSourceReader {

  override def readData(spark: SparkSession,
                        datasourceConfig: DataSourceConfigEntry,
                        options: Map[String, String]): DataFrame = {
    val sourceConfig = datasourceConfig.asInstanceOf[JdbcSourceConfigEntry]

    options + (JDBC_DRIVER_CLASS -> sourceConfig.driver)
    options + (JDBC_URL-> sourceConfig.url)
    options + ("user"-> sourceConfig.user)
    options + ("password"-> sourceConfig.passwd)

    if (sourceConfig.fetchSize.isDefined) {
      options + (JDBC_BATCH_FETCH_SIZE-> sourceConfig.fetchSize.get)
    }

    if (sourceConfig.table != null) {
      options + (JDBC_TABLE_NAME -> sourceConfig.table)
      if (sourceConfig.partitionColumn.isDefined) {
        options + (JDBC_PARTITION_COLUMN -> sourceConfig.partitionColumn)
        options + (JDBC_UPPER_BOUND -> sourceConfig.upperBound)
        options + (JDBC_LOWER_BOUND -> sourceConfig.lowerBound)
        options +(JDBC_NUM_PARTITIONS-> sourceConfig.readParallel)
      }
    } else {
      options + (JDBC_QUERY_STRING-> sourceConfig.statement)
    }
    spark.read
      .format("jdbc")
      .options(options)
      .load()
  }
}
