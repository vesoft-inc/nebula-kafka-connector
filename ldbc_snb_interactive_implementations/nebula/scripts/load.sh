#!/bin/bash

set -eu
# Directory of this script
SCRIPT_DIR=$(dirname $(readlink -f "$0"))
# Directory of this project
PROJECT_DIR=$(dirname ${SCRIPT_DIR})
echo  ${SCRIPT_DIR}
source ${SCRIPT_DIR}/env.sh

for i in generate.sh convert.sh convert_person.sh split.sh;do
	${SCRIPT_DIR}/$i
	if [ $? -ne 0 ];then
		echo "ERROR: $i failed"
		exit 1
	fi
done

cd ${SCRIPT_DIR}
data_dir=${LDBC_HOME}/test_data

# replace nebula, graph_name and data path
# {{graph_addresses}}
# {{graph_name}}
sed "s/{{graph_name}}/${GRAPH_NAME}/g" importer_template.yaml > importer_${GRAPH_NAME}.yaml
sed -i "s/{{graph_addresses}}/${NEBULA_ADDRESS}/g" importer_${GRAPH_NAME}.yaml
sed -i "s#{{data_folder}}#${data_dir}#g" importer_${GRAPH_NAME}.yaml

if [ -x ${SCRIPT_DIR}/nebula-importer ];then
  echo "nebula-importer: ${SCRIPT_DIR}/nebula-importer is existed, ignore download"
else
        wget "http://192.168.15.12/softs/nebula-importer"
        chmod +x ${SCRIPT_DIR}/nebula-importer
fi

# import data
${SCRIPT_DIR}/nebula-importer --config importer_${GRAPH_NAME}.yaml