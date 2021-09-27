#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
import_sample v1
register_features v1

case='featctl list feature works'
expected='Name,Group,Revision,Status,Category,ValueType,Description,RevisionsLimit,CreateTime,ModifyTime
price,device,v1,disabled,batch,int(11),device average price,3,2021-09-27T08:24:26Z,2021-09-27T08:24:26Z
model,device,v1,disabled,batch,varchar(32),device model name,3,2021-09-27T08:24:26Z,2021-09-27T08:24:26Z
'
actual=$(featctl list feature)
ignore_time() { xsv select 1-8 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected")" "$(ignore_time "$actual")"
