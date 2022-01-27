#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

# clean up the tmp file
trap 'command rm -rf entity_rows.csv entity_rows_with_values.csv' EXIT INT TERM HUP

# push stream feature
t0=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')

oomcli push --entity-key 1 --group user-click --feature last_5_click_posts=0,1,2,3,4 --feature number_of_user_starred_posts=10 >> /dev/null 2>&1
t1=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1

# import sample data to offline store
import_student_sample
t2=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1

oomcli push --entity-key 1 --group user-click --feature last_5_click_posts=1,2,3,4,5 --feature number_of_user_starred_posts=10 >> /dev/null 2>&1
t3=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1

oomcli push --entity-key 2 --group user-click --feature last_5_click_posts=2,3,4,5,6 --feature number_of_user_starred_posts=11 >> /dev/null 2>&1
t4=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1

oomcli push --entity-key 3 --group user-click --feature last_5_click_posts=3,4,5,6,7 --feature number_of_user_starred_posts=12 >> /dev/null 2>&1
t5=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')

oomcli snapshot user-click

oomcli_join_historical_feature() {
    case='oomcli join historical-feature'

cat <<-EOF > entity_rows.csv
entity_key,unix_milli
1,$t0
1,$t1
1,$t2
1,$t3
2,$t0
2,$t1
2,$t2
2,$t3
EOF
    expected="
entity_key,unix_milli,student.name,student.gender,student.age
1,$t0,,,
1,$t1,,,
1,$t2,lian,m,18
1,$t3,lian,m,18
2,$t0,,,
2,$t1,,,
2,$t2,gao,m,20
2,$t3,gao,m,20
    "

    actual=$(oomcli join \
        --feature student.name,student.gender,student.age \
        --input-file entity_rows.csv \
        --output csv
    )

    sorted_expected=$(echo "$expected"|sort)
    sorted_actual=$(echo "$actual"|sort)

    assert_eq "$case" "$sorted_expected" "$sorted_actual"
}

oomcli_join_historical_feature_with_real_time_feature_values() {
    case='oomcli join historical-feature with real-time feature values'

cat <<-EOF > entity_rows_with_values.csv
entity_key,unix_milli,value_1,value_2
1,$t0,1,2
1,$t1,3,4
1,$t2,5,6
1,$t3,7,8
2,$t0,1,2
2,$t1,3,4
2,$t2,5,6
2,$t3,7,8
EOF

    expected="
entity_key,unix_milli,value_1,value_2,student.name,student.gender,student.age
1,$t0,1,2,,,
1,$t1,3,4,,,
1,$t2,5,6,lian,m,18
1,$t3,7,8,lian,m,18
2,$t0,1,2,,,
2,$t1,3,4,,,
2,$t2,5,6,gao,m,20
2,$t3,7,8,gao,m,20
"

    actual=$(oomcli join \
        --feature student.name,student.gender,student.age \
        --input-file entity_rows_with_values.csv \
        --output csv
    )

    sorted_expected=$(echo "$expected"|sort)
    sorted_actual=$(echo "$actual"|sort)

    assert_eq "$case" "$sorted_expected" "$sorted_actual"
}

oomcli_join_stream_historical_feature() {
    case='oomcli join stream historical-feature'
    oomcli snapshot user-click

    cat <<-EOF > entity_rows.csv
entity_key,unix_milli
1,$t0
1,$t1
1,$t2
1,$t3
1,$t4
1,$t5
2,$t0
2,$t1
2,$t2
2,$t3
2,$t4
2,$t5
3,$t0
3,$t1
3,$t2
3,$t3
3,$t4
3,$t5
EOF

    expected="
entity_key,unix_milli,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$t0,,
1,$t1,\"0,1,2,3,4\",10
1,$t2,\"0,1,2,3,4\",10
1,$t3,\"1,2,3,4,5\",10
1,$t4,\"1,2,3,4,5\",10
1,$t5,\"1,2,3,4,5\",10
2,$t0,,
2,$t1,,
2,$t2,,
2,$t3,,
2,$t4,\"2,3,4,5,6\",11
2,$t5,\"2,3,4,5,6\",11
3,$t0,,
3,$t1,,
3,$t2,,
3,$t3,,
3,$t4,,
3,$t5,\"3,4,5,6,7\",12"

    actual=$(oomcli join \
        --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts \
        --input-file entity_rows.csv \
        --output csv
    )
    sorted_expected=$(echo "$expected"|sort)
    sorted_actual=$(echo "$actual"|sort)
    assert_eq "$case" "$sorted_expected" "$sorted_actual"
}

