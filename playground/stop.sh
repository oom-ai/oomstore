#!/usr/bin/env bash
set -euo pipefail
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

[ "$(docker ps -q -f name="$CONTAINER_NAME")" ] || { info <<< "$CONTAINER_NAME not running" && exit; }

info <<< "$(docker stop "$CONTAINER_NAME") stoped"
