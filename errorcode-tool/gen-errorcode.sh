#!/bin/bash

errorCodeFilePath=source_file/errorcode
errorClassFilePath=source_file/errorcode_class
errorMessageFilePath=source_file/errormessage
codeDescFilePath=source_file/code_description.md


export CLASSPATH=$CLASSPATH:.
rm -rf ErrorCodeGenerate.class
javac ErrorCodeGenerate.java
java ErrorCodeGenerate $errorCodeFilePath $errorClassFilePath $errorMessageFilePath $codeDescFilePath
