/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.writer

import com.vesoft.nebula.spark.common.nebula.VidType
import com.vesoft.nebula.spark.common.writer.NebulaExecutor
import com.vesoft.nebula.spark.common.{NebulaEdge, NebulaEdges, NebulaOptions, WriteMode}
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.sources.v2.writer.{DataWriter, WriterCommitMessage}
import org.apache.spark.sql.types.StructType
import org.slf4j.LoggerFactory

import scala.collection.mutable.ListBuffer

class NebulaEdgeWriter(nebulaOptions: NebulaOptions,
                       srcIndex: Int,
                       dstIndex: Int,
                       schema: StructType)
    extends NebulaWriter(nebulaOptions)
    with DataWriter[InternalRow] {

  private val LOG      = LoggerFactory.getLogger(this.getClass)
  private val edgeDesc = graphProvider.getEdgeDesc(nebulaOptions.graphName, nebulaOptions.label)

  val fieldTypeMap: Map[String, String] = edgeDesc.properties

  /** buffer to save batch edges */
  var edges: ListBuffer[NebulaEdge] = new ListBuffer()

  private val isSourceIdStringType = edgeDesc.srcNodePkDataType == VidType.STRING
  private val isTargetIdStringType = edgeDesc.dstNodePkDataType == VidType.STRING

  /**
    * write one edge record to buffer
    */
  override def write(row: InternalRow): Unit = {
    val srcId = NebulaExecutor.extraID(schema, row, srcIndex, isSourceIdStringType)
    val dstId = NebulaExecutor.extraID(schema, row, dstIndex, isTargetIdStringType)

    val values =
      if (nebulaOptions.writeMode == WriteMode.DELETE) {
        // delete mode does not need property.
        Map[String, String]()
      } else {
        NebulaExecutor.assignEdgeValues(schema,
                                        row,
                                        srcIndex,
                                        dstIndex,
                                        nebulaOptions.srcAsProp,
                                        nebulaOptions.dstAsProp,
                                        fieldTypeMap)
      }
    val nebulaEdge = NebulaEdge(srcId, dstId, values)
    edges.append(nebulaEdge)
    if (edges.size >= nebulaOptions.batch) {
      execute()
    }
  }

  /**
    * submit buffer edges to nebula
    */
  def execute(): Unit = {
    writeEdges(edges)
    edges.clear()
  }

  def writeEdges(edges: ListBuffer[NebulaEdge]): Unit = {
    val exec   = getGql(edges.toList)
    val result = submit(exec)
    if (result.isSucceeded) {
      if (!nebulaOptions.disableWriteLog) {
        LOG.info(
          s"batch write for ${nebulaOptions.label} succeed. batch size(${edges.size}), latency(${result.getLatency})")
      }
    } else {
      // re-execute the vertices one by one
      LOG.warn(
        s"write edge ${nebulaOptions.label} failed error message: ${result.getGqlStatus}, " +
          s"mow retry writing one by one.\n ${exec}")
      edges.par.foreach { edge =>
        writeEdge(edge)
      }
    }
  }

  def writeEdge(edge: NebulaEdge): Unit = {
    val exec   = getGql(List(edge))
    val result = submit(exec)
    if (result.isSucceeded) {
      if (!nebulaOptions.disableWriteLog) {
        LOG.info(s"write ${nebulaOptions.label}, batch size(1), latency(${result.getLatency}ms)")
      }
      return
    }
    // retry the write execution for RAFT_RPC_FAILURE(6029), RAFT_LEADER_CHANGED(6026), RAFT_BUFFER_OVERFLOW(6019)
    var executeResult = result
    var retry         = 0
    while (retry < nebulaOptions.executionRetry &&
           (executeResult.getGqlStatus.contains("6029")
           || executeResult.getGqlStatus.contains("6026")
           || executeResult.getGqlStatus.contains("6019"))) {
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
    LOG.error(s"write edge failed for ${executeResult.getGqlStatus}, statement:\n ${exec}")
    failedExecs.append(exec)
  }

  private def getGql(edges: List[NebulaEdge]): String = {
    val nebulaEdges = NebulaEdges(nebulaOptions.label, edges)
    val exec = nebulaOptions.writeMode match {
      case WriteMode.INSERT =>
        NebulaExecutor.toExecuteSentence(nebulaOptions.graphName, nebulaOptions.label, nebulaEdges)
      case _ =>
        throw new IllegalArgumentException(s"write mode ${nebulaOptions.writeMode} not supported.")
    }
    exec
  }

  override def commit(): WriterCommitMessage = {
    if (edges.nonEmpty) {
      execute()
    }
    graphProvider.close()
    NebulaCommitMessage.apply(failedExecs.toList)
  }

  override def abort(): Unit = {
    LOG.error("insert edge task abort.")
    graphProvider.close()
  }
}
