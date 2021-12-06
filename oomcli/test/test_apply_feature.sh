#!/usr/bin/env bash
set -euo pipefail

SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store

trap 'rm -f "$TMPFILE"' EXIT
TMPFILE=$(mktemp) || exit 1
cat > "$TMPFILE" <<EOF
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
kind: Feature
name: model
group-name: device
category: batch
db-value-type: varchar(16)
description: 'description'
---
kind: Feature
name: price
group-name: device
category: batch
db-value-type: int
description: 'description'
EOF

oomcli apply -f "$TMPFILE"
feature_expected='ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,model,device,user,batch,string,description
2,price,device,user,batch,int64,description
'
feature_actual=$(oomcli get meta feature -o csv)
assert_eq "oomcli get meta feature" "$(sort <<< "$feature_expected")" "$(sort <<< "$feature_actual")"

cat > "$TMPFILE" <<EOF
kind: Entity
name: user
length: 8
description: user ID
---
kind: Group
name: account
entity-name: user
category: batch
description: user account info
---
kind: Group
name: transaction_stats
entity-name: user
category: batch
description: user transaction statistics
---
items:
  - kind: Feature
    name: credit_score
    group-name: account
    db-value-type: int
    description: "credit_score description"
  - kind: Feature
    name: account_age_days
    group-name: account
    db-value-type: int
    description: "account_age_days description"
  - kind: Feature
    name: has_2fa_installed
    group-name: account
    db-value-type: bool
    description: "has_2fa_installed description"
  - kind: Feature
    name: transaction_count_7d
    group-name: transaction_stats
    db-value-type: int
    description: "transaction_count_7d description"
  - kind: Feature
    name: transaction_count_30d
    group-name: transaction_stats
    db-value-type: int
    description: "transaction_count_30d description"
EOF

init_store
oomcli apply -f "$TMPFILE"
feature_expected='
ID,NAME,GROUP,ENTITY,CATEGORY,VALUE-TYPE,DESCRIPTION
1,credit_score,account,user,batch,int64,credit_score description
2,account_age_days,account,user,batch,int64,account_age_days description
3,has_2fa_installed,account,user,batch,bool,has_2fa_installed description
4,transaction_count_7d,transaction_stats,user,batch,int64,transaction_count_7d description
5,transaction_count_30d,transaction_stats,user,batch,int64,transaction_count_30d description
'
feature_actual=$(oomcli get meta feature -o csv)
assert_eq "oomcli apply multiple features" "$(sort <<< "$feature_expected")" "$(sort <<< "$feature_actual")"
