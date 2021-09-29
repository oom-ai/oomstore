#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
import_sample v1
register_features v1

case='featctl export works'
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
trap 'command rm -rf price.csv' EXIT INT TERM HUP
featctl export -g device -n price --output-file price.csv
actual=$(cat price.csv)
assert_eq "$case" "$expected" "$actual"
