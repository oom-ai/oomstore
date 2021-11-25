#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomctl list feature works'
expected='Name,Group,Entity,Category,DBValueType,ValueType,Description,OnlineRevisionID
model,phone,device,batch,varchar(32),string,model,<NULL>
price,phone,device,batch,int,int64,price,<NULL>
'
actual=$(oomctl list feature -o csv)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"
