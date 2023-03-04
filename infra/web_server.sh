#!/bin/bash

podman run -d --name web-server \
    -p 8080:80 \
    -v /home/arisygdc/Projects/gopath/src/ikmv2/infra/site-available:/etc/apache2/sites-available \
    -v /home/arisygdc/Projects/gopath/src/ikmv2/infra/site-enabled:/etc/apache2/sites-enabled \
    -v /home/arisygdc/Projects/gopath/src/ikmv2/frontend:/var/www/html \
    docker.io/php:8.2-apache