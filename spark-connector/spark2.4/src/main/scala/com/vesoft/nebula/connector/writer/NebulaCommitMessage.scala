
package com.vesoft.nebula.connector.writer

import org.apache.spark.sql.sources.v2.writer.WriterCommitMessage

case class NebulaCommitMessage(executeStatements: List[String]) extends WriterCommitMessage
