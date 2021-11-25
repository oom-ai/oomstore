#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

case='oomctl list group works'
expected='GroupName,GroupID,EntityName,Description
phone,1,device,phone
'
actual=$(oomctl list group -o csv)
ignore_time() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"
