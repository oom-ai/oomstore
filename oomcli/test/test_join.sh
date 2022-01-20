#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

# clean up the tmp file
trap 'command rm -rf entity_rows.csv entity_rows_with_values.csv' EXIT INT TERM HUP

# import sample data to offline store
import_student_sample 80
bt1=50
bt2=100

# push stream feature
st0=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
oomcli push --entity-key 1 --group user-click --features last_5_click_posts=0,1,2,3,4 --features number_of_user_starred_posts=10 >> /dev/null 2>&1
st1=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1
oomcli push --entity-key 1 --group user-click --features last_5_click_posts=1,2,3,4,5 --features number_of_user_starred_posts=10 >> /dev/null 2>&1
st2=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1
oomcli push --entity-key 2 --group user-click --features last_5_click_posts=2,3,4,5,6 --features number_of_user_starred_posts=11 >> /dev/null 2>&1
st3=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
sleep 1
oomcli push --entity-key 3 --group user-click --features last_5_click_posts=3,4,5,6,7 --features number_of_user_starred_posts=12 >> /dev/null 2>&1
st4=$(perl -MTime::HiRes -e 'printf("%.0f\n",Time::HiRes::time()*1000)')
oomcli snapshot user-click

oomcli_join_historical_feature() {
    case='oomcli join historical-feature'

cat <<-EOF > entity_rows.csv
entity_key,unix_milli
1,$bt1
2,$bt1
1,$bt2
2,$bt2
EOF
    expected="
entity_key,unix_milli,student.name,student.gender,student.age
1,$bt1,,,
2,$bt1,,,
1,$bt2,lian,m,18
2,$bt2,gao,m,20
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
1,$bt1,1,2
2,$bt1,3,4
1,$bt2,5,6
2,$bt2,7,8
EOF

    expected="
entity_key,unix_milli,value_1,value_2,student.name,student.gender,student.age
1,$bt1,1,2,,,
2,$bt1,3,4,,,
1,$bt2,5,6,lian,m,18
2,$bt2,7,8,gao,m,20
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
1,$st0
1,$st1
1,$st2
1,$st3
2,$st1
2,$st2
2,$st3
2,$st4
3,$st1
3,$st2
3,$st3
3,$st4
EOF

    expected="
entity_key,unix_milli,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$st0,,
1,$st1,\"0,1,2,3,4\",10
1,$st2,\"1,2,3,4,5\",10
1,$st3,\"1,2,3,4,5\",10
2,$st1,,
2,$st2,,
2,$st3,\"2,3,4,5,6\",11
2,$st4,\"2,3,4,5,6\",11
3,$st1,,
3,$st2,,
3,$st3,,
3,$st4,\"3,4,5,6,7\",12"

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
1,$st0,1,2
1,$st1,2,3
1,$st2,3,4
1,$st3,4,5
2,$st1,1,2
2,$st2,2,3
2,$st3,3,4
2,$st4,1,2
3,$st1,2,3
3,$st2,3,4
3,$st3,4,5
3,$st4,5,6
EOF

    expected="
entity_key,unix_milli,value_1,value_2,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$st0,1,2,,
1,$st1,2,3,\"0,1,2,3,4\",10
1,$st2,3,4,\"1,2,3,4,5\",10
1,$st3,4,5,\"1,2,3,4,5\",10
2,$st1,1,2,,
2,$st2,2,3,,
2,$st3,3,4,\"2,3,4,5,6\",11
2,$st4,1,2,\"2,3,4,5,6\",11
3,$st1,2,3,,
3,$st2,3,4,,
3,$st3,4,5,,
3,$st4,5,6,\"3,4,5,6,7\",12"

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
1,$st0
1,$st1
1,$st2
1,$st3
2,$st1
2,$st2
2,$st3
2,$st4
3,$st1
3,$st2
3,$st3
3,$st4
EOF
    expected="
entity_key,unix_milli,student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$st0,lian,m,18,,
1,$st1,lian,m,18,\"0,1,2,3,4\",10
1,$st2,lian,m,18,\"1,2,3,4,5\",10
1,$st3,lian,m,18,\"1,2,3,4,5\",10
2,$st1,gao,m,20,,
2,$st2,gao,m,20,,
2,$st3,gao,m,20,\"2,3,4,5,6\",11
2,$st4,gao,m,20,\"2,3,4,5,6\",11
3,$st1,zhang,f,17,,
3,$st2,zhang,f,17,,
3,$st3,zhang,f,17,,
3,$st4,zhang,f,17,\"3,4,5,6,7\",12"

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
1,$st0,1,2
1,$st1,2,3
1,$st2,3,4
1,$st3,4,5
2,$st1,1,2
2,$st2,2,3
2,$st3,3,4
2,$st4,4,5
3,$st1,1,2
3,$st2,2,3
3,$st3,3,4
3,$st4,4,5
EOF
    expected="
entity_key,unix_milli,value1,value2,student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,$st0,1,2,lian,m,18,,
1,$st1,2,3,lian,m,18,\"0,1,2,3,4\",10
1,$st2,3,4,lian,m,18,\"1,2,3,4,5\",10
1,$st3,4,5,lian,m,18,\"1,2,3,4,5\",10
2,$st1,1,2,gao,m,20,,
2,$st2,2,3,gao,m,20,,
2,$st3,3,4,gao,m,20,\"2,3,4,5,6\",11
2,$st4,4,5,gao,m,20,\"2,3,4,5,6\",11
3,$st1,1,2,zhang,f,17,,
3,$st2,2,3,zhang,f,17,,
3,$st3,3,4,zhang,f,17,,
3,$st4,4,5,zhang,f,17,\"3,4,5,6,7\",12"

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
