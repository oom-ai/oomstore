#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
import_sample v1
register_features v1

case='featctl update works'
# import v2 data
import_sample v2
# update active revision to v2
featctl set -g device -n price --revision v2
expected='
Name,Revision
price,v2
model,v1
'
actual=$(featctl list feature | xsv select Name,Revision)
assert_eq "$case" "$expected" "$actual"
