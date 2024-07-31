
package com.vesoft.nebula.utils

import org.scalatest.funsuite.AnyFunSuite

class NebulaUtilsSuite extends AnyFunSuite {

  test("testIsNumeric"){
    assert(NebulaUtils.isNumeric("123"))
    assert(NebulaUtils.isNumeric("-1"))
    assert(NebulaUtils.isNumeric("0"))
    assert(NebulaUtils.isNumeric("000000"))
    assert(NebulaUtils.isNumeric("-0"))
    assert(NebulaUtils.isNumeric("3452389426358734242346542345243134325342"))
    assert(NebulaUtils.isNumeric("-1247359123461792549218431924235432413"))

    assert(!NebulaUtils.isNumeric("-"))
    assert(!NebulaUtils.isNumeric("a123"))
    assert(!NebulaUtils.isNumeric("abc"))
    assert(!NebulaUtils.isNumeric("1 234"))
    assert(!NebulaUtils.isNumeric("12345a"))
  }

}
