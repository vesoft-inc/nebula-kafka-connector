/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration

/**
  * source data pre-process rules
  * */
sealed trait DataPreProcessConfig

case class ConcatConfig(oldFields: List[String], newFiled: String, sep: String = "_")
    extends DataPreProcessConfig {
  override def toString: String =
    s"ConcatConfig{oldFields:$oldFields, newFiled:$newFiled, sep:$sep}"
}

case class SeparatorConfig(oldFiled: String, newFields: List[String], sep: String = "_")
    extends DataPreProcessConfig {
  override def toString: String =
    s"SeparatorConfig{oldField:$oldFiled, newFields:$newFields, sep:$sep}"
}

case class FilterConfig(conditions: List[String]) extends DataPreProcessConfig {
  override def toString: String = s"FilterConfig{conditions:$conditions}"
}

case class NonValueConfig(value: String) extends DataPreProcessConfig {
  override def toString: String = s"NonValueConfig{value:$value}"
}
