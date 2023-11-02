/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import org.slf4j.{Logger, LoggerFactory}

import scala.collection.mutable.ListBuffer

class NebulaConnectionConfig(graphAddress: String,
                             timeout: Int,
                             connectionRetry: Int,
                             executeRetry: Int)
    extends Serializable {
  def getGraphAddress    = graphAddress
  def getTimeout         = timeout
  def getConnectionRetry = connectionRetry
  def getExecRetry       = executeRetry

}

object NebulaConnectionConfig {
  class ConfigBuilder {
    private val LOG = LoggerFactory.getLogger(this.getClass)

    protected var graphAddress: String = _
    protected var timeout: Int         = 6000
    protected var connectionRetry: Int = 1
    protected var executeRetry: Int    = 1

    /**
      * set nebula graph server address, multi addresses is split by English comma
      */
    def withGraphAddress(graphAddress: String): ConfigBuilder = {
      this.graphAddress = graphAddress
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
      * set connectionRetry, connectionRetry is optional
      */
    def withConnectionRetry(connectionRetry: Int): ConfigBuilder = {
      this.connectionRetry = connectionRetry
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
      * check if the connection config is valid
      */
    def check(): Unit = {
      assert(timeout > 0, "timeout must be larger than 0")
      assert(connectionRetry > 0 && executeRetry > 0, "retry must be larger than 0.")
    }

    /**
      * build NebulaConnectionConfig
      */
    def build(): NebulaConnectionConfig = {
      check()
      new NebulaConnectionConfig(graphAddress, timeout, connectionRetry, executeRetry)
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
                        user: String,
                        passwd: String,
                        batch: Int,
                        writeMode: String)
    extends Serializable {
  def getGraphName = graphName
  def getBatch     = batch
  def getUser      = user
  def getPasswd    = passwd
  def getWriteMode = writeMode
}

/**
  * subclass of WriteNebulaConfig to config vertex when write dataframe into nebula graph
  *
  * @param graphName: nebula graph name
  * @param nodeType: node type name
  * @param vidField: field in dataframe to indicate vertexId
  * @param batch: amount of one batch when write into nebula graph
  */
class WriteNebulaVertexConfig(graphName: String,
                              nodeType: String,
                              vidField: String,
                              batch: Int,
                              vidAsProp: Boolean,
                              user: String,
                              passwd: String,
                              writeMode: String)
    extends WriteNebulaConfig(graphName, user, passwd, batch, writeMode) {
  def getNodeType  = nodeType
  def getVidField  = vidField
  def getVidAsProp = vidAsProp

}

/**
  * object WriteNebulaVertexConfig
  * */
object WriteNebulaVertexConfig {

  private val LOG: Logger = LoggerFactory.getLogger(this.getClass)

  class WriteVertexConfigBuilder {

    var graphName: String = _
    var nodeType: String  = _
    var vidField: String  = _
    var batchSize: Int    = 512
    var user: String      = "root"
    var passwd: String    = "nebula"
    var writeMode: String = "insert"

    /** whether set vid as property */
    var vidAsProp: Boolean = false

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
    def withVidField(vidField: String): WriteVertexConfigBuilder = {
      this.vidField = vidField
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
      * set whether vid as prop, default is false
      */
    def withVidAsProp(vidAsProp: Boolean): WriteVertexConfigBuilder = {
      this.vidAsProp = vidAsProp
      this
    }

    /**
      * set user name for nebula graph
      */
    def withUser(user: String): WriteVertexConfigBuilder = {
      this.user = user
      this
    }

    /**
      * set password for nebula graph's user
      */
    def withPasswd(passwd: String): WriteVertexConfigBuilder = {
      this.passwd = passwd
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
      * check and get WriteNebulaVertexConfig
      */
    def build(): WriteNebulaVertexConfig = {
      check()
      new WriteNebulaVertexConfig(graphName,
                                  nodeType,
                                  vidField,
                                  batchSize,
                                  vidAsProp,
                                  user,
                                  passwd,
                                  writeMode)
    }

    private def check(): Unit = {
      assert(graphName != null && graphName.nonEmpty, s"config graphName is empty.")

      assert(vidField != null && vidField.nonEmpty, "config vidField is empty.")
      assert(batchSize > 0, s"config batch must be positive, your batch is $batchSize.")

      assert(user != null && user.nonEmpty, "user is empty")
      assert(passwd != null && passwd.nonEmpty, "passwd is empty")
      try {
        WriteMode.withName(writeMode.toLowerCase())
      } catch {
        case e: Throwable =>
          assert(false, s"optional write mode: insert or update, your write mode is $writeMode")
      }
      if (writeMode.equalsIgnoreCase(WriteMode.UPDATE.toString)) {
        assert(batchSize <= 512, "the maximum number of statements for Nebula is 512")
      }

      LOG.info(
        s"NebulaWriteVertexConfig={graphName=$graphName,nodeType=$nodeType,vidField=$vidField," +
          s"batch=$batchSize,writeMode=$writeMode}")
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
  * @param edgeName: edge name
  * @param srcFiled: field in dataframe to indicate src vertex id
  * @param dstField: field in dataframe to indicate dst vertex id
  * @param batch: amount of one batch when write into nebula graph
  */
class WriteNebulaEdgeConfig(graphName: String,
                            edgeName: String,
                            srcFiled: String,
                            dstField: String,
                            batch: Int,
                            srcAsProp: Boolean,
                            dstAsProp: Boolean,
                            user: String,
                            passwd: String,
                            writeMode: String)
    extends WriteNebulaConfig(graphName, user, passwd, batch, writeMode) {
  def getEdgeName = edgeName
  def getSrcFiled = srcFiled
  def getDstField = dstField

  def getSrcAsProp = srcAsProp
  def getDstAsProp = dstAsProp

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

    var graphName: String  = _
    var edgeName: String   = _
    var srcIdField: String = _
    var dstIdField: String = _
    var batchSize: Int     = 512
    var user: String       = "root"
    var passwd: String     = "nebula"

    /** whether srcId as property */
    var srcAsProp: Boolean = false

    /** whether dstId as property */
    var dstAsProp: Boolean = false

    /** write mode for nebula, insert or update */
    var writeMode: String = WriteMode.INSERT.toString

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
    def withEdge(edgeName: String): WriteEdgeConfigBuilder = {
      this.edgeName = edgeName
      this
    }

    /**
      * set which field in dataframe as nebula edge's src id
      */
    def withSrcIdField(srcIdField: String): WriteEdgeConfigBuilder = {
      this.srcIdField = srcIdField
      this
    }

    /**
      * set which field in dataframe as nebula edge's dst id
      */
    def withDstIdField(dstIdField: String): WriteEdgeConfigBuilder = {
      this.dstIdField = dstIdField
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
    def withSrcAsProperty(srcAsProp: Boolean): WriteEdgeConfigBuilder = {
      this.srcAsProp = srcAsProp
      this
    }

    /**
      * set whether dst id as property
      */
    def withDstAsProperty(dstAsProp: Boolean): WriteEdgeConfigBuilder = {
      this.dstAsProp = dstAsProp
      this
    }

    /**
      * set user name for nebula graph
      */
    def withUser(user: String): WriteEdgeConfigBuilder = {
      this.user = user
      this
    }

    /**
      * set password for nebula graph's user
      */
    def withPasswd(passwd: String): WriteEdgeConfigBuilder = {
      this.passwd = passwd
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
      * check configs and get WriteNebulaEdgeConfig
      */
    def build(): WriteNebulaEdgeConfig = {
      check()
      new WriteNebulaEdgeConfig(graphName,
                                edgeName,
                                srcIdField,
                                dstIdField,
                                batchSize,
                                srcAsProp,
                                dstAsProp,
                                user,
                                passwd,
                                writeMode)
    }

    private def check(): Unit = {
      assert(graphName != null && !graphName.isEmpty, s"config graphName is empty.")

      assert(srcIdField != null && !srcIdField.isEmpty, "config srcIdField is empty.")
      assert(dstIdField != null && !dstIdField.isEmpty, "config dstIdField is empty.")

      assert(batchSize > 0, s"config batch must be positive, your batch is $batchSize.")
      assert(user != null && user.nonEmpty, "user is empty")
      assert(passwd != null && passwd.nonEmpty, "passwd is empty")
      try {
        WriteMode.withName(writeMode.toLowerCase)
      } catch {
        case e: Throwable =>
          assert(false, s"optional write mode: insert or update, your write mode is $writeMode")
      }
      if (writeMode.equalsIgnoreCase(WriteMode.UPDATE.toString)) {
        assert(batchSize <= 512, "the maximum number of statements for Nebula is 512")
      }
      assert(edgeName != null && edgeName.nonEmpty, s"config edgeName is empty.")
      LOG.info(
        s"NebulaWriteEdgeConfig={graphName=$graphName,edgeName=$edgeName,srcField=$srcIdField," +
          s"dstField=$dstIdField,writeMode=$writeMode}")
    }
  }

  def builder(): WriteEdgeConfigBuilder = {
    new WriteEdgeConfigBuilder
  }
}
