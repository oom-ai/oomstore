#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample

case='get all'
expected='device,price
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
#actual=$(featctl get historical-feature --group phone --feature price)
#assert_eq "$case" "$expected" "$actual"
#
#case='get with limit'
#expected='device,price
#1,3999
#2,5299
#3,3999
#4,1999
#5,999
#'
#actual=$(featctl get historical-feature --group phone --feature price --limit 5)
#assert_eq "$case" "$expected" "$actual"
