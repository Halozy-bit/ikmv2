#!/bin/bash

AppName=web-server

if [ ! "$(podman ps -a -q -f name=$AppName)" ]
then 
    podman run -d --name $AppName \
    -p 8080:80 \
    -v /home/arisygdc/Projects/gopath/src/ikmv2/infra/apache/site-available:/etc/apache2/sites-available \
    -v /home/arisygdc/Projects/gopath/src/ikmv2/infra/apache/site-enabled:/etc/apache2/sites-enabled \
    -v /home/arisygdc/Projects/gopath/src/ikmv2/frontend:/var/www/html \
    docker.io/php:8.2-apache

    exit 0
fi

podman start $AppName