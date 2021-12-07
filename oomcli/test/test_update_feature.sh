#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli update feature works'
oomcli update feature price --description "new description"
expected='
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION,DB-VALUE-TYPE,ONLINE-REVISION-ID
1,price,phone,device,batch,int64,new description,int,<NULL>
'
actual=$(oomcli get meta feature price -o csv --wide)
ignore_time() { cut -d ',' -f 1-9 <<<"$1"; }
assert_eq "$case"  "$expected" "$(ignore_time "$actual")"
