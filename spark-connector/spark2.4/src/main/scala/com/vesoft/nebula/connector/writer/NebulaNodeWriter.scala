
package com.vesoft.nebula.connector.writer

import com.vesoft.nebula.driver.graph.ErrorCode
import com.vesoft.nebula.spark.common.writer.NebulaExecutor
import com.vesoft.nebula.spark.common.{NebulaNode, NebulaNodes, NebulaOptions, WriteMode}
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.sources.v2.writer.{DataWriter, WriterCommitMessage}
import org.apache.spark.sql.types.StructType
import org.slf4j.LoggerFactory

import scala.collection.mutable.ListBuffer

class NebulaNodeWriter(nebulaOptions: NebulaOptions, schema: StructType)
  extends NebulaWriter(nebulaOptions)
    with DataWriter[InternalRow] {

  private val LOG = LoggerFactory.getLogger(this.getClass)
  private val nodeDesc = graphProvider.getNodeDesc(nebulaOptions.graphName, nebulaOptions.label)

  val fieldTypeMap: Map[String, String] = nodeDesc.properties
  private val pkNames: List[String] = nodeDesc.nodePkNames

  /** buffer to save batch vertices */
  private var vertices: ListBuffer[NebulaNode] = new ListBuffer()

  /**
   * write one node row to buffer
   */
  override def write(row: InternalRow): Unit = {
    // check the node primary key value's validation, the pkName must exist in row,
    // it's already checked in NebulaDataSourceNodeWriter.createWriterFactory
    val pkIndicesInSparkRow: Map[String,Int] = schema.fields.toList.map(field => field.name).zip(schema.fields.indices).toMap
    pkIndicesInSparkRow.foreach((pkIndex)=>{
      if (row.isNullAt(pkIndex._2)) {
        LOG.warn(s">>>> record has null value at index ${pkIndex._2} for primary key ${pkIndex._1}, ignore it. record:$row")
        return
      }
    })

    val values =
      if (nebulaOptions.writeMode == WriteMode.DELETE) {
        NebulaExecutor.assignNodePkValues(schema, row, fieldTypeMap, pkNames)
      } else {
        NebulaExecutor.assignNodePropValues(schema,
          row,
          fieldTypeMap)
      }
    val nebulaNode = NebulaNode(values)
    vertices.append(nebulaNode)
    if (vertices.size >= nebulaOptions.batchSize) {
      execute()
    }
  }

  /**
   * submit buffer vertices to nebula
   */
  def execute(): Unit = {
    writeNodes(vertices)
    vertices.clear()
  }

  def writeNodes(vertices: ListBuffer[NebulaNode]): Unit = {
    val exec = getGql(vertices.toList)
    val result = submit(exec)
    if (result.isSucceeded) {
      if (!nebulaOptions.disableWriteLog) {
        LOG.info(
          s"batch write for ${nebulaOptions.label} succeed. batch size(${vertices.size}), latency(${result.getLatency})")
      }
    } else {
      if (vertices.size == 1
        && result.getErrorCode != ErrorCode.LEADER_CHANGED
        && !result.getErrorCode.isRpcError
        && !result.getErrorCode.isRaftError) {
        failedExecs.append(exec)
        LOG.error(s"write edge ${nebulaOptions.label} failed: ${result.getErrorMessage}.")
        return
      }
      // re-execute the vertices one by one
      LOG.warn(
        s"write node ${nebulaOptions.label} failed: ${result.getErrorMessage}, " +
          s"now retry writing one by one.")
      vertices.par.foreach { node => writeNode(node) }
    }
  }

  def writeNode(node: NebulaNode): Unit = {
    val exec = getGql(List(node))
    val result = submit(exec)
    if (result.isSucceeded) {
      if (!nebulaOptions.disableWriteLog) {
        LOG.info(s"write ${nebulaOptions.label}, batch size(1), latency(${result.getLatency}ms)")
      }
      return
    }
    if (result.getErrorCode == ErrorCode.NODE_ALREADY_EXIST) {
      LOG.warn(
        s"write ${nebulaOptions.label} failed, the node already exists. ")
      return
    }
    // retry the write execution for RPC_ERROR(NN), LEADER_CHANGED(ND005), RAFT_ERROR(NA)
    var executeResult = result
    var retry = 0
    while (retry < nebulaOptions.executionRetry &&
      (executeResult.getErrorCode.isRpcError
        || executeResult.getErrorCode == ErrorCode.LEADER_CHANGED
        || executeResult.getErrorCode.isRaftError)) {
      retry += 1
      Thread.sleep(nebulaOptions.executionRetryInterval)
      executeResult = submit(exec)
      if (executeResult.isSucceeded) {
        if (!nebulaOptions.disableWriteLog) {
          LOG.info(
            s"write node ${nebulaOptions.label}, batch size(1), latency(${executeResult.getLatency}ms)")
        }
        return
      }
    }
    LOG.error(s"write node ${nebulaOptions.label} failed: ${executeResult.getErrorMessage}.")
    failedExecs.append(exec)
  }

  private def getGql(nebulaVertices: List[NebulaNode]): String = {
    val nebulaNodes = NebulaNodes(nebulaOptions.label, nebulaVertices, pkNames, fieldTypeMap)
    val exec = nebulaOptions.writeMode match {
      case WriteMode.INSERT =>
        NebulaExecutor.toInsertSentence(nebulaOptions.graphName, nebulaNodes, "")
      case WriteMode.INSERTREPLACE =>
        NebulaExecutor.toInsertSentence(nebulaOptions.graphName, nebulaNodes, "OR REPLACE")
      case WriteMode.INSERTIGNORE=>
        NebulaExecutor.toInsertSentence(nebulaOptions.graphName, nebulaNodes, "OR IGNORE")
      case WriteMode.INSERTUPDATE=>
        NebulaExecutor.toInsertSentence(nebulaOptions.graphName, nebulaNodes, "OR UPDATE")
      case WriteMode.UPDATE =>
        NebulaExecutor.toUpdateSentence(nebulaOptions.graphName, nebulaOptions.label, nebulaNodes)
      case WriteMode.DELETE =>
        NebulaExecutor.toDeleteSentence(nebulaOptions.graphName, nebulaOptions.label, nebulaNodes, "DELETE")
      case WriteMode.DETACHDELETE =>
        NebulaExecutor.toDeleteSentence(nebulaOptions.graphName, nebulaOptions.label, nebulaNodes, "DETACH DELETE")
      case _ =>
        throw new IllegalArgumentException(s"write mode ${nebulaOptions.writeMode} not supported.")
    }
    exec
  }

  override def commit(): WriterCommitMessage = {
    if (vertices.nonEmpty) {
      execute()
    }
    graphProvider.close()
    NebulaCommitMessage(failedExecs.toList)
  }

  override def abort(): Unit = {
    LOG.error("insert node task abort.")
    graphProvider.close()
  }
}
