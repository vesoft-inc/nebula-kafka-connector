/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration

sealed trait SchemaConfigs {

  /** nebula node type or edge tape */
  def name: String

  /** source data fields */
  def sourceFields: List[String]

  /** nebula property fields */
  def nebulaFields: List[String]

  /** node or edge amount of one write request */
  def batchSize: Int

}

/**
  * note type config
  */
case class NodeConfig(override val name: String,
                      override val sourceFields: List[String],
                      override val nebulaFields: List[String],
                      override val batchSize: Int)
    extends SchemaConfigs {}

/**
  * edge type config
  */
case class EdgeConfig(override val name: String,
                      override val sourceFields: List[String],
                      override val nebulaFields: List[String],
                      override val batchSize: Int)
    extends SchemaConfigs {}
