
package com.vesoft.nebula.common.configuration

import scala.collection.mutable.ListBuffer

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

  def savePreProcessConfig(preProcessConfigs: List[DataPreProcessConfig]): String = {
    val concatConfigs: ListBuffer[ConcatConfig]       = new ListBuffer[ConcatConfig]
    val separatorConfigs: ListBuffer[SeparatorConfig] = new ListBuffer[SeparatorConfig]
    val filterConfigs: ListBuffer[FilterConfig]       = new ListBuffer[FilterConfig]
    var nonValueConf: NonValueConfig                  = NonValueConfig(null)

    for (preProcessConfig <- preProcessConfigs) {
      preProcessConfig match {
        case concatConfig: ConcatConfig       => concatConfigs += concatConfig
        case separatorConfig: SeparatorConfig => separatorConfigs += separatorConfig
        case filterConfig: FilterConfig       => filterConfigs += filterConfig
        case nonValueConfig: NonValueConfig   => nonValueConf = nonValueConfig
      }
    }
    val preProcessConfigString: StringBuilder = new StringBuilder()
    val space                                 = "         "

    preProcessConfigString.append(s"${space}concat:[")

    for (concatConfig <- concatConfigs) {
      preProcessConfigString.append(concatConfig.saveSchema())
      preProcessConfigString.append("\n")
    }
    preProcessConfigString.append(s"\n${space}]")

    preProcessConfigString.append(s"\n${space}separate:[")
    for (separatorConfig <- separatorConfigs) {
      preProcessConfigString.append(separatorConfig.saveSchema())
      preProcessConfigString.append("\n")
    }
    preProcessConfigString.append(s"\n$space]")

    for (filterConfig <- filterConfigs) {
      preProcessConfigString.append(filterConfig.saveSchema())
    }
    preProcessConfigString.append(nonValueConf.saveSchema())
    preProcessConfigString.toString()
  }

  def saveSchema(schemas: List[SchemaConfig]): (String, String) = {
    val nodeSchemas: ListBuffer[NodeConfig] = new ListBuffer[NodeConfig]()
    val edgeSchemas: ListBuffer[EdgeConfig] = new ListBuffer[EdgeConfig]()
    for (schema <- schemas) {
      schema match {
        case nodeConfig: NodeConfig => nodeSchemas += nodeConfig
        case edgeConfig: EdgeConfig => edgeSchemas += edgeConfig
      }
    }

    val space                            = "        "
    val doubleSpace                      = "           "
    val nodeSchemasString: StringBuilder = new StringBuilder()
    for (node <- nodeSchemas) {
      nodeSchemasString.append(s"""
                                  |${space}{
                                  |${doubleSpace}name: \"${node.name}\"
                                  |${doubleSpace}primaryKey: \"${node.vid}\"
                                  |${doubleSpace}sourceFields: [${node.sourceFields.mkString(",")}]
                                  |${doubleSpace}nebulaFields: [${node.nebulaFields.mkString(",")}]
                                  |${doubleSpace}batchSize: ${node.batchSize}
                                  |${doubleSpace}writeParallel: ${node.partition}
                                  |$space}
                                  |""".stripMargin)
    }

    val edgeSchemasString: StringBuilder = new StringBuilder()
    for (edge <- edgeSchemas) {
      edgeSchemasString.append(s"""
                                  |$space{
                                  |${doubleSpace}name: \"${edge.name}\"
                                  |${doubleSpace}sourceKey: \"${edge.src}\"
                                  |${doubleSpace}targetKey: \"${edge.dst}\"
                                  |${doubleSpace}sourceFields: [${edge.sourceFields.mkString(",")}]
                                  |${doubleSpace}nebulaFields: [${edge.nebulaFields.mkString(",")}]
                                  |${doubleSpace}batchSize: ${edge.batchSize}
                                  |${doubleSpace}writeParallel: ${edge.partition}
                                  |${space}}""".stripMargin)
    }
    (nodeSchemasString.toString(), edgeSchemasString.toString)
  }
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
