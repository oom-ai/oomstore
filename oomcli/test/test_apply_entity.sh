#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

apply_single_complex_entity() {
    init_store

    cat <<EOF | oomcli apply -f /dev/stdin
kind: Entity
name: user
length: 8
description: 'description'
batch-features:
- group: device
  description: a description
  features:
  - name: model
    value-type: string
    description: 'description'
  - name: price
    value-type: int64
    description: 'description'
- group: user
  description: a description
  features:
  - name: age
    value-type: int64
    description: 'description'
  - name: gender
    value-type: int64
    description: 'description'
EOF

    entity_expected='
ID,NAME,LENGTH,DESCRIPTION
1,user,8,description
'
    entity_actual=$(oomcli get meta entity -o csv)
    assert_eq "apply_single_complex_entity: check entity" "$(sort <<< "$entity_expected")" "$(sort <<< "$entity_actual")"

    group_expected='
ID,NAME,ENTITY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,device,user,a description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
2,user,user,a description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
'
    group_actual=$(oomcli get meta group -o csv --wide)
    filter() { cut -d ',' -f 1-4 <<<"$1"; }
    assert_eq "apply_single_complex_entity: check group" "$(filter "$group_expected" | sort)" "$(filter "$group_actual" | sort)"

    feature_expected='
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,model,device,user,batch,string,description
2,price,device,user,batch,int64,description
3,age,user,user,batch,int64,description
4,gender,user,user,batch,int64,description
'
    feature_actual=$(oomcli get meta feature -o csv)
    assert_eq "apply_single_complex_entity: check feature" "$(sort <<< "$feature_expected")" "$(sort <<< "$feature_actual")"
}

apply_multiple_files_of_entity() {
    init_store

    cat <<EOF | oomcli apply -f /dev/stdin
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


  entity_expected='
ID,NAME,LENGTH,DESCRIPTION
1,user,8,description
2,device,16,description
3,test,32,description
'
    entity_actual=$(oomcli get meta entity -o csv)
    assert_eq "apply_multiple_files_of_entity: oomcli get meta entity" "$entity_expected" "$entity_actual"

    group_expected='
ID,NAME,ENTITY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,student,user,student feature group,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
'
    group_actual=$(oomcli get meta group -o csv --wide)
    filter() { cut -d ',' -f 1-4 <<<"$1"; }
    assert_eq "oapply_multiple_files_of_entity: check group" "$(filter "$group_expected"| sort)" "$(filter "$group_actual" | sort)"
}

apply_entity_items() {
    init_store

    cat <<EOF | oomcli apply -f /dev/stdin
items:
  - kind: Entity
    name: user
    length: 8
    description: user ID
    batch-features:
      - group: account
        description: user account info
        features:
          - name: credit_score
            value-type: int64
            description: credit_score description
          - name: account_age_days
            value-type: int64
            description: account_age_days description
          - name: has_2fa_installed
            value-type: bool
            description: has_2fa_installed description
      - group: transaction_stats
        description: user transaction statistics
        features:
          - name: transaction_count_7d
            value-type: int64
            description: transaction_count_7d description
          - name: transaction_count_30d
            value-type: int64
            description: transaction_count_30d description
  - kind: Entity
    name: device
    length: 8
    description: device info
    batch-features:
      - group: phone
        description: phone info
        features:
          - name: model
            value-type: string
            description: model description
          - name: price
            value-type: int64
            description: price description
EOF

      entity_expected='
ID,NAME,LENGTH,DESCRIPTION
1,user,8,user ID
2,device,8,device info
'
      entity_actual=$(oomcli get meta entity -o csv)
      assert_eq "apply_entity_items: oomcli apply mutiple entity: check entity" "$(sort <<< "$entity_expected")" "$(sort <<< "$entity_actual")"

    group_expected='
ID,NAME,ENTITY,DESCRIPTION
1,account,user,user account info
2,transaction_stats,user,user transaction statistics
3,phone,device,phone info
'
      group_actual=$(oomcli get meta group -o csv)
      assert_eq "apply_entity_items: oomcli apply multiple entity: check group" "$group_expected" "$group_actual"

      feature_expected='
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,credit_score,account,user,batch,int64,credit_score description
2,account_age_days,account,user,batch,int64,account_age_days description
3,has_2fa_installed,account,user,batch,bool,has_2fa_installed description
4,transaction_count_7d,transaction_stats,user,batch,int64,transaction_count_7d description
5,transaction_count_30d,transaction_stats,user,batch,int64,transaction_count_30d description
6,model,phone,device,batch,string,model description
7,price,phone,device,batch,int64,price description
'
      feature_actual=$(oomcli get meta feature -o csv)
      assert_eq "apply_entity_items: oomcli apply multiple entity: check feature" "$feature_expected" "$feature_actual"
}

apply_single_complex_entity
apply_multiple_files_of_entity
apply_entity_items
