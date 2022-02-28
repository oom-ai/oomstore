#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

import_student_sample
oomcli push --entity-key 1 --group user-click --feature last_5_click_posts=1,2,3,4,5 --feature number_of_user_starred_posts=10

oomcli_export_temporary_table() {
    case="temporary table in export"

    t0=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
    oomcli export --feature student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts --unix-milli $t0 -o csv 2>&1 >> /dev/null

    t1=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
    oomcli export --feature student.name,student.gender,student.age,user-click.last_5_click_posts,user-click.number_of_user_starred_posts --unix-milli $t0 -o csv 2>&1 >> /dev/null

    t2=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}

    actual=$(oomcli gc --unix-milli $t0 2>&1 |wc -l)
    assert_eq "$case" 0 "$actual"

    actual=$(oomcli gc --unix-milli $t1 2>&1 |wc -l)
    assert_eq "$case" 2 "$actual"

    actual=$(oomcli gc --unix-milli $t2 2>&1 |wc -l)
    assert_eq "$case" 3 "$actual"

    oomcli gc --unix-milli $t0 --force
    actual=$(oomcli gc --unix-milli $t2 2>&1 |wc -l)
    assert_eq "$case" 3 "$actual"

    oomcli gc --unix-milli $t1 --force
    actual=$(oomcli gc --unix-milli $t2 2>&1 |wc -l)
    assert_eq "$case" 2 "$actual"

    oomcli gc --unix-milli $t2 --force
    actual=$(oomcli gc --unix-milli $t2 2>&1 |wc -l)
    assert_eq "$case" 0 "$actual"
}

oomcli_join_temporary_table() {
    t0=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}

    # clean up the tmp file
    trap 'command rm -rf entity_rows.csv entity_rows.csv' EXIT INT TERM HUP
    cat <<-EOF > entity_rows.csv
entity_key,unix_milli
1,$t0
2,$t0
EOF

    case="temporary table in join"
    oomcli join \
        --feature student.name,student.gender,student.age \
        --input-file entity_rows.csv \
        --output csv > /dev/null

    oomcli join \
        --feature student.name,student.gender,student.age \
        --input-file entity_rows.csv \
        --output csv > /dev/null

    oomcli join \
        --feature student.name,student.gender,student.age \
        --input-file entity_rows.csv \
        --output csv > /dev/null


    t3=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
    actual=$(oomcli gc --unix-milli $t3 2>&1 |wc -l)
    assert_eq "$case" 0 "$actual"
}

main() {
    oomcli_export_temporary_table
    oomcli_join_temporary_table
}

main
