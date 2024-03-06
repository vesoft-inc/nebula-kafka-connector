/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.reader

import com.vesoft.nebula.spark.common.{NebulaOptions, NebulaUtils}
import org.apache.spark.TaskContext
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.sources.v2.reader.{DataSourceReader, InputPartition}
import org.apache.spark.sql.types.StructType
import org.slf4j.LoggerFactory

import java.util
import scala.collection.JavaConverters._

/**
 * Base class of Nebula Source Reader
 */
abstract class NebulaSourceReader(nebulaOptions: NebulaOptions) extends DataSourceReader {
  private val LOG = LoggerFactory.getLogger(this.getClass)

  private var datasetSchema: StructType = _

  override def readSchema(): StructType = {
    if (datasetSchema == null) {
      datasetSchema = NebulaUtils.getSchema(nebulaOptions)
      LOG.info(s"dataset's schema: $datasetSchema")
    }
    datasetSchema
  }

}

/**
 * DataSourceReader for Nebula Node
 */
class NebulaDataSourceNodeReader(nebulaOptions: NebulaOptions)
  extends NebulaSourceReader(nebulaOptions) {

  override def planInputPartitions(): util.List[InputPartition[InternalRow]] = {
    val partitionNum = nebulaOptions.partitionNums
    val partitions = for (index <- 1 to partitionNum)
      yield {
        new NebulaNodePartition(index, nebulaOptions, readSchema())
      }
    partitions.map(_.asInstanceOf[InputPartition[InternalRow]]).asJava
  }
}

/**
 * DataSourceReader for Nebula Edge
 */
class NebulaDataSourceEdgeReader(nebulaOptions: NebulaOptions)
  extends NebulaSourceReader(nebulaOptions) {

  override def planInputPartitions(): util.List[InputPartition[InternalRow]] = {
    val partitionNum = nebulaOptions.partitionNums
    val partitions = for (index <- 1 to partitionNum)
      yield new NebulaEdgePartition(index, nebulaOptions, readSchema())

    partitions.map(_.asInstanceOf[InputPartition[InternalRow]]).asJava
  }
}
