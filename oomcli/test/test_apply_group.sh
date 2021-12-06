#!/usr/bin/env bash
set -euo pipefail

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

apply_single_complex_group() {
    init_store

    cat <<EOF | oomcli apply -f /dev/stdin
kind: Entity
name: user
length: 8
description: 'description'
---
kind: Group
name: device
entity-name: user
category: batch
description: 'description'
features:
- name: model
  db-value-type: varchar(16)
  description: 'description'
- name: price
  db-value-type: int
  description: 'description'
- name: radio
  db-value-type: int
  description: 'description'
EOF

    group_expected='
ID,NAME,ENTITY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,device,user,description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
'
    group_actual=$(oomcli get meta group -o csv --wide)
    filter() { cut -d ',' -f 1-4 <<<"$1"; }
    assert_eq "apply_single_complex_group: check group" "$(filter "$group_expected" | sort)" "$(filter "$group_actual" | sort)"

    feature_expected='
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,model,device,user,batch,string,description
2,price,device,user,batch,int64,description
3,radio,device,user,batch,int64,description
'
    feature_actual=$(oomcli get meta feature -o csv)
    assert_eq "apply_single_complex_group: check feature" "$(sort <<< "$feature_expected")" "$(sort <<< "$feature_actual")"
}

apply_multiple_files_of_group() {
    init_store

    cat <<EOF | oomcli apply -f /dev/stdin
kind: Entity
name: user
length: 8
description: 'description'
---
kind: Group
name: device
entity-name: user
category: batch
description: 'description'
---
kind: Group
name: account
entity-name: user
category: batch
description: 'description'
EOF

    group_expected='
ID,NAME,ENTITY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,device,user,description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
2,account,user,description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
'
    group_actual=$(oomcli get meta group -o csv --wide)
    filter() { cut -d ',' -f 1-4 <<<"$1"; }
    assert_eq "apply_multiple_files_of_group: check group" "$(filter "$group_expected" | sort)" "$(filter "$group_actual" | sort)"
}

apply_group_items() {
    init_store

    cat <<EOF | oomcli apply -f /dev/stdin
kind: Entity
name: user
length: 8
description: 'description'
---
kind: Entity
name: device
length: 8
description: 'description'
---
items:
  - kind: Group
    name: account
    entity-name: user
    category: batch
    description: user account info
    features:
      - name: state
        db-value-type: varchar(32)
        description: ""
      - name: credit_score
        db-value-type: int
        description: credit_score description
      - name: account_age_days
        db-value-type: int
        description: account_age_days description
      - name: has_2fa_installed
        db-value-type: bool
        description: has_2fa_installed description
  - kind: Group
    name: transaction_stats
    entity-name: user
    category: batch
    description: user transaction statistics
    features:
      - name: transaction_count_7d
        db-value-type: int
        description: transaction_count_7d description
      - name: transaction_count_30d
        db-value-type: int
        description: transaction_count_30d description
  - kind: Group
    name: phone
    entity-name: device
    category: batch
    description: phone info
    features:
      - name: model
        db-value-type: varchar(32)
        description: model description
      - name: price
        db-value-type: int
        description: price description
EOF

    group_expected='
ID,NAME,ENTITY,DESCRIPTION
1,account,user,user account info
2,transaction_stats,user,user transaction statistics
3,phone,device,phone info
'
    group_actual=$(oomcli get meta group -o csv --wide)
    filter() { cut -d ',' -f 1-4 <<<"$1"; }
    assert_eq "apply_single_complex_group: check features" "$(filter "$group_expected" | sort)" "$(filter "$group_actual" | sort)"

    feature_expected='
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,state,account,user,batch,string,
2,credit_score,account,user,batch,int64,credit_score description
3,account_age_days,account,user,batch,int64,account_age_days description
4,has_2fa_installed,account,user,batch,bool,has_2fa_installed description
5,transaction_count_7d,transaction_stats,user,batch,int64,transaction_count_7d description
6,transaction_count_30d,transaction_stats,user,batch,int64,transaction_count_30d description
7,model,phone,device,batch,string,model description
8,price,phone,device,batch,int64,price description
'
    feature_actual=$(oomcli get meta feature -o csv)
    assert_eq "apply_single_complex_group: check group" "$feature_expected" "$feature_actual"
}

apply_single_complex_group
apply_multiple_files_of_group
apply_group_items
