#!/bin/bash

errorCodeFilePath=errorcode
errorClassFilePath=errorcode_class
errorMessageFilePath=errormessage


export CLASSPATH=$CLASSPATH:.
rm -rf ErrorCodeGenerate.class
javac ErrorCodeGenerate.java
java ErrorCodeGenerate $errorCodeFilePath $errorClassFilePath $errorMessageFilePath
