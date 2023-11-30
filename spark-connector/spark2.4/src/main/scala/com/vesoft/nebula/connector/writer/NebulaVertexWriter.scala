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

import scala.collection.mutable
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
    val nebulaVertices = NebulaVertices(nebulaOptions.label, vertices.toList)
    val exec = nebulaOptions.writeMode match {
      case WriteMode.INSERT =>
        NebulaExecutor.toExecuteSentence(nebulaOptions.graphName,
                                         nebulaOptions.label,
                                         nebulaVertices)
      case _ =>
        throw new IllegalArgumentException(s"write mode ${nebulaOptions.writeMode} not supported.")
    }
    vertices.clear()
    submit(exec)
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
