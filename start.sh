#!/usr/bin/bash

set -e

MYSQL_VERSION=5.7.19
MYSQL_NODE_NAME=deepsea-mysql

trap "docker stop $MYSQL_NODE_NAME && docker rm $MYSQL_NODE_NAME" SIGINT SIGTERM

mkdir -p docker_mysql_volume
if [[ ! $(docker container ls --format {{.Names}} | grep $MYSQL_NODE_NAME) ]]; then
    docker run \
        --name $MYSQL_NODE_NAME \
        -v $PWD/docker_mysql_volume:/var/lib/mysql \
        -e MYSQL_ROOT_PASSWORD=deepsea \
        -d \
        -p 3307:3306 \
        mysql:$MYSQL_VERSION
fi

mysql -uroot -h127.0.0.1 -P3307 -pdeepsea -e "CREATE DATABASE IF NOT EXISTS test_deepsea;"

./build/deepsea-api -kubeconfig ~/.kube/config