#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli get meta features works'
expected='NAME,GROUP,ENTITY,CATEGORY,DB-VALUE-TYPE,VALUE-TYPE,DESCRIPTION,ONLINE-REVISION-ID
model,phone,device,batch,varchar(32),string,model,<NULL>
price,phone,device,batch,int,int64,price,<NULL>
'
actual=$(oomcli get meta feature -o csv --wide)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"

case='oomcli get simplified meta features works'
expected='NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE
model,phone,device,batch,string
price,phone,device,batch,int64
'
actual=$(oomcli get meta feature -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get one meta feature works'
expected='NAME,GROUP,ENTITY,CATEGORY,DB-VALUE-TYPE,VALUE-TYPE,DESCRIPTION,ONLINE-REVISION-ID
model,phone,device,batch,varchar(32),string,model,<NULL>
'
actual=$(oomcli get meta feature model -o csv --wide)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"


case='oomcli get meta feature: one feature'
expected='
kind: Feature
name: model
group-name: phone
db-value-type: varchar(32)
description: varchar(32)
'

actual=$(oomcli get meta feature model -o yaml)
assert_eq "$case" "$expected" "$actual"

case='oomcli get meta feature: multiple features'
expected='
items:
    - kind: Feature
      name: price
      group-name: phone
      db-value-type: int
      description: int
    - kind: Feature
      name: model
      group-name: phone
      db-value-type: varchar(32)
      description: varchar(32)
'

actual=$(oomcli get meta feature -o yaml)
assert_eq "$case" "$expected" "$actual"
