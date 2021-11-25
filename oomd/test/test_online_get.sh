#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomd

case="query single feature"
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
actual=$(grpcurl -protoset ../../proto/oomd.protoset -plaintext -d @ localhost:50051 oomd.OomD/OnlineGet <<EOF
{
    "entity_key": "19",
    "feature_names": ["state"]
}
EOF
)
assert_json_eq "$case" "$expected" "$actual"

case="query multiple features"
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
actual=$(grpcurl -protoset ../../proto/oomd.protoset -plaintext -d @ localhost:50051 oomd.OomD/OnlineGet <<EOF
{
    "entity_key": "48",
    "feature_names": ["state", "credit_score", "transaction_count_7d", "transaction_count_30d"]
}
EOF
)
assert_json_eq "$case" "$expected" "$actual"
