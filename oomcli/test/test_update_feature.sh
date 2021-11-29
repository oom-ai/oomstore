#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli update feature works'
oomcli update feature price --description "new description"
expected='
Name:             price
Group:            phone
Entity:           device
Category:         batch
DBValueType:      int
ValueType:        int64
Description:      new description
OnlineRevisionID: <NULL>
'
actual=$(oomcli describe feature price)
ignore() { grep -Ev '^(CreateTime|ModifyTime)' <<<"$1"; }
assert_eq "$case"  "$expected" "$(ignore "$actual")"
