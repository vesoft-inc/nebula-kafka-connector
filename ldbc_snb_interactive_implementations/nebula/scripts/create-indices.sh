#!/usr/bin/env bash

set -eu
set -o pipefail

cd "$( cd "$( dirname "${BASH_SOURCE[0]:-${(%):-%x}}" )" >/dev/null 2>&1 && pwd )"

docker exec --interactive ${NEBULA_CONTAINER_NAME} cypher-shell < ../ddl/indices.nebula
