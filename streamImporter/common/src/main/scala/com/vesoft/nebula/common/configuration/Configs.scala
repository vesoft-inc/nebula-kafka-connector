
package com.vesoft.nebula.common.configuration

import com.vesoft.nebula.common.reader.hdfsreader.HdfsSourceConfigEntry
import com.vesoft.nebula.common.reader.hivereader.HiveSourceConfigEntry
import com.vesoft.nebula.common.reader.jdbcreader.JdbcSourceConfigEntry
import com.vesoft.nebula.common.reader.ossreader.OSSSourceConfigEntry
import com.vesoft.nebula.common.reader.s3reader.S3SourceConfigEntry

case class Configs(nebulaGraphConfigEntry: NebulaGraphConfigEntry,
                   mqClusterConfigEntry: MQClusterConfigEntry,
                   errorConfigEntry: ErrorConfigEntry,
                   sourceConfigEntrys: List[DataSourceConfigEntry]) {
  override def toString: String = {
    s"Configs{$nebulaGraphConfigEntry, " +
      s"$mqClusterConfigEntry, " +
      s"$errorConfigEntry, " +
      s"${sourceConfigEntrys.map(_.toString)}}"
  }

  def saveToFile: String = {
    val tabSpace = "  "
    s"{\n ${tabSpace}${nebulaGraphConfigEntry.saveSchema} \n" +
      s"\n${tabSpace}${mqClusterConfigEntry.saveSchema} \n" +
      s"\n${tabSpace}${errorConfigEntry.saveSchema} \n" +
      s"\n${tabSpace}sources:[" +
      s"\n${saveSources(sourceConfigEntrys)}" +
      s"\n${tabSpace}]" +
      s"\n}"
  }

  private[this] def saveSources(sourceConfigs: List[DataSourceConfigEntry]): String = {
    val sourceConfigsString: StringBuilder = new StringBuilder()
    for (sourceConfig <- sourceConfigs) {
      sourceConfig match {
        case hdfsSourceConfigEntry: HdfsSourceConfigEntry =>
          sourceConfigsString.append(hdfsSourceConfigEntry.saveToFile)
        case s3SourceConfigEntry: S3SourceConfigEntry =>
          sourceConfigsString.append(s3SourceConfigEntry.saveToFile)
        case ossSourceConfigEntry: OSSSourceConfigEntry =>
          sourceConfigsString.append(ossSourceConfigEntry.saveToFile)
        case hiveSourceConfigEntry: HiveSourceConfigEntry =>
          sourceConfigsString.append(hiveSourceConfigEntry.saveToFile)
        case jdbcSourceConfigEntry: JdbcSourceConfigEntry =>
          sourceConfigsString.append(jdbcSourceConfigEntry.saveToFile)
      }
    }
    sourceConfigsString.toString()
  }
}

/**
  * NebulaGraph config
  * used to query the specific graph's schema information
  */
case class NebulaGraphConfigEntry(graphAddress: String,
                                  graphName: String,
                                  user: String,
                                  passwd: String,
                                  connectTimeout: Int,
                                  requestTimeout: Int,
                                  retryIntervalTime: Int,
                                  mode: SinkCategory.Value,
                                  generateDDL: Boolean = false) {
  def check(): Unit = {
    require(graphAddress != null && graphAddress.nonEmpty, "graphAddr cannot be null")
    require(graphName != null && graphName.nonEmpty, "graph name cannot be null")
    require(user != null && user.nonEmpty, "NebulaGraph user cannot be null")
    require(passwd != null && passwd.nonEmpty, "NebulaGraph passwd cannot be null")
    require((new ValidateUtil).validateServer(graphAddress), "graph address is not valid")
    require(connectTimeout >= 0, "graph connect timeout cannot be less than 0")
    require(requestTimeout >= 0, "graph request timeout cannot be less than 0")
    require(retryIntervalTime >= 0, "graph interval time between retrys cannot less than 0")
  }
  override def toString: String =
    s"NebulaGraphConfigEntry{graphAddress:$graphAddress, graph:$graphName, user:$user, passwd:****}"

  def saveSchema: String = {
    val space       = "  "
    val doubleSpace = "    "
    s"""
       |${space}nebula:{
       |${doubleSpace}graphAddr: \"${graphAddress}\"
       |${doubleSpace}graphName: \"${graphName}\"
       |${doubleSpace}user: \"${user}\"
       |${doubleSpace}passwd: \"${passwd}\"
       |${doubleSpace}mode: \"${mode}\"
       |${doubleSpace}connectTimeout: $connectTimeout
       |${doubleSpace}requestTimeout: $requestTimeout
       |${doubleSpace}retryIntervalTime: $retryIntervalTime
       |$space}
       |""".stripMargin
  }
}

/**
  * RedPanda MQ config
  *
  */
case class MQClusterConfigEntry(server: String, topic: String, replic: Int) {
  def check(): Unit = {
    require((new ValidateUtil).validateServer(server), "mq server address is not valid")
    require(topic != null && topic.nonEmpty, "mq topic cannot be null")
  }
  override def toString: String =
    s"MQClusterConfigEntry{server:$server, topic:$topic, replic=$replic}"

  def saveSchema: String = {
    val space       = "  "
    val doubleSpace = "    "
    s"""
       |${space}mq:{
       |${doubleSpace}server: \"${server}\"
       |${doubleSpace}topic: \"${topic}\"
       |${doubleSpace}replic: $replic
       |$space}
       |""".stripMargin
  }
}

case class ErrorConfigEntry(path: String, maxRecords: Int) {
  def check(): Unit = {
    require(maxRecords >= 0, "maxRecords can not be less than 0")
  }
  override def toString: String = s"ErrorConfigEntry{path:$path, maxRecords:$maxRecords}"

  def saveSchema: String = {
    val space       = "  "
    val doubleSpace = "    "
    s"""
       |${space}error:{
       |${doubleSpace}path: \"${path}\"
       |${doubleSpace}maxRecords: $maxRecords
       |$space}
       |""".stripMargin
  }
}
