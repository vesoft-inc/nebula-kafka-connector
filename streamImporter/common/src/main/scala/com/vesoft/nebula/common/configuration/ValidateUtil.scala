
package com.vesoft.nebula.common.configuration

import com.google.common.net.{InetAddresses, InternetDomainName}
import com.vesoft.nebula.common.connect.GraphProvider
import org.apache.log4j.Logger

import java.net.InetAddress

class ValidateUtil {
  private[this] val LOG = Logger.getLogger(this.getClass)
  def validateServer(server: String): Boolean = {
    if (server == null || server.isEmpty) return false
    for (addr <- server.split(",")) {
      val hostAndPort: Array[String] = addr.split(":")
      if (hostAndPort.length != 2) {
        return false
      }
      val host: String = hostAndPort(0)
      val port: Int    = hostAndPort(1).toInt
      // get all host name
      val inetAddresses: Array[InetAddress] = InetAddress.getAllByName(host)
      for (inetAddress <- inetAddresses) {
        val ip: String = inetAddress.getHostAddress
        if (!(InetAddresses.isInetAddress(ip) || InetAddresses.isUriInetAddress(ip) || InternetDomainName
              .isValid(ip)) || (port <= 0 || port >= 65535)) {
          return false
        }
      }
    }
    true
  }

  def validateSourceSchema(configs: Configs): Unit = {
    val nebulaGraphConfigEntry = configs.nebulaGraphConfigEntry
    val graphProvider = new GraphProvider(
      nebulaGraphConfigEntry.graphAddress,
      nebulaGraphConfigEntry.user,
      nebulaGraphConfigEntry.passwd,
      nebulaGraphConfigEntry.connectTimeout,
      nebulaGraphConfigEntry.requestTimeout,
      nebulaGraphConfigEntry.retryIntervalTime
    )
    val graphName = nebulaGraphConfigEntry.graphName
    val sources   = configs.sourceConfigEntrys
    for (source <- sources) {
      val schemaConfigs = source.schemaConfigs
      for (schema <- schemaConfigs) {
        schema match {
          case nodeConfig: NodeConfig =>
            val schemaMap     = graphProvider.getNodeSchemas(graphName, nodeConfig.name)
            val nodePropNames = schemaMap.keySet.toList
            require(
              nodeConfig.nebulaFields.size == nodeConfig.sourceFields.size,
              s"sourceFields must has the same number elements as nebulaFields for node ${nodeConfig.name}"
            )
            for (nebulaField <- nodeConfig.nebulaFields) {
              require(
                nodePropNames.contains(nebulaField),
                s"field $nebulaField in config nebulaFields does not exist in " +
                  s"node ${nodeConfig.name}, available property names are $nodePropNames."
              )
            }
          case edgeConfig: EdgeConfig =>
            val schemaMap     = graphProvider.getNodeSchemas(graphName, edgeConfig.name)
            val edgePropNames = schemaMap.keySet.toList
            require(
              edgeConfig.nebulaFields.size == edgeConfig.sourceFields.size,
              s"sourceFields must has the same number elements as nebulaFields for edge ${edgeConfig.name}"
            )
            for (nebulaField <- edgeConfig.nebulaFields) {
              require(
                edgePropNames.contains(nebulaField),
                s"field $nebulaField in config nebulaFields does not exist in " +
                  s"edge ${edgeConfig.name}, available property names are $edgePropNames."
              )
            }
        }
      }
    }

    // TODO validate if the source fields and vid field/ src field/ dst field really exist in Source
  }

}
