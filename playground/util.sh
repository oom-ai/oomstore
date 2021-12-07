#!/usr/bin/env bash
info() {
    while IFS='' read -r line || [[ -n "$line" ]]; do
        printf "%b[info]%b %s\n" '\e[0;32m\033[1m' '\e[0m' "$line" >&2;
    done
}
erro() {
    while IFS='' read -r line || [[ -n "$line" ]]; do
        printf "%b[erro]%b %s\n" '\e[0;31m\033[1m' '\e[0m' "$line" >&2;
    done
}

usage() { erro <<< "Usage: $(basename "$0") [postgres|mysql]"; }

DB=${1:-postgres}
CONTAINER_NAME="oomstore_playground_$DB"
