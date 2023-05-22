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

# replace nebula, space and data path
# {{graph_addresses}}
# {{graph_name}}
sed "s/{{graph_name}}/${GRAPH_NAME}/g" importer_template_node.yaml > importer_${GRAPH_NAME}_node.yaml
sed -i "s/{{graph_addresses}}/${NEBULA_ADDRESS}/g" importer_${GRAPH_NAME}_node.yaml
sed -i "s#{{data_folder}}#${data_dir}#g" importer_${GRAPH_NAME}_node.yaml

sed "s/{{graph_name}}/${GRAPH_NAME}/g" importer_template_edge.yaml > importer_${GRAPH_NAME}_edge.yaml
sed -i "s/{{graph_addresses}}/${NEBULA_ADDRESS}/g" importer_${GRAPH_NAME}_edge.yaml
sed -i "s#{{data_folder}}#${data_dir}#g" importer_${GRAPH_NAME}_edge.yaml

if [ -x ${SCRIPT_DIR}/nebula-importer ];then
  echo "nebula-importer: ${SCRIPT_DIR}/nebula-importer is existed, ignore download"
else
        wget "https://nebula-graph-ent.oss-accelerate.aliyuncs.com/general/nebula-importer/nebula-importer-v5-amd64?OSSAccessKeyId=LTAI5tPwfHdcUx2NtmZPmqUh&Expires=1684136792&Signature=UCCMXqjFtKNt6XzU2%2FVB6NP4AgY%3D" -O "nebula-importer"
        chmod +x ${SCRIPT_DIR}/nebula-importer
fi

# import data
echo "import node"
${SCRIPT_DIR}/nebula-importer --config importer_${GRAPH_NAME}_node.yaml
echo "import edge"
${SCRIPT_DIR}/nebula-importer --config importer_${GRAPH_NAME}_edge.yaml