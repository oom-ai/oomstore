#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

oomcli apply -f ./data/user_click.yaml

case="push stream feature"
arg='
{
    "entity_key": "1",
    "group_name": "user-click",
    "feature_names": ["last_5_click_posts", "number_of_user_starred_posts"],
    "feature_values": [
    {
        "stringValue": "1,2,3"
    },
    {
        "int64Value": 10
    }
]
}
'
expected='
{
  "status": {}
}
'

# wait informer refresh
sleep 1
actual=$(testgrpc Push <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="query stream feature"
arg='
{
    "entity_key": "1",
    "feature_full_names": ["user-click.last_5_click_posts", "user-click.number_of_user_starred_posts"]
}
'
expected='
{
  "status": {

  },
  "result": {
    "map": {
      "user-click.last_5_click_posts": {
        "stringValue": "1,2,3"
      },
      "user-click.number_of_user_starred_posts": {
        "int64Value": "10"
      }
    }
  }
}
'
actual=$(testgrpc OnlineGet <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
