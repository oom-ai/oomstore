#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

case='oomcli get meta group works'
expected='ID,NAME,ENTITY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,phone,device,phone,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
'
actual=$(oomcli get meta group -o csv --wide)
ignore_time() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"

case='oomcli get simplified group works'
expected='ID,NAME,ENTITY,DESCRIPTION
1,phone,device,phone
'
actual=$(oomcli get meta group -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get one group works'
expected='ID,NAME,ENTITY,DESCRIPTION
1,phone,device,phone
'
actual=$(oomcli get meta group phone -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"