oomcli_join_stream_historical_feature_with_real_time_feature_value() {
    case='oomcli join stream historical-feature with real-time feature values'

    cat <<-EOF > entity_rows_with_values.csv
entity_key,unix_milli,value_1,value_2
1,$t0,1,2
1,$t1,3,4
1,$t2,5,6
1,$t3,7,8
1,$t4,9,10
1,$t5,11,12
2,$t0,1,2
2,$t1,3,4
2,$t2,5,6
2,$t3,7,8
2,$t4,9,10
2,$t5,11,12
3,$t0,1,2
3,$t1,3,4
3,$t2,5,6
3,$t3,7,8
3,$t4,9,10
3,$t5,11,12
EOF

    expected="
entity_key,unix_milli,value_1,value_2,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$t0,1,2,,
1,$t1,3,4,\"0,1,2,3,4\",10
1,$t2,5,6,\"0,1,2,3,4\",10
1,$t3,7,8,\"1,2,3,4,5\",10
1,$t4,9,10,\"1,2,3,4,5\",10
1,$t5,11,12,\"1,2,3,4,5\",10
2,$t0,1,2,,
2,$t1,3,4,,
2,$t2,5,6,,
2,$t3,7,8,,
2,$t4,9,10,\"2,3,4,5,6\",11
2,$t5,11,12,\"2,3,4,5,6\",11
3,$t0,1,2,,
3,$t1,3,4,,
3,$t2,5,6,,
3,$t3,7,8,,
3,$t4,9,10,,
3,$t5,11,12,\"3,4,5,6,7\",12"

    actual=$(oomcli join \
        --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts \
        --input-file entity_rows_with_values.csv \
        --output csv
    )
    sorted_expected=$(echo "$expected"|sort)
    sorted_actual=$(echo "$actual"|sort)

    assert_eq "$case" "$sorted_expected" "$sorted_actual"
}

oomcli_join_batch_and_stream_historical_feature() {
   case='oomcli join batch and stream historical-feature'

   cat <<-EOF > entity_rows.csv
entity_key,unix_milli
1,$t0
1,$t1
1,$t2
1,$t3
1,$t4
1,$t5
2,$t0
2,$t1
2,$t2
2,$t3
2,$t4
2,$t5
3,$t0
3,$t1
3,$t2
3,$t3
3,$t4
3,$t5
EOF
    expected="
entity_key,unix_milli,student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$t0,,,,,
1,$t1,,,,\"0,1,2,3,4\",10
1,$t2,lian,m,18,\"0,1,2,3,4\",10
1,$t3,lian,m,18,\"1,2,3,4,5\",10
1,$t4,lian,m,18,\"1,2,3,4,5\",10
1,$t5,lian,m,18,\"1,2,3,4,5\",10
2,$t0,,,,,
2,$t1,,,,,
2,$t2,gao,m,20,,
2,$t3,gao,m,20,,
2,$t4,gao,m,20,\"2,3,4,5,6\",11
2,$t5,gao,m,20,\"2,3,4,5,6\",11
3,$t0,,,,,
3,$t1,,,,,
3,$t2,zhang,f,17,,
3,$t3,zhang,f,17,,
3,$t4,zhang,f,17,,
3,$t5,zhang,f,17,\"3,4,5,6,7\",12"

    actual=$(oomcli join \
        --feature student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts \
        --input-file entity_rows.csv \
        --output csv
    )

    sorted_expected=$(echo "$expected"|sort)
    sorted_actual=$(echo "$actual"|sort)

    assert_eq "$case" "$sorted_expected" "$sorted_actual"
}

oomcli_join_batch_and_stream_historical_feature_with_real_time_feature_value() {
    case='oomcli join batch and stream historical-feature with real time feature value'

   cat <<-EOF > entity_rows.csv
entity_key,unix_milli,value1,value2
1,$t0,1,2
1,$t1,3,4
1,$t2,5,6
1,$t3,7,8
1,$t4,9,10
1,$t5,11,12
2,$t0,1,2
2,$t1,3,4
2,$t2,5,6
2,$t3,7,8
2,$t4,9,10
2,$t5,11,12
3,$t0,1,2
3,$t1,3,4
3,$t2,5,6
3,$t3,7,8
3,$t4,9,10
3,$t5,11,12
EOF
    expected="
entity_key,unix_milli,value1,value2,student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$t0,1,2,,,,,
1,$t1,3,4,,,,\"0,1,2,3,4\",10
1,$t2,5,6,lian,m,18,\"0,1,2,3,4\",10
1,$t3,7,8,lian,m,18,\"1,2,3,4,5\",10
1,$t4,9,10,lian,m,18,\"1,2,3,4,5\",10
1,$t5,11,12,lian,m,18,\"1,2,3,4,5\",10
2,$t0,1,2,,,,,
2,$t1,3,4,,,,,
2,$t2,5,6,gao,m,20,,
2,$t3,7,8,gao,m,20,,
2,$t4,9,10,gao,m,20,\"2,3,4,5,6\",11
2,$t5,11,12,gao,m,20,\"2,3,4,5,6\",11
3,$t0,1,2,,,,,
3,$t1,3,4,,,,,
3,$t2,5,6,zhang,f,17,,
3,$t3,7,8,zhang,f,17,,
3,$t4,9,10,zhang,f,17,,
3,$t5,11,12,zhang,f,17,\"3,4,5,6,7\",12"

    actual=$(oomcli join \
        --feature student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts \
        --input-file entity_rows.csv \
        --output csv
    )

    sorted_expected=$(echo "$expected"|sort)
    sorted_actual=$(echo "$actual"|sort)

    assert_eq "$case" "$sorted_expected" "$sorted_actual"

}

main() {
    oomcli_join_historical_feature
    oomcli_join_historical_feature_with_real_time_feature_values
    oomcli_join_stream_historical_feature
    oomcli_join_stream_historical_feature_with_real_time_feature_value
    oomcli_join_batch_and_stream_historical_feature
    oomcli_join_batch_and_stream_historical_feature_with_real_time_feature_value
}

main
