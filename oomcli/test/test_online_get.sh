#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
revisionID=$(import_device_sample)
sync "phone" $revisionID

case="query single feature"
expected='
device,phone.model
1,xiaomi-mix3
'
actual=$(oomcli get online --feature phone.model -k 1 -o csv)
assert_eq "$case" "$expected" "$actual"


case="query multiple features 1"
expected='
device,phone.model,phone.price
6,apple-iphone11,4999
'
actual=$(oomcli get online --feature phone.model,phone.price -k 6 -o csv)
assert_eq "$case" "$expected" "$actual"

case="query multiple features 2"
expected='
device,phone.price,phone.model
6,4999,apple-iphone11
'
actual=$(oomcli get online --feature phone.price,phone.model -k 6 -o csv)
assert_eq "$case" "$expected" "$actual"

case="query multiple entity and features"
expected='
device,phone.price,phone.model
1,3999,xiaomi-mix3
6,4999,apple-iphone11
'
actual=$(oomcli get online --feature phone.price,phone.model -k 1,6 -o csv)
assert_eq "$case" "$expected" "$actual"
