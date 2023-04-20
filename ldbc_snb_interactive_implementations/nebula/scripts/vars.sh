cd "$( cd "$( dirname "${BASH_SOURCE[0]:-${(%):-%x}}" )" >/dev/null 2>&1 && pwd )"
cd ..

export NEBULA_GRAPH_CONTAINER_NAME=nebula-graph
export NEBULA_META_CONTAINER_NAME=nebula-meta
export NEBULA_STORAGE_CONTAINER_NAME=nebula-storage

export NEBULA_CONTAINER_ROOT=`pwd`/scratch
export NEBULA_DATA_DIR=${NEBULA_CONTAINER_ROOT}/data
export NEBULA_ENV_VARS=${NEBULA_ENV_VARS:-}
export NEBULA_HEADER_DIR=`pwd`/headers
export NEBULA_VERSION=${NEBULA_VERSION:-5.0.0}

if [[ `uname -m` == "arm64" ]]; then
    export NEBULA_DOCKER_PLATFORM_FLAG="--platform linux/arm64"
else
    export NEBULA_DOCKER_PLATFORM_FLAG=""
fi
