/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import com.vesoft.nebula.spark.common.WriteNebulaVertexConfig.WriteVertexConfigBuilder
import org.slf4j.{Logger, LoggerFactory}

class NebulaConnectionConfig(graphAddress: String,
                             user: String,
                             passwd: String,
                             timeout: Int,
                             executeRetry: Int,
                             executeRetryIntervalMs: Int)
    extends Serializable {
  def getGraphAddress = graphAddress

  def getUser = user

  def getPasswd              = passwd
  def getTimeout             = timeout
  def getExecRetry           = executeRetry
  def getExecRetryIntervalMs = executeRetryIntervalMs

}

object NebulaConnectionConfig {
  class ConfigBuilder {
    private val LOG = LoggerFactory.getLogger(this.getClass)

    protected var graphAddress: String        = _
    protected var user: String                = _
    protected var passwd: String              = _
    protected var timeout: Int                = 6000
    protected var executeRetry: Int           = 3
    protected var executeRetryIntervalMs: Int = 0

    /**
      * set nebula graph server address, multi addresses is split by English comma
      */
    def withGraphAddress(graphAddress: String): ConfigBuilder = {
      this.graphAddress = graphAddress
      this
    }

    /**
      * set user name for nebula graph
      */
    def withUser(user: String): ConfigBuilder = {
      this.user = user
      this
    }

    /**
      * set password for nebula graph's user
      */
    def withPasswd(passwd: String): ConfigBuilder = {
      this.passwd = passwd
      this
    }

    /**
      * set timeout, timeout is optional， unit: second
      */
    def withTimeoutSec(timeout: Int): ConfigBuilder = {
      this.timeout = timeout
      this
    }

    /**
      * set executeRetry, executeRetry is optional
      */
    def withExecuteRetry(executeRetry: Int): ConfigBuilder = {
      this.executeRetry = executeRetry
      this
    }

    /**
      * set executeRetryIntervalMs, executeRetryIntervalMs is optional
      */
    def withExecuteRetryIntervalMs(executeRetryIntervalMs: Int): ConfigBuilder = {
      this.executeRetryIntervalMs = executeRetryIntervalMs
      this
    }

    /**
      * check if the connection config is valid
      */
    def check(): Unit = {
      assert(graphAddress != null && graphAddress.nonEmpty, "graph address cannot be blank.")
      assert(user != null && user.nonEmpty, "user cannot be blank.")
      assert(passwd != null && passwd.nonEmpty, "password cannot be blank.")
      assert(timeout > 0, "timeout must be larger than 0.")
      assert(executeRetry >= 0, "retry must be equal or larger than 0.")
    }

    /**
      * build NebulaConnectionConfig
      */
    def build(): NebulaConnectionConfig = {
      check()
      new NebulaConnectionConfig(graphAddress,
                                 user,
                                 passwd,
                                 timeout,
                                 executeRetry,
                                 executeRetryIntervalMs)
    }
  }

  def builder(): ConfigBuilder = {
    new ConfigBuilder
  }

}

/**
  * Base config needed when write dataframe into nebula graph
  */
class WriteNebulaConfig(graphName: String,
                        batchSize: Int,
                        writeMode: String,
                        disableWriteLog: Boolean)
    extends Serializable {
  def getGraphName: String       = graphName
  def getBatchSize: Int          = batchSize
  def getWriteMode: String       = writeMode
  def isDisableWriteLog: Boolean = disableWriteLog
}

/**
  * subclass of WriteNebulaConfig to config vertex when write dataframe into nebula graph
  *
  * @param graphName: nebula graph name
  * @param nodeType: node type name
  * @param pkField: field in dataframe to indicate vertexId
  * @param batchSize: amount of one batch when write into nebula graph
  * @param writeMode: write mode, insert / update / delete
  *  @param disableWriteLog: disable the log print for write result, such as batch size and latency
  */
class WriteNebulaVertexConfig(graphName: String,
                              nodeType: String,
                              pkField: String,
                              batchSize: Int,
                              writeMode: String,
                              disableWriteLog: Boolean)
    extends WriteNebulaConfig(graphName, batchSize, writeMode, disableWriteLog) {
  def getNodeType = nodeType
  def getPkField  = pkField
}

/**
  * object WriteNebulaVertexConfig
  * */
object WriteNebulaVertexConfig {

