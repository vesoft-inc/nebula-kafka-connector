
package com.vesoft.nebula.common.configuration

import com.vesoft.nebula.common.reader.hdfsreader.HdfsSourceConfigEntry
import org.scalatest.BeforeAndAfterAll
import org.scalatest.funsuite.AnyFunSuite

class ConfigsResolveSuite extends AnyFunSuite {

  test("resolve config file") {
    val configFilePath = "src/test/resources/import.conf"
    val configs        = ConfigsResolve.parse(configFilePath)

    val nebulaGraphConfigEntry = configs.nebulaGraphConfigEntry
    val mqClusterConfigEntry   = configs.mqClusterConfigEntry
    val errorConfigEntry       = configs.errorConfigEntry
    val sourceConfigEntrys     = configs.sourceConfigEntrys

    assert(!nebulaGraphConfigEntry.generateDDL)
    assert(
      "127.0.0.1:9669,127.0.0.1:9670,127.0.0.1:9671".equalsIgnoreCase(
        nebulaGraphConfigEntry.graphAddress))
    assert("nba".equalsIgnoreCase(nebulaGraphConfigEntry.graphName))
    assert("root".equalsIgnoreCase(nebulaGraphConfigEntry.user))
    assert("IMPORT".equalsIgnoreCase(nebulaGraphConfigEntry.mode.toString))

    assert("cluster_nba".equalsIgnoreCase(mqClusterConfigEntry.topic))
    assert("127.0.0.1:9092".equalsIgnoreCase(mqClusterConfigEntry.server))

    assert("/tmp/errors".equalsIgnoreCase(errorConfigEntry.path))
    assert(errorConfigEntry.maxRecords==10)

    assert(sourceConfigEntrys.size==2)
    val sourceConfigEntry = sourceConfigEntrys.head
    assert(sourceConfigEntry.category == SourceCategory.HDFS)
    val hdfsSource = sourceConfigEntry.asInstanceOf[HdfsSourceConfigEntry]
    assert(hdfsSource.fileFormat == FileFormatCategory.CSV)
    assert("hdfs://127.0.0.1:9000/tmp/a.csv".equalsIgnoreCase(hdfsSource.path))
    assert(hdfsSource.header)
    assert(",".equalsIgnoreCase(hdfsSource.separator))
    assert(hdfsSource.readParallel==10)
    assert(hdfsSource.preProcessConfigs.size==4)
    assert(hdfsSource.schemaConfigs.size==2)
    val schemaConfigs = hdfsSource.schemaConfigs
    for (schemaConfig <- schemaConfigs) {
      schemaConfig match {
        case nodeConfig: NodeConfig =>
          val nodeConfig = schemaConfig.asInstanceOf[NodeConfig]
          assert("person".equalsIgnoreCase(nodeConfig.name))
          assert("a".equalsIgnoreCase(nodeConfig.vid))
          assert(nodeConfig.sourceFields.contains("c"))
          assert(nodeConfig.nebulaFields.contains("c"))
          assert(nodeConfig.batchSize==10)
          assert(nodeConfig.partition==10)
        case edgeConfig: EdgeConfig =>
          val edgeConfig = schemaConfig.asInstanceOf[EdgeConfig]
          assert("friend".equalsIgnoreCase(edgeConfig.name))
          assert("a".equalsIgnoreCase(edgeConfig.src))
          assert("b".equalsIgnoreCase(edgeConfig.dst))
          assert(edgeConfig.batchSize==5)
          assert(edgeConfig.partition==10)
      }
    }

    for (preProcessConfig <- sourceConfigEntry.preProcessConfigs) {
      preProcessConfig match {
        case concatConfig: ConcatConfig =>
          assert("_".equalsIgnoreCase(concatConfig.sep))
          assert("ab".equalsIgnoreCase(concatConfig.newFiled))
          assert(concatConfig.oldFields.size == 2)
        case separateConfig: SeparatorConfig =>
          assert("&".equalsIgnoreCase(separateConfig.sep))
          assert("aaa".equalsIgnoreCase(separateConfig.oldField))
          assert(separateConfig.newFields.size == 2)
        case filterConfig: FilterConfig =>
          assert(filterConfig.conditions.size == 1)
          assert("aa>=ss".equalsIgnoreCase(filterConfig.conditions.head))
        case nonValueConfig: NonValueConfig =>
          assert(nonValueConfig.value.isEmpty)
      }
    }
  }
}
