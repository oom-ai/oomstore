#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

case="api returns ok"
arg=$(cat <<-EOF
{
    "group": "account",
    "input_file": "./data/account_10.csv",
    "delimiter": ","
}
EOF
)
expected='
{
  "revisionId": 3
}
'
actual=$(testgrpc Import <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="data actually imported"
expected='
user,account.state,account.credit_score,account.account_age_days,account.has_2fa_installed
10,Idaho,693,212,true
'
oomcli sync --group-name account --revision-id 3
actual=$(oomcli get online --feature account.state,account.credit_score,account.account_age_days,account.has_2fa_installed -k 10 -o csv)
assert_eq "$case" "$expected" "$actual"
