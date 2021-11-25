#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomctl describe entity works'
expected='
Name:        device
Length:      32
Description: device
CreateTime:
ModifyTime:
'
actual=$(oomctl describe entity device)
ignore_time() { grep -Ev '^(CreateTime|ModifyTime)' <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected")" "$(ignore_time "$actual")"