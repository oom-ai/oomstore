#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

# clean up the tmp file
trap 'command rm -rf entity_rows.csv entity_rows_with_values.csv' EXIT INT TERM HUP

# import sample data to offline store
import_sample 80

t1=50
t2=100
cat <<-EOF > entity_rows.csv
entity_key,unix_milli
1,$t1
2,$t1
1,$t2
2,$t2
EOF

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


cat <<-EOF > entity_rows_with_values.csv
entity_key,unix_milli,value_1,value_2
1,$t1,1,2
2,$t1,3,4
1,$t2,5,6
2,$t2,7,8
EOF

case='oomcli join historical-feature with real-time feature values'
expected="
entity_key,unix_milli,value_1,value_2,model,price
1,$t1,1,2,,
2,$t1,3,4,,
1,$t2,5,6,xiaomi-mix3,3999
2,$t2,7,8,huawei-p40,5299
"

actual=$(oomcli join \
    --feature model,price \
    --input-file entity_rows_with_values.csv \
    --output csv
)

sorted_expected=$(echo "$expected"|sort)
sorted_actual=$(echo "$actual"|sort)

assert_eq "$case" "$sorted_expected" "$sorted_actual"
