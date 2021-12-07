#!/usr/bin/env bash
set -euo pipefail
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

[ "$(docker ps -q -f name="$CONTAINER_NAME")" ] && info <<< "$CONTAINER_NAME already running" && exit

case $DB in
    postgres)
        docker run --rm -d --name "$CONTAINER_NAME" \
            --health-cmd='pg_isready' \
            -e POSTGRES_PASSWORD=postgres \
            -e POSTGRES_USER=postgres \
            -p 5432:5432 \
            postgres:14.0-alpine 2>&1 | info
        docker exec "$CONTAINER_NAME" sh -c 'while ! pg_isready; do sleep 1; done' 2>&1 |info
        info <<< "username: postgres"
        info <<< "password: postgres"
        ;;
    mysql)
        docker run --rm -d --name "$CONTAINER_NAME" \
            -e MYSQL_ALLOW_EMPTY_PASSWORD=1 \
            -p 3306:3306 \
            mysql:5.7 | info
        docker exec "$CONTAINER_NAME" sh -c 'while ! mysqladmin ping -h localhost; do sleep 1; done' 2>&1 | info
        info <<< "username: root"
        info <<< "password:"
        ;;
    *)
        usage && exit 1
esac
