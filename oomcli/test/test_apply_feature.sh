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
feature_expected='ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,model,device,user,batch,string,description
2,price,device,user,batch,int64,description
'
feature_actual=$(oomcli get meta feature -o csv)
assert_eq "oomcli get meta feature" "$(sort <<< "$feature_expected")" "$(sort <<< "$feature_actual")"
