#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

oomcli apply -f ./data/driver_stats.yaml

for i in {1..5}; do
    import_sample driver_stats "./data/driver_stats_revision_$i.csv" "$i"
done

case="api returns correct result"
arg=$(cat <<-EOF
{
  "feature_names": [
    "conv_rate",
    "acc_rate",
    "avg_daily_trips"
  ],
  "entity_row": {
    "entity_key": "1",
    "unix_milli": 3
  }
}
{
  "entity_row": {
    "entity_key": "7",
    "unix_milli": 0
  }
}
EOF
)

expected='
{
  "status": {},
  "header": [
    "entity_key",
    "unix_milli",
    "conv_rate",
    "acc_rate",
    "avg_daily_trips"
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
{
  "status": {},
  "joinedRow": "TODO"
}
'
actual=$(testgrpc ChannelJoin <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