  private val LOG: Logger = LoggerFactory.getLogger(this.getClass)

  class WriteVertexConfigBuilder {
    private var graphName: String        = _
    private var nodeType: String         = _
    private var writeMode: String        = "insert"
    private var disableWriteLog: Boolean = false
    private var pkField: String          = _
    private var batchSize: Int           = 512

    /**
      * set graph name
      */
    def withGraphName(graphName: String): WriteVertexConfigBuilder = {
      this.graphName = graphName
      this
    }

    /**
      * set tag name
      */
    def withNodeType(nodeType: String): WriteVertexConfigBuilder = {
      this.nodeType = nodeType
      this
    }

    /**
      * set which field in dataframe as nebula tag's id
      */
    def withPrimaryKeyField(pkField: String): WriteVertexConfigBuilder = {
      this.pkField = pkField
      this
    }

    /**
      * set data amount for one batch, default is 512
      */
    def withBatchSize(batchSize: Int): WriteVertexConfigBuilder = {
      this.batchSize = batchSize
      this
    }

    /**
      * set nebula write mode for nebula tag, INSERT or UPDATE
      */
    def withWriteMode(writeMode: WriteMode.Value): WriteVertexConfigBuilder = {
      this.writeMode = writeMode.toString
      this
    }

    /**
      * set disableWriteLog, default is false
      */
    def withDisableWriteLog(disableWriteLog: Boolean): WriteVertexConfigBuilder = {
      this.disableWriteLog = disableWriteLog
      this
    }

    /**
      * check and get WriteNebulaVertexConfig
      */
    def build(): WriteNebulaVertexConfig = {
      check()
      new WriteNebulaVertexConfig(graphName,
                                  nodeType,
                                  pkField,
                                  batchSize,
                                  writeMode,
                                  disableWriteLog)
    }

    /**
      * check the validation for {@link WriteNebulaVertexConfig}
      */
    private def check(): Unit = {
      assert(graphName != null && graphName.nonEmpty, s"config graphName can not be empty.")
      assert(nodeType != null && nodeType.nonEmpty, "config nodeType can not be empty")
      assert(pkField != null && pkField.nonEmpty, "config primaryKeyField can not be empty.")
      assert(batchSize > 0, s"config batchSize must be positive, your batchSize is $batchSize.")

      try {
        WriteMode.withName(writeMode.toLowerCase())
      } catch {
        case e: Throwable =>
          assert(false, s"optional write mode: insert or update, your write mode is $writeMode")
      }
      if (!writeMode.equalsIgnoreCase(WriteMode.INSERT.toString)) {
        assert(false, s"the writeMode is ${writeMode}, for now just INSERT is supported.")
      }

      LOG.info(
        s"NebulaWriteVertexConfig={graphName=$graphName,nodeType=$nodeType,pkField=$pkField," +
          s"batchSize=$batchSize,writeMode=$writeMode,disableWriteLog=$disableWriteLog}")
    }
  }

  def builder(): WriteVertexConfigBuilder = {
    new WriteVertexConfigBuilder
  }
}

/**
  * subclass of WriteNebulaConfig to config edge when write dataframe into nebula graph
  *
  * @param graphName: nebula graph name
  * @param edgeType: edge name
  * @param srcPkFiled: field in dataframe to indicate src vertex primary key
  * @param dstPkField: field in dataframe to indicate dst vertex primary key
  * @param batchSize: amount of one batch when write into nebula graph
  * @param srcPkAsProp: whether use src node primary key as edge's property
  * @param dstPkAsProp: whether use dst node primary key as edge's property
  * @param writeMode: write mode, insert / update / delete
  * @param disableWriteLog: disable the log print for write result, such as batch size and latency
  */
class WriteNebulaEdgeConfig(graphName: String,
                            edgeType: String,
                            srcPkFiled: String,
                            dstPkField: String,
                            batchSize: Int,
                            srcPkAsProp: Boolean,
                            dstPkAsProp: Boolean,
                            writeMode: String,
                            disableWriteLog: Boolean)
    extends WriteNebulaConfig(graphName, batchSize, writeMode, disableWriteLog) {
  def getEdgeType: String   = edgeType
  def getSrcPkFiled: String = srcPkFiled
  def getDstPkField: String = dstPkField
  def getSrcAsProp: Boolean = srcPkAsProp
  def getDstAsProp: Boolean = dstPkAsProp
}

