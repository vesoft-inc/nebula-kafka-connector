#!/bin/bash

# A demo script to play with ngctl in a local machine

CMD=../ngctl
INSTALL_PATH=""
PACKAGE_PATH=""

if [ -z $1 ]; then
    echo "bash demo.sh <install_path> <package_path>. Please provide the <install_path>."
    exit 1
else
    INSTALL_PATH=$1
fi

if [ -z $2 ]; then
    echo "bash demo.sh <install_path> <package_path>. Please provide the <package_path>."
    exit 1
else
    PACKAGE_PATH=$2
fi

cp ./single.yaml ./config.yaml

sed -i "s@\$INSTALL_PATH@$INSTALL_PATH@g" ./config.yaml

sed -i "s@\$PACKAGE_PATH@$PACKAGE_PATH@g" ./config.yaml

echo "INSTALL_PATH set to: $INSTALL_PATH"

echo "PACKAGE_PATH set to: $PACKAGE_PATH"

$CMD supercluster create --config ./config.yaml --with_install

$CMD supercluster login -P 49559 -u root

$CMD cluster create -c testcluster

$CMD host add -f ./config.yaml -c testcluster

$CMD service add -t storaged -H 127.0.0.1 -P 49779 -c testcluster

$CMD service add -t graphd -H 127.0.0.1 -P 49669 -c testcluster

$CMD service start -t storaged -H 127.0.0.1 -P 49779 -f ./config.yaml

$CMD service start -t graphd -H 127.0.0.1 -P 49669 -f ./config.yaml

$CMD service show -c testcluster

$CMD service stop -t storaged -H 127.0.0.1 -P 49779 -f ./config.yaml

$CMD service stop -t graphd -H 127.0.0.1 -P 49669 -f ./config.yaml

$CMD service drop -t storaged -H 127.0.0.1 -P 49779 -c testcluster

$CMD service drop -t graphd -H 127.0.0.1 -P 49669 -c testcluster

$CMD host drop -f ./config.yaml -c testcluster

$CMD host show -c testcluster

$CMD supercluster stop --config ./config.yaml

# logout to the supercluster
$CMD supercluster logout

rm ./config.yaml
