#!/usr/bin/env bash

docker run --rm -d --name onestore \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_USER=postgres \
    -p 5432:5432 \
    postgres:14.0-alpine >/dev/null

docker logs -f onestore
