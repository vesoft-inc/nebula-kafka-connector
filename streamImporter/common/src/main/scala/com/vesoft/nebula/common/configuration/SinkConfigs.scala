
package com.vesoft.nebula.common.configuration

/**
  * sink config
  * */
sealed trait DataSinkConfigEntry {
  def category: SinkCategory.Value

}
