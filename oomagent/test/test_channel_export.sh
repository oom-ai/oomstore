#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

import_sample account "./data/account_10.csv"

case1() {
    prefix="export some features"
    arg='
    {
        "feature_names": ["state"],
        "revision_id": 3
    }
    '
    actual=$(testgrpc ChannelExport <<<"$arg")

    case="$prefix - first response returns header"
    actual_header=$(jq -s '.[0].header' <<< "$actual")
    expected_header='["user","state"]'
    assert_json_eq "$case" "$expected_header" "$actual_header"

    case="$prefix - returns correct rows"
    actual_rows=$(jq -c ".row" <<< "$actual" | sort)
    expected_rows='
    [{"stringValue":"1"},{"stringValue":"Nevada"}]
    [{"stringValue":"2"},{"stringValue":"South Carolina"}]
    [{"stringValue":"3"},{"stringValue":"New Jersey"}]
    [{"stringValue":"4"},{"stringValue":"Ohio"}]
    [{"stringValue":"5"},{"stringValue":"California"}]
    [{"stringValue":"6"},{"stringValue":"North Carolina"}]
    [{"stringValue":"7"},{"stringValue":"North Dakota"}]
    [{"stringValue":"8"},{"stringValue":"West Virginia"}]
    [{"stringValue":"9"},{"stringValue":"Alabama"}]
    [{"stringValue":"10"},{"stringValue":"Idaho"}]
    '
    assert_json_eq "$case" "$(sort <<<"$expected_rows")" "$actual_rows"
}

case2() {
    prefix="export some features"
    arg='{"revision_id":3}'
    actual=$(testgrpc ChannelExport <<<"$arg")

    case="$prefix - first response returns header"
    actual_header=$(jq -s '.[0].header' <<< "$actual")
    expected_header='["user","state","credit_score","account_age_days","has_2fa_installed"]'
    assert_json_eq "$case" "$expected_header" "$actual_header"

    case="$prefix - returns correct rows"
    actual_rows=$(jq -c ".row" <<< "$actual" | sort)
    expected_rows='
    [{"stringValue":"1"},{"stringValue":"Nevada"},{"int64Value":"530"},{"int64Value":"242"},{"boolValue":true}]
    [{"stringValue":"2"},{"stringValue":"South Carolina"},{"int64Value":"520"},{"int64Value":"268"},{"boolValue":false}]
    [{"stringValue":"3"},{"stringValue":"New Jersey"},{"int64Value":"655"},{"int64Value":"84"},{"boolValue":false}]
    [{"stringValue":"4"},{"stringValue":"Ohio"},{"int64Value":"677"},{"int64Value":"119"},{"boolValue":true}]
    [{"stringValue":"5"},{"stringValue":"California"},{"int64Value":"566"},{"int64Value":"289"},{"boolValue":false}]
    [{"stringValue":"6"},{"stringValue":"North Carolina"},{"int64Value":"533"},{"int64Value":"155"},{"boolValue":true}]
    [{"stringValue":"7"},{"stringValue":"North Dakota"},{"int64Value":"605"},{"int64Value":"334"},{"boolValue":true}]
    [{"stringValue":"8"},{"stringValue":"West Virginia"},{"int64Value":"664"},{"int64Value":"282"},{"boolValue":false}]
    [{"stringValue":"9"},{"stringValue":"Alabama"},{"int64Value":"577"},{"int64Value":"150"},{"boolValue":true}]
    [{"stringValue":"10"},{"stringValue":"Idaho"},{"int64Value":"693"},{"int64Value":"212"},{"boolValue":true}]
    '
    assert_json_eq "$case" "$(sort <<<"$expected_rows")" "$actual_rows"
}

case1
case2
