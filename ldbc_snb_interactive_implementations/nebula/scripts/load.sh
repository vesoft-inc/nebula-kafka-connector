#!/bin/bash

set -eu
# Directory of this script
SCRIPT_DIR=$(dirname $(readlink -f "$0"))
# Directory of this project
PROJECT_DIR=$(dirname ${SCRIPT_DIR})
echo  ${SCRIPT_DIR}
source ${SCRIPT_DIR}/env.sh

