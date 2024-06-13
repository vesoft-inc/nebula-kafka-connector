/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula

import com.vesoft.nebula.spark.common.utils.SparkValidate
import com.vesoft.nebula.spark.common.{DataTypeEnum, GqlNebulaConfig, NebulaConnectionConfig, NebulaOptions, OperaType, ReadNebulaConfig, WriteNebulaConfig, WriteNebulaEdgeConfig, WriteNebulaNodeConfig}
import org.apache.spark.sql.{DataFrame, DataFrameReader, DataFrameWriter, Row, SaveMode}

package object connector {

  /**
   * spark writer for nebula graph
   */
  implicit class NebulaDataFrameWriter(writer: DataFrameWriter[Row]) {

    private var connectionConfig: NebulaConnectionConfig = _
    private var writeNebulaConfig: WriteNebulaConfig = _

    /**
     * config nebula connection
     *
     * @param connectionConfig  connection parameters
     * @param writeNebulaConfig write parameters for node or edge
     */
    def nebula(connectionConfig: NebulaConnectionConfig,
               writeNebulaConfig: WriteNebulaConfig): NebulaDataFrameWriter = {
      SparkValidate.validate("2.4.*")
      this.connectionConfig = connectionConfig
      this.writeNebulaConfig = writeNebulaConfig
      this
    }

    /**
     * write dataframe into nebula node type
     */
    def writeVertices(): Unit = {
      assert(connectionConfig != null && writeNebulaConfig != null,
        "nebula config is not set, please call nebula() before writeVertices")
      val writeConfig = writeNebulaConfig.asInstanceOf[WriteNebulaNodeConfig]
      val dfWriter = writer
        .format(classOf[NebulaDataSource].getName)
        .mode(SaveMode.Overwrite)
        .option(NebulaOptions.TYPE, DataTypeEnum.NODE.toString)
        .option(NebulaOptions.OPERATE_TYPE, OperaType.WRITE.toString)
        .option(NebulaOptions.USER_NAME, connectionConfig.getUser)
        .option(NebulaOptions.AUTHOPTIONS, connectionConfig.getAuthOptions)
        .option(NebulaOptions.GRAPH_ADDRESS, connectionConfig.getGraphAddress)
        .option(NebulaOptions.TIMEOUT, connectionConfig.getTimeout)
        .option(NebulaOptions.EXECUTION_RETRY, connectionConfig.getExecRetry)
        .option(NebulaOptions.EXECUTION_RETRY_INTERVAL, connectionConfig.getExecRetryIntervalMs)
        .option(NebulaOptions.GRAPH_NAME, writeConfig.getGraphName)
        .option(NebulaOptions.LABEL, writeConfig.getNodeType)
        .option(NebulaOptions.BATCH_SIZE, writeConfig.getBatchSize)
        .option(NebulaOptions.WRITE_MODE, writeConfig.getWriteMode)
        .option(NebulaOptions.DISABLE_WRITE_LOG, writeConfig.isDisableWriteLog)
      dfWriter.save()
    }

    /**
     * write dataframe into nebula edge type
     */
    def writeEdges(): Unit = {

      assert(connectionConfig != null && writeNebulaConfig != null,
        "nebula config is not set, please call nebula() before writeEdges")
      val writeConfig = writeNebulaConfig.asInstanceOf[WriteNebulaEdgeConfig]
      val dfWriter = writer
        .format(classOf[NebulaDataSource].getName)
        .mode(SaveMode.Overwrite)
        .option(NebulaOptions.TYPE, DataTypeEnum.EDGE.toString)
        .option(NebulaOptions.OPERATE_TYPE, OperaType.WRITE.toString)
        .option(NebulaOptions.USER_NAME, connectionConfig.getUser)
        .option(NebulaOptions.AUTHOPTIONS, connectionConfig.getAuthOptions)
        .option(NebulaOptions.GRAPH_ADDRESS, connectionConfig.getGraphAddress)
        .option(NebulaOptions.TIMEOUT, connectionConfig.getTimeout)
        .option(NebulaOptions.EXECUTION_RETRY, connectionConfig.getExecRetry)
        .option(NebulaOptions.EXECUTION_RETRY_INTERVAL, connectionConfig.getExecRetryIntervalMs)
        .option(NebulaOptions.GRAPH_NAME, writeConfig.getGraphName)
        .option(NebulaOptions.LABEL, writeConfig.getEdgeType)
        .option(NebulaOptions.SRC_PK_FIELD, writeConfig.getSrcPkFiled)
        .option(NebulaOptions.DST_PK_FIELD, writeConfig.getDstPkField)
        .option(NebulaOptions.BATCH_SIZE, writeConfig.getBatchSize)
        .option(NebulaOptions.SRC_PK_AS_PROP, writeConfig.getSrcAsProp)
        .option(NebulaOptions.DST_PK_AS_PROP, writeConfig.getDstAsProp)
        .option(NebulaOptions.WRITE_MODE, writeConfig.getWriteMode)
        .option(NebulaOptions.DISABLE_WRITE_LOG, writeConfig.isDisableWriteLog)
      dfWriter.save()
    }
  }

  /**
   * spark reader for nebula graph
   */
  implicit class NebulaDataFrameReader(reader: DataFrameReader) {
    var connectionConfig: NebulaConnectionConfig = _
    var readConfig: ReadNebulaConfig = _

    def nebula(connectionConfig: NebulaConnectionConfig, readConfig: ReadNebulaConfig): NebulaDataFrameReader = {
      SparkValidate.validate("2.4.*")
      this.connectionConfig = connectionConfig
      this.readConfig = readConfig
      this
    }

    /**
     * Reading nodes from NebulaGraph
     *
     * @return DataFrame
     */
    def loadNode(): DataFrame = {
      assert(connectionConfig != null && readConfig != null,
        "nebula config is not set, please call nebula() before loadVerticesToDF")
      val dfReader = reader
        .format(classOf[NebulaDataSource].getName)
        .option(NebulaOptions.TYPE, DataTypeEnum.NODE.toString)
        .option(NebulaOptions.OPERATE_TYPE, OperaType.READ.toString)
        .option(NebulaOptions.TIMEOUT, connectionConfig.getTimeout)
        .option(NebulaOptions.EXECUTION_RETRY, connectionConfig.getExecRetry)
        .option(NebulaOptions.EXECUTION_RETRY_INTERVAL, connectionConfig.getExecRetryIntervalMs)
        .option(NebulaOptions.GRAPH_ADDRESS, connectionConfig.getGraphAddress)
        .option(NebulaOptions.USER_NAME, connectionConfig.getUser)
        .option(NebulaOptions.AUTHOPTIONS, connectionConfig.getAuthOptions)
        .option(NebulaOptions.GRAPH_NAME, readConfig.getGraphName)
        .option(NebulaOptions.LABEL, readConfig.getTypeName)
        .option(NebulaOptions.RETURN_COLS, readConfig.getReturnColsString)
        .option(NebulaOptions.BATCH_SIZE, readConfig.getBatchSize)
        .option(NebulaOptions.PARTITION_NUMBER, readConfig.getPartitionNum)

      dfReader.load()
    }


    def loadEdge(): DataFrame = {
      assert(connectionConfig != null && readConfig != null,
        "nebula config is not set, please call nebula() before loadVerticesToDF")
      val dfReader = reader
        .format(classOf[NebulaDataSource].getName)
        .option(NebulaOptions.TYPE, DataTypeEnum.EDGE.toString)
        .option(NebulaOptions.OPERATE_TYPE, OperaType.READ.toString)
        .option(NebulaOptions.TIMEOUT, connectionConfig.getTimeout)
        .option(NebulaOptions.EXECUTION_RETRY, connectionConfig.getExecRetry)
        .option(NebulaOptions.EXECUTION_RETRY_INTERVAL, connectionConfig.getExecRetryIntervalMs)
        .option(NebulaOptions.GRAPH_ADDRESS, connectionConfig.getGraphAddress)
        .option(NebulaOptions.USER_NAME, connectionConfig.getUser)
        .option(NebulaOptions.AUTHOPTIONS, connectionConfig.getAuthOptions)
        .option(NebulaOptions.GRAPH_NAME, readConfig.getGraphName)
        .option(NebulaOptions.LABEL, readConfig.getTypeName)
        .option(NebulaOptions.RETURN_COLS, readConfig.getReturnColsString)
        .option(NebulaOptions.BATCH_SIZE, readConfig.getBatchSize)
        .option(NebulaOptions.PARTITION_NUMBER, readConfig.getPartitionNum)

      dfReader.load()
    }
  }


  /**
   * spark reader for nebula gql query
   */
  implicit class NebulaDataFrameGqlReader(reader: DataFrameReader) {
    var connectionConfig: NebulaConnectionConfig = _
    var gqlConfig: GqlNebulaConfig = _

    def gql(connectionConfig: NebulaConnectionConfig, gqlConfig: GqlNebulaConfig): NebulaDataFrameGqlReader = {
      SparkValidate.validate("2.4.*")
      this.connectionConfig = connectionConfig
      this.gqlConfig = gqlConfig
      this
    }

    /**
     * Reading gql result from NebulaGraph
     *
     * @return DataFrame
     */
    def load(): DataFrame = {
      assert(connectionConfig != null && gqlConfig != null,
        "nebula gql config is not set, please call gql() before load()")
      val dfReader = reader
        .format(classOf[NebulaDataSource].getName)
        .option(NebulaOptions.TYPE, DataTypeEnum.GQL.toString)
        .option(NebulaOptions.OPERATE_TYPE, OperaType.READ.toString)
        .option(NebulaOptions.TIMEOUT, connectionConfig.getTimeout)
        .option(NebulaOptions.EXECUTION_RETRY, connectionConfig.getExecRetry)
        .option(NebulaOptions.EXECUTION_RETRY_INTERVAL, connectionConfig.getExecRetryIntervalMs)
        .option(NebulaOptions.GRAPH_ADDRESS, connectionConfig.getGraphAddress)
        .option(NebulaOptions.USER_NAME, connectionConfig.getUser)
        .option(NebulaOptions.AUTHOPTIONS, connectionConfig.getAuthOptions)
        .option(NebulaOptions.GQL, gqlConfig.getGql)
      dfReader.load()
    }
  }

}
