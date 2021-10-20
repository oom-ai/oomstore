#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample

case='featctl describe feature works'
expected='
Name:          price
Group:         phone
Entity:        device
Category:      batch
ValueType:     int
Description:
Revision:
DataTable:
CreateTime:
ModifyTime:
'
actual=$(featctl describe feature price)
ignore_time() { grep -Ev '^(CreateTime|ModifyTime|Revision|DataTable)' <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected")" "$(ignore_time "$actual")"
