/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration

case class Configs(nebulaGraphConfigEntry: NebulaGraphConfigEntry,
                   MQClusterConfigEntry: MQClusterConfigEntry,
                   errorConfigEntry: ErrorConfigEntry,
                   sourceConfigEntrys: List[DataSourceConfigEntry])

/**
  * NebulaGraph config
  * used to query the specific graph's schema information
  */
case class NebulaGraphConfigEntry(graphAddress: String,
                                  graphName: String,
                                  user: String,
                                  passwd: String,
                                  connectTimeout:Int,
                                  requestTimeout:Int,
                                  retryIntervalTime:Int,
                                  mode: SinkCategory.Value,
                                  generateDDL: Boolean = false) {
  def check(): Unit = {
    require(graphAddress != null && graphAddress.nonEmpty, "graphAddr cannot be null")
    require(graphName != null && graphName.nonEmpty, "graph name cannot be null")
    require(user != null && user.nonEmpty, "NebulaGraph user cannot be null")
    require(passwd != null && passwd.nonEmpty, "NebulaGraph passwd cannot be null")
    require((new ValidateUtil).validateServer(graphAddress), "graph address is not valid")
    require(connectTimeout >=0, "graph connect timeout cannot be less than 0")
    require(requestTimeout >=0, "graph request timeout cannot be less than 0")
    require(retryIntervalTime >=0, "graph interval time between retrys cannot less than 0")
  }
  override def toString: String =
    s"NebulaGraphConfigEntry{graphAddress:$graphAddress, graph:$graphName, user:$user, passwd:****}"
}

/**
  * RedPanda MQ config
  *
  */
case class MQClusterConfigEntry(server: String, topic: String) {
  def check(): Unit = {
    require((new ValidateUtil).validateServer(server), "mq server address is not valid")
    require(topic != null && topic.nonEmpty, "mq topic cannot be null")
  }
  override def toString: String = s"MQClusterConfigEntry{server:$server, topic:$topic}"
}

case class ErrorConfigEntry(path: String, maxRecords: Int) {
  def check(): Unit = {
    require(maxRecords >=0, "maxRecords can not be less than 0")
  }
  override def toString: String = s"ErrorConfigEntry{path:$path, maxRecords:$maxRecords}"
}
