#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='featctl list revision works'
expected='
Revision,RevisionID,GroupName,DataTable,Description,CreateTime,ModifyTime
1634700104,1,phone,phone_1634700104,test data,2021-10-20T03:21:44Z,2021-10-20T03:21:44Z
'
actual=$(featctl list revision --group phone -o csv)
ignore_time() { cut -d ',' -f 2,3,5 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected")" "$(ignore_time "$actual")"
