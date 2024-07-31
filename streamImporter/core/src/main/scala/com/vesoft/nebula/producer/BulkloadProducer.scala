
package com.vesoft.nebula.producer

import com.alibaba.fastjson.JSON
import com.alibaba.fastjson.serializer.SerializerFeature
import com.vesoft.nebula.common.configuration.{EdgeConfig, MQClusterConfigEntry, NebulaGraphConfigEntry, NodeConfig}
import com.vesoft.nebula.common.connect.GraphProvider
import com.vesoft.nebula.entity.{Edge, EdgeToWrite, NebulaEdgesRequest, NebulaNodesRequest, NodeToWrite, Vertex}
import com.vesoft.nebula.utils.RedpandaSink
import org.apache.log4j.Logger
import org.apache.spark.broadcast.Broadcast
import org.apache.spark.sql.Dataset
import org.apache.spark.util.LongAccumulator

import scala.collection.JavaConverters.{mapAsJavaMapConverter, seqAsJavaListConverter}
import scala.collection.mutable.ListBuffer

class BulkloadProducerForVertex(data: Dataset[(Int, Vertex)],
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
    val schema   = graphProvider.getNodeSchemas(nebulaGraphConfigEntry.graphName, nodeConfig.name)
    val graphId  = graphProvider.getGraphId(nebulaGraphConfigEntry.graphName)
    val nodeType = graphProvider.getNodeType(nebulaGraphConfigEntry.graphName, nodeConfig.name)
    val nodeTypeId =
      graphProvider.getNodeTypeId(nebulaGraphConfigEntry.graphName, nodeConfig.name)

    data
      .map(row => (row._1, List(row._2)))
      .rdd
      .reduceByKey((v1, v2) => v1 ++ v2)
      .mapPartitions(iter => {
        iter.grouped(nodeConfig.batchSize).foreach {
          batch =>
            batch.map(vertices => {
              val primaryKeys = vertices._2.map(v => v.vertexID)
              val primaryKey2InternalId: Map[String, String] =
                graphProvider.generateInternalIds(nodeType, primaryKeys)
              val validVertices: ListBuffer[Vertex] = new ListBuffer[Vertex]
              // validate the internal id generation for the vertex primaryKey
              vertices._2.foreach(v => {
                if (primaryKey2InternalId.contains(v.vertexID)) {
                  v.vertexID = primaryKey2InternalId(v.vertexID)
                  validVertices.append(v)
                } else {
                  // TODO(Anqi) record the error vertex record for not generated vertex internal id
                  failureRecords.add(1)
                }
              })

              val nebulaNodesRequest = new NebulaNodesRequest
              nebulaNodesRequest.setGraphId(graphId)
              nebulaNodesRequest.setNodeTypeId(nodeTypeId)
              nebulaNodesRequest.setPartId(vertices._1)
              val nodesList: ListBuffer[NodeToWrite] = new ListBuffer[NodeToWrite]
              validVertices.foreach(v =>
                // TODO(Anqi) convert the property value to Nebula Data type
                nodesList.append(new NodeToWrite(v.vertexID.toLong, v.values.asJava)))
              nebulaNodesRequest.setNodes(nodesList.toList.asJava)

              val messageJson = JSON.toJSONString(nebulaNodesRequest, SerializerFeature.PrettyFormat)
              val rm = kafkaProducer.value
                .send(mqClusterConfigEntry.topic, vertices._1, "vertices", messageJson)
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

class BulkloadProducerForEdge(data: Dataset[(Int, Edge)],
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
    val schema  = graphProvider.getEdgeSchemas(nebulaGraphConfigEntry.graphName, edgeConfig.name)
    val graphId = graphProvider.getGraphId(nebulaGraphConfigEntry.graphName)
    val edgeTypeId =
      graphProvider.getEdgeTypeId(nebulaGraphConfigEntry.graphName, edgeConfig.name)

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
              val nebulaEdgesRequest = new NebulaEdgesRequest
              nebulaEdgesRequest.setGraphId(graphId)
              nebulaEdgesRequest.setPartId(edges._1)
              nebulaEdgesRequest.setEdgeTypeId(edgeTypeId)
              val edgeList: ListBuffer[EdgeToWrite] = new ListBuffer[EdgeToWrite]
              validEdges.foreach(
                e =>
                  edgeList.append(
                    // TODO(Anqi) convert the property value to Nebula Data type
                    new EdgeToWrite(e.sourceID.toLong, e.targetID.toLong, e.values.asJava)))
              nebulaEdgesRequest.setEdges(edgeList.toList.asJava)
              val message = JSON.toJSONString(nebulaEdgesRequest, SerializerFeature.PrettyFormat)
              val rm =
                kafkaProducer.value.send(mqClusterConfigEntry.topic, edges._1, "edges", message)
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
