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

    case="$group: api returns ok"
    actual=$(testgrpc Export <<<"$arg")
    assert_json_eq "$case" '{"status":{}}' "$actual"

    case="$group: result is correct"
    # sort result ignoring header
    actual=$(head -1 "$output" && tail -n +2 "$output" | sort -t, -n -k 1,1)
    assert_eq "$case" "$expected" "$actual"
}

##################
#  test group 1  #
##################
group="export some features"
arg=$(cat <<-EOF
{
    "feature_names": ["state"],
    "output_file_path": "$output",
    "revision_id": "3"
}
EOF
)
expected='
user,state
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
    "output_file_path": "$output",
    "revision_id": "3"
}
EOF
)
expected=$(cat ./data/account_10.csv)
group_test "$group" "$arg" "$expected"
)
