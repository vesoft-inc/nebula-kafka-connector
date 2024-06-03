/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import org.apache.commons.lang.StringUtils
import org.apache.spark.sql.catalyst.util.CaseInsensitiveMap

import java.util.Properties

class NebulaOptions(@transient val parameters: CaseInsensitiveMap[String]) extends Serializable {

  import NebulaOptions._

  def this(parameters: Map[String, String], operaType: OperaType.Value) =
    this(CaseInsensitiveMap(parameters))

  def this(hostAndPorts: String,
           graphName: String,
           dataType: String,
           label: String,
           parameters: Map[String, String]) = {
    this(
      CaseInsensitiveMap(
        parameters ++ Map(
          NebulaOptions.GRAPH_NAME -> graphName,
          NebulaOptions.TYPE -> dataType,
          NebulaOptions.LABEL -> label
        ))
    )
  }

  val operaType = OperaType.withName(parameters(OPERATE_TYPE))

  /**
   * Return property with all options
   */
  val asProperties: Properties = {
    val properties = new Properties()
    parameters.originalMap.foreach { case (k, v) => properties.setProperty(k, v) }
    properties
  }

  val timeout: Int =
    parameters.getOrElse(TIMEOUT, DEFAULT_CONNECTION_TIMEOUT_SECONDS).toString.toInt
  val executionRetry: Int =
    parameters.getOrElse(EXECUTION_RETRY, DEFAULT_EXECUTION_RETRY).toString.toInt
  val executionRetryInterval: Int =
    parameters.getOrElse(EXECUTION_RETRY_INTERVAL, DEFAULT_EXECUTION_RETRY_INTERVAL).toString.toInt
  val user: String = parameters.getOrElse(USER_NAME, DEFAULT_USER_NAME)
  val passwd: String = parameters.getOrElse(PASSWD, DEFAULT_PASSWD)
  val rateLimit: Long = parameters.getOrElse(RATE_LIMIT, DEFAULT_RATE_LIMIT).toString.toLong

  require(parameters.isDefinedAt(TYPE), s"Option '$TYPE' is required")
  val dataType: String = parameters(TYPE)
  require(
    DataTypeEnum.validDataType(dataType),
    s"Option '$TYPE' is illegal, it should be '${DataTypeEnum.NODE}' or '${DataTypeEnum.EDGE}' or `${DataTypeEnum.GQL}`")

  /** nebula common parameters */
  require(parameters.isDefinedAt(GRAPH_ADDRESS),
    s"option $GRAPH_ADDRESS is required and can not be blank")
  var graphAddress = parameters(GRAPH_ADDRESS)
  var graphName: String = _
  var label: String = _
  if (!dataType.equals(DataTypeEnum.GQL.toString)) {
    require(parameters.isDefinedAt(GRAPH_NAME) && StringUtils.isNotBlank(parameters(GRAPH_NAME)),
      s"Option '$GRAPH_NAME' is required and can not be blank")
    graphName = parameters(GRAPH_NAME)
    require(parameters.isDefinedAt(LABEL) && StringUtils.isNotBlank(parameters(LABEL)),
      s"Option '$LABEL' is required and can not be blank")
    label = parameters(LABEL)
  }

  var batchSize: Int = parameters.getOrElse(BATCH_SIZE, DEFAULT_BATCH_SIZE).toString.toInt

  /** write parameters */

  var srcPkField: String = _
  var dstPkField: String = _
  var srcPkAsProp: Boolean = _
  var dstPkAsProp: Boolean = _
  var writeMode: WriteMode.Value = _
  var disableWriteLog: Boolean = _

  if (operaType == OperaType.WRITE) {
    srcPkField = parameters.getOrElse(SRC_PK_FIELD, null)
    dstPkField = parameters.getOrElse(DST_PK_FIELD, null)
    srcPkAsProp = parameters.getOrElse(SRC_PK_AS_PROP, false).toString.toBoolean
    dstPkAsProp = parameters.getOrElse(DST_PK_AS_PROP, false).toString.toBoolean
    writeMode =
      WriteMode.withName(parameters.getOrElse(WRITE_MODE, DEFAULT_WRITE_MODE).toString.toLowerCase)
    disableWriteLog = parameters.getOrElse(DISABLE_WRITE_LOG, false).toString.toBoolean
  }

  /** read parameters */
  var partitionNums: Int = parameters.getOrElse(PARTITION_NUMBER, "1").toInt
  val returnCols = parameters.getOrElse(RETURN_COLS, null)

  def getReturnCols: List[String] = {
    if (returnCols.equals("$null")) {
      null
    } else if (returnCols.isEmpty) {
      List()
    } else {
      returnCols.split(",").toList
    }
  }

  /** gql parameters */
  var gql: String = parameters.getOrElse(GQL, null)
}

object NebulaOptions {

  /** nebula common config */
  val GRAPH_NAME: String = "graphName"
  val GRAPH_ADDRESS: String = "graphAddress"
  val TYPE: String = "type"
  val LABEL: String = "label"
  val BATCH_SIZE: String = "batchSize"

  /** connection config */
  val TIMEOUT: String = "timeout"
  val EXECUTION_RETRY: String = "executionRetry"
  val EXECUTION_RETRY_INTERVAL: String = "executionRetryInterval"
  val USER_NAME: String = "user"
  val PASSWD: String = "passwd"

  val OPERATE_TYPE: String = "operateType"

  /** gql config */
  val GQL: String = "gql"

  /** write config */
  val RATE_LIMIT: String = "rateLimit"
  val SRC_PK_FIELD = "srcPkField"
  val DST_PK_FIELD = "dstPkField"
  val SRC_PK_AS_PROP: String = "srcPkAsProp"
  val DST_PK_AS_PROP: String = "dstPkAsProp"
  val WRITE_MODE: String = "writeMode"
  val DISABLE_WRITE_LOG: String = "disableWriteLog"

  /** read config */
  val PARTITION_NUMBER: String = "partitionNumber"
  val RETURN_COLS: String = "returnCols"

  val DEFAULT_TIMEOUT_SECONDS: Int = 10
  val DEFAULT_CONNECTION_TIMEOUT_SECONDS: Int = 3
  val DEFAULT_CONNECTION_RETRY: Int = 3
  val DEFAULT_EXECUTION_RETRY: Int = 3
  val DEFAULT_EXECUTION_RETRY_INTERVAL: Int = 0
  val DEFAULT_USER_NAME: String = "root"
  val DEFAULT_PASSWD: String = "nebula"

  val DEFAULT_RATE_LIMIT: Long = 1024L
  val DEFAULT_BATCH_SIZE: Int = 1000

  val DEFAULT_WRITE_MODE = WriteMode.INSERT

  val EMPTY_STRING: String = ""
}
