package com.vesoft.nebula.common.reader.jdbcreader

import com.vesoft.nebula.common.reader.DataSourceReader
import org.apache.spark.sql.DataFrame

/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

class JdbcReader extends DataSourceReader {

  override def readSchema(): String = ???

  override def readData(): DataFrame = ???
}
