#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

# clean up the tmp file
trap 'command rm -rf entity_rows.csv' EXIT INT TERM HUP

before_unix_time=$(date +%s)
echo "1,${before_unix_time}" >> entity_rows.csv
echo "2,${before_unix_time}" >> entity_rows.csv
sleep 1
import_sample > /dev/null
after_unix_time=$(date +%s)
echo "1,${after_unix_time}" >> entity_rows.csv
echo "2,${after_unix_time}" >> entity_rows.csv

case='featctl join historical-feature'
expected="
entity_key,unix_time,model,price
1,${before_unix_time},,
1,${after_unix_time},xiaomi-mix3,3999
2,${before_unix_time},,
2,${after_unix_time},huawei-p40,5299
"

actual=$(featctl join historical-feature \
    --feature model,price \
    --input-file entity_rows.csv \
    --output csv
    )

assert_eq "$case" "$expected" "$actual"
