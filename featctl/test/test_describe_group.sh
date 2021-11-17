#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='featctl describe group works'
expected='
Name:             phone
Entity:           device
Description:      phone
OnlineRevisionID: <NULL>
CreateTime:       2021-11-17T11:24:03Z
ModifyTime:       2021-11-17T11:24:03Z
'
actual=$(featctl describe group phone)
ignore_time() { grep -Ev '^(CreateTime|ModifyTime)' <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected")" "$(ignore_time "$actual")"
