#!/bin/bash

errorCodeFilePath=source_file/errorcode
errorClassFilePath=source_file/errorcode_class
errorMessageFilePath=source_file/errormessage
codeDescFilePath=source_file/code_description.md


export CLASSPATH=$CLASSPATH:.
rm -rf ErrorCodeGenerate.class
javac ErrorCodeGenerate.java
java ErrorCodeGenerate $errorCodeFilePath $errorClassFilePath $errorMessageFilePath $codeDescFilePath

# replace the errorcode file
cp error.go ../golang/pkg/errors
cp ErrorCode.java ../java/client/src/main/java/com/vesoft/nebula/driver/graph
cp _error_code.py ../python/src/nebulagraph_python

# format go
cd ../golang
go fmt ./...
