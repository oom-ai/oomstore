#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli get meta features works'
expected='Name,Group,Entity,Category,DBValueType,ValueType,Description,OnlineRevisionID
model,phone,device,batch,varchar(32),string,model,<NULL>
price,phone,device,batch,int,int64,price,<NULL>
'
actual=$(oomcli get meta feature -o csv)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"

case='oomcli get one meta feature works'
expected='Name,Group,Entity,Category,DBValueType,ValueType,Description,OnlineRevisionID
model,phone,device,batch,varchar(32),string,model,<NULL>
'
actual=$(oomcli get meta feature model -o csv)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"
