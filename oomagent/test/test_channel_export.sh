#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

import_sample account "./data/account_10.csv"
oomcli push --entity-key 1 --group user_fake_stream --features f1=10
oomcli snapshot user_fake_stream

unix_milli=${3:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
echo $unix_milli

case1() {
    case="export and no features"
    arg=$(cat <<-EOF
    {
       "features": [],
       "unix_milli": $unix_milli
    }
EOF
    )
    actual=$(testgrpc ChannelExport <<<"$arg")
    expected=''
    assert_json_eq "$case" "$expected" "$actual"
}

case2() {
    prefix="export some features"
    arg=$(cat <<-EOF
    {
       "features": ["account.state"],
       "unix_milli": $unix_milli
    }
EOF
    )
    actual=$(testgrpc ChannelExport <<<"$arg")

    case="$prefix - first response returns header"
    actual_header=$(jq -s '.[0].header' <<< "$actual")
    expected_header='["user","account.state"]'
    assert_json_eq "$case" "$expected_header" "$actual_header"

    case="$prefix - returns correct rows"
    actual_rows=$(jq -c ".row" <<< "$actual" | sort)
    expected_rows='
    [{"string":"1"},{"string":"Nevada"}]
    [{"string":"2"},{"string":"South Carolina"}]
    [{"string":"3"},{"string":"New Jersey"}]
    [{"string":"4"},{"string":"Ohio"}]
    [{"string":"5"},{"string":"California"}]
    [{"string":"6"},{"string":"North Carolina"}]
    [{"string":"7"},{"string":"North Dakota"}]
    [{"string":"8"},{"string":"West Virginia"}]
    [{"string":"9"},{"string":"Alabama"}]
    [{"string":"10"},{"string":"Idaho"}]
    '
    assert_json_eq "$case" "$(sort <<<"$expected_rows")" "$actual_rows"
}

case3() {
    prefix="export some features"
    arg=$(cat <<-EOF
    {
       "features": ["account.state","account.credit_score","account.account_age_days","account.has_2fa_installed", "user_fake_stream.f1"],
       "unix_milli": $unix_milli
    }
EOF
)
    actual=$(testgrpc ChannelExport <<<"$arg")

    case="$prefix - first response returns header"
    actual_header=$(jq -s '.[0].header' <<< "$actual")
    expected_header='["user","account.state","account.credit_score","account.account_age_days","account.has_2fa_installed","user_fake_stream.f1"]'
    assert_json_eq "$case" "$expected_header" "$actual_header"

    case="$prefix - returns correct rows"
    actual_rows=$(jq -c ".row" <<< "$actual" | sort)
    expected_rows='
    [{"string":"1"},{"string":"Nevada"},{"int64":"530"},{"int64":"242"},{"bool":true},{"int64": "10"}]
    [{"string":"2"},{"string":"South Carolina"},{"int64":"520"},{"int64":"268"},{"bool":false},{}]
    [{"string":"3"},{"string":"New Jersey"},{"int64":"655"},{"int64":"84"},{"bool":false},{}]
    [{"string":"4"},{"string":"Ohio"},{"int64":"677"},{"int64":"119"},{"bool":true},{}]
    [{"string":"5"},{"string":"California"},{"int64":"566"},{"int64":"289"},{"bool":false},{}]
    [{"string":"6"},{"string":"North Carolina"},{"int64":"533"},{"int64":"155"},{"bool":true},{}]
    [{"string":"7"},{"string":"North Dakota"},{"int64":"605"},{"int64":"334"},{"bool":true},{}]
    [{"string":"8"},{"string":"West Virginia"},{"int64":"664"},{"int64":"282"},{"bool":false},{}]
    [{"string":"9"},{"string":"Alabama"},{"int64":"577"},{"int64":"150"},{"bool":true},{}]
    [{"string":"10"},{"string":"Idaho"},{"int64":"693"},{"int64":"212"},{"bool":true},{}]
    '
    assert_json_eq "$case" "$(sort <<<"$expected_rows")" "$actual_rows"
}

case1
case2
case3
