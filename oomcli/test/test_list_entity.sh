#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli list entity works'
expected='Name,Length,Description,CreateTime,ModifyTime
device,32,device,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
user,64,user,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
'
actual=$(oomcli list entity -o csv)
ignore_time() { cut -d ',' -f 1-3 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"