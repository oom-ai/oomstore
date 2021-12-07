#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli get meta revision works'
expected='
ID,REVISION,GROUP,DATA-TABLE,DESCRIPTION,ANCHORED,CREATE-TIME,MODIFY-TIME
1,1638519905556,phone,offline_1_1,test data,true,2021-12-03T08:25:05Z,2021-12-03T08:25:05Z
'
actual=$(oomcli get meta revision --group phone -o csv --wide)
filter() { cut -d ',' -f 1,3-6 <<<"$1"; }
assert_eq "$case" "$(filter "$expected")" "$(filter "$actual")"
