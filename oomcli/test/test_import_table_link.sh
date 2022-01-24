#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_device_sample

case='oomcli import using table link'
oomcli import \
    --group phone \
    --table-link "offline_batch_1_2" \
    --description 'linked table' > /dev/null
expected='
ID,REVISION,GROUP,SNAPSHOT-TABLE,CDC-TABLE,DESCRIPTION
1,0,phone,offline_stream_snapshot_1_0,
2,1639047117470,phone,offline_batch_1_2,,test data
3,1639047117552,phone,offline_batch_1_2,,linked table
'
filter() { cut -d, -f 1,3,4,5; }
actual=$(oomcli get meta revision -o csv)
assert_eq "$case" "$(filter <<<"$expected")" "$(filter <<<"$actual")"
