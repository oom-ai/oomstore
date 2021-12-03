#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli get meta entity works'
expected='ID,NAME,LENGTH,DESCRIPTION,CREATE-TIME,MODIFY-TIME
1,device,32,device,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
2,user,64,user,2021-10-19T06:56:07Z,2021-10-19T06:56:07Z
'
actual=$(oomcli get meta entity -o csv)
ignore_time() { cut -d ',' -f 1-4 <<<"$1"; }
assert_eq "$case" "$(ignore_time "$expected" | sort)" "$(ignore_time "$actual" | sort)"


case='oomcli get simplified meta entity works'
expected='ID,NAME,LENGTH,DESCRIPTION
1,device,32,device
2,user,64,user
'
actual=$(oomcli get meta entity -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get meta entity -o yaml: one entity'
expected='
kind: Entity
name: device
length: 32
description: device
batch-features:
    - group: phone
      description: phone
      features:
        - name: price
          db-value-type: int
          description: price
        - name: model
          db-value-type: varchar(32)
          description: model
'

actual=$(oomcli get meta entity device -o yaml)
assert_eq "$case" "$expected" "$actual"

case='oomcli get meta entity -o yaml: multiple entities'
expected='
items:
    - kind: Entity
      name: device
      length: 32
      description: device
      batch-features:
        - group: phone
          description: phone
          features:
            - name: price
              db-value-type: int
              description: price
            - name: model
              db-value-type: varchar(32)
              description: model
    - kind: Entity
      name: user
      length: 64
      description: user
      batch-features:
        - group: student
          description: student
          features:
            - name: age
              db-value-type: int
              description: age
'

actual=$(oomcli get meta entity -o yaml)
assert_eq "$case" "$expected" "$actual"
