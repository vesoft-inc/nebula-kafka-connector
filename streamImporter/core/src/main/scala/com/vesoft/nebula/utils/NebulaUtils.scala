/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.utils

object NebulaUtils {

  /**
    * check if a string is numeric, true if str is numeric, false else
    * @param str string
    * @return Boolean
    */
  def isNumeric(str: String): Boolean = str.matches("-?\\d+")


  def escapeUtil(str:String): String = {
    str
      .replaceAll("\\\\", "\\\\\\\\")
      .replaceAll("\t", "\\\\t")
      .replaceAll("\n", "\\\\n")
      .replaceAll("\"", "\\\\\"")
      .replaceAll("\'", "\\\\'")
      .replaceAll("\r", "\\\\r")
      .replaceAll("\b", "\\\\b")
  }
}
