#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
import_sample v1
register_features v1

case='featctl describe feature works'
expected='
Name:           price
Group:          device
Revision:       v1
Status:         disabled
Category:       batch
ValueType:      int(11)
Description:    device average price
RevisionsLimit: 3
CreateTime:     2021-09-28T05:59:15Z
ModifyTime:     2021-09-28T05:59:15Z
'
actual=$(featctl describe feature -g device -n price)
ignore_time() { grep -Ev '^(CreateTime|ModifyTime)' <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected")" "$(ignore_time "$actual")"
