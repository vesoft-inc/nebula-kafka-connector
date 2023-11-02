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

import scala.collection.mutable
import scala.collection.mutable.ListBuffer

class NebulaEdgeWriter(nebulaOptions: NebulaOptions,
                       srcIndex: Int,
                       dstIndex: Int,
                       schema: StructType)
    extends NebulaWriter(nebulaOptions)
    with DataWriter[InternalRow] {

  private val LOG = LoggerFactory.getLogger(this.getClass)

  val fieldTypeMap: Map[String, String] =
    graphProvider.getEdgeSchema(nebulaOptions.graphName, nebulaOptions.label)

  /** buffer to save batch edges */
  var edges: ListBuffer[NebulaEdge] = new ListBuffer()

  val (sourceIdType, targetIdType) =
    graphProvider.getIdsType(nebulaOptions.graphName, nebulaOptions.label)
  private val isSourceIdStringType = sourceIdType == VidType.STRING
  private val isTargetIdStringType = targetIdType == VidType.STRING

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
    val nebulaEdges = NebulaEdges(nebulaOptions.label, edges.toList)
    val exec = nebulaOptions.writeMode match {
      case WriteMode.INSERT =>
        NebulaExecutor.toExecuteSentence(nebulaOptions.graphName, nebulaOptions.label, nebulaEdges)
      case _ =>
        throw new IllegalArgumentException(s"write mode ${nebulaOptions.writeMode} not supported.")
    }
    edges.clear()
    submit(exec)
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
