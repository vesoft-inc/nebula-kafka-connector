
package com.vesoft.nebula.spark.connector.writer

import com.vesoft.nebula.spark.common.writer.NodeWriter
import com.vesoft.nebula.spark.common.NebulaOptions
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.sources.v2.writer.{DataWriter, WriterCommitMessage}
import org.apache.spark.sql.types.StructType

class NebulaNodeWriter(nebulaOptions: NebulaOptions, schema: StructType)
  extends NodeWriter(nebulaOptions, schema)
    with DataWriter[InternalRow] {


  /**
   * write one node row to buffer
   */
  override def write(row: InternalRow): Unit = {
    writeRow(row)
  }

  override def commit(): WriterCommitMessage = {
    if (nodes.nonEmpty) {
      execute()
    }
    graphProvider.close()
    NebulaCommitMessage(failedExecs.toList)
  }

  override def abort(): Unit = {
    abortWriter()
  }
}
