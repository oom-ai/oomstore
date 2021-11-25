#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

case='oomctl import using absolute path'
oomctl import \
    --group phone \
    --input-file "$(pwd)/device.csv" \
    --description 'test data' > /dev/null
actual=$?
expected=0
assert_eq "$case" "$expected" "$actual"