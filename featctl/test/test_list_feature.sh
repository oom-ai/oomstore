#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample

case='featctl list feature works'
expected='Name,Group,Entity,Category,DBValueType,ValueType,Description,Revision,DataTable,CreateTime,ModifyTime
model,phone,device,batch,varchar(32),string,,1634626568,phone_1634626568,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
price,phone,device,batch,int,int32,1634626568,phone_1634626568,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
'
actual=$(featctl list feature)
ignore_time() { cut -d ',' -f 1-6 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"
