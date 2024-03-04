/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.mock

import com.vesoft.nebula.client.graph.net.NebulaClient
import org.apache.log4j.Logger

class NebulaGraphMock {
  private[this] val LOG = Logger.getLogger(this.getClass)

  val graphAddr = "192.168.8.6:3820"
  val user = "root"
  val passwd = "nebula"
  @transient val client: NebulaClient = NebulaClient.builder(graphAddr, user, passwd).build()

  def mockReadGraph(): Unit = {
    try {
      val createGraphType = "CREATE GRAPH TYPE IF NOT EXISTS spark_read_type AS {(node_player LABEL player {col1 INT PRIMARY KEY, col2 STRING, col3 FLOAT, col4 bool, col5 DOUBLE, col6 INT64, col7 local time, col8 local datetime, col9 zoned time})," +
        "(node_player)-[edge_follow LABEL follow {col1 INT, col2 FLOAT64, col3 STRING, col4 DOUBLE}]->(node_player)}"

      var result = client.execute(createGraphType)
      if (!result.isSucceeded) {
        LOG.error(s"create graph type spark_read_type failed: ${result.getErrorMessage}")
        System.exit(1)
      }
      Thread.sleep(3000)


      val createGraph = "CREATE GRAPH IF NOT EXISTS spark_read spark_read_type"
      result = client.execute(createGraph)
      if (!result.isSucceeded) {
        LOG.error("create graph spark_read failed: $result.getErrorMessage")
        System.exit(1)
      }
      Thread.sleep(3000)

      val insertNode = "USE spark_read INSERT NODE node_player " +
        "({col1:1, col2:\"Tim\", col3: 11.0, col4: true, col5: 7.99, col6: 10, col7: local_time(\"12:00:00\"), col8:local_datetime(\"2024-01-01T12:00:00.123\"), col9: zoned_time(\"12:00:00.000+0800\")})," +
        "({col1:2, col2:\"Tom\", col3: 12.0, col4: false, col5: 8.18, col6: 11, col7: local_time(\"13:00:00.222\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0830\")})," +
        "({col1:3, col2:\"Bob\", col3: 13.1, col4: true, col5: 9.76, col6: 12, col7: local_time(\"10:10:20.123\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0200\")})," +
        "({col1:4, col2:\"Jena\", col3: 14.5, col4: true, col5: 0.55, col6: 13, col7: local_time(\"12:00:01.999\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0230\")})," +
        "({col1:5, col2:\"Alex\", col3: 15.6, col4: true, col5: 1.42, col6: 14, col7: local_time(\"12:00:00\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0900\")})," +
        "({col1:6, col2:\"James\", col3: 16.7, col4: true, col5: 2.34, col6: 15, col7: local_time(\"12:00:00\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0900\")})," +
        "({col1:7, col2:\"Emma\", col3: 17.8, col4: true, col5: 3.65, col6: 16, col7: local_time(\"12:00:00\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0800\")})," +
        "({col1:8, col2:\"Michael\", col3: 18.1, col4: true, col5: 4.12, col6: 17, col7: local_time(\"12:00:00\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0000\")})," +
        "({col1:9, col2:\"Sarah\", col3: 19.2, col4: false, col5: 15.01, col6: 18, col7: local_time(\"12:00:00\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0800\")})," +
        "({col1:10, col2:\"Jessica\", col3: 20.6, col4: false, col5: 21.56, col6: 19, col7: local_time(\"12:00:00\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0800\")})," +
        "({col1:-1, col2:\"William\", col3: 21.5, col4: false, col5: 45.43, col6: 20, col7: null, col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0800\")})," +
        "({col1:-2, col2:\"Robert\", col3: 22.8, col4: false, col5: 21.43, col6: 21, col7: local_time(\"12:00:00\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0800\")})," +
        "({col1:-3, col2:\"David\", col3: 23.9, col4: false, col5: 72.15, col6: 22, col7: local_time(\"12:00:00\"), col8:local_datetime(\"2024-01-01T12:00:00\"), col9: zoned_time(\"12:00:00.000+0800\")})"
      result = client.execute(insertNode)
      if (!result.isSucceeded) {
        LOG.error(s"insert node for graph spark_read failed: ${result.getErrorMessage}")
        System.exit(1)
      }

      val insertEdge = "USE spark_read INSERT EDGE edge_follow " +
        "({col1:1})-[{col1:90, col2: 66.8, col3: \"A\", col4: 1.0}]->({col1:2})," +
        "({col1:2})-[{col1:90, col2: 66.8, col3: \"B\", col4: 1.0}]->({col1:3})," +
        "({col1:4})-[{col1:90, col2: 66.8, col3: \"C\", col4: 1.0}]->({col1:5})," +
        "({col1:5})-[{col1:90, col2: 66.8, col3: \"D\", col4: 1.0}]->({col1:1})," +
        "({col1:6})-[{col1:90, col2: 66.8, col3: \"E\", col4: 1.0}]->({col1:8})," +
        "({col1:7})-[{col1:90, col2: 66.8, col3: \"F\", col4: 1.0}]->({col1:9})," +
        "({col1:8})-[{col1:90, col2: 66.8, col3: \"G\", col4: 1.0}]->({col1:10})," +
        "({col1:9})-[{col1:90, col2: 66.8, col3: \"H\", col4: 1.0}]->({col1:-1})," +
        "({col1:-2})-[{col1:90, col2: 66.8, col3: null, col4: 1.0}]->({col1:-1})," +
        "({col1:-3})-[{col1:90, col2: 66.8, col3: \"\", col4: 1.0}]->({col1:1})"
      result = client.execute(insertEdge)
      if (!result.isSucceeded) {
        LOG.error("insert edge for graph {} failed: {}", "spark_read", result.getErrorMessage)
        System.exit(1)
      }
    } catch {
      case e: Exception => LOG.error("mock read graph failed", e)
        System.exit(1)
    } finally {
      client.close()
    }
  }

}
