/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common

import org.apache.hadoop.conf.Configuration
import org.apache.hadoop.fs.Path
import org.apache.log4j.Logger

import scala.collection.mutable.ArrayBuffer

object ErrorHandler {
  @transient
  private[this] val LOG = Logger.getLogger(this.getClass)

  def save(buffer: ArrayBuffer[String], path: String): Unit = {
    LOG.info("create error output path: $path")
    val targetPath = new Path(path)
    val fileSystem = targetPath.getFileSystem(new Configuration())
    val errors = if (fileSystem.exists(targetPath)) {
      // For kafka, the error ngql need to append to a same file instead of overwrite
      fileSystem.append(targetPath)
    } else {
      fileSystem.create(targetPath)
    }

    try {
      for (error <- buffer) {
        errors.write(error.getBytes)
        errors.writeBytes("\n")
      }
    } finally {
      errors.close()
    }
  }
}
