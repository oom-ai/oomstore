#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

case="it works"
arg='
{
    "group": "account",
    "revision_id": "3"
}
'
expected='{}'
import_sample account ./data/account_10.csv
actual=$(testgrpc Sync <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
