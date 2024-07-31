
package com.vesoft.nebula.common.reader.jdbcreader

import com.vesoft.nebula.common.configuration.{
  DataBaseServerSourceConfigEntry,
  DataPreProcessConfig,
  SchemaConfig,
  SourceCategory
}
import org.apache.log4j.Logger

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
  */
case class JdbcSourceConfigEntry(override val category: SourceCategory.Value,
                                 override val readParallel: Int,
                                 override val statement: String,
                                 override val schemaConfigs: List[SchemaConfig],
                                 url: String,
                                 driver: String,
                                 user: String,
                                 passwd: String,
                                 table: String,
                                 partitionColumn: Option[String] = None,
                                 lowerBound: Option[Long] = None,
                                 upperBound: Option[Long] = None,
                                 fetchSize: Option[Long] = None,
                                 override val preProcessConfigs: List[DataPreProcessConfig])
    extends DataBaseServerSourceConfigEntry {
  private[this] val LOG = Logger.getLogger(this.getClass)

  def check(): Unit = {
    require(url != null && url.nonEmpty, "url for jdbc source cannot be null.")
    require(driver != null && driver.nonEmpty, "driver for jdbc source cannot be null.")
    require(user != null && user.nonEmpty, "user for jdbc source cannot be null.")
    require(passwd != null && passwd.nonEmpty, "passwd for jdbc source cannot be null.")

    require(table == null || statement == null,
            "statement and table cannot be used at the same time.")
    require(table != null || statement != null, "either sentence or table should be config.")

    if (statement != null && partitionColumn.isDefined) {
      LOG.warn("statement is configured, partitionColumn will be ignored.")
    }

    if (table != null && partitionColumn.isDefined) {
      require(
        lowerBound.isDefined && upperBound.isDefined,
        "partitionColumn is configured, lowerBound and upperBound must be configured, refer " +
          "https://spark.apache.org/docs/latest/sql-data-sources-jdbc.html#data-source-option"
      )
    }
  }

  override def toString: String =
    s"JdbcSourceConfigEntry{category:$category, " +
      s"readParallel:$readParallel, " +
      s"statement:$statement, " +
      s"url:$url, " +
      s"driver:$driver, " +
      s"user:$user, " +
      s"passwd:***, " +
      s"table:$table, " +
      s"partitionColumn:${if (partitionColumn.isDefined) partitionColumn.get}, " +
      s"lowerBound:${if (lowerBound.isDefined) lowerBound.get}," +
      s"upperBound:${if (upperBound.isDefined) upperBound.get}, " +
      s"fetchSize:${if (fetchSize.isDefined) fetchSize.get}, " +
      s"schemaConfigs:${schemaConfigs.map(_.toString)}" +
      s"preProcessConfigs:${preProcessConfigs.map(_.toString())}" +
      s" }"

  def saveToFile: String = {
    val space                    = "  "
    val doubleSpace              = "    "
    val tripleSpace              = "      "
    val (nodeString, edgeString) = saveSchema(schemaConfigs)
    val preProcessConfigString   = savePreProcessConfig(preProcessConfigs)

    val optionalColumn =
      if (partitionColumn.isDefined)
        doubleSpace + "partitionColumn: \"" + partitionColumn + "\""
      else
        ""
    val optionalLowerBound =
      if (lowerBound.isDefined) doubleSpace + "lowerBound: " + lowerBound else ""
    val optionalUpperBound =
      if (upperBound.isDefined) doubleSpace + "upperBound: " + upperBound else ""
    val optionalFetchSize = if (fetchSize.isDefined) doubleSpace + "fetchSize: " + fetchSize else ""

    s"""
       |$space{
       |${doubleSpace}type: \"${category.toString}\"
       |${doubleSpace}url: \"$url\"
       |${doubleSpace}driver: \"$driver\"
       |${doubleSpace}user: \"$user\"
       |${doubleSpace}passwd: \"$passwd\"
       |${doubleSpace}table: \"$table\"
       |${doubleSpace}statement: \"$statement\"
       |${doubleSpace}readParallel: $readParallel
       |$optionalColumn
       |$optionalLowerBound
       |$optionalUpperBound
       |$optionalFetchSize
       |${doubleSpace}pre_processes:{
       |${preProcessConfigString}
       |$tripleSpace }
       |${doubleSpace}nodetypes: [
       |${nodeString}
       |${doubleSpace}]
       |${doubleSpace}edgetypes: [
       |${edgeString}
       |${doubleSpace}]
       |$space}""".stripMargin
  }
}
