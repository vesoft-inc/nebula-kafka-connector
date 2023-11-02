/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import org.scalatest.BeforeAndAfterAll
import org.scalatest.funsuite.AnyFunSuite

class NebulaConfigSuite extends AnyFunSuite with BeforeAndAfterAll {
  test("test NebulaConnectionConfig") {
    assertThrows[AssertionError](NebulaConnectionConfig.builder().withTimeoutSec(1).build())
    assertThrows[AssertionError](NebulaConnectionConfig.builder().withTimeoutSec(-1).build())
    NebulaConnectionConfig
      .builder()
      .withTimeoutSec(1)
      .build()
  }

  test("test WriteNebulaConfig") {
    var writeNebulaConfig: WriteNebulaVertexConfig = null
    writeNebulaConfig = WriteNebulaVertexConfig
      .builder()
      .withGraphName("test")
      .withNodeType("tag")
      .withVidField("vid")
      .build()

    assert(!writeNebulaConfig.getVidAsProp)
    assert(writeNebulaConfig.getGraphName.equals("test"))
  }

  test("wrong batch size for update") {
    assertThrows[AssertionError](
      WriteNebulaVertexConfig
        .builder()
        .withGraphName("test")
        .withNodeType("tag")
        .withVidField("vId")
        .withWriteMode(WriteMode.UPDATE)
        .withBatchSize(513)
        .build())
    assertThrows[AssertionError](
      WriteNebulaEdgeConfig
        .builder()
        .withGraphName("test")
        .withEdge("edge")
        .withSrcIdField("src")
        .withDstIdField("dst")
        .withWriteMode(WriteMode.UPDATE)
        .withBatchSize(513)
        .build())
  }

  test("test wrong policy") {
    assertThrows[AssertionError](
      WriteNebulaVertexConfig
        .builder()
        .withGraphName("test")
        .withNodeType("tag")
        .withVidField("vId")
        .build())
  }

  test("test wrong batch") {
    assertThrows[AssertionError](
      WriteNebulaVertexConfig
        .builder()
        .withGraphName("test")
        .withNodeType("tag")
        .withVidField("vId")
        .withBatchSize(-1)
        .build())
  }

  test("test ReadNebulaConfig") {
    ReadNebulaConfig
      .builder()
      .withGraphName("test")
      .withLabel("tagName")
      .withNoColumn(true)
      .withReturnCols(List("col"))
      .build()
  }
}
