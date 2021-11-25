#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

prepare_store
prepare_oomd

case="it works"
expected='
{
  "status": {}
}
'
import_sample account ./data/account_10.csv
actual=$(grpcurl -protoset ../../proto/oomd.protoset -plaintext -d @ localhost:50051 oomd.OomD/Sync <<EOF
{
    "revision_id": "3"
}
EOF
)
assert_json_eq "$case" "$expected" "$actual"
