#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
import_sample v1
register_features v1

case='featctl import using absolute path'
featctl import \
    --group device \
    --revision v2 \
    --schema-template "$(pwd)/schema.sql" \
    --input-file "$(pwd)/device.csv" \
    --description 'test data'
actual=$?
expected=0
assert_eq "$case" "$expected" "$actual"
