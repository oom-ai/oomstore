#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcl get meta entity works'
expected='NAME,LENGTH,DESCRIPTION,CREATE-TIME,MODIFY-TIME
device,32,device,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
user,64,user,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
'
actual=$(oomcli get meta entity -o csv)
ignore_time() { cut -d ',' -f 1-3 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"
