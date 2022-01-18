#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli get meta entity works'
expected='ID,NAME,DESCRIPTION,CREATE-TIME,MODIFY-TIME
1,device,device,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
2,user,user,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
'
actual=$(oomcli get meta entity -o csv)
ignore_time() { cut -d ',' -f 1-3 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"


case='oomcli get simplified meta entity works'
expected='ID,NAME,DESCRIPTION
1,device,device
2,user,user
'
actual=$(oomcli get meta entity -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get meta entity -o yaml: one entity'
expected='
kind: Entity
name: device
description: device
groups:
    - name: phone
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

actual=$(oomcli get meta entity device -o yaml)
assert_eq "$case" "$expected" "$actual"

case='oomcli get meta entity -o yaml: multiple entities'
expected='
items:
    - kind: Entity
      name: device
      description: device
      groups:
        - name: phone
          category: batch
          description: phone
          features:
            - name: price
              value-type: int64
              description: price
            - name: model
              value-type: string
              description: model
    - kind: Entity
      name: user
      description: user
      groups:
        - name: student
          category: batch
          description: student
          features:
            - name: age
              value-type: int64
              description: age
        - name: user-click
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

actual=$(oomcli get meta entity -o yaml)
assert_eq "$case" "$expected" "$actual"
