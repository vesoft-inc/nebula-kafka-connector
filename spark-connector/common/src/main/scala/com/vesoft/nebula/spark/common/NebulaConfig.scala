/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import com.alibaba.fastjson.JSONObject
import com.alibaba.fastjson.serializer.JSONObjectCodec
import org.slf4j.{Logger, LoggerFactory}

import scala.collection.mutable
import scala.collection.mutable.ListBuffer

class NebulaConnectionConfig(graphAddress: String,
                             user: String,
                             passwd: String,
                             authOptions: Map[String, Any],
                             timeout: Int,
                             executeRetry: Int,
                             executeRetryIntervalMs: Int)
  extends Serializable {
  def getGraphAddress = graphAddress

  def getUser = user


  def getAuthOptions = {
    val newAuthOptions = new mutable.HashMap[String, Any]()
    if (passwd != null) {
      newAuthOptions.put("password", passwd)
    }
    newAuthOptions ++= authOptions
    val json = new JSONObject
    newAuthOptions.foreach(x => json.put(x._1, x._2))
    json.toJSONString
  }

  def getTimeout = timeout

  def getExecRetry = executeRetry

  def getExecRetryIntervalMs = executeRetryIntervalMs

}

object NebulaConnectionConfig {
  class ConfigBuilder {
    private val LOG = LoggerFactory.getLogger(this.getClass)

    protected var graphAddress: String = _
    protected var user: String = _
    protected var passwd: String = _
    protected var authOptions = new mutable.HashMap[String, Any]
    protected var timeout: Int = 5
    protected var executeRetry: Int = 3
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
     * set auth options for nebula graph's user
     */
    def withAuthOptions(authOptions: Map[String, Any]): ConfigBuilder = {
      if (authOptions != null) {
        this.authOptions ++= authOptions
      }
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
      assert((passwd != null && passwd.nonEmpty) || authOptions.nonEmpty,
        "password and authOptions cannot be blank at the same time.")
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
        authOptions.toMap,
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
  def getGraphName: String = graphName

  def getBatchSize: Int = batchSize

  def getWriteMode: String = writeMode

  def isDisableWriteLog: Boolean = disableWriteLog
}

/**
 * subclass of WriteNebulaConfig to config node when write dataframe into nebula graph
 *
 * @param graphName       : nebula graph name
 * @param nodeType        : node type name
 * @param batchSize       : amount of one batch when write into nebula graph
 * @param writeMode       : write mode, insert / update / delete
 * @param disableWriteLog : disable the log print for write result, such as batch size and latency
 */
class WriteNebulaNodeConfig(graphName: String,
                            nodeType: String,
                            batchSize: Int,
                            writeMode: String,
                            disableWriteLog: Boolean)
  extends WriteNebulaConfig(graphName, batchSize, writeMode, disableWriteLog) {
  def getNodeType = nodeType

}

/**
 * object WriteNebulaNodeConfig
 * */
object WriteNebulaNodeConfig {

  private val LOG: Logger = LoggerFactory.getLogger(this.getClass)

  class WriteNodeConfigBuilder {
    private var graphName: String = _
    private var nodeType: String = _
    private var writeMode: String = "insert"
    private var disableWriteLog: Boolean = false
    private var batchSize: Int = 512

    /**
     * set graph name
     */
    def withGraphName(graphName: String): WriteNodeConfigBuilder = {
      this.graphName = graphName
      this
    }

    /**
     * set tag name
     */
    def withNodeType(nodeType: String): WriteNodeConfigBuilder = {
      this.nodeType = nodeType
      this
    }

    /**
     * set data amount for one batch, default is 512
     */
    def withBatchSize(batchSize: Int): WriteNodeConfigBuilder = {
      this.batchSize = batchSize
      this
    }

    /**
     * set nebula write mode for nebula tag, INSERT or UPDATE
     */
    def withWriteMode(writeMode: WriteMode.Value): WriteNodeConfigBuilder = {
      this.writeMode = writeMode.toString
      this
    }

    /**
     * set disableWriteLog, default is false
     */
    def withDisableWriteLog(disableWriteLog: Boolean): WriteNodeConfigBuilder = {
      this.disableWriteLog = disableWriteLog
      this
    }

    /**
     * check and get WriteNebulaNodeConfig
     */
    def build(): WriteNebulaNodeConfig = {
      check()
      new WriteNebulaNodeConfig(graphName,
        nodeType,
        batchSize,
        writeMode,
        disableWriteLog)
    }

    /**
     * check the validation for {@link WriteNebulaNodeConfig}
     */
    private def check(): Unit = {
      assert(graphName != null && graphName.nonEmpty, s"config graphName can not be empty.")
      assert(nodeType != null && nodeType.nonEmpty, "config nodeType can not be empty")
      assert(batchSize > 0, s"config batchSize must be positive, your batchSize is $batchSize.")

      try {
        WriteMode.withName(writeMode.toLowerCase())
      } catch {
        case e: Throwable =>
          assert(false, s"optional write mode: insert or update, your write mode is $writeMode")
      }
      if (writeMode.equalsIgnoreCase(WriteMode.UPDATE.toString)) {
        assert(false, s"the writeMode is ${writeMode}, for now just INSERT and DELETE is supported.")
      }

      LOG.info(
        s"NebulaWriteNodeConfig={graphName=$graphName,nodeType=$nodeType," +
          s"batchSize=$batchSize,writeMode=$writeMode,disableWriteLog=$disableWriteLog}")
    }
  }

