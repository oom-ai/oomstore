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
actual=$(featctl get online-feature --feature model -k 1 -o csv)
assert_eq "$case" "$expected" "$actual"


case="query multiple features"
expected='
device,model,price
6,apple-iphone11,4999
'
actual=$(featctl get online-feature --feature model,price -k 6 -o csv)
assert_eq "$case" "$expected" "$actual"
