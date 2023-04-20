#!/usr/bin/env bash

set -eu
set -o pipefail

cd "$( cd "$( dirname "${BASH_SOURCE[0]:-${(%):-%x}}" )" >/dev/null 2>&1 && pwd )"
cd ..

. scripts/vars.sh

if [ ! -d ${NEBULA_CSV_DIR} ]; then
    echo "Nebula CSV directory does not exist. \${NEBULA_CSV_DIR} is set to: ${NEBULA_CSV_DIR}"
    exit 1
fi

echo "==============================================================================="
echo "Loading the Nebula database"
echo "-------------------------------------------------------------------------------"
echo "NEBULA_CONTAINER_ROOT: ${NEBULA_CONTAINER_ROOT}"
echo "NEBULA_VERSION: ${NEBULA_VERSION}"
echo "NEBULA_CONTAINER_NAME: ${NEBULA_CONTAINER_NAME}"
echo "NEBULA_ENV_VARS: ${NEBULA_ENV_VARS}"
echo "NEBULA_DATA_DIR (on the host machine):"
echo "  ${NEBULA_DATA_DIR}"
echo "NEBULA_CSV_DIR (on the host machine):"
echo "  ${NEBULA_CSV_DIR}"
echo "==============================================================================="

scripts/stop.sh
scripts/delete-database.sh
scripts/import.sh
scripts/start.sh
scripts/create-indices.sh
