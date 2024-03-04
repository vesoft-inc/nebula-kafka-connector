/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

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
