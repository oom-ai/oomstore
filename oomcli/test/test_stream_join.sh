#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

# clean up the tmp file
trap 'command rm -rf entity_rows.csv entity_rows_with_values.csv' EXIT INT TERM HUP

case='oomcli join stream historical-feature'
oomcli push --entity-key 1 --group user-click --features last_5_click_posts=0,1,2,3,4 --features number_of_user_starred_posts=10 >> /dev/null 2>&1
oomcli push --entity-key 1 --group user-click --features last_5_click_posts=1,2,3,4,5 --features number_of_user_starred_posts=10 >> /dev/null 2>&1
oomcli push --entity-key 2 --group user-click --features last_5_click_posts=2,3,4,5,6 --features number_of_user_starred_posts=11 >> /dev/null 2>&1
oomcli push --entity-key 3 --group user-click --features last_5_click_posts=3,4,5,6,7 --features number_of_user_starred_posts=12 >> /dev/null 2>&1
oomcli snapshot user-click

t1=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
cat <<-EOF > entity_rows.csv
entity_key,unix_milli
1,$t1
2,$t1
3,$t1
EOF

expected="
entity_key,unix_milli,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$t1,\"1,2,3,4,5\",10
2,$t1,\"2,3,4,5,6\",11
3,$t1,\"3,4,5,6,7\",12
"
actual=$(oomcli join \
    --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts \
    --input-file entity_rows.csv \
    --output csv
)
sorted_expected=$(echo "$expected"|sort)
sorted_actual=$(echo "$actual"|sort)

assert_eq "$case" "$sorted_expected" "$sorted_actual"

case='oomcli join stream historical-feature with real-time feature values'
cat <<-EOF > entity_rows_with_values.csv
entity_key,unix_milli,value_1,value_2
1,$t1,1,2
2,$t1,3,4
3,$t1,5,6
EOF

expected="
entity_key,unix_milli,value_1,value_2,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$t1,1,2,\"1,2,3,4,5\",10
2,$t1,3,4,\"2,3,4,5,6\",11
3,$t1,5,6,\"3,4,5,6,7\",12
"
actual=$(oomcli join \
    --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts \
    --input-file entity_rows_with_values.csv \
    --output csv
)
sorted_expected=$(echo "$expected"|sort)
sorted_actual=$(echo "$actual"|sort)

assert_eq "$case" "$sorted_expected" "$sorted_actual"
