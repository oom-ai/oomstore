#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomctl update entity works'
oomctl update entity device --description "new description"
expected='
Name:        device
Length:      32
Description: new description
'
actual=$(oomctl describe entity device)
ignore() { grep -Ev '^(CreateTime|ModifyTime)' <<<"$1"; }
assert_eq "$case" "$expected" "$(ignore "$actual")"