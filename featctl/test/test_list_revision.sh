#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
import_sample v1
register_features v1

case='featctl list revision works'
expected='
Group,Revision,Source,Description
device,v1,device_v1,test data
'
actual=$(featctl list revision --group device)
ignore_time() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected")" "$(ignore_time "$actual")"
