#!/usr/bin/env bash

set -eu
set -o pipefail

cd "$( cd "$( dirname "${BASH_SOURCE[0]:-${(%):-%x}}" )" >/dev/null 2>&1 && pwd )"
cd ..

. scripts/vars.sh

# make sure directories exist
mkdir -p ${NEO4J_CONTAINER_ROOT}/{logs,import,plugins}

if [ ! -d ${NEO4J_CSV_DIR} ]; then
    echo "Neo4j CSV directory does not exist. \${NEO4J_CSV_DIR} is set to: ${NEO4J_CSV_DIR}"
    exit 1
fi

# start with a fresh data dir (required by the CSV importer)
mkdir -p ${NEO4J_DATA_DIR}
rm -rf ${NEO4J_DATA_DIR}/*

# We run a $(find) command in the shell of the host machine to list the `part-*.csv` or (if compression is applied) `part-*.csv.gz` files
# in each entity's directory under ${NEO4_CSV_DIR}/initial_snapshot/{static,dynamic}/${ENTITY}.
#
# The paths are mapped to the directory structure inside the container (/import/initial_snapshot/...) and concatenated to a comma-separated string, e.g.
# ",/import/initial_snapshot/static/Place/part1.csv,/import/initial_snapshot/static/Place/part2.csv"
#
# This string is prepended with the path of the header file in the container's /headers directory, yielding e.g.:
# "/headers/static/Place.csv,/import/initial_snapshot/static/Place/part1.csv,/import/initial_snapshot/static/Place/part2.csv"
#
# The average path length in the container is around 110 characters. Depending on the maximum length of the argument list, which is usually between 131072 and 2097152 [1],
# this is sufficient for between 1200 and 19000 files, respectively.
#
# It's also important to consider the length of MAX_ARG_STRLEN that limits the length of a single argument which is usually (?) 131072. [2]
#
# For reference, SF300 has about 12000 part-*.csv files.
#
# [1] https://serverfault.com/a/163390/573198
# [2] https://unix.stackexchange.com/a/120842/315847

NEO4J_PART_FIND_PATTERN="part-*.csv*"
NEO4J_HEADER_EXTENSION=".csv"

if [ "$(uname)" == "Darwin" ]; then
    FIND_COMMAND=gfind
    if ! command -v ${FIND_COMMAND} > /dev/null; then
        echo "Command '${FIND_COMMAND}' not found. Install it with 'brew install findutils'"
        exit 1
    fi
else
    FIND_COMMAND=find
fi

docker run --rm \
    --user="$(id -u):$(id -g)" \
    --volume=${NEO4J_DATA_DIR}:/data \
    --volume=${NEO4J_CSV_DIR}:/import \
    --volume=${NEO4J_HEADER_DIR}:/headers \
    ${NEO4J_ENV_VARS} \
    ${NEO4J_DOCKER_PLATFORM_FLAG} \
    neo4j:${NEO4J_VERSION} \
    neo4j-admin database import full \
    --id-type=INTEGER \
    --ignore-empty-strings=true \
    --bad-tolerance=0 \
    --delimiter '|'
