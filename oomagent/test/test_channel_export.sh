#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

import_sample account "./data/account_10.csv"

case="export some features"
arg='
{
    "feature_names": ["state"],
    "revision_id": 3
}
'
expected='
{"status":{},"header":["user","state"],"row":[{"stringValue":"1"},{"stringValue":"Nevada"}]}
{"status":{},"row":[{"stringValue":"2"},{"stringValue":"South Carolina"}]}
{"status":{},"row":[{"stringValue":"3"},{"stringValue":"New Jersey"}]}
{"status":{},"row":[{"stringValue":"4"},{"stringValue":"Ohio"}]}
{"status":{},"row":[{"stringValue":"5"},{"stringValue":"California"}]}
{"status":{},"row":[{"stringValue":"6"},{"stringValue":"North Carolina"}]}
{"status":{},"row":[{"stringValue":"7"},{"stringValue":"North Dakota"}]}
{"status":{},"row":[{"stringValue":"8"},{"stringValue":"West Virginia"}]}
{"status":{},"row":[{"stringValue":"9"},{"stringValue":"Alabama"}]}
{"status":{},"row":[{"stringValue":"10"},{"stringValue":"Idaho"}]}
'
actual=$(testgrpc ChannelExport <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="export all features"
arg='
{
    "revision_id": 3
}
'
expected='
{"status":{},"header":["user","state","credit_score","account_age_days","has_2fa_installed"],"row":[{"stringValue":"1"},{"stringValue":"Nevada"},{"int64Value":"530"},{"int64Value":"242"},{"boolValue":true}]}
{"status":{},"row":[{"stringValue":"2"},{"stringValue":"South Carolina"},{"int64Value":"520"},{"int64Value":"268"},{"boolValue":false}]}
{"status":{},"row":[{"stringValue":"3"},{"stringValue":"New Jersey"},{"int64Value":"655"},{"int64Value":"84"},{"boolValue":false}]}
{"status":{},"row":[{"stringValue":"4"},{"stringValue":"Ohio"},{"int64Value":"677"},{"int64Value":"119"},{"boolValue":true}]}
{"status":{},"row":[{"stringValue":"5"},{"stringValue":"California"},{"int64Value":"566"},{"int64Value":"289"},{"boolValue":false}]}
{"status":{},"row":[{"stringValue":"6"},{"stringValue":"North Carolina"},{"int64Value":"533"},{"int64Value":"155"},{"boolValue":true}]}
{"status":{},"row":[{"stringValue":"7"},{"stringValue":"North Dakota"},{"int64Value":"605"},{"int64Value":"334"},{"boolValue":true}]}
{"status":{},"row":[{"stringValue":"8"},{"stringValue":"West Virginia"},{"int64Value":"664"},{"int64Value":"282"},{"boolValue":false}]}
{"status":{},"row":[{"stringValue":"9"},{"stringValue":"Alabama"},{"int64Value":"577"},{"int64Value":"150"},{"boolValue":true}]}
{"status":{},"row":[{"stringValue":"10"},{"stringValue":"Idaho"},{"int64Value":"693"},{"int64Value":"212"},{"boolValue":true}]}
'
actual=$(testgrpc ChannelExport <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
