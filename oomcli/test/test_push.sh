#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

case="push stream feature"
expected='
user,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,"1,2,3,4,5",10
'
oomcli push --entity-key 1 --group user-click --features last_5_click_posts=1,2,3,4,5 --features number_of_user_starred_posts=10

actual=$(oomcli get online --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts -k 1 -o csv)
assert_eq "$case" "$expected" "$actual"
