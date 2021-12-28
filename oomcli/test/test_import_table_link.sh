#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample

case='oomcli import using table link'
oomcli import \
    --group phone \
    --table-link "offline_batch_1_1" \
    --description 'linked table' > /dev/null
expected='
ID,REVISION,GROUP,SNAPSHOT-TABLE,CDC-TABLE,DESCRIPTION
1,1639047117470,phone,offline_batch_1_1,,test data
2,1639047117552,phone,offline_batch_1_1,,linked table
'
filter() { cut -d, -f 1,3,4,5; }
actual=$(oomcli get meta revision -o csv)
assert_eq "$case" "$(filter <<<"$expected")" "$(filter <<<"$actual")"
