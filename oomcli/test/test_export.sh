#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

t0=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
oomcli push --entity-key 1 --group user-click --feature last_5_click_posts=1,2,3,4,5 --feature number_of_user_starred_posts=10
t1=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}

oomcli push --entity-key 1 --group user-click --feature last_5_click_posts=2,3,4,5,6 --feature number_of_user_starred_posts=11
t2=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}

import_student_sample
t3=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}

oomcli push --entity-key 2 --group user-click --feature last_5_click_posts=1,2,3,4,5 --feature number_of_user_starred_posts=10
t4=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}

oomcli_export_no_register_feature() {
    case="export all no register feature"
    actual=$(oomcli export --feature a.b,a.c --unix-milli $t0 2>&1 || true)
    expected='Error: failed exporting features: invalid feature names [a.b a.c]'
    assert_eq "$case" "$expected" "$actual"
}

oomcli_export_has_no_register_feature() {
    case="export has no register feature"
    actual=$(oomcli export --feature user-click.last_5_click_posts,user-click.a --unix-milli $t0 2>&1 ||true)
    expected='Error: failed exporting features: invalid feature names [user-click.a]'
    assert_eq "$case" "$expected" "$actual"
}

oomcli_export_push_feature() {
    case="push feature"
    expected='user,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,"1,2,3,4,5",10
'
    actual=$(oomcli export --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts --unix-milli $t1 -o csv)
    assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"
}

oomcli_export_update_feature() {
    case="update feature"
    expected='user,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,"2,3,4,5,6",11
'
    actual=$(oomcli export --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts --unix-milli $t2 -o csv)
    assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"
}

oomcli_export_batch() {
    case='export batch feature'
    expected='user,student.name,student.gender,student.age
1,lian,m,18
2,gao,m,20
3,zhang,f,17
4,dong,m,25
5,tang,f,18
6,chen,m,25
7,he,f,19
'
    actual=$(oomcli export --feature student.name,student.gender,student.age --unix-milli $t3 -o csv)
    assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"
}

oomcli_export_batch_and_stream_before_import() {
    case="export batch and stream feature before import"
    expected='user,student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,,,,"2,3,4,5,6",11
'
    actual=$(oomcli export --feature student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts --unix-milli $t2 -o csv)
    assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"
}

oomcli_export_batch_and_stream() {
    case="export batch and stream feature"
    expected='user,student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,lian,m,18,"2,3,4,5,6",11
2,gao,m,20,,
3,zhang,f,17,,
4,dong,m,25,,
5,tang,f,18,,
6,chen,m,25,,
7,he,f,19,,
'
    actual=$(oomcli export --feature student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts --unix-milli $t3 -o csv)
    assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"
}

main() {
    oomcli_export_no_register_feature
    oomcli_export_has_no_register_feature
    oomcli_export_push_feature
    oomcli_export_update_feature
    oomcli_export_batch
    oomcli_export_batch_and_stream_before_import
    oomcli_export_batch_and_stream
}

main
