
package com.vesoft.nebula.common.configuration

/**
  * source data pre-process rules
  * */
sealed trait DataPreProcessConfig

case class ConcatConfig(oldFields: List[String], newFiled: String, sep: String = "_")
    extends DataPreProcessConfig {
  override def toString: String =
    s"ConcatConfig{oldFields:$oldFields, newFiled:$newFiled, sep:$sep}"

  def saveSchema(): String = {
    val space       = "          "
    val doubleSpace = "            "
    s"""
       |$space{"
       |${doubleSpace}oldFields: [${oldFields.mkString(",")}]
       |${doubleSpace}newField: \"$newFiled\"
       |${doubleSpace}sep: \"$sep\"
       |$space}""".stripMargin
  }
}

case class SeparatorConfig(oldField: String, newFields: List[String], sep: String = "_")
    extends DataPreProcessConfig {
  override def toString: String =
    s"SeparatorConfig{oldField:$oldField, newFields:$newFields, sep:$sep}"

  def saveSchema(): String = {
    val space       = "          "
    val doubleSpace = "            "
    s"""
       |$space{
       |${doubleSpace}oldField: \"$oldField\"
       |${doubleSpace}newFields: [${newFields.mkString(",")}]
       |${doubleSpace}sep: \"$sep\"
       |$space}
       |""".stripMargin
  }
}

case class FilterConfig(conditions: List[String]) extends DataPreProcessConfig {
  override def toString: String = s"FilterConfig{conditions:$conditions}"

  def saveSchema(): String = {
    val space = "            "
    s"\n${space}filter: [${conditions.mkString(",")}]"
  }
}

case class NonValueConfig(value: String) extends DataPreProcessConfig {
  override def toString: String = s"NonValueConfig{value:$value}"

  def saveSchema(): String = {
    val space = "          "
    if (value == null) {
      ""
    } else {
      s"\n${space}nonValue:${value}"
    }
  }
}
