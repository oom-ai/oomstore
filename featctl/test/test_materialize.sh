#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample

case='featctl materialize feature values from offline store to online store'
featctl materialize phone

expected='
device,model,price
1,xiaomi-mix3,3999
'
actual=$(featctl get online-feature --feature price,model -k 1)
assert_eq "$case" "$expected" "$actual"
