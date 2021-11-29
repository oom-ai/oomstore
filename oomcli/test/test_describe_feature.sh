#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli describe feature works'
expected='
Name:             price
Group:            phone
Entity:           device
Category:         batch
DBValueType:      int
ValueType:        int64
Description:      price
OnlineRevisionID: <NULL>
CreateTime:       2021-11-17T11:16:37Z
ModifyTime:       2021-11-17T11:16:37Z
'
actual=$(oomcli describe feature price)
ignore_time() { grep -Ev '^(CreateTime|ModifyTime)' <<<"$1"; }

assert_eq "$case" "$(ignore_time "$expected")" "$(ignore_time "$actual")"
