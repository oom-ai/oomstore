#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli update feature works'
oomcli update feature price --description "new description"
expected='
NAME,GROUP,ENTITY,CATEGORY,DB-VALUE-TYPE,VALUE-TYPE,DESCRIPTION,ONLINE-REVISION-ID
price,phone,device,batch,int,int64,new description,<NULL>
'
actual=$(oomcli get meta feature price -o csv)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "$case"  "$expected" "$(ignore_time "$actual")"
