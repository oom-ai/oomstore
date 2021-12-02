#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

trap 'rm -f "$TMPFILE"' EXIT
TMPFILE=$(mktemp) || exit 1

init_store

cat > "$TMPFILE" <<EOF
kind: Entity
name: user
length: 8
description: 'description'
batch-features:
- group: student
  description: student feature group
---
kind: Entity
name: device
length: 16
description: 'description'
---
kind: Entity
name: test
length: 32
description: 'description'
EOF

oomcli apply -f "$TMPFILE"

entity_expected='ID,NAME,LENGTH,DESCRIPTION
1,user,8,description
2,device,16,description
3,test,32,description
'
entity_actual=$(oomcli get meta entity -o csv)
#ignore_time() { cut -d ',' -f 1-3 <<<"$1"; }
assert_eq "oomcli get meta entity" "$(sort <<< "$entity_expected")" "$(sort <<< "$entity_actual")"

group_expected='
ID,NAME,ENTITY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,student,user,student feature group,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
'
group_actual=$(oomcli get meta group -o csv --wide)
filter() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "oomcli get meta group" "$(filter "$group_expected"| sort)" "$(filter "$group_actual" | sort)"


init_store

cat > "$TMPFILE" <<EOF
kind: Entity
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

oomcli apply -f "$TMPFILE"

entity_expected='ID,NAME,LENGTH,DESCRIPTION
1,user,8,description
'
entity_actual=$(oomcli get meta entity -o csv)
assert_eq "oomcli get meta entity" "$(sort <<< "$entity_expected")" "$(sort <<< "$entity_actual")"

group_expected='
ID,NAME,ENTITY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,device,user,a description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
2,user,user,a description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z

'
group_actual=$(oomcli get meta group -o csv --wide)
filter() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "oomcli get meta group" "$(filter "$group_expected" | sort)" "$(filter "$group_actual" | sort)"

feature_expected='ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,model,device,user,batch,string,description
2,price,device,user,batch,int64,description
3,age,user,user,batch,int64,description
4,gender,user,user,batch,int64,description
'
feature_actual=$(oomcli get meta feature -o csv)
assert_eq "oomcli get meta feature" "$(sort <<< "$feature_expected")" "$(sort <<< "$feature_actual")"
