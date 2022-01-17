#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomagent

case="query single feature"
arg='
{
    "entity_keys": ["19", "50", "78"],
    "feature_names": ["account.state"]
}
'
expected='
{
  "result": {
    "19": {
      "map": {
        "account.state": {
          "stringValue": "Illinois"
        }
      }
    },
    "50": {
      "map": {
        "account.state": {
          "stringValue": "Hawaii"
        }
      }
    },
    "78": {
      "map": {
        "account.state": {
          "stringValue": "Tennessee"
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
    "feature_names": ["account.state", "account.credit_score", "transaction_stats.transaction_count_7d", "transaction_stats.transaction_count_30d"]
}
'
expected='
{
  "result": {
    "48": {
      "map": {
        "account.credit_score": {
          "int64Value": "708"
        },
        "account.state": {
          "stringValue": "Indiana"
        },
        "transaction_stats.transaction_count_30d": {
          "int64Value": "45"
        },
        "transaction_stats.transaction_count_7d": {
          "int64Value": "5"
        }
      }
    },
    "74": {
      "map": {
        "account.credit_score": {
          "int64Value": "703"
        },
        "account.state": {
          "stringValue": "Ohio"
        },
        "transaction_stats.transaction_count_30d": {
          "int64Value": "25"
        },
        "transaction_stats.transaction_count_7d": {
          "int64Value": "8"
        }
      }
    }
  }
}
'
actual=$(testgrpc OnlineMultiGet <<<"$arg")
assert_json_eq "$case" "$expected" "$actual"
