#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli update entity works'
oomcli update entity device --description "new description"
expected='
NAME,LENGTH,DESCRIPTION
device,32,new description
'
actual=$(oomcli get meta entity device -o csv)
ignore() { cut -d ',' -f 1-3 <<<"$1"; }
assert_eq "$case" "$expected" "$(ignore "$actual")"
