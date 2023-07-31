/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.core

import com.vesoft.nebula.common.configuration.{EdgeConfig, IDType, NebulaGraphConfigEntry}
import com.vesoft.nebula.common.connect.GraphProvider
import com.vesoft.nebula.entity.Edge
import com.vesoft.nebula.utils.{NebulaUtils, PartitionUtils}
import org.apache.log4j.Logger
import org.apache.spark.sql.{DataFrame, Dataset, Encoders, Row}
import org.apache.spark.util.LongAccumulator

import scala.collection.mutable
import scala.collection.mutable.ListBuffer
import scala.util.control.Breaks.{break, breakable}

class EdgeProcessor(data: DataFrame,
                    nebulaGraphConfigEntry: NebulaGraphConfigEntry,
                    edgeConfig: EdgeConfig,
                    failureRecords: LongAccumulator)
    extends Serializable {
  @transient
  private[this] lazy val LOG = Logger.getLogger(this.getClass)

  /**
    * convert DataSet[Row] to DataSet[Edge]
    *
    */
  def process(): Dataset[(Int, Edge)] = {
    val graphProvider = new GraphProvider(
      nebulaGraphConfigEntry.graphAddress,
      nebulaGraphConfigEntry.user,
      nebulaGraphConfigEntry.passwd,
      nebulaGraphConfigEntry.connectTimeout,
      nebulaGraphConfigEntry.requestTimeout,
      nebulaGraphConfigEntry.retryIntervalTime
    )

    val (sourceIdType, targetIdType) =
      graphProvider.getEdgeNodesIdTypes(nebulaGraphConfigEntry.graphName, edgeConfig.name)

    val propertyDataType =
      graphProvider.getEdgeSchemas(nebulaGraphConfigEntry.graphName, edgeConfig.name)

    // 校验：过滤脏数据， id数据类型校验，属性数据类型校验，属性格式校验（尤其是时间格式）。 记录错误数据，以及条数
    // 校验的同时做数据转换，转成vertex
    import data.sparkSession.implicits._
    data.mapPartitions(iter =>
      checkAndConvertEdge(iter, sourceIdType, targetIdType, propertyDataType))
  }

  /**
    *
    * @param iter
    * @param idType
    * @param propertyDataType the mapping of source field name to nebula data type
    *
    * */
  def checkAndConvertEdge(iter: Iterator[Row],
                          sourceIdType: IDType.Value,
                          targetIdType: IDType.Value,
                          propertyDataType: Map[String, String]): Iterator[(Int, Edge)] = {
    val edges = iter.map(row => {
      val sourceIndex = row.schema.fieldIndex(edgeConfig.src)
      val targetIndex = row.schema.fieldIndex(edgeConfig.dst)
      breakable({
        if (sourceIndex < 0 || row.isNullAt(sourceIndex) || targetIndex < 0 || row.isNullAt(
              targetIndex)) {
          LOG.warn(
            s"source or target primary key for edge must exist and cannot be null, but your data is $row")
          // TODO recode the dirty record
          failureRecords.add(1)
          break()
        }
      })

      var primarySrcKeyValue = row.get(sourceIndex).toString.trim
      var primaryDstKeyValue = row.get(targetIndex).toString.trim
      breakable({
        if (sourceIdType == IDType.INT && !NebulaUtils.isNumeric(primarySrcKeyValue)) {
          LOG.warn(
            s"source primary key has wrong data type, source primary key for ${edgeConfig.name} is INT, but your value is $primarySrcKeyValue")
          // TODO recode the dirty record
          failureRecords.add(1)
          break()
        }
        if (targetIdType == IDType.INT && !NebulaUtils.isNumeric(primaryDstKeyValue)) {
          LOG.warn(
            s"target primary key has wrong data type, target primary key for ${edgeConfig.name} is INT, but your value is $primaryDstKeyValue")
          // TODO recode the dirty record
          failureRecords.add(1)
          break()
        }
      })

      // convert the row into Edge
      if (sourceIdType == IDType.STRING) {
        primarySrcKeyValue = NebulaUtils.escapeUtil(primarySrcKeyValue).mkString("\"", "", "\"")
      }
      if (targetIdType == IDType.STRING) {
        primaryDstKeyValue = NebulaUtils.escapeUtil(primaryDstKeyValue).mkString("\"", "", "\"")
      }
      val properties: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
      for (i <- 0 to edgeConfig.sourceFields.size) {
        properties.put(edgeConfig.nebulaFields(i), extraValue(row, edgeConfig.sourceFields(i)))
      }
      val edge = Edge(primarySrcKeyValue, primaryDstKeyValue, properties.toMap)
      (PartitionUtils.getBucketIdForEdge(edge), edge)
    })
    edges
  }

  def extraValue(row: Row, field: String): String = {
    val index = row.schema.fieldIndex(field)
    if (row.isNullAt(index)) null else row.get(index).toString.trim
  }
}
