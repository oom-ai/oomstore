#!/usr/bin/env bash
set -euo pipefail

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store

trap 'rm -f "$TMPFILE"' EXIT
TMPFILE=$(mktemp) || exit 1
cat > "$TMPFILE" <<EOF
kind: entity
name: user
length: 8
description: 'description'
---
kind: group
name: device
entity-name: user
category: batch
description: 'description'
---
kind: feature
name: model
group-name: device
category: batch
db-value-type: varchar(16)
description: 'description'
---
kind: feature
name: price
group-name: device
category: batch
db-value-type: int
description: 'description'
EOF

featctl apply -f "$TMPFILE"
feature_expected='Name,Group,Entity,Category,DBValueType,ValueType,Description,OnlineRevisionID
model,device,user,batch,varchar(16),string,description,<NULL>
price,device,user,batch,int,int64,description,<NULL>
'
feature_actual=$(featctl list feature -o csv)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "featctl list feature" "$(sort <<< "$feature_expected")" "$(ignore_time "$feature_actual" | sort)"
