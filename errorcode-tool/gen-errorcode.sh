#!/bin/bash

errorCodeFilePath=errorcode
errorClassFilePath=errorcode_class
errorMessageFilePath=errormessage
codeDescFilePath=code_description.md


export CLASSPATH=$CLASSPATH:.
rm -rf ErrorCodeGenerate.class
javac ErrorCodeGenerate.java
java ErrorCodeGenerate $errorCodeFilePath $errorClassFilePath $errorMessageFilePath $codeDescFilePath
