
package com.vesoft.nebula.spark.common.exception


/**
  * An exception thrown if a required option is missing form [[NebulaOptions]]
  */
class IllegalOptionException(message: String, cause: Throwable = null)
    extends IllegalArgumentException(message, cause)

/**
  * An exception thrown if nebula execution occur rpc exception.
  */
class NebulaRPCException(message: String, cause: Throwable = null)
    extends RuntimeException(message, cause)
