#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
import_sample v1
register_features v1

case='featctl export all'
expected='entity_key,price
1,3999
2,5299
3,3999
4,1999
5,999
6,4999
7,5999
8,6500
9,4500
'
actual=$(featctl export -g device -n price)
assert_eq "$case" "$expected" "$actual"

case='featctl export with limit'
expected='entity_key,price
1,3999
2,5299
3,3999
4,1999
5,999
'
actual=$(featctl export -g device -n price --limit 5)
assert_eq "$case" "$expected" "$actual"
