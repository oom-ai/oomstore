#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

oomcli apply -f ./data/driver_stats.yaml

for i in {1..5}; do
    import_sample driver_stats "./data/driver_stats_revision_$i.csv" "$i"
done

successful_cases() {
  oomcli push --entity-key 1 --group fake_stream --feature f1=10
  t1=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
  oomcli snapshot fake_stream

  arg=$(cat <<-EOF
{
  "join_features": [
    "driver_stats.conv_rate",
    "driver_stats.acc_rate",
    "driver_stats.avg_daily_trips",
    "fake_stream.f1"
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
{
  "entity_row": {
    "entity_key": "1",
    "unix_milli": $t1
  }
}
EOF
)

  case="first response contains header"
  actual=$(testgrpc ChannelJoin <<<"$arg")
  actual_header=$(jq -s '.[0].header' <<< "$actual")
  expected_header='
[
    "entity_key",
    "unix_milli",
    "driver_stats.conv_rate",
    "driver_stats.acc_rate",
    "driver_stats.avg_daily_trips",
    "fake_stream.f1"
]
'
  assert_json_eq "$case" "$expected_header" "$actual_header"

  case="api returns correct joined rows"
  actual_rows=$(jq -c ".joinedRow" <<< "$actual" | sort)
  expected_rows=$(cat <<-EOF
[{"string":"1"},{"int64":"3"},{"double":0.556},{"double":0.465},{"int64":"464"},{}]
[{"string":"7"},{"int64":"1"},{"double":0.758},{"double":0.02},{"int64":"389"},{}]
[{"string":"7"},{"int64":"0"},{},{},{},{}]
[{"string":"1"},{"int64":"$t1"},{"double":0.146},{"double":0.031},{"int64":"286"},{"int64": "10"}]
EOF
)

assert_json_eq "$case" "$(sort <<<"$expected_rows")" "$actual_rows"
}


failed_cases() {
  case="invalid feature name"
  arg=$(cat <<-EOF
{
  "join_features": [
    "invalid_feature_name"
  ],
  "entity_row": {
    "entity_key": "1",
    "unix_milli": 3
  }
}
EOF
)

  actual=$(testgrpc ChannelJoin <<<"$arg" 2>&1 || flag=true)
  expected=$(cat <<-EOF
ERROR:
  Code: Internal
  Message: invalid full feature name: 'invalid_feature_name'
EOF
)
  assert_eq "$case" "$expected" "$actual"

  case="empty entity rows"
  arg=$(cat <<-EOF
{
  "join_features": [
    "driver_stats.conv_rate"
  ]
}
EOF
)
  actual=$(testgrpc ChannelJoin <<<"$arg" 2>&1 || flag=true)
  expected=$(cat <<-EOF
EOF
)
  assert_eq "$case" "$expected" "$actual"


   case="invalid request"
  arg=$(cat <<-EOF
EOF
)
  actual=$(testgrpc ChannelJoin <<<"$arg" 2>&1 || flag=true)
  expected=$(cat <<-EOF
ERROR:
  Code: InvalidArgument
  Message: invalid request: empty feature
EOF
)
  assert_eq "$case" "$expected" "$actual"
}

main() {
  successful_cases
  failed_cases
}

main
