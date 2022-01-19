#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

oomcli push --entity-key 1 --group user-click --features last_5_click_posts=1,2,3,4,5 --features number_of_user_starred_posts=10
t1=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}

oomcli push --entity-key 1 --group user-click --features last_5_click_posts=2,3,4,5,6 --features number_of_user_starred_posts=11
t2=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}

oomcli push --entity-key 2 --group user-click --features last_5_click_posts=1,2,3,4,5 --features number_of_user_starred_posts=10
t3=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}

case="push feature"
expected='user,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,"1,2,3,4,5",10
'
actual=$(oomcli export --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts --unix-milli $t1 -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case="update feature"
expected='user,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,"2,3,4,5,6",11
'
actual=$(oomcli export --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts --unix-milli $t2 -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case="push new entity key feature "
expected='user,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,"2,3,4,5,6",11
2,"1,2,3,4,5",10
'
actual=$(oomcli export --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts --unix-milli $t3 -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"
