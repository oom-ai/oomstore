#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='featctl update group works'
featctl update group phone --description "new description"
expected='
Name:                     phone
Entity:                   device
Description:              new description
Online Revision:          <NULL>
Offline Latest Revision:  <NULL>
Offline Latest DataTable: <NULL>
CreateTime:
ModifyTime:
'
actual=$(featctl describe group phone)
ignore() { grep -Ev '^(CreateTime|ModifyTime|Online Revision|Offline Latest Revision|Offline Latest DataTable)' <<<"$1"; }
assert_eq "$case" "$(ignore "$expected")" "$(ignore "$actual")"
