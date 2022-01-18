#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

case="query single feature"
arg='
{
    "entity_key": "19",
    "features": ["account.state"]
}
'
expected='
{
  "result": {
    "map": {
      "account.state": {
        "string": "Illinois"
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
    "features": ["account.state", "account.credit_score", "transaction_stats.transaction_count_7d", "transaction_stats.transaction_count_30d"]
}
'
expected='
{
  "result": {
    "map": {
      "account.credit_score": {
        "int64": "708"
      },
      "account.state": {
        "string": "Indiana"
      },
      "transaction_stats.transaction_count_30d": {
        "int64": "45"
      },
      "transaction_stats.transaction_count_7d": {
        "int64": "5"
      }
    }
  }
}
'
actual=$(testgrpc OnlineGet <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
