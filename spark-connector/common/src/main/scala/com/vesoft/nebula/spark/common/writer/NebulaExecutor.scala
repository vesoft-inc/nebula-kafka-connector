/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common.writer

import com.vesoft.nebula.spark.common.{NebulaEdges, NebulaNodes, NebulaUtils}
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.types.StructType
import org.slf4j.LoggerFactory

import scala.collection.mutable

object NebulaExecutor {

  private val LOG = LoggerFactory.getLogger(this.getClass)

  /**
   * used to extra edge's src node primary key,dst node primary key
   *
   * @param schema
   * @param record
   * @param index
   * @param policy
   * @param isPkStringType true if primary key is String data type
   */
  def extraID(schema: StructType,
              record: InternalRow,
              index: Int,
              isPkStringType: Boolean): String = {
    val types = schema.fields.map(field => field.dataType)
    if (record.isNullAt(index)) {
      return null
    }
    val vid = record.get(index, types(index)).toString
    if (isPkStringType) {
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
   * deal with node property values
   *
   * @param schema
   * @param record
   * @param fieldTypeMap
   * @return Map of property name to property value
   * */
  def assignNodePropValues(schema: StructType,
                           record: InternalRow,
                           fieldTypeMap: Map[String, String]): Map[String, String] = {
    val properties: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    for {
      index <- schema.fields.indices
    } yield {
      properties += (schema.fields(index).name -> extraValue(record, schema, index, fieldTypeMap))
    }
    properties.toMap
  }

  /**
   * deal with node pk value
   * this function just need to deal with the primary key, not all the properties
   *
   * @param schema       spark dataframe schema
   * @param record       spark dataframe row
   * @param fieldTypeMap field name and data type map in NebulaGraph
   * @return Map of pk name to property value
   */
  def assignNodePkValues(schema: StructType,
                         record: InternalRow,
                         fieldTypeMap: Map[String, String],
                         pkName: String): Map[String, String] = {
    val properties: mutable.HashMap[String, String] = new mutable.HashMap[String, String]()
    val index = schema.fieldIndex(pkName)
    properties += (pkName -> extraValue(record, schema, index, fieldTypeMap))
    properties.toMap
  }

  /**
   * deal with edge property values
   *
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
   * @param record       DataFrame internal row
   * @param schema       DataFrame schema
   * @param index        the position of row columns
   * @param fieldTypeMap property name -> property datatype in nebula
   */
  private[this] def extraValue(record: InternalRow,
                               schema: StructType,
                               index: Int,
                               fieldTypeMap: Map[String, String]): String = {
    if (record.isNullAt(index)) return null

    val types = schema.fields.map(field => field.dataType)
    val propValue = record.get(index, types(index))
    val propValueTypeClassName = propValue.getClass.getName
    val simpleName = propValueTypeClassName.substring(propValueTypeClassName.lastIndexOf(".") + 1,
      propValueTypeClassName.length)

    val fieldName = schema.fields(index).name
    fieldTypeMap(fieldName) match {
      case "STRING" =>
        NebulaUtils.escapeUtil(propValue.toString).mkString("\"", "", "\"")
      case "DATE" => "date(\"" + propValue + "\")"
      case "LOCAL DATETIME" => "local_datetime(\"" + propValue + "\")"
      case "LOCAL TIME" => "local_time(\"" + propValue + "\")"
      case "ZONED TIME" => "zoned_time(\"" + propValue + " \")"
      case "ZONED DATETIME" => "zoned_datetime(\"" + propValue + " \")"
      case _ =>
        if (simpleName.equalsIgnoreCase("UTF8String")) propValue.toString
        else propValue.toString
    }
  }

  /**
   * construct insert statement for node
   * TABLE t {id,firstName,lastName, tag_name} =
   * (1, "f1", "l1", "tag1"),
   * (2, "f2", "l2", "tag2"),
   * (3, "f3", "l3", "tag3")
   * USE nba
   * FOR r IN t
   * INSERT (a@Person{id:r.id,firstName:r.firstName,lastName:r.lastName})
   */
  def toExecuteSentence(graphName: String, nodes: NebulaNodes, mode:String): String = {
    s"""
       |TABLE t {${nodes.propNamesStr}} =
       |${nodes.getNodesStr}
       |USE `$graphName`
       |FOR r IN t
       |INSERT $mode (@`${nodes.nodeType}`{${nodes.propNamesWithTableStr}})
       |""".stripMargin
  }

  /**
   * construct insert statement for edge
   */
  def toExecuteSentence(graphName: String, edgeType: String, edges: NebulaEdges): String = {
    s"USE `$graphName` INSERT EDGE `$edgeType` ${edges.getEdgesStr}"
  }

  /**
   * construct update statement for node
   */
  def toUpdateSentence(graphName: String, nodeType: String, nodes: NebulaNodes): String = {
    null
  }


  /**
   * construct update statement for edge
   */
  def toUpdateSentence(graphName: String, edgeType: String, edges: NebulaEdges): String = {
    null
  }


  /**
   * construct delete statement for node
   * USE graphName MATCH((a@nodeType where a.id in [ 1,2,3,4,5]-[r]-(b)) DETACH DELETE a
   */
  def toDeleteSentence(graphName: String, nodeType: String, nodes: NebulaNodes): String = {
    val nodePks = nodes.values.map(node => node.values(nodes.pkName)).mkString(",")
    s"USE `$graphName` MATCH(a@$nodeType where a.${nodes.pkName} in [$nodePks]) DETACH DELETE a "
  }


  /**
   * construct delete statement for edge
   * USE graphName MATCH (a@Person{id:2})-[e@edgeType]->(b@Person{id:3}) DELETE e
   */
  def toDeleteSentence(graphName: String, edgeType: String, edges: NebulaEdges): String = {
    val edge = edges.values.iterator.next()
    val edgeMatchStatement = s"MATCH(a@${edges.srcType}{${edge.srcPkName}:${edge.srcId}})-[e@$edgeType]-(b@${edges.dstType}{${edge.dstPkName}:${edge.dstId}}) DELETE e"
    s"USE `$graphName` $edgeMatchStatement"
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
