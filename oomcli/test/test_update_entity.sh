#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_device_sample > /dev/null

case='oomcli update entity works'
oomcli update entity device --description "new description"
expected='
ID,NAME,DESCRIPTION
1,device,new description
'
actual=$(oomcli get meta entity device -o csv)
assert_eq "$case" "$expected" "$actual"
