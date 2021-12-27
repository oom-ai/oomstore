#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli update group works'
oomcli update group phone --description "new description"
expected='
ID,NAME,ENTITY,CATEGORY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,phone,device,batch,new description,,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
'
filter() { cut -d ',' -f 1-5 <<<"$1"; }
actual=$(oomcli get meta group phone -o csv)
assert_eq "$case" "$(filter "$expected")" "$(filter "$actual")"
