#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomd

case="query single feature"
arg='
{
    "entity_key": "19",
    "feature_names": ["state"]
}
'
expected='
{
  "status": {

  },
  "result": {
    "map": {
      "state": {
        "stringValue": "Illinois"
      }
    }
  }
}
'
actual=$(testgrpc OnlineGet <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="query multiple features"
arg='
{
    "entity_key": "48",
    "feature_names": ["state", "credit_score", "transaction_count_7d", "transaction_count_30d"]
}
'
expected='
{
  "status": {},
  "result": {
    "map": {
      "credit_score": {
        "int64Value": "708"
      },
      "state": {
        "stringValue": "Indiana"
      },
      "transaction_count_30d": {
        "int64Value": "45"
      },
      "transaction_count_7d": {
        "int64Value": "5"
      }
    }
  }
}
'
actual=$(testgrpc OnlineGet <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
