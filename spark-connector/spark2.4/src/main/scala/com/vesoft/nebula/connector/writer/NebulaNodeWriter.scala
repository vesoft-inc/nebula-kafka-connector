/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.writer

import com.vesoft.nebula.spark.common.exception.IllegalOptionException
import com.vesoft.nebula.spark.common.nebula.VidType
import com.vesoft.nebula.spark.common.writer.NebulaExecutor
import com.vesoft.nebula.spark.common.{NebulaOptions, NebulaNode, NebulaNodes, WriteMode}
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
  val pkName: String = nodeDesc.nodePkName

  /** buffer to save batch vertices */
  var vertices: ListBuffer[NebulaNode] = new ListBuffer()

  /**
   * write one node row to buffer
   */
  override def write(row: InternalRow): Unit = {
    // check the node primary key value's validation, the pkName must exist in row,
    // it's already checked in NebulaDataSourceNodeWriter.createWriterFactory
    val pkIndexInSparkRow: Int = schema.fields.toList.map(field => field.name).zip(schema.fields.indices).toMap.get(pkName).get
    if (row.isNullAt(pkIndexInSparkRow)) {
      LOG.warn(s">>>> record has null value at index $pkIndexInSparkRow for primary key $pkName, ignore it. record:$row")
      return
    }

    val values =
      if (nebulaOptions.writeMode == WriteMode.DELETE) {
        // delete mode does not need property.
        Map[String, String]()
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
          s"batch write for ${nebulaOptions.label} succeed. batch size(${vertices.size}, latency(${result.getLatency}))")
      }
    } else {
      // re-execute the vertices one by one
      LOG.warn(
        s"write node ${nebulaOptions.label} failed error message: ${result.getErrorMessage}, " +
          s"mow retry writing one by one.\n ${exec}")
      vertices.par.foreach { node =>
        writeNode(node)
      }
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
    if (result.getErrorCode.equals("ND002")) {
      LOG.warn(
        s"write ${nebulaOptions.label} failed, the node already exists. node: ${node.toString}")
      return
    }
    // retry the write execution for RPC_ERROR(NN), LEADER_CHANGED(ND005), RAFT_ERROR(NA)
    var executeResult = result
    var retry = 0
    while (retry < nebulaOptions.executionRetry &&
      (executeResult.getErrorCode.contains("NN")
        || executeResult.getErrorCode.equals("ND005")
        || executeResult.getErrorCode.contains("NA"))) {
      retry += 1
      Thread.sleep(nebulaOptions.executionRetryInterval)
      executeResult = submit(exec)
      if (executeResult.isSucceeded) {
        if (!nebulaOptions.disableWriteLog) {
          LOG.info(
            s"write ${nebulaOptions.label}, batch size(1), latency(${executeResult.getLatency}ms)")
        }
        return
      }
    }
    LOG.error(s"write node failed for ${executeResult.getErrorMessage}, statement:\n ${exec}")
    failedExecs.append(exec)
  }

  private def getGql(nebulaVertices: List[NebulaNode]): String = {
    val nebulaNodes = NebulaNodes(nebulaOptions.label, nebulaVertices)
    val exec = nebulaOptions.writeMode match {
      case WriteMode.INSERT =>
        NebulaExecutor.toExecuteSentence(nebulaOptions.graphName, nebulaOptions.label, nebulaNodes)
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