  def builder(): WriteNodeConfigBuilder = {
    new WriteNodeConfigBuilder
  }
}

/**
 * subclass of WriteNebulaConfig to config edge when write dataframe into nebula graph
 *
 * @param graphName       : nebula graph name
 * @param edgeType        : edge name
 * @param srcPkFiled      : field in dataframe to indicate src node primary key
 * @param dstPkField      : field in dataframe to indicate dst node primary key
 * @param batchSize       : amount of one batch when write into nebula graph
 * @param srcPkAsProp     : whether use src node primary key as edge's property
 * @param dstPkAsProp     : whether use dst node primary key as edge's property
 * @param writeMode       : write mode, insert / update / delete
 * @param disableWriteLog : disable the log print for write result, such as batch size and latency
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
  def getEdgeType: String = edgeType

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
    var graphName: String = _
    var writeMode: String = WriteMode.INSERT.toString
    var disableWriteLog: Boolean = false

    private var edgeType: String = _
    private var srcPkField: String = _
    private var dstPkField: String = _
    private var srcPkAsProp: Boolean = false
    private var dstPkAsProp: Boolean = false
    private var batchSize: Int = 512

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
      if (writeMode.equalsIgnoreCase(WriteMode.UPDATE.toString)) {
        assert(false, s"the writeMode is ${writeMode}, for now just INSERT and DELETE is supported.")
      }
      // the batch size must be 1 for DELETE edge
      if (writeMode.equalsIgnoreCase(WriteMode.DELETE.toString)) {
        LOG.info("the write mode is DELETE for edge, batch size is automatically adjusted to 1")
        batchSize = 1
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


class ReadNebulaConfig(graphName: String, typeName: String, returnCols: ListBuffer[String], partitionNum: Int, batchSize: Int) extends Serializable {
  def getGraphName: String = graphName

  def getTypeName: String = typeName

  def getReturnCols: ListBuffer[String] = returnCols

  def getReturnColsString: String = {
    if (returnCols == null) {
      return "$null"
    }
    returnCols.mkString(",")
  }

  def getPartitionNum: Int = partitionNum

  def getBatchSize: Int = batchSize
}


/**
 * config for reading NebulaGraph
 */
object ReadNebulaConfig {
  private val LOG: Logger = LoggerFactory.getLogger(this.getClass)

  def builder(): ReadConfigBuilder = {
    new ReadConfigBuilder
  }

  class ReadConfigBuilder {
    private var graphName: String = _
    private var typeName: String = _
    private var returnCols: ListBuffer[String] = _
    private var partitionNum: Int = 10
    private var batchSize: Int = 2000

    /**
     * config the graph name for reading
     */
    def withGraphName(graphName: String): ReadConfigBuilder = {
      this.graphName = graphName
      this
    }

    /**
     * config the type name for reading
     */
    def withTypeName(typeName: String): ReadConfigBuilder = {
      this.typeName = typeName
      this
    }

    /**
     * config the property names for reading
     * if you want to read all the properties, please ignore withReturnCols, or config returnCols as null.
     */
    def withReturnCols(returnCols: List[String]): ReadConfigBuilder = {
      if (returnCols != null) {
        this.returnCols = new ListBuffer[String]
        for (col: String <- returnCols) {
          this.returnCols.append(col)
        }
      }
      this
    }


    /**
     * config the partition num for spark, default is 10
     */
    def withPartitionNum(partitionNum: Int): ReadConfigBuilder = {
      this.partitionNum = partitionNum
      this
    }

    /**
     * config the batch size for each scan request, default is 2000
     */
    def withBatchSize(batchSize: Int): ReadConfigBuilder = {
      this.batchSize = batchSize
      this
    }


    /**
     * build a ReadNebulaConfig
     */
    def build(): ReadNebulaConfig = {
      check()
      new ReadNebulaConfig(graphName, typeName, returnCols, partitionNum, batchSize)
    }

    /**
     * check the validation for configs
     */
    private def check(): Unit = {
      assert(graphName != null && graphName.nonEmpty, "config graphName can't be empty.")
      assert(typeName != null && typeName.nonEmpty, "config typeName can't be empty.")
      assert(partitionNum > 0, s"config partitionNum must be positive, your value is $partitionNum.")
      assert(batchSize > 0, s"config batchSize must be positive, your value is $batchSize.")
    }
  }
}

class GqlNebulaConfig(gql: String) extends Serializable {
  def getGql: String = gql
}

/**
 * config for reading NebulaGraph through gql
 */
object GqlNebulaConfig {
  private val LOG: Logger = LoggerFactory.getLogger(this.getClass)

  def builder(): GqlConfigBuilder = {
    new GqlConfigBuilder
  }

  class GqlConfigBuilder {
    private var gql: String = _

    def withGql(gql: String): GqlConfigBuilder = {
      this.gql = gql;
      this
    }


    /**
     * build a ReadNebulaConfig
     */
    def build(): GqlNebulaConfig = {
      check()
      new GqlNebulaConfig(gql)
    }

    private def check(): Unit = {
      assert(gql != null, "config gql can't be empty.")
    }
  }
}

