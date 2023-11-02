/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import com.vesoft.nebula.client.graph.data.{NDateTime, NTime}
import org.apache.spark.sql.catalyst.InternalRow
import org.apache.spark.sql.types._
import org.apache.spark.unsafe.types.UTF8String
import org.slf4j.LoggerFactory

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
          prop match {
            case wrapper: NDateTime =>
              row.update(pos, UTF8String.fromString(wrapper.toString))
            case wrapper: NTime =>
              row.update(pos, UTF8String.fromString(wrapper.toString))
            case _ =>
              row.update(pos, UTF8String.fromString(String.valueOf(prop)))
          }
    }
  }

  /**
    * check if a str is numic
    * @param str string
    *
    * @return true if str is numic
    */
  def isNumic(str: String): Boolean =
    str.matches("-?\\d+")

  /**
    * escape the string which contains escape str
    * @param str string
    *
    * @return escaped string
    */
  def escapeUtil(str: String): String =
    str
      .replaceAll("\\\\", "\\\\\\\\")
      .replaceAll("\t", "\\\t")
      .replaceAll("\n", "\\\n")
      .replaceAll("\"", "\\\"")
      .replaceAll("\'", "\\\'")
      .replaceAll("\r", "\\\r")
      .replaceAll("\b", "\\\b")

}
