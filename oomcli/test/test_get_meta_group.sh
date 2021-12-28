#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

case='oomcli get meta group works'
expected='ID,NAME,ENTITY,CATEGORY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,phone,device,batch,phone,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
2,student,user,batch,student
3,user-click,user,stream,user click post feature
'
actual=$(oomcli get meta group -o csv --wide)
ignore_time() { cut -d ',' -f 1-5 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"

case='oomcli get simplified group works'
expected='ID,NAME,ENTITY,CATEGORY,DESCRIPTION
1,phone,device,batch,phone
2,student,user,batch,student
3,user-click,user,stream,user click post feature
'
actual=$(oomcli get meta group -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get one group works'
expected='ID,NAME,ENTITY,CATEGORY,DESCRIPTION
1,phone,device,batch,phone
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
      value-type: int64
      description: price
    - name: model
      value-type: string
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
          value-type: int64
          description: price
        - name: model
          value-type: string
          description: model
    - kind: Group
      name: student
      entity-name: user
      category: batch
      description: student
      features:
        - name: age
          value-type: int64
          description: age
    - kind: Group
      name: user-click
      entity-name: user
      category: stream
      description: user click post feature
'
actual=$(oomcli get meta group -o yaml)
assert_eq "$case" "$expected" "$actual"
