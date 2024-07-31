
package com.vesoft.nebula.core

import com.vesoft.nebula.common.configuration.{EdgeConfig, NodeConfig, SchemaConfig}
import org.apache.spark.sql.{DataFrame, Row}
import org.apache.spark.util.LongAccumulator

import scala.collection.mutable.ListBuffer

class Processor(data: DataFrame, configs: List[SchemaConfig], dirtyRecords: LongAccumulator)
    extends Serializable {

  def filterDirtyData(): DataFrame = {
    val notNullFields: ListBuffer[String] = new ListBuffer[String]
    notNullFields.appendAll(getIdFields ++ getNonEmptyPropertyNames)

    // val condition: String = notNullFields.map(field => s"$field is not null").mkString(" and ")
    // data.filter(condition)
    data.filter(row => { filterFun(notNullFields.toList, row) })
  }

  def filterFun(notNullFields: List[String], row: Row): Boolean = {
    val dirtyValue: ListBuffer[Boolean] = new ListBuffer[Boolean]
    notNullFields.foreach(field => {
      val index = row.fieldIndex(field)
      dirtyValue.append(row.isNullAt(index))
    })
    val isDirty = dirtyValue.collectFirst { case value if value => true }.getOrElse(false)
    if (isDirty) {
      dirtyRecords.add(1)
      false
    } else {
      true
    }
  }

  /**
    * get Node Id field and Edge srcId field and dstId field
    */
  def getIdFields: List[String] = {
    val fields: ListBuffer[String] = new ListBuffer
    configs.foreach {
      case nodeConfig: NodeConfig =>
        fields.append(nodeConfig.vid)
      case edgeConfig: EdgeConfig =>
        fields.append(edgeConfig.src)
        fields.append(edgeConfig.dst)
    }
    fields.toList.distinct
  }

  /**
    * get non-empty NebulaGraph property names
    * TODO(Anqi) get non-empty properties through schema
    * */
  def getNonEmptyPropertyNames: List[String] = {
    List()
  }
}
