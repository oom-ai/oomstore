#!/usr/bin/env bash
set -euo pipefail

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store

trap 'rm -f "$TMPFILE"' EXIT
TMPFILE=$(mktemp) || exit 1
cat > "$TMPFILE" <<EOF
kind: Entity
name: user
length: 8
description: 'description'
---
kind: Group
name: device
entity-name: user
category: batch
description: 'description'
---
kind: Feature
name: model
group-name: device
category: batch
db-value-type: varchar(16)
description: 'description'
---
kind: Feature
name: price
group-name: device
category: batch
db-value-type: int
description: 'description'
EOF

oomcli apply -f "$TMPFILE"
feature_expected='Name,Group,Entity,Category,DBValueType,ValueType,Description,OnlineRevisionID
model,device,user,batch,varchar(16),string,description,<NULL>
price,device,user,batch,int,int64,description,<NULL>
'
feature_actual=$(oomcli get meta feature -o csv)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "oomcli get meta feature" "$(sort <<< "$feature_expected")" "$(ignore_time "$feature_actual" | sort)"
