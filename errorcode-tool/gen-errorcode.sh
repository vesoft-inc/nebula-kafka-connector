#!/bin/bash

errorCodeFilePath=errorcode
client=$1

export CLASSPATH=$CLASSPATH:.
rm -rf ErrorCodeGenerate.class
javac ErrorCodeGenerate.java
java ErrorCodeGenerate $errorCodeFilePath golang > errorcode_golang.txt
java ErrorCodeGenerate $errorCodeFilePath java > errorcode_java.txt
