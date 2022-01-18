#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

case="query single feature"
arg='
{
    "entity_keys": ["19", "50", "78"],
    "features": ["account.state"]
}
'
expected='
{
  "result": {
    "19": {
      "map": {
        "account.state": {
          "string": "Illinois"
        }
      }
    },
    "50": {
      "map": {
        "account.state": {
          "string": "Hawaii"
        }
      }
    },
    "78": {
      "map": {
        "account.state": {
          "string": "Tennessee"
        }
      }
    }
  }
}
'
actual=$(testgrpc OnlineMultiGet <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"

case="query multiple features"
arg='
{
    "entity_keys": ["48", "74"],
    "features": ["account.state", "account.credit_score", "transaction_stats.transaction_count_7d", "transaction_stats.transaction_count_30d"]
}
'
expected='
{
  "result": {
    "48": {
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
    },
    "74": {
      "map": {
        "account.credit_score": {
          "int64": "703"
        },
        "account.state": {
          "string": "Ohio"
        },
        "transaction_stats.transaction_count_30d": {
          "int64": "25"
        },
        "transaction_stats.transaction_count_7d": {
          "int64": "8"
        }
      }
    }
  }
}
'
actual=$(testgrpc OnlineMultiGet <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
