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
  value-type: string
  description: 'description'
- name: price
  value-type: int64
  description: 'description'
- name: radio
  value-type: int64
  description: 'description'
EOF

    group_expected='
ID,NAME,ENTITY,CATEGORY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,device,user,batch,description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
'
    group_actual=$(oomcli get meta group -o csv --wide)
    filter() { cut -d ',' -f 1-5 <<<"$1"; }
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
ID,NAME,ENTITY,CATEGORY,DESCRIPTION,ONLINE-REVISION-ID,CREATE-TIME,MODIFY-TIME
1,device,user,batch,description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
2,account,user,batch,description,<NULL>,2021-11-30T07:51:03Z,2021-11-30T08:19:13Z
'
    group_actual=$(oomcli get meta group -o csv --wide)
    filter() { cut -d ',' -f 1-5 <<<"$1"; }
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
        value-type: string
        description: ""
      - name: credit_score
        value-type: int64
        description: credit_score description
      - name: account_age_days
        value-type: int64
        description: account_age_days description
      - name: has_2fa_installed
        value-type: bool
        description: has_2fa_installed description
  - kind: Group
    name: transaction_stats
    entity-name: user
    category: batch
    description: user transaction statistics
    features:
      - name: transaction_count_7d
        value-type: int64
        description: transaction_count_7d description
      - name: transaction_count_30d
        value-type: int64
        description: transaction_count_30d description
  - kind: Group
    name: phone
    entity-name: device
    category: batch
    description: phone info
    features:
      - name: model
        value-type: string
        description: model description
      - name: price
        value-type: int64
        description: price description
EOF

    group_expected='
ID,NAME,ENTITY,CATEGORY,DESCRIPTION
1,account,user,batch,user account info
2,transaction_stats,user,batch,user transaction statistics
3,phone,device,batch,phone info
'
    group_actual=$(oomcli get meta group -o csv --wide)
    filter() { cut -d ',' -f 1-5 <<<"$1"; }
    assert_eq "apply_single_complex_group: check group" "$(filter "$group_expected" | sort)" "$(filter "$group_actual" | sort)"

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
