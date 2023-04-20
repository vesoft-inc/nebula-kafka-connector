#!/usr/bin/env bash

set -eu
set -o pipefail

cd "$( cd "$( dirname "${BASH_SOURCE[0]:-${(%):-%x}}" )" >/dev/null 2>&1 && pwd )"
cd ..

. scripts/vars.sh

if [ ! -d ${NEBULA_DATA_DIR} ]; then
    echo "Nebula data directory does not exist"
    exit 1
fi

docker run --rm \
    --user="$(id -u):$(id -g)" \
    --publish=3713:3713 \
    --publish=13713:13713 \
    --detach \
    --ulimit nofile=40000:40000 \
    ${NEBULA_ENV_VARS} \
    --volume=${NEBULA_DATA_DIR}:/data:z \
    --volume=${NEBULA_CONTAINER_ROOT}/logs:/logs:z \
    --volume=${NEBULA_CONTAINER_ROOT}/plugins:/plugins:z \
    --env NEBULA_AUTH=none \
    --name ${NEBULA_CONTAINER_NAME} \
    ${NEBULA_DOCKER_PLATFORM_FLAG} \
    nebula:${NEBULA_VERSION}

echo -n "Waiting for the nebula database to start ."
until docker exec --interactive --tty ${NEBULA_CONTAINER_NAME} cypher-shell "RETURN 'Database has started successfully' AS message" > /dev/null 2>&1; do
    echo -n " ."
    sleep 1
done
echo
echo "Database started."