/**
  * object WriteNebulaEdgeConfig
  */
object WriteNebulaEdgeConfig {

  private val LOG: Logger = LoggerFactory.getLogger(WriteNebulaEdgeConfig.getClass)

  /**
    * a builder to create {@link WriteNebulaEdgeConfig}
    */
  class WriteEdgeConfigBuilder {
    var graphName: String        = _
    var writeMode: String        = WriteMode.INSERT.toString
    var disableWriteLog: Boolean = false

    private var edgeType: String     = _
    private var srcPkField: String   = _
    private var dstPkField: String   = _
    private var srcPkAsProp: Boolean = false
    private var dstPkAsProp: Boolean = false
    private var batchSize: Int       = 512

    /**
      * set graph name
      */
    def withGraphName(graphName: String): WriteEdgeConfigBuilder = {
      this.graphName = graphName
      this
    }

    /**
      * set edge type name
      */
    def withEdge(edgeType: String): WriteEdgeConfigBuilder = {
      this.edgeType = edgeType
      this
    }

    /**
      * set which field in dataframe as nebula edge's src id
      */
    def withSrcPkField(srcPkField: String): WriteEdgeConfigBuilder = {
      this.srcPkField = srcPkField
      this
    }

    /**
      * set which field in dataframe as nebula edge's dst id
      */
    def withDstPkField(dstIdField: String): WriteEdgeConfigBuilder = {
      this.dstPkField = dstIdField
      this
    }

    /**
      * set data amount for one batch, default is 512
      */
    def withBatchSize(batchSize: Int): WriteEdgeConfigBuilder = {
      this.batchSize = batchSize
      this
    }

    /**
      * set whether src id as property
      */
    def withSrcPkAsProperty(srcPkAsProp: Boolean): WriteEdgeConfigBuilder = {
      this.srcPkAsProp = srcPkAsProp
      this
    }

    /**
      * set whether dst id as property
      */
    def withDstPkAsProperty(dstPkAsProp: Boolean): WriteEdgeConfigBuilder = {
      this.dstPkAsProp = dstPkAsProp
      this
    }

    /**
      * set write mode for nebula edge, INSERT or UPDATE
      */
    def withWriteMode(writeMode: WriteMode.Value): WriteEdgeConfigBuilder = {
      this.writeMode = writeMode.toString
      this
    }

    /**
      * set disableWriteLog, default is false
      */
    def withDisableWriteLog(disableWriteLog: Boolean): WriteEdgeConfigBuilder = {
      this.disableWriteLog = disableWriteLog
      this
    }

    /**
      * check configs and get WriteNebulaEdgeConfig
      */
    def build(): WriteNebulaEdgeConfig = {
      check()
      new WriteNebulaEdgeConfig(graphName,
                                edgeType,
                                srcPkField,
                                dstPkField,
                                batchSize,
                                srcPkAsProp,
                                dstPkAsProp,
                                writeMode,
                                disableWriteLog)
    }

    private def check(): Unit = {
      assert(graphName != null && graphName.nonEmpty, s"config graphName can not be empty.")
      assert(edgeType != null && edgeType.nonEmpty, "config edgeType can not be empty")
      assert(srcPkField != null && srcPkField.nonEmpty, "config srcPkField can not be empty.")
      assert(dstPkField != null && dstPkField.nonEmpty, "config dstPkField can not be empty.")
      assert(batchSize > 0, s"config batchSize must be positive, your batchSize is $batchSize.")

      try {
        WriteMode.withName(writeMode.toLowerCase)
      } catch {
        case e: Throwable =>
          assert(false, s"optional write mode: insert or update, your write mode is $writeMode")
      }
      if (!writeMode.equalsIgnoreCase(WriteMode.INSERT.toString)) {
        assert(false, s"the writeMode is ${writeMode}, for now just INSERT is supported.")
      }
      LOG.info(
        s"NebulaWriteEdgeConfig={graphName=$graphName,edgeType=$edgeType,srcPkField=$srcPkField," +
          s"dstPkField=$dstPkField,batchSize=$batchSize,srcPkAsProp=$srcPkField," +
          s"dstPkAsProp=$dstPkField,writeMode=$writeMode,disableWriteLog=$disableWriteLog}")
    }
  }

  def builder(): WriteEdgeConfigBuilder = {
    new WriteEdgeConfigBuilder
  }
}
