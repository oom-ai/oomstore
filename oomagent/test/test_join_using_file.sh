#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

oomcli apply -f ./data/driver_stats.yaml

for i in {1..5}; do
    import_sample driver_stats "./data/driver_stats_revision_$i.csv" "$i"
done

# We cannot trap multiple times in one process so we step into a subshell
(
trap 'command rm -rf $TEMP' EXIT INT TERM HUP
TEMP="$(mktemp -dt "$(basename "$0")".XXXXXX)"

output="$TEMP/result.csv"
case="api returns ok"
arg=$(cat <<-EOF
{
  "feature_full_names": [
    "driver_stats.conv_rate",
    "driver_stats.acc_rate",
    "driver_stats.avg_daily_trips"
  ],
  "input_file_path": "./data/driver_stats_label.csv",
  "output_file_path": "$output"
}
EOF
)
expected='
{
  "status": {}
}
'
actual=$(testgrpc Join <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="result is correct"
expected='
entity_key,unix_milli,driver_stats.conv_rate,driver_stats.acc_rate,driver_stats.avg_daily_trips
1,0,,,
1,3,0.556,0.465,464
1,4,0.377,0.991,329
2,3,0.934,0.082,646
3,3,0.892,0.222,148
4,4,0.413,0.445,107
5,2,0.259,0.954,833
5,2,0.259,0.954,833
5,2,0.259,0.954,833
5,2,0.259,0.954,833
5,3,0.794,0.887,371
5,4,0.567,0.846,714
5,5,0.751,0.939,281
6,1,0.289,0.481,169
6,3,0.951,0.189,433
7,0,,,
7,4,0.272,0.247,233
7,4,0.272,0.247,233
7,4,0.272,0.247,233
7,4,0.272,0.247,233
8,4,0.532,0.476,725
9,3,0.073,0.532,948
10,4,0.96,0.582,909
10,4,0.96,0.582,909
10,4,0.96,0.582,909
10,4,0.96,0.582,909
'
# sort result ignoring header
actual=$(head -1 "$output" && tail -n +2 "$output" | sort -t, -n -k 1,1)
assert_eq "$case" "$expected" "$actual"
)
