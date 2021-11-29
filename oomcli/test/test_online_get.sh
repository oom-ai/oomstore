#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
revisionID=$(import_sample)
sync $revisionID

case="query single feature"
expected='
device,model
1,xiaomi-mix3
'
actual=$(oomcli get online --feature model -k 1 -o csv)
assert_eq "$case" "$expected" "$actual"


case="query multiple features 1"
expected='
device,model,price
6,apple-iphone11,4999
'
actual=$(oomcli get online --feature model,price -k 6 -o csv)
assert_eq "$case" "$expected" "$actual"

case="query multiple features 2"
expected='
device,price,model
6,4999,apple-iphone11
'
actual=$(oomcli get online --feature price,model -k 6 -o csv)
assert_eq "$case" "$expected" "$actual"
