#!/usr/bin/env bash

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features

case='oomcli get meta group works'
expected='ID,NAME,ENTITY,CATEGORY,SNAPSHOT-INTERVAL,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,phone,device,batch,0s,phone,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
2,student,user,batch,0s,student
3,user-click,user,stream,1s,user click post feature
'
actual=$(oomcli get meta group -o csv --wide)
ignore_time() { cut -d ',' -f 1-5 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"

case='oomcli get simplified group works'
expected='ID,NAME,ENTITY,CATEGORY,SNAPSHOT-INTERVAL,DESCRIPTION
1,phone,device,batch,0s,phone
2,student,user,batch,0s,student
3,user-click,user,stream,1s,user click post feature
'
actual=$(oomcli get meta group -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get one group works'
expected='ID,NAME,ENTITY,CATEGORY,SNAPSHOT-INTERVAL,DESCRIPTION
1,phone,device,batch,0s,phone
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
      - name: name
        value-type: string
        description: name
      - name: gender
        value-type: string
        description: gender
      - name: age
        value-type: int64
        description: age
  - kind: Group
    name: user-click
    entity-name: user
    category: stream
    snapshot-interval: 1s
    description: user click post feature
    features:
      - name: last_5_click_posts
        value-type: string
        description: user last 5 click posts
      - name: number_of_user_starred_posts
        value-type: int64
        description: number of posts that users starred today
'
actual=$(oomcli get meta group -o yaml)
assert_eq "$case" "$expected" "$actual"
