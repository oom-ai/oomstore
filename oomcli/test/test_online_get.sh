#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
revisionID=$(import_device_sample)
sync "phone" $revisionID

case="query single feature"
expected='
device,phone.model
1,xiaomi-mix3
'
actual=$(oomcli get online --feature phone.model -k 1 -o csv)
assert_eq "$case" "$expected" "$actual"

case="sync stream feature with no import"
actual=$(oomcli sync --group-name user-click 2>&1 || true)
expected=$(cat <<-EOF
syncing features ...
Error: failed sync features: group user-click doesn't have any revision
EOF
)
assert_eq "$case" "$expected" "$actual"

import_user_click
case="sync stream feature"
oomcli sync --group-name user-click

case="query stream feature"
actual=$(oomcli get online -k 1,2 --feature user-click.last_5_click_posts,user-click.number_of_user_starred_posts -o csv)
expected=$(cat <<-EOF
user,user-click.last_5_click_posts,user-click.number_of_user_starred_posts
1,"1,2",10
2,"2,3",10
EOF
)
assert_eq "$case" "$expected" "$actual"

case="query multiple features 1"
expected=$(cat <<-EOF
device,phone.model,phone.price
6,apple-iphone11,4999
EOF
)

actual=$(oomcli get online --feature phone.model,phone.price -k 6 -o csv)
assert_eq "$case" "$expected" "$actual"

case="query multiple features 2"
expected=$(cat <<-EOF
device,phone.price,phone.model
6,4999,apple-iphone11
EOF
)

actual=$(oomcli get online --feature phone.price,phone.model -k 6 -o csv)
assert_eq "$case" "$expected" "$actual"

case="query multiple entity and features"
expected=$(cat <<-EOF
device,phone.price,phone.model
1,3999,xiaomi-mix3
6,4999,apple-iphone11
EOF
)
actual=$(oomcli get online --feature phone.price,phone.model -k 1,6 -o csv)
assert_eq "$case" "$expected" "$actual"
