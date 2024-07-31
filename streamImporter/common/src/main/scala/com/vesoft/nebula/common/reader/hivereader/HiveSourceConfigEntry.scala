
package com.vesoft.nebula.common.reader.hivereader

import com.vesoft.nebula.common.configuration.{
  DataBaseServerSourceConfigEntry,
  DataPreProcessConfig,
  SchemaConfig,
  SourceCategory
}

/**
  * hive source config
  *
  * @param category source category
  * @param sentence source query statement
  * @param schemaConfigs list schema configs for this source data
  * @param preProcessConfigs list preProcessConfig for processing the source data
  */
case class HiveSourceConfigEntry(override val category: SourceCategory.Value,
                                 override val readParallel: Int,
                                 override val statement: String,
                                 uris: String = null,
                                 override val schemaConfigs: List[SchemaConfig],
                                 override val preProcessConfigs: List[DataPreProcessConfig])
    extends DataBaseServerSourceConfigEntry {

  def check(): Unit = {
    require(statement != null && statement.nonEmpty, "statement for hive source cannot be null.")
  }

  override def toString: String =
    s"HiveSourceConfigEntry{category:$category, " +
      s"readParallel:$readParallel, " +
      s"statement:$statement, " +
      s"uris:$uris, " +
      s"schemaConfigs:${schemaConfigs.map(_.toString)}, " +
      s"preProcessConfigs:${preProcessConfigs.map(_.toString())}}"

  def saveToFile: String = {
    val space                    = "  "
    val doubleSpace              = "    "
    val tripleSpace              = "      "
    val (nodeString, edgeString) = saveSchema(schemaConfigs)
    val preProcessConfigString   = savePreProcessConfig(preProcessConfigs)

    s"""$space{
       |${doubleSpace}type: \"${category.toString}\"
       |${doubleSpace}uris: \"$uris\"
       |${doubleSpace}statement: \"$statement\"
       |${doubleSpace}readParallel: $readParallel
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
