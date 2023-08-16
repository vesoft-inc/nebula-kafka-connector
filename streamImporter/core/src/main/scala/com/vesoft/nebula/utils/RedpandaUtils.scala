/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.utils

import com.vesoft.nebula.common.configuration.MQClusterConfigEntry
import org.apache.kafka.clients.admin.{AdminClient, KafkaAdminClient, NewTopic}
import org.apache.kafka.clients.producer.{KafkaProducer, ProducerRecord, RecordMetadata}

import java.util
import java.util.Properties

class RedpandaProvider() extends Serializable {

  /**
    * create topic for redpanda
    * @param mqClusterConfigEntry redpanda serveres, topic, replic for topic
    * @param partition partition for topic, the same with NebulaGraph graph's total bucket number
    */
  def createTopic(mqClusterConfigEntry: MQClusterConfigEntry, partition: Int): Unit = {
    val bootstrapServers = mqClusterConfigEntry.server
    val topic            = mqClusterConfigEntry.topic
    val props            = new Properties()
    props.put("bootstrap.servers", bootstrapServers)
    props.put("acks", "all")
    props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer")
    props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer")
    val adminClient: KafkaAdminClient = AdminClient.create(props).asInstanceOf[KafkaAdminClient]
    val topics: util.List[NewTopic]   = new util.ArrayList[NewTopic]()
    topics.add(new NewTopic(topic, partition, mqClusterConfigEntry.replic.toShort))
    adminClient.createTopics(topics)
    adminClient.close()
  }
}

class RedpandaSink[K, V](createProducer: () => KafkaProducer[K, V]) extends Serializable {
  lazy val producer = createProducer()

  /**
    * send value into specific topic and partition
    */
  def send(topic: String,
           partition: Int,
           key: K,
           value: V): java.util.concurrent.Future[RecordMetadata] =
    producer.send(new ProducerRecord[K, V](topic, partition, key, value))
}

object RedpandaSink {
  import scala.collection.JavaConversions._
  def apply[K, V](config: Map[String, Object]): RedpandaSink[K, V] = {
    val createProducerFunc = () => {
      val producer = new KafkaProducer[K, V](config)
      sys.addShutdownHook {
        producer.close()
      }
      producer
    }
    new RedpandaSink(createProducerFunc)
  }
  def apply[K, V](config: java.util.Properties): RedpandaSink[K, V] = apply(config.toMap)
}
