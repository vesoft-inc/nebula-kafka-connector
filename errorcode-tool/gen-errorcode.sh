#!/bin/bash

errorCodeFilePath=errorcode
client=$1

export CLASSPATH=$CLASSPATH:.
javac ErrorCodeGenerate.java
java ErrorCodeGenerate $errorCodeFilePath $client
