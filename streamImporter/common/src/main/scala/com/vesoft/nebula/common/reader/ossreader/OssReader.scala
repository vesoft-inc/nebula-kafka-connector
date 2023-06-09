/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.reader.ossreader

import com.vesoft.nebula.common.reader.DataSourceReader
import org.apache.spark.sql.DataFrame

class OssReader extends DataSourceReader{

  override def readSchema(): String = ???

  override def readData(): DataFrame = ???
}
