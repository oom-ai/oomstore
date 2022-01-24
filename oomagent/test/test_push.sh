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
    "group": "user-click",
    "feature_values": {
        "last_5_click_posts": {
            "string": "1,2,3"
        },
        "number_of_user_starred_posts": {
            "int64": 10
        }
    }
}
'
expected='{}'

# wait informer refresh
sleep 1
actual=$(testgrpc Push <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="query stream feature"
arg='
{
    "entity_key": "1",
    "features": ["user-click.last_5_click_posts", "user-click.number_of_user_starred_posts"]
}
'
expected='
{
  "result": {
    "map": {
      "user-click.last_5_click_posts": {
        "string": "1,2,3"
      },
      "user-click.number_of_user_starred_posts": {
        "int64": "10"
      }
    }
  }
}
'
actual=$(testgrpc OnlineGet <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="push stream feature with difference value"
arg='
{
    "entity_key": "1",
    "group": "user-click",
    "feature_values": {
        "last_5_click_posts": {
            "string": "2,3,4"
        },
        "number_of_user_starred_posts": {
            "int64": 11
        }
    }
}
'
expected='{}'

# wait informer refresh
sleep 1
actual=$(testgrpc Push <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="query stream feature"
arg='
{
    "entity_key": "1",
    "features": ["user-click.last_5_click_posts", "user-click.number_of_user_starred_posts"]
}
'
expected='
{
  "result": {
    "map": {
      "user-click.last_5_click_posts": {
        "string": "2,3,4"
      },
      "user-click.number_of_user_starred_posts": {
        "int64": "11"
      }
    }
  }
}
'
actual=$(testgrpc OnlineGet <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
