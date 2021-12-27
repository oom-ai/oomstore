#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='oomcli get meta features works'
expected='
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION,ONLINE-REVISION-ID
1,price,phone,device,batch,int64,price,<NULL>
2,model,phone,device,batch,string,model,<NULL>
3,age,student,user,batch,int64,age,<NULL>
'
actual=$(oomcli get meta feature -o csv --wide)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
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
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION,ONLINE-REVISION-ID
2,model,phone,device,batch,string,model,<NULL>
'
actual=$(oomcli get meta feature phone.model -o csv --wide)
ignore_time() { cut -d ',' -f 1-8 <<<"$1"; }
assert_eq "$case" "$(sort <<< "$expected")" "$(ignore_time "$actual" | sort)"


case='oomcli get meta feature in yaml: one feature'
expected='
kind: Feature
name: model
group-name: phone
value-type: string
description: model
'

actual=$(oomcli get meta feature phone.model -o yaml)
assert_eq "$case" "$expected" "$actual"

case='oomcli get meta feature: multiple features'
expected='
items:
    - kind: Feature
      name: price
      group-name: phone
      value-type: int64
      description: price
    - kind: Feature
      name: model
      group-name: phone
      value-type: string
      description: model
    - kind: Feature
      name: age
      group-name: student
      value-type: int64
      description: age
'

actual=$(oomcli get meta feature -o yaml)
assert_eq "$case" "$expected" "$actual"
