
package com.vesoft.nebula.common.configuration

sealed trait SchemaConfig {

  /** nebula node type or edge tape */
  def name: String

  /** source data fields */
  def sourceFields: List[String]

  /** nebula property fields */
  def nebulaFields: List[String]

  /** node or edge amount of one write request */
  def batchSize: Int

  /** data partition for label $name to write into NebulaGraph */
  def partition: Int

}

/**
  * note type config
  */
case class NodeConfig(override val name: String,
                      override val sourceFields: List[String],
                      override val nebulaFields: List[String],
                      override val batchSize: Int,
                      override val partition: Int,
                      vid: String)
    extends SchemaConfig {
  override def toString: String =
    s"NodeConfig{name:$name, sourceFields:$sourceFields, nebulaFields:$nebulaFields, batchSize:$batchSize, vid:$vid, partition:$partition}"
}

/**
  * edge type config
  */
case class EdgeConfig(override val name: String,
                      override val sourceFields: List[String],
                      override val nebulaFields: List[String],
                      override val batchSize: Int,
                      override val partition: Int,
                      src: String,
                      dst: String)
    extends SchemaConfig {
  override def toString: String =
    s"EdgeConfig{name:$name, sourceFields:$sourceFields, nebulaFields:$nebulaFields, batchSize:$batchSize, src:$src, dst:$dst, partition:$partition}"
}
