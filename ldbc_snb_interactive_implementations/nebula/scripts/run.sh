#!/usr/bin/env bash

operation=${1:-}

set -eu
# Directory of this script
SCRIPT_DIR=$(dirname $(readlink -f "$0"))
# Directory of this project
PROJECT_DIR=$(dirname ${SCRIPT_DIR})
# target data
TARGET_DIR=${PROJECT_DIR}/target

if [[ -z "${operation}" ]]; then
        echo "Usage: $0 <validate|benchmark>"
        exit 1
fi

if [[ "${operation}" != "validate" && "${operation}" != "benchmark" ]]; then
        echo "Usage: $0 <validate|benchmark>"
        exit 1
fi

if [[ "${operation}" == "validate" ]]; then
        echo "Running LDBC validate."
        echo "java -cp nebula-1.0.0.jar org.ldbcouncil.snb.driver.Client -P ${PROJECT_DIR}/conf/validate.properties"
        java -cp ${TARGET_DIR}/nebula-1.0.0.jar org.ldbcouncil.snb.driver.Client -P ${PROJECT_DIR}/conf/validate.properties
else
        echo "Running LDBC benchmark for product."
        echo "java -cp -cp ${TARGET_DIR}/nebula-1.0.0.jar org.ldbcouncil.snb.driver.Client -P ${PROJECT_DIR}/conf/benchmark.properties"
        java -cp ${TARGET_DIR}/nebula-1.0.0.jar org.ldbcouncil.snb.driver.Client -P ${PROJECT_DIR}/conf/benchmark.properties
fi

echo "Finish"
