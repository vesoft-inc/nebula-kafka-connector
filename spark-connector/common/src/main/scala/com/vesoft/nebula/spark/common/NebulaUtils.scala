/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import com.vesoft.nebula.spark.common.nebula.GraphProvider
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.types._
import org.apache.spark.unsafe.types.UTF8String
import org.slf4j.LoggerFactory

import scala.collection.mutable.ListBuffer

object NebulaUtils {
  private val LOG = LoggerFactory.getLogger(this.getClass)


  type NebulaValueGetter = (Any, InternalRow, Int) => Unit

  /**
   * make getter
   *
   * @param schema Spark DataFrame schema
   * @return list of NebulaValueGetter
   */
  def makeGetters(schema: StructType): Array[NebulaValueGetter] = {
    schema.fields.map(field => makeGetter(field.dataType))
  }

  private def makeGetter(dataType: DataType): NebulaValueGetter = {
    dataType match {
      case BooleanType =>
        (prop: Any, row: InternalRow, pos: Int) =>
          row.setBoolean(pos, prop.asInstanceOf[Boolean])
      case TimestampType | LongType =>
        (prop: Any, row: InternalRow, pos: Int) =>
          row.setLong(pos, prop.asInstanceOf[Long])
      case FloatType | DoubleType =>
        (prop: Any, row: InternalRow, pos: Int) =>
          row.setDouble(pos, prop.asInstanceOf[Double])
      case IntegerType =>
        (prop: Any, row: InternalRow, pos: Int) =>
          row.setInt(pos, prop.asInstanceOf[Int])
      case _ =>
        (prop: Any, row: InternalRow, pos: Int) =>
          row.update(pos, UTF8String.fromString(String.valueOf(prop)))
    }
  }

  /**
   * check if a str is numic
   *
   * @param str string
   * @return true if str is numic
   */
  def isNumic(str: String): Boolean =
    str.matches("-?\\d+")

  /**
   * escape the string which contains escape str
   *
   * @param str string
   * @return escaped string
   */
  def escapeUtil(str: String): String =
    str
      .replaceAll("\\\\", "\\\\\\\\")
      .replaceAll("\t", "\\\\t")
      .replaceAll("\n", "\\\\n")
      .replaceAll("\"", "\\\\\"")
      .replaceAll("\'", "\\\\'")
      .replaceAll("\r", "\\\\r")
      .replaceAll("\b", "\\\\b")


  /**
   * return the dataset's schema.
   * schema includes configured cols in returnCols, if returnCols is null, return all the properties.
   * if returnCols is empty, return no properties but just vid for node and srcId, dstId for edge.
   *
   * for node, the pk name always be the first position of schema.
   * for edge, the schema fields are: src pk name, dst pk name, edge properties name
   *
   * @param nebulaOptions operations for schema
   * @return StructType
   */
  def getSchema(nebulaOptions: NebulaOptions): StructType = {
    var returnCols = nebulaOptions.getReturnCols
    val graphProvider = new GraphProvider(nebulaOptions.graphAddress, nebulaOptions.user, nebulaOptions.passwd, nebulaOptions.timeout)
    val isNodeType = DataTypeEnum.NODE.toString.equalsIgnoreCase(nebulaOptions.dataType)

    val fields: ListBuffer[StructField] = new ListBuffer[StructField]
    if (isNodeType) {
      val nodeDesc = graphProvider.getNodeDesc(nebulaOptions.graphName, nebulaOptions.label)
      val pk = nodeDesc.nodePkName
      fields.append(DataTypes.createStructField(pk, DataTypes.StringType, false))
      // if returnCols is null, read all the property of node type/edge type
      if (returnCols == null) {
        returnCols = nodeDesc.properties.keySet.toList
      }
      // add node returnCols name to Spark schema's fields
      for (propName <- returnCols) {
        if (!propName.equals(pk)) {
          fields.append(DataTypes.createStructField(propName, DataTypes.StringType, true))
        }
      }
      new StructType(fields.toArray)
    } else {
      val edgeDesc = graphProvider.getEdgeDesc(nebulaOptions.graphName, nebulaOptions.label)
      fields.append(DataTypes.createStructField(edgeDesc.srcNodePkName, DataTypes.StringType, false))
      fields.append(DataTypes.createStructField(edgeDesc.dstNodePkName, DataTypes.StringType, false))
      // if returnCols is null, read all the property of node type/edge type
      if (returnCols == null) {
        returnCols = edgeDesc.properties.keySet.toList
      }
      // add edge returnCols name to Spark schema's fields
      for (propName <- returnCols) {
        // if edge property has the same name with src/dst node's pk name, rename it with suffix $
        val finalPropName =
          if (propName.equals(edgeDesc.srcNodePkName) || propName.equals(edgeDesc.dstNodePkName)) {
            propName + "$"
          } else {
            propName
          }
        fields.append(DataTypes.createStructField(finalPropName, DataTypes.StringType, true))
      }
      new StructType(fields.toArray)
    }
  }
}
