
package com.vesoft.nebula.common.reader.hdfsreader

import com.vesoft.nebula.common.configuration.{
  DataPreProcessConfig,
  FileDataSourceConfigEntry,
  FileFormatCategory,
  SchemaConfig,
  SourceCategory
}

/**
  * HDFS file source config， support for csv,json
  * TODO support kerberos
  */
case class HdfsSourceConfigEntry(override val category: SourceCategory.Value,
                                 override val readParallel: Int,
                                 override val fileFormat: FileFormatCategory.Value,
                                 override val path: String,
                                 separator: String,
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
      s"separator:'$separator', " +
      s"header:$header, " +
      s"schemaConfigs:${schemaConfigs.map(_.toString)}, " +
      s"preProcessConfigs:${preProcessConfigs.map(_.toString())}}"

  def saveToFile: String = {
    val space                    = "  "
    val doubleSpace              = "    "
    val tripleSpace              = "      "
    val (nodeString, edgeString) = saveSchema(schemaConfigs)
    val preProcessConfigString   = savePreProcessConfig(preProcessConfigs)

    s"""
       |$space{
       |${doubleSpace}type: \"${category.toString}\"
       |${doubleSpace}format: \"${fileFormat.toString}\"
       |${doubleSpace}path: \"${path}\"
       |${doubleSpace}header: $header
       |${doubleSpace}separator: \"$separator\"
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
       |$space}\n
       |""".stripMargin
  }
}
