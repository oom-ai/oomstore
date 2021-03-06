#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

oomcli apply -f ./data/driver_stats.yaml

import_sample account "./data/account_10.csv"

(
trap 'command rm -rf $TEMP' EXIT INT TERM HUP
TEMP="$(mktemp -dt "$(basename "$0")".XXXXXX)"
output="$TEMP/result.csv"

group_test() {
    local group="$1"
    local arg="$2"
    local expected="$3"

    actual=$(testgrpc Export <<<"$arg")

    case="$group: result is correct"
    # sort result ignoring header
    actual=$(head -1 "$output" && tail -n +2 "$output" | sort -t, -n -k 1,1)
    assert_eq "$case" "$expected" "$actual"
}

##################
#  test group 1  #
##################
group="export some features"
unix_milli=${3:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
arg=$(cat <<-EOF
{
    "features": ["account.state"],
    "unix_milli": $unix_milli,
    "output_file": "$output"
}
EOF
)
expected='
user,account.state
1,Nevada
2,South Carolina
3,New Jersey
4,Ohio
5,California
6,North Carolina
7,North Dakota
8,West Virginia
9,Alabama
10,Idaho
'
group_test "$group" "$arg" "$expected"

##################
#  test group 2  #
##################
group="export all features"
arg=$(cat <<-EOF
{
    "features": ["account.state","account.credit_score","account.account_age_days","account.has_2fa_installed"],
    "output_file": "$output",
    "unix_milli": $unix_milli
}
EOF
)
expected='
user,account.state,account.credit_score,account.account_age_days,account.has_2fa_installed
1,Nevada,530,242,true
2,South Carolina,520,268,false
3,New Jersey,655,84,false
4,Ohio,677,119,true
5,California,566,289,false
6,North Carolina,533,155,true
7,North Dakota,605,334,true
8,West Virginia,664,282,false
9,Alabama,577,150,true
10,Idaho,693,212,true
'
group_test "$group" "$arg" "$expected"
)
