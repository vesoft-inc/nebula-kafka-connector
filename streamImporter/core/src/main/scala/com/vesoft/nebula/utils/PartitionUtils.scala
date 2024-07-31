
package com.vesoft.nebula.utils

import com.vesoft.nebula.entity.{Edge, Vertex}
import org.apache.spark.sql.{Dataset, SparkSession}

object PartitionUtils {

  def repartitionForVertex(spark: SparkSession, data: Dataset[Vertex]): Dataset[(Int, Vertex)] = {
    import spark.implicits._
    data.mapPartitions(iter => {
      iter.map(v => {
        (getBucketIdForVertex(v), v)
      })
    })
  }

  def getBucketIdForVertex(vertex: Vertex): Int = {
    val vID = vertex.vertexID
    vID.toInt % 65535
  }

  def repartitionForEdge(spark: SparkSession, data: Dataset[Edge]): Dataset[(Int, Edge)] = {
    import spark.implicits._
    data.mapPartitions(iter => {
      iter.map(e => {
        (getBucketIdForEdge(e), e)
      })
    })
  }

  def getBucketIdForEdge(edge: Edge): Int = {
    val srcID = edge.sourceID
    val dstID = edge.targetID
    srcID.toInt * dstID.toInt % 65535
  }

}
