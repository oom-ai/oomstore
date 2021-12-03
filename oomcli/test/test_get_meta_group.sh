#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

case='oomcli get meta group works'
expected='ID,NAME,ENTITY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,phone,device,phone,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
2,student,user,student
'
actual=$(oomcli get meta group -o csv --wide)
ignore_time() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"

case='oomcli get simplified group works'
expected='ID,NAME,ENTITY,DESCRIPTION
1,phone,device,phone
2,student,user,student
'
actual=$(oomcli get meta group -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get one group works'
expected='ID,NAME,ENTITY,DESCRIPTION
1,phone,device,phone
'
actual=$(oomcli get meta group phone -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get meta group -o yaml: one group'
expected='
kind: Group
name: phone
entity-name: device
category: batch
description: phone
features:
    - name: price
      db-value-type: int
      description: price
    - name: model
      db-value-type: varchar(32)
      description: model
'

actual=$(oomcli get meta group phone -o yaml)
assert_eq "$case" "$expected" "$actual"

case='oomcli get meta group -o yaml: multiple groups'
expected='
items:
    - kind: Group
      name: phone
      entity-name: device
      category: batch
      description: phone
      features:
        - name: price
          db-value-type: int
          description: price
        - name: model
          db-value-type: varchar(32)
          description: model
    - kind: Group
      name: student
      entity-name: user
      category: batch
      description: student
      features:
        - name: age
          db-value-type: int
          description: age
'
actual=$(oomcli get meta group -o yaml)
assert_eq "$case" "$expected" "$actual"
