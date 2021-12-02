#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli get meta features works'
expected='ID,NAME,GROUP,ENTITY,CATEGORY,DB-VALUE-TYPE,VALUE-TYPE,DESCRIPTION,ONLINE-REVISION-ID
1,price,phone,device,batch,int,int64,price,<NULL>
2,model,phone,device,batch,varchar(32),string,model,<NULL>
'
actual=$(oomcli get meta feature -o csv --wide)
ignore_time() { cut -d ',' -f 1-9 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"

case='oomcli get simplified meta features works'
expected='ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,price,phone,device,batch,int64,price
2,model,phone,device,batch,string,model
'
actual=$(oomcli get meta feature -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get one meta feature works'
expected='ID,NAME,GROUP,ENTITY,CATEGORY,DB-VALUE-TYPE,VALUE-TYPE,DESCRIPTION,ONLINE-REVISION-ID
2,model,phone,device,batch,varchar(32),string,model,<NULL>
'
actual=$(oomcli get meta feature model -o csv --wide)
ignore_time() { cut -d ',' -f 1-9 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"
