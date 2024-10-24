
package com.vesoft.nebula.spark.common.nebula


case class NodeDesc(nodeTypeName: String,
                    nodePkNames: List[String],
                    properties: Map[String, String])

case class EdgeDesc(edgeTypeName: String,
                    srcNodeTypeName: String,
                    srcNodePkNames: List[String],
                    srcNodePkDataTypeMap: Map[String, String],
                    dstNodeTypeName: String,
                    dstNodePkNames: List[String],
                    dstNodePkDataTypeMap: Map[String, String],
                    properties: Map[String, String])
