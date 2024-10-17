#!/bin/bash

# A demo script to play with ngctl in a local machine

CMD=../ngctl
INSTALL_PATH=""
PACKAGE_PATH=""

if [[ -z $1 ]]; then
    echo "bash demo.sh <install_path> <package_path>. Please provide the <install_path>."
    exit 1
else
    INSTALL_PATH=$1
fi

if [[ -z $2 ]]; then
    echo "bash demo.sh <install_path> <package_path>. Please provide the <package_path>."
    exit 1
else
    PACKAGE_PATH=$2
fi

cp ./single.yaml ./config.yaml

sed -i "s@\$INSTALL_PATH@${INSTALL_PATH}@g" ./config.yaml

sed -i "s@\$PACKAGE_PATH@${PACKAGE_PATH}@g" ./config.yaml

echo "INSTALL_PATH set to: ${INSTALL_PATH}"

echo "PACKAGE_PATH set to: ${PACKAGE_PATH}"

../ngctl metad create --config ./config.yaml

../ngctl metad login -P 49559 -u root

../ngctl srvgrp create -c testsrvgrp

../ngctl host add -f ./config.yaml -c testsrvgrp

../ngctl host install -f ./config.yaml -c testsrvgrp

../ngctl service add -t storaged -H 127.0.0.1 -P 49779 -c testsrvgrp

../ngctl service add -t graphd -H 127.0.0.1 -P 49669 -c testsrvgrp

../ngctl service start -t storaged -H 127.0.0.1 -P 49779 -f ./config.yaml

../ngctl service start -t graphd -H 127.0.0.1 -P 49669 -f ./config.yaml

../ngctl service show -c testsrvgrp

../ngctl service stop -t storaged -H 127.0.0.1 -P 49779 -f ./config.yaml

../ngctl service stop -t graphd -H 127.0.0.1 -P 49669 -f ./config.yaml

../ngctl service drop -t storaged -H 127.0.0.1 -P 49779 -c testsrvgrp

../ngctl service drop -t graphd -H 127.0.0.1 -P 49669 -c testsrvgrp

../ngctl host drop -f ./config.yaml -c testsrvgrp

../ngctl host uninstall -f ./config.yaml -c testsrvgrp

../ngctl host show -c testsrvgrp

../ngctl metad stop --config ./config.yaml

# logout to the metad
../ngctl metad logout

