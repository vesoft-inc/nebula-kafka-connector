/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.common.statistics

import java.util.concurrent.ConcurrentHashMap

class Communication {
  // 读取失败条数，读取总条数，写入失败条数，写入总条数，脏数据条数
  val counter             = new ConcurrentHashMap[String, Long]()
  var maxFailedSize: Long = 0

  /**
    *
    * get the specific type record size
    */
  def getCounter(key: String): Long = {
    if (counter.contains(key)) {
      counter.get(key)
    } else {
      0
    }
  }

  /**
    * increase the value for specific type
    */
  def increaseCounter(key: String, value: Long): Unit = {
    if (counter.contains()) {
      val count = counter.get(key)
      counter.put(key, count + value)
    } else {
      counter.put(key, value)
    }
  }

  /**
    * if one specific type's record size beyond the maxFailedSize
    */
  def beyondMaxFailedSize(key: String): Boolean = {
    if (counter.contains(key)) {
      return counter.get(key) >= maxFailedSize
    }
    false
  }

  /**
    * set max field record size
    */
  def setMaxFailedSize(maxFailedRecordSize: Long) = {
    this.maxFailedSize = maxFailedRecordSize
  }
}

class StatisticType {
  val READ_FAILED_RECORD_SIZE   = "read_failed_record_size"
  val DIRTY_DATA_RECORD_SIZE    = "dirty_data_record_size"
  val WRITE_FAILED_RECORD_SIZE  = "write_failed_record_size"
  val READ_SUCCEED_RECORD_SIZE  = "read_succeed_record_size"
  val WRITE_SUCCEED_RECORD_SIZE = "write_succeed_record_size"
}
