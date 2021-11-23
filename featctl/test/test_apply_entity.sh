#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

trap 'rm -f "$TMPFILE"' EXIT
TMPFILE=$(mktemp) || exit 1

init_store

cat > "$TMPFILE" <<EOF
kind: entity
name: user
length: 8
description: 'description'
batch-features:
- group: student
  description: student feature group
---
kind: entity
name: device
length: 16
description: 'description'
---
kind: entity
name: test
length: 32
description: 'description'
EOF

featctl apply -f "$TMPFILE"

entity_expected='Name,Length,Description
user,8,description
device,16,description
test,32,description
'
entity_actual=$(featctl list entity -o csv)
ignore_time() { cut -d ',' -f 1-3 <<<"$1"; }
assert_eq "featctl list entity" "$(sort <<< "$entity_expected")" "$(ignore_time "$entity_actual" | sort)"

group_expected='
GroupName,GroupID,EntityName,Description
student,1,user,student feature group
'
group_actual=$(featctl list group -o csv)
ignore_time() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "featctl list group" "$(sort <<< "$group_expected")" "$(ignore_time "$group_actual" | sort)"


init_store

cat > "$TMPFILE" <<EOF
kind: entity
name: user
length: 8
description: 'description'
batch-features:
- group: device
  description: a description
  features:
  - name: model
    db-value-type: varchar(16)
    description: 'description'
  - name: price
    db-value-type: int
    description: 'description'
- group: user
  description: a description
  features:
  - name: age
    db-value-type: int
    description: 'description'
  - name: gender
    db-value-type: int
    description: 'description'
EOF

featctl apply -f "$TMPFILE"

entity_expected='Name,Length,Description
user,8,description
'
entity_actual=$(featctl list entity -o csv)
ignore_time() { cut -d ',' -f 1-3 <<<"$1"; }
assert_eq "featctl list entity" "$(sort <<< "$entity_expected")" "$(ignore_time "$entity_actual" | sort)"

group_expected='
GroupName,GroupID,EntityName,Description
device,1,user,a description
user,2,user,a description
'
group_actual=$(featctl list group -o csv)
ignore_time() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "featctl list group" "$(sort <<< "$group_expected")" "$(ignore_time "$group_actual" | sort)"

feature_expected='Name,Group,Entity,Category,DBValueType,ValueType,Description,OnlineRevisionID
model,device,user,batch,varchar(16),string,description,<NULL>
price,device,user,batch,int,int64,description,<NULL>
age,user,user,batch,int,int64,description,<NULL>
gender,user,user,batch,int,int64,description,<NULL>
'
feature_actual=$(featctl list feature -o csv)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "featctl list feature" "$(sort <<< "$feature_expected")" "$(ignore_time "$feature_actual" | sort)"
