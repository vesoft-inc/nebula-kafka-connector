/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.reader

import org.apache.spark.sql.DataFrame

trait DataSourceReader {
  def readSchema(): String

  def readData(): DataFrame
}


