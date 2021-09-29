#!/usr/bin/env bash

docker run --rm -d --name onestore \
    -e MYSQL_ALLOW_EMPTY_PASSWORD=1 \
    -p 4000:3306 \
    mysql:5.7 >/dev/null

docker logs -f onestore
