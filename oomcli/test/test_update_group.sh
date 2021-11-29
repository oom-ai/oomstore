#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli update group works'
oomcli update group phone --description "new description"
expected='
Name:             phone
Entity:           device
Description:      new description
OnlineRevisionID: <NULL>
'
actual=$(oomcli describe group phone)
ignore() { grep -Ev '^(CreateTime|ModifyTime)' <<<"$1"; }
assert_eq "$case" "$expected" "$(ignore "$actual")"
