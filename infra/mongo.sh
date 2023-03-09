#!/bin/bash

AppName=mongo-ikm

if [ ! "$(podman ps -a -q -f name=$AppName)" ]
then
    podman run --name $AppName -d \
    -p 27017:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=user \
    -e MONGO_INITDB_ROOT_PASSWORD=secret \
    docker.io/mongo:6.0.4-jammy

    exit 0
fi

podman start $AppName