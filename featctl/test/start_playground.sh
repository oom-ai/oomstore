#!/usr/bin/env bash

info() { printf "$(date -Is) %b[info]%b %s\n" '\e[0;32m\033[1m' '\e[0m' "$*" >&2; }

docker run --rm -d --name onestore \
    -p 2379:2379 \
    -p 9090:9090 \
    -p 3000:3000 \
    -p 4000:4000 \
    -p 3930:3930 \
    aiinfra/tidb-playground >/dev/null

docker logs -f onestore | sed -E 's/[0-9]+.[0-9]+.[0-9]+.[0-9]+/127.0.0.1/g'
