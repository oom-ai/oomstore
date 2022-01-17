#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

oomcli apply -f ./data/driver_stats.yaml

for i in {1..5}; do
    import_sample driver_stats "./data/driver_stats_revision_$i.csv" "$i"
done

arg='
{
  "join_features": [
    "driver_stats.conv_rate",
    "driver_stats.acc_rate",
    "driver_stats.avg_daily_trips"
  ],
  "entity_row": {
    "entity_key": "1",
    "unix_milli": 3
  }
}
{
  "entity_row": {
    "entity_key": "7",
    "unix_milli": 1
  }
}
{
  "entity_row": {
    "entity_key": "7",
    "unix_milli": 0
  }
}
'

actual=$(testgrpc ChannelJoin <<<"$arg")

case="first response contains header"
actual_header=$(jq -s '.[0].header' <<< "$actual")
expected_header='
[
    "entity_key",
    "unix_milli",
    "driver_stats.conv_rate",
    "driver_stats.acc_rate",
    "driver_stats.avg_daily_trips"
]
'
assert_json_eq "$case" "$expected_header" "$actual_header"

case="api returns correct joined rows"
actual_rows=$(jq -c ".joinedRow" <<< "$actual" | sort)
expected_rows='
[{"stringValue":"1"},{"int64Value":"3"},{"doubleValue":0.556},{"doubleValue":0.465},{"int64Value":"464"}]
[{"stringValue":"7"},{"int64Value":"1"},{"doubleValue":0.758},{"doubleValue":0.02},{"int64Value":"389"}]
[{"stringValue":"7"},{"int64Value":"0"},{"nullValue":"NULL_VALUE"},{"nullValue":"NULL_VALUE"},{"nullValue":"NULL_VALUE"}]
'
assert_json_eq "$case" "$(sort <<<"$expected_rows")" "$actual_rows"
