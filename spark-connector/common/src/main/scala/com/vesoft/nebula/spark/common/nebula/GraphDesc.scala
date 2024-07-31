
package com.vesoft.nebula.spark.common.nebula


case class NodeDesc(nodeTypeName: String,
                    nodePkName: String,
                    nodePkDataType: VidType.Value,
                    properties: Map[String, String])

case class EdgeDesc(edgeTypeName: String,
                    srcNodeTypeName: String,
                    srcNodePkName:String,
                    srcNodePkDataType: VidType.Value,
                    dstNodeTypeName: String,
                    dstNodePkName:String,
                    dstNodePkDataType: VidType.Value,
                    properties: Map[String, String])
