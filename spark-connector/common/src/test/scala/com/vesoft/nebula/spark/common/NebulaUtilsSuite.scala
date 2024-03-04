/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.spark.common

import org.apache.spark.sql.types.{LongType, StructType, StructField}
import org.scalatest.funsuite.AnyFunSuite

class NebulaUtilsSuite extends AnyFunSuite{
  test("makeGetters") {
    val schema = StructType(
      List(
        StructField("col1", LongType, nullable = false),
        StructField("col2", LongType, nullable = true)
      ))
    assert(NebulaUtils.makeGetters(schema).length == 2)
  }

  test("isNumic") {
    assert(NebulaUtils.isNumic("123"))
    assert(NebulaUtils.isNumic("-123"))
    assert(!NebulaUtils.isNumic(""))
    assert(!NebulaUtils.isNumic("-"))
    assert(!NebulaUtils.isNumic("1.0"))
    assert(!NebulaUtils.isNumic("a123"))
    assert(!NebulaUtils.isNumic("123b"))
  }

  test("escapeUtil") {
    assert(NebulaUtils.escapeUtil("123").equals("123"))
    // a\bc -> a\\bc
    assert(NebulaUtils.escapeUtil("a\bc").equals("a\\bc"))
    // a\tbc -> a\\tbc
    assert(NebulaUtils.escapeUtil("a\tbc").equals("a\\tbc"))
    // a\nbc -> a\\nbc
    assert(NebulaUtils.escapeUtil("a\nbc").equals("a\\nbc"))
    // a\"bc -> a\\"bc
    assert(NebulaUtils.escapeUtil("a\"bc").equals("a\\\"bc"))
    // a\'bc -> a\\'bc
    assert(NebulaUtils.escapeUtil("a\'bc").equals("a\\'bc"))
    // a\rbc -> a\\rbc
    assert(NebulaUtils.escapeUtil("a\rbc").equals("a\\rbc"))
    // a\bbc -> a\\bbc
    assert(NebulaUtils.escapeUtil("a\bbc").equals("a\\bbc"))
  }


}
