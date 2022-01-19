#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

# clean up the tmp file
trap 'command rm -rf entity_rows.csv entity_rows_with_values.csv' EXIT INT TERM HUP

case='oomcli join stream historical-feature'
t0=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
oomcli push --entity-key 1 --group user-click --features last_5_click_posts=0,1,2,3,4 --features number_of_user_starred_posts=10 >> /dev/null 2>&1
t1=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1
oomcli push --entity-key 1 --group user-click --features last_5_click_posts=1,2,3,4,5 --features number_of_user_starred_posts=10 >> /dev/null 2>&1
t2=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1
oomcli push --entity-key 2 --group user-click --features last_5_click_posts=2,3,4,5,6 --features number_of_user_starred_posts=11 >> /dev/null 2>&1
t3=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1
oomcli push --entity-key 3 --group user-click --features last_5_click_posts=3,4,5,6,7 --features number_of_user_starred_posts=12 >> /dev/null 2>&1
t4=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')

oomcli snapshot user-click

cat <<-EOF > entity_rows.csv
entity_key,unix_milli
1,$t0
1,$t1
1,$t2
1,$t3
2,$t1
2,$t2
2,$t3
2,$t4
3,$t1
3,$t2
3,$t3
3,$t4
EOF

expected="
entity_key,unix_milli,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$t0,,
1,$t1,\"0,1,2,3,4\",10
1,$t2,\"1,2,3,4,5\",10
1,$t3,\"1,2,3,4,5\",10
2,$t1,,
2,$t2,,
2,$t3,\"2,3,4,5,6\",11
2,$t4,\"2,3,4,5,6\",11
3,$t1,,
3,$t2,,
3,$t3,,
3,$t4,\"3,4,5,6,7\",12"
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
1,$t0,1,2
1,$t1,2,3
1,$t2,3,4
1,$t3,4,5
2,$t1,1,2
2,$t2,2,3
2,$t3,3,4
2,$t4,1,2
3,$t1,2,3
3,$t2,3,4
3,$t3,4,5
3,$t4,5,6
EOF

expected="
entity_key,unix_milli,value_1,value_2,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$t0,1,2,,
1,$t1,2,3,\"0,1,2,3,4\",10
1,$t2,3,4,\"1,2,3,4,5\",10
1,$t3,4,5,\"1,2,3,4,5\",10
2,$t1,1,2,,
2,$t2,2,3,,
2,$t3,3,4,\"2,3,4,5,6\",11
2,$t4,1,2,\"2,3,4,5,6\",11
3,$t1,2,3,,
3,$t2,3,4,,
3,$t3,4,5,,
3,$t4,5,6,\"3,4,5,6,7\",12"

actual=$(oomcli join \
    --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts \
    --input-file entity_rows_with_values.csv \
    --output csv
)
sorted_expected=$(echo "$expected"|sort)
sorted_actual=$(echo "$actual"|sort)

assert_eq "$case" "$sorted_expected" "$sorted_actual"
