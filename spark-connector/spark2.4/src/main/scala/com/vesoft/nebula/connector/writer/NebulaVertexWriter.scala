/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.writer

import com.vesoft.nebula.spark.common.nebula.VidType
import com.vesoft.nebula.spark.common.writer.NebulaExecutor
import com.vesoft.nebula.spark.common.{NebulaOptions, NebulaVertex, NebulaVertices, WriteMode}
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.sources.v2.writer.{DataWriter, WriterCommitMessage}
import org.apache.spark.sql.types.StructType
import org.slf4j.LoggerFactory

import scala.collection.mutable.ListBuffer

class NebulaVertexWriter(nebulaOptions: NebulaOptions, vertexIndex: Int, schema: StructType)
    extends NebulaWriter(nebulaOptions)
    with DataWriter[InternalRow] {

  private val LOG = LoggerFactory.getLogger(this.getClass)

  val fieldTypeMap: Map[String, String] =
    graphProvider.getNodeSchema(nebulaOptions.graphName, nebulaOptions.label)

  /** buffer to save batch vertices */
  var vertices: ListBuffer[NebulaVertex] = new ListBuffer()

  private val isIdStringType: Boolean = graphProvider.getIdType(
    nebulaOptions.graphName,
    nebulaOptions.label) == VidType.STRING

  /**
    * write one vertex row to buffer
    */
  override def write(row: InternalRow): Unit = {
    val vertex =
      NebulaExecutor.extraID(schema, row, vertexIndex, isIdStringType)
    val values =
      if (nebulaOptions.writeMode == WriteMode.DELETE) {
        // delete mode does not need property.
        Map[String, String]()
      } else {
        NebulaExecutor.assignVertexPropValues(schema,
                                              row,
                                              vertexIndex,
                                              nebulaOptions.vidAsProp,
                                              fieldTypeMap)
      }
    val nebulaVertex = NebulaVertex(vertex, values)
    vertices.append(nebulaVertex)
    if (vertices.size >= nebulaOptions.batch) {
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

  def writeNodes(vertices: ListBuffer[NebulaVertex]): Unit = {
    val exec   = getGql(vertices.toList)
    val result = submit(exec)
    if (result.isSucceeded) {
      if (!nebulaOptions.disableWriteLog) {
        LOG.info(
          s"batch write for ${nebulaOptions.label} succeed. batch size(${vertices.size}, latency(${result.getLatency}))")
      }
    } else {
      // re-execute the vertices one by one
      LOG.warn(
        s"write node ${nebulaOptions.label} failed error message: ${result.getGqlStatus}, " +
          s"mow retry writing one by one.\n ${exec}")
      vertices.par.foreach { vertex =>
        writeNode(vertex)
      }
    }
  }

  def writeNode(vertex: NebulaVertex): Unit = {
    val exec   = getGql(List(vertex))
    val result = submit(exec)
    if (result.isSucceeded) {
      if (!nebulaOptions.disableWriteLog) {
        LOG.info(s"write ${nebulaOptions.label}, batch size(1), latency(${result.getLatency}ms)")
      }
      return
    }
    if (result.getGqlStatus.contains("(9,8000)")) {
      LOG.warn(
        s"write ${nebulaOptions.label} failed, the node already exists. node: ${vertex.toString}")
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
    LOG.error(s"write node failed for ${executeResult.getGqlStatus}, statement:\n ${exec}")
    failedExecs.append(exec)
  }

  private def getGql(nebulaVertices: List[NebulaVertex]): String = {
    val nebulaNodes = NebulaVertices(nebulaOptions.label, nebulaVertices)
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
    LOG.error("insert vertex task abort.")
    graphProvider.close()
  }
}
