#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample

case='featctl update group works'
featctl update group phone --description "new description"
expected='
Name:          phone
Entity:        device
Description:   new description
Revision:
DataTable:
CreateTime:
ModifyTime:
'
actual=$(featctl describe group phone)
ignore() { grep -Ev '^(CreateTime|ModifyTime|Revision|DataTable)' <<<"$1"; }
assert_eq "$case" "$(ignore "$expected")" "$(ignore "$actual")"
