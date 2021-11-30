#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

case="api returns ok"
arg=$(cat <<-EOF
{
    "group_name": "account",
    "row": "$(base64 <./data/account_10.csv | tr -d '\n\r')"
}
EOF
)
expected='
{
  "status": {},
  "revisionId": "3"
}
'
actual=$(testgrpc ChannelImport <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="data actually imported"
expected='
user,state,credit_score,account_age_days,has_2fa_installed
10,Idaho,693,212,true
'
oomcli sync -r 3
actual=$(oomcli get online --feature state,credit_score,account_age_days,has_2fa_installed -k 10 -o csv)
assert_eq "$case" "$expected" "$actual"