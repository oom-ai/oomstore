#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli get meta features works'
expected='
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION,DB-VALUE-TYPE,ONLINE-REVISION-ID
1,price,phone,device,batch,int64,price,int,<NULL>
2,model,phone,device,batch,string,model,varchar(32),<NULL>
3,age,student,user,batch,int64,age,int,<NULL>
'
actual=$(oomcli get meta feature -o csv --wide)
ignore_time() { cut -d ',' -f 1-9 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"

case='oomcli get simplified meta features works'
expected='ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,price,phone,device,batch,int64,price
2,model,phone,device,batch,string,model
3,age,student,user,batch,int64,age
'
actual=$(oomcli get meta feature -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"

case='oomcli get one meta feature works'
expected='
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION,DB-VALUE-TYPE,ONLINE-REVISION-ID
2,model,phone,device,batch,string,model,varchar(32),<NULL>
'
actual=$(oomcli get meta feature model -o csv --wide)
ignore_time() { cut -d ',' -f 1-9 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"


case='oomcli get meta feature: one feature'
expected='
kind: Feature
name: model
group-name: phone
db-value-type: varchar(32)
description: model
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
      description: price
    - kind: Feature
      name: model
      group-name: phone
      db-value-type: varchar(32)
      description: model
    - kind: Feature
      name: age
      group-name: student
      db-value-type: int
      description: age
'

actual=$(oomcli get meta feature -o yaml)
assert_eq "$case" "$expected" "$actual"
