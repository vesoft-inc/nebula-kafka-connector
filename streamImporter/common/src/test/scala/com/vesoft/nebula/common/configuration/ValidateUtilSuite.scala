/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.configuration

import org.scalatest.funsuite.AnyFunSuite

class ValidateUtilSuite extends AnyFunSuite {
  test("test server address validation"){
     val server = "127.0.0.1:9669,127.0.0.1:9670"
    val validateUtil = new ValidateUtil
    assert(validateUtil.validateServer(server))

    // test wrong server address
    assert(!validateUtil.validateServer("127.0.0.1"))
    assert(!validateUtil.validateServer("127.0.0.0:0000"))
    assert(!validateUtil.validateServer("127.0.0.1:1，127.0.0.1:2"))
    assert(!validateUtil.validateServer("127.0.0.1:65536"))
  }



}
