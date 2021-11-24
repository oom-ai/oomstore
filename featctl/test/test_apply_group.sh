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
kind: Group
name: account
entity-name: user
category: batch
description: 'description'
EOF

featctl apply -f "$TMPFILE"

group_expected='
GroupName,GroupID,EntityName,Description
device,1,user,description
account,2,user,description
'
group_actual=$(featctl list group -o csv)
ignore_time() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "featctl list group" "$(sort <<< "$group_expected")" "$(ignore_time "$group_actual" | sort)"


init_store

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
features:
- name: model
  db-value-type: varchar(16)
  description: 'description'
- name: price
  db-value-type: int
  description: 'description'
- name: radio
  db-value-type: int
  description: 'description'
EOF

featctl apply -f "$TMPFILE"

group_expected='
GroupName,GroupID,EntityName,Description
device,1,user,description
'
group_actual=$(featctl list group -o csv)
ignore_time() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "featctl list group" "$(sort <<< "$group_expected")" "$(ignore_time "$group_actual" | sort)"

feature_expected='Name,Group,Entity,Category,DBValueType,ValueType,Description,OnlineRevisionID
model,device,user,batch,varchar(16),string,description,<NULL>
price,device,user,batch,int,int64,description,<NULL>
radio,device,user,batch,int,int64,description,<NULL>
'
feature_actual=$(featctl list feature -o csv)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "featctl list feature" "$(sort <<< "$feature_expected")" "$(ignore_time "$feature_actual" | sort)"
