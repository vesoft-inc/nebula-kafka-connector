/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.producer

import com.alibaba.fastjson.JSONObject
import com.vesoft.nebula.common.configuration.{
  EdgeConfig,
  MQClusterConfigEntry,
  NebulaGraphConfigEntry,
  NodeConfig
}
import com.vesoft.nebula.common.connect.GraphProvider
import com.vesoft.nebula.entity.GQLTemplate.BATCH_INSERT_TEMPLATE
import com.vesoft.nebula.entity.{Edge, Vertex}
import com.vesoft.nebula.utils.RedpandaSink
import org.apache.log4j.Logger
import org.apache.spark.broadcast.Broadcast
import org.apache.spark.sql.Dataset
import org.apache.spark.util.LongAccumulator

import java.util
import scala.collection.mutable
import scala.collection.mutable.ListBuffer

class ImportProducerForVertex(data: Dataset[(Int, Vertex)],
                              nebulaGraphConfigEntry: NebulaGraphConfigEntry,
                              mqClusterConfigEntry: MQClusterConfigEntry,
                              nodeConfig: NodeConfig,
                              failureRecords: LongAccumulator,
                              kafkaProducer: Broadcast[RedpandaSink[String, String]]) {

  @transient
  private[this] lazy val LOG = Logger.getLogger(this.getClass)

  import data.sparkSession.implicits._
  def produceVertex(): Unit = {
    val graphProvider = new GraphProvider(
      nebulaGraphConfigEntry.graphAddress,
      nebulaGraphConfigEntry.user,
      nebulaGraphConfigEntry.passwd,
      nebulaGraphConfigEntry.connectTimeout,
      nebulaGraphConfigEntry.requestTimeout,
      nebulaGraphConfigEntry.retryIntervalTime
    )
    val schema = graphProvider.getNodeSchemas(nebulaGraphConfigEntry.graphName, nodeConfig.name)

    data
      .map(row => (row._1, List(row._2)))
      .rdd
      .reduceByKey((v1, v2) => v1 ++ v2)
      .mapPartitions(iter => {
        iter.grouped(nodeConfig.batchSize).foreach {
          batch =>
            batch.map(vertices => {
              val primaryKeys = vertices._2.map(v => v.vertexID)
              val nodeType =
                graphProvider.getNodeType(nebulaGraphConfigEntry.graphName, nodeConfig.name)
              val primaryKey2InternalId: Map[String, String] =
                graphProvider.generateInternalIds(nodeType, primaryKeys)
              val validVertices: ListBuffer[Vertex] = new ListBuffer[Vertex]
              vertices._2.foreach(v => {
                if (primaryKey2InternalId.contains(v.vertexID)) {
                  v.vertexID = primaryKey2InternalId(v.vertexID)
                  validVertices.append(v)
                } else {
                  // TODO(Anqi) record the error vertex record for not generated vertex internal id
                  failureRecords.add(1)
                }
              })
              val verticesValues = vertices._2.map(v => v.getVertexString(schema)).mkString(",")
              val statement = BATCH_INSERT_TEMPLATE
                .format(nebulaGraphConfigEntry.graphName, "NODE", nodeConfig.name, verticesValues)
              // send the statement to MQ
              val dataMap: util.HashMap[String, Object] = new util.HashMap[String, Object]()
              dataMap.put("value", statement)
              val messageJson = new JSONObject(dataMap)
              val rm = kafkaProducer.value
                .send(mqClusterConfigEntry.topic, vertices._1, "statement", messageJson.toString)
              val meta = rm.get()
              LOG.info("topic name={},partition={},offset={}",
                       meta.topic(),
                       meta.partition(),
                       meta.offset())
            })
        }
        null
      })
  }
}

class ImportProducerForEdge(data: Dataset[(Int, Edge)],
                            nebulaGraphConfigEntry: NebulaGraphConfigEntry,
                            mqClusterConfigEntry: MQClusterConfigEntry,
                            edgeConfig: EdgeConfig,
                            failureRecords: LongAccumulator,
                            kafkaProducer: Broadcast[RedpandaSink[String, String]]) {
  @transient
  private[this] lazy val LOG = Logger.getLogger(this.getClass)

  import data.sparkSession.implicits._
  def produceEdge(): Unit = {
    val graphProvider = new GraphProvider(
      nebulaGraphConfigEntry.graphAddress,
      nebulaGraphConfigEntry.user,
      nebulaGraphConfigEntry.passwd,
      nebulaGraphConfigEntry.connectTimeout,
      nebulaGraphConfigEntry.requestTimeout,
      nebulaGraphConfigEntry.retryIntervalTime
    )
    val schema = graphProvider.getEdgeSchemas(nebulaGraphConfigEntry.graphName, edgeConfig.name)
    data
      .map(row => (row._1, List(row._2)))
      .rdd
      .reduceByKey((e1, e2) => e1 ++ e2)
      .mapPartitions(iter => {
        iter.grouped(edgeConfig.batchSize).foreach {
          batch =>
            batch.map(edges => {
              val srcPrimaryKeys = edges._2.map(e => e.sourceID)
              val dstPrimaryKeys = edges._2.map(e => e.targetID)

              // query the internal ID for primarykey
              val (srcNodeType, dstNodeType) =
                graphProvider.getNodesType(nebulaGraphConfigEntry.graphName, edgeConfig.name)
              val srcPrimaryKey2Internal: Map[String, String] =
                graphProvider.getInternalIds(srcNodeType, srcPrimaryKeys)
              val dstPrimaryKeys2Internal =
                graphProvider.getInternalIds(dstNodeType, dstPrimaryKeys)

              val validEdges: ListBuffer[Edge] = new ListBuffer()
              edges._2.foreach(e => {
                if (srcPrimaryKey2Internal.contains(e.sourceID) && dstPrimaryKeys2Internal.contains(
                      e.targetID)) {
                  e.sourceID = srcPrimaryKey2Internal(e.sourceID)
                  e.targetID = dstPrimaryKeys2Internal(e.targetID)
                  validEdges.append(e)
                } else {
                  // TODO(Anqi) record the error edge record for lack of vertex internal id
                  failureRecords.add(1)
                }
              })

              val edgeValues = validEdges.map(e => e.getEdgeString(schema)).mkString(",")
              val statement = BATCH_INSERT_TEMPLATE
                .format(nebulaGraphConfigEntry.graphName, "EDGE", edgeConfig.name, edgeValues)
              // send the statement to MQ
              val dataMap: util.HashMap[String, Object] = new util.HashMap[String, Object]()
              dataMap.put("value", statement)
              val messageJson = new JSONObject(dataMap)
              val rm = kafkaProducer.value
                .send(mqClusterConfigEntry.topic, edges._1, "statement", messageJson.toString)
              val meta = rm.get()
              LOG.info("topic name={},partition={},offset={}",
                       meta.topic(),
                       meta.partition(),
                       meta.offset())
            })
        }
        null
      })
  }
}
