#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomd

case="it works"
arg='
{
    "revision_id": "3"
}
'
expected='
{
  "status": {}
}
'
import_sample account ./data/account_10.csv
actual=$(testgrpc Sync <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
