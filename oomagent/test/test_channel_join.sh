#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

oomcli apply -f ./data/driver_stats.yaml

for i in {1..5}; do
    import_sample driver_stats "./data/driver_stats_revision_$i.csv" "$i"
done

case1="api returns correct result"
arg1='
{
  "feature_full_names": [
    "driver_stats.conv_rate",
    "driver_stats.acc_rate",
    "driver_stats.avg_daily_trips"
  ],
  "entity_row": {
    "entity_key": "1",
    "unix_milli": 3
  }
}
'
expected1='
{
  "status": {},
  "header": [
    "entity_key",
    "unix_milli",
    "driver_stats.conv_rate",
    "driver_stats.acc_rate",
    "driver_stats.avg_daily_trips"
  ],
  "joinedRow": [
    {
      "stringValue": "1"
    },
    {
      "int64Value": "3"
    },
    {
      "doubleValue": 0.556
    },
    {
      "doubleValue": 0.465
    },
    {
      "int64Value": "464"
    }
  ]
}
'

case2="no header in subsequent request"
arg2='
{
  "entity_row": {
    "entity_key": "7",
    "unix_milli": 1
  }
}
'
expected2='
{
  "status": {},
  "joinedRow": [
    {
      "stringValue": "7"
    },
    {
      "int64Value": "1"
    },
    {
      "doubleValue": 0.758
    },
    {
      "doubleValue": 0.02
    },
    {
      "int64Value": "389"
    }
  ]
}
'

case3="handle null value correctly"
arg3='
{
  "entity_row": {
    "entity_key": "7",
    "unix_milli": 0
  }
}
'
expected3='
{
  "status": {},
  "joinedRow": [
    {
      "stringValue": "7"
    },
    {
      "int64Value": "0"
    },
    {
      "nullValue": "NULL_VALUE"
    },
    {
      "nullValue": "NULL_VALUE"
    },
    {
      "nullValue": "NULL_VALUE"
    }
  ]
}
'

arg="$arg1"
arg="$arg$arg2"
 arg="$arg$arg3"

actual=$(testgrpc ChannelJoin <<<"$arg")

assert_json_eq "$case1" "$expected1" "$(jq -s '.[0]' <<< "$actual")"
assert_json_eq "$case2" "$expected2" "$(jq -s '.[1]' <<< "$actual")"
assert_json_eq "$case3" "$expected3" "$(jq -s '.[2]' <<< "$actual")"
