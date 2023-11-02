/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula

import com.vesoft.nebula.spark.common.utils.SparkValidate
import com.vesoft.nebula.spark.common.{DataTypeEnum, NebulaConnectionConfig, NebulaOptions, OperaType, WriteNebulaConfig, WriteNebulaEdgeConfig, WriteNebulaVertexConfig}
import org.apache.spark.sql.{ DataFrameWriter, Row, SaveMode}

package object connector {

  /**
    * spark writer for nebula graph
    */
  implicit class NebulaDataFrameWriter(writer: DataFrameWriter[Row]) {

    private var connectionConfig: NebulaConnectionConfig = _
    private var writeNebulaConfig: WriteNebulaConfig     = _

    /**
      * config nebula connection
      * @param connectionConfig connection parameters
      * @param writeNebulaConfig write parameters for vertex or edge
      */
    def nebula(connectionConfig: NebulaConnectionConfig,
               writeNebulaConfig: WriteNebulaConfig): NebulaDataFrameWriter = {
      SparkValidate.validate("2.4.*")
      this.connectionConfig = connectionConfig
      this.writeNebulaConfig = writeNebulaConfig
      this
    }

    /**
      * write dataframe into nebula vertex
      */
    def writeVertices(): Unit = {
      assert(connectionConfig != null && writeNebulaConfig != null,
             "nebula config is not set, please call nebula() before writeVertices")
      val writeConfig = writeNebulaConfig.asInstanceOf[WriteNebulaVertexConfig]
      val dfWriter = writer
        .format(classOf[NebulaDataSource].getName)
        .mode(SaveMode.Overwrite)
        .option(NebulaOptions.TYPE, DataTypeEnum.VERTEX.toString)
        .option(NebulaOptions.OPERATE_TYPE, OperaType.WRITE.toString)
        .option(NebulaOptions.GRAPH_NAME, writeConfig.getGraphName)
        .option(NebulaOptions.LABEL, writeConfig.getNodeType)
        .option(NebulaOptions.USER_NAME, writeConfig.getUser)
        .option(NebulaOptions.PASSWD, writeConfig.getPasswd)
        .option(NebulaOptions.VERTEX_FIELD, writeConfig.getVidField)
        .option(NebulaOptions.BATCH, writeConfig.getBatch)
        .option(NebulaOptions.VID_AS_PROP, writeConfig.getVidAsProp)
        .option(NebulaOptions.WRITE_MODE, writeConfig.getWriteMode)
        .option(NebulaOptions.GRAPH_ADDRESS, connectionConfig.getGraphAddress)
        .option(NebulaOptions.TIMEOUT, connectionConfig.getTimeout)
        .option(NebulaOptions.CONNECTION_RETRY, connectionConfig.getConnectionRetry)
        .option(NebulaOptions.EXECUTION_RETRY, connectionConfig.getExecRetry)

      dfWriter.save()
    }

    /**
      * write dataframe into nebula edge
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
        .option(NebulaOptions.GRAPH_NAME, writeConfig.getGraphName)
        .option(NebulaOptions.USER_NAME, writeConfig.getUser)
        .option(NebulaOptions.PASSWD, writeConfig.getPasswd)
        .option(NebulaOptions.LABEL, writeConfig.getEdgeName)
        .option(NebulaOptions.SRC_VERTEX_FIELD, writeConfig.getSrcFiled)
        .option(NebulaOptions.DST_VERTEX_FIELD, writeConfig.getDstField)
        .option(NebulaOptions.BATCH, writeConfig.getBatch)
        .option(NebulaOptions.SRC_AS_PROP, writeConfig.getSrcAsProp)
        .option(NebulaOptions.DST_AS_PROP, writeConfig.getDstAsProp)
        .option(NebulaOptions.WRITE_MODE, writeConfig.getWriteMode)
        .option(NebulaOptions.GRAPH_ADDRESS, connectionConfig.getGraphAddress)
        .option(NebulaOptions.TIMEOUT, connectionConfig.getTimeout)
        .option(NebulaOptions.CONNECTION_RETRY, connectionConfig.getConnectionRetry)
        .option(NebulaOptions.EXECUTION_RETRY, connectionConfig.getExecRetry)

      dfWriter.save()
    }
  }

}
