/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import org.scalatest.BeforeAndAfterAll
import org.scalatest.funsuite.AnyFunSuite

class NebulaConfigSuite extends AnyFunSuite with BeforeAndAfterAll {
  test("test NebulaConnectionConfig") {
    val config = NebulaConnectionConfig
      .builder()
      .withGraphAddress("127.0.0.1:9669")
      .withUser("root")
      .withPasswd("nebula")
      .withTimeoutSec(1)
      .withExecuteRetry(2)
      .withExecuteRetryIntervalMs(1000)
      .build()
    assert(config.getGraphAddress.equals("127.0.0.1:9669"))
    assert(config.getUser.equals("root"))
    assert(config.getPasswd.equals("nebula"))
    assert(config.getExecRetry == 2)
    assert(config.getExecRetryIntervalMs == 1000)
    assert(config.getTimeout == 1)
  }

  test("test wrong connection config") {
    assertThrows[AssertionError](NebulaConnectionConfig.builder().withTimeoutSec(1).build())
    assertThrows[AssertionError](NebulaConnectionConfig.builder().withTimeoutSec(-1).build())
    assertThrows[AssertionError](
      NebulaConnectionConfig
        .builder()
        .withGraphAddress("127.0.0.1:9669")
        .withUser("")
        .withPasswd("nebula")
        .build())
    assertThrows[AssertionError](
      NebulaConnectionConfig
        .builder()
        .withGraphAddress("127.0.0.1:9669")
        .withUser("root")
        .withPasswd("")
        .build())
    assertThrows[AssertionError](
      NebulaConnectionConfig
        .builder()
        .withGraphAddress("")
        .withUser("root")
        .withPasswd("nebula")
        .build())
  }

  test("test WriteNebulaConfig") {
    var writeNebulaConfig: WriteNebulaVertexConfig = null
    writeNebulaConfig = WriteNebulaVertexConfig
      .builder()
      .withGraphName("test")
      .withNodeType("tag")
      .withPrimaryKeyField("vid")
      .build()

    assert(writeNebulaConfig.getGraphName.equals("test"))
  }

  test("test wrong batch") {
    assertThrows[AssertionError](
      WriteNebulaVertexConfig
        .builder()
        .withGraphName("test")
        .withNodeType("tag")
        .withPrimaryKeyField("vId")
        .withBatchSize(-1)
        .build())
  }

}
