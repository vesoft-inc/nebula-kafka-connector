
package com.vesoft.nebula.core

import com.vesoft.nebula.common.configuration.{Configs, IDType, NebulaGraphConfigEntry, NodeConfig}
import com.vesoft.nebula.common.connect.GraphProvider
import com.vesoft.nebula.entity.Vertex
import com.vesoft.nebula.utils.{NebulaUtils, PartitionUtils}
import org.apache.log4j.Logger
import org.apache.spark.sql.{DataFrame, Dataset, Encoders, Row}
import org.apache.spark.util.LongAccumulator

import scala.collection.mutable
import scala.collection.mutable.ListBuffer
import scala.util.control.Breaks.{break, breakable}

class VertexProcessor(data: DataFrame,
                      nebulaGraphConfigEntry: NebulaGraphConfigEntry,
                      nodeConfig: NodeConfig,
                      failureRecords: LongAccumulator)
    extends Serializable {
  @transient
  private[this] lazy val LOG = Logger.getLogger(this.getClass)

  /**
    * convert DataSet[Row] to DataSet[Vertex]
    *
    */
  def process(): Dataset[(Int, Vertex)] = {
    val graphProvider = new GraphProvider(
      nebulaGraphConfigEntry.graphAddress,
      nebulaGraphConfigEntry.user,
      nebulaGraphConfigEntry.passwd,
      nebulaGraphConfigEntry.connectTimeout,
      nebulaGraphConfigEntry.requestTimeout,
      nebulaGraphConfigEntry.retryIntervalTime
    )

    val vertexIdType =
      graphProvider.getNodeIdType(nebulaGraphConfigEntry.graphName, nodeConfig.name)

    val propertyDataType =
      graphProvider.getNodeSchemas(nebulaGraphConfigEntry.graphName, nodeConfig.name)

    // 校验：过滤脏数据， id数据类型校验，属性数据类型校验，属性格式校验（尤其是时间格式）。 记录错误数据，以及条数
    // 校验的同时做数据转换，转成vertex
    import data.sparkSession.implicits._
    data.mapPartitions(iter => checkAndConvertVertex(iter, vertexIdType, propertyDataType))
  }

  /**
    *
    * @param iter
    * @param idType
    * @param propertyDataType the mapping of source field name to nebula data type
    *
    * */
  def checkAndConvertVertex(iter: Iterator[Row],
                            idType: IDType.Value,
                            propertyDataType: Map[String, String]): Iterator[(Int, Vertex)] = {
    val vertices = iter.map(row => {
      val index = row.schema.fieldIndex(nodeConfig.vid)
      breakable({
        if (index < 0 || row.isNullAt(index)) {
          LOG.warn(s"node primary key must exist and cannot be null, but your data is $row")
          // TODO recode the dirty record
          failureRecords.add(1)
          break
        }
      })

      var primaryKeyValue = row.get(index).toString.trim
      breakable({
        if (idType == IDType.INT && !NebulaUtils.isNumeric(primaryKeyValue)) {
          LOG.warn(
            s"node primary key has wrong data type, primary key for ${nodeConfig.name} is INT, but your value is $primaryKeyValue")
          // TODO recode the dirty record
          failureRecords.add(1)
          break
        }
      })

      // convert the row into Vertex
      if (idType == IDType.STRING) {
        primaryKeyValue = NebulaUtils.escapeUtil(primaryKeyValue).mkString("\"", "", "\"")
      }

      val properties: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
      for (i <- 0 to nodeConfig.sourceFields.size) {
        properties.put(nodeConfig.nebulaFields(i), extraValue(row, nodeConfig.sourceFields(i)))
      }
      val vertex = Vertex(primaryKeyValue, properties.toMap)
      (PartitionUtils.getBucketIdForVertex(vertex), vertex)
    })
    vertices
  }

  def extraValue(row: Row, field: String): String = {
    val index = row.schema.fieldIndex(field)
    if (row.isNullAt(index)) null else row.get(index).toString.trim
  }

}
