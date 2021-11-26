#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

# clean up the tmp file
trap 'command rm -rf entity_rows.csv' EXIT INT TERM HUP

before_unix_time_ms=$(date +%s000)
echo "1,${before_unix_time_ms}" >> entity_rows.csv
echo "2,${before_unix_time_ms}" >> entity_rows.csv
sleep 1
import_sample > /dev/null
sleep 1
after_unix_time_ms=$(date +%s000)
echo "1,${after_unix_time_ms}" >> entity_rows.csv
echo "2,${after_unix_time_ms}" >> entity_rows.csv

case='oomctl join historical-feature'
expected="
entity_key,unix_time,model,price
1,${after_unix_time_ms},xiaomi-mix3,3999
2,${after_unix_time_ms},huawei-p40,5299
1,${before_unix_time_ms},,
2,${before_unix_time_ms},,
"

actual=$(oomctl join \
    --feature model,price \
    --input-file entity_rows.csv \
    --output csv
)

sorted_expected=$(echo "$expected"|sort)
sorted_actual=$(echo "$actual"|sort)

assert_eq "$case" "$sorted_expected" "$sorted_actual"
