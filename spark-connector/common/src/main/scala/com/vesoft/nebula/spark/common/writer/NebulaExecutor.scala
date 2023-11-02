/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common.writer

import com.vesoft.nebula.spark.common.{NebulaEdges, NebulaUtils, NebulaVertices}
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.types.StructType
import org.slf4j.LoggerFactory

import scala.collection.mutable

object NebulaExecutor {

  private val LOG = LoggerFactory.getLogger(this.getClass)

  /**
    * deal with vertex ID, used to extra vertex's id and edge's srcId,dstId
    * @param schema
    * @param record
    * @param index
    * @param policy
    * @param isVidStringType true if vid_type is Fix_String
    */
  def extraID(schema: StructType,
              record: InternalRow,
              index: Int,
              isVidStringType: Boolean): String = {
    val types = schema.fields.map(field => field.dataType)
    if (record.isNullAt(index)) {
      LOG.warn(s">>>> record has null value at index $index, ignore it. record:$record")
    }
    val vid = record.get(index, types(index)).toString
    if (isVidStringType) {
      NebulaUtils.escapeUtil(vid).mkString("\"", "", "\"")
    } else {
      if (!NebulaUtils.isNumic(vid)) {
        LOG.error(
          s">>>> record has wrong value (expect numeric value) at index $index, ignore it. record:$record")
      }
      vid
    }
  }

  /**
    * deal with vertex property values
    * @param schema
    * @param record
    * @param vertexIndex
    * @param fieldTypeMap
    *
    * @return Map of property name to property value
    * */
  def assignVertexPropValues(schema: StructType,
                             record: InternalRow,
                             vertexIndex: Int,
                             vidAsProp: Boolean,
                             fieldTypeMap: Map[String, String]): Map[String, String] = {
    val properties: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    for {
      index <- schema.fields.indices
      if vidAsProp || index != vertexIndex
    } yield {
      properties += (schema.fields(index).name -> extraValue(record, schema, index, fieldTypeMap))
    }
    properties.toMap
  }

  /**
    * deal with edge property values
    * @param schema
    * @param record
    * @param srcIndex
    * @param dstIndex
    * @param fieldTypeMap
    */
  def assignEdgeValues(schema: StructType,
                       record: InternalRow,
                       srcIndex: Int,
                       dstIndex: Int,
                       srcAsProp: Boolean,
                       dstAsProp: Boolean,
                       fieldTypeMap: Map[String, String]): Map[String, String] = {
    val properties: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    for {
      index <- schema.fields.indices
      if (srcAsProp || index != srcIndex) && (dstAsProp || index != dstIndex)
    } yield {
      properties += (schema.fields(index).name -> extraValue(record, schema, index, fieldTypeMap))
    }
    properties.toMap
  }

  /**
    * get and convert property value
    *
    * @param record DataFrame internal row
    * @param schema DataFrame schema
    * @param index  the position of row columns
    * @param fieldTypeMap property name -> property datatype in nebula
    */
  private[this] def extraValue(record: InternalRow,
                               schema: StructType,
                               index: Int,
                               fieldTypeMap: Map[String, String]): String = {
    if (record.isNullAt(index)) return null

    val types                  = schema.fields.map(field => field.dataType)
    val propValue              = record.get(index, types(index))
    val propValueTypeClassName = propValue.getClass.getName
    val simpleName = propValueTypeClassName.substring(propValueTypeClassName.lastIndexOf(".") + 1,
                                                      propValueTypeClassName.length)

    val fieldName = schema.fields(index).name
    fieldTypeMap(fieldName) match {
      case "STRING" =>
        NebulaUtils.escapeUtil(propValue.toString).mkString("\"", "", "\"")
      case "DATE"          => "date(\"" + propValue + "\")"
      case "LOCALDATETIME" => "localdatetime(\"" + propValue + "\")"
      case "LOCALTIME"        => "localtime(\"" + propValue + "\")"
      case _ =>
        if (simpleName.equalsIgnoreCase("UTF8String")) propValue.toString
        else propValue.toString
    }
  }

  /**
    * construct insert statement for vertex
    */
  def toExecuteSentence(graphName: String, nodeType: String, vertices: NebulaVertices): String = {
    s"USE $graphName INSERT NODE `$nodeType` ${vertices.getVerticesStr}"
  }

  /**
    * construct insert statement for edge
    */
  def toExecuteSentence(graphName: String, edgeType: String, edges: NebulaEdges): String = {
    s"USE $graphName INSERT EDGE `$edgeType` ${edges.getEdgesStr}"
  }

  /**
    * escape nebula property name, add `` for each property.
    *
    * @param nebulaFields nebula property name list
    * @return escaped nebula property name list
    */
  def escapePropName(nebulaFields: List[String]): List[String] =
    nebulaFields.map(key => s"`$key`")

}
