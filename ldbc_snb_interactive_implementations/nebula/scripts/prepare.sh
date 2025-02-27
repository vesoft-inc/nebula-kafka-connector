#!/bin/bash

set -eu
# Directory of this script
SCRIPT_DIR=$(dirname $(readlink -f "$0"))
# Directory of this project
PROJECT_DIR=$(dirname ${SCRIPT_DIR})
source ${SCRIPT_DIR}/env.sh
# rename files
cd ${LDBC_HOME}
for file in `ls dynamic/*_0_0.*`; do
    new_name=$(echo "$file" | sed 's/_0_0//g')
    mv "$file" "$new_name"
done

for file in `ls static/*_0_0.*`; do
    new_name=$(echo "$file" | sed 's/_0_0//g')
    mv "$file" "$new_name"
done

cd static
ls *.csv|xargs sed -i 's/.id//g'
cd ../dynamic
ls *.csv|xargs sed -i 's/.id//g'