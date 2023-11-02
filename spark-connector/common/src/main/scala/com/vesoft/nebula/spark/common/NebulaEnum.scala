/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

object DataTypeEnum extends Enumeration {

  type DataType = Value
  val VERTEX = Value("vertex")
  val EDGE   = Value("edge")

  def validDataType(dataType: String): Boolean =
    values.exists(_.toString.equalsIgnoreCase(dataType))
}

object OperaType extends Enumeration {

  type Operation = Value
  val READ  = Value("read")
  val WRITE = Value("write")
}

object WriteMode extends Enumeration {

  type Mode = Value
  val INSERT = Value("insert")
  val UPDATE = Value("update")
  val DELETE = Value("delete")
}
