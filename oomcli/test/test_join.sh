#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

# clean up the tmp file
trap 'command rm -rf entity_rows.csv' EXIT INT TERM HUP

t1=50
t2=100
cat <<-EOF > entity_rows.csv
1,$t1
2,$t1
1,$t2
2,$t2
EOF

import_sample 80
case='oomcli join historical-feature'
expected="
entity_key,unix_milli,model,price
1,$t1,,
2,$t1,,
1,$t2,xiaomi-mix3,3999
2,$t2,huawei-p40,5299
"

actual=$(oomcli join \
    --feature model,price \
    --input-file entity_rows.csv \
    --output csv
)

sorted_expected=$(echo "$expected"|sort)
sorted_actual=$(echo "$actual"|sort)

assert_eq "$case" "$sorted_expected" "$sorted_actual"
