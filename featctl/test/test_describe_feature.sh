#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample > /dev/null

case='featctl describe feature works'
expected='
Name:                     price
Group:                    phone
Entity:                   device
Category:                 batch
DBValueType:              int
ValueType:                int32
Description:
Online Revision:          <NULL>
Offline Latest Revision:  <NULL>
Offline Latest DataTable: <NULL>
CreateTime:
ModifyTime:
'
actual=$(featctl describe feature price)
ignore_time() { grep -Ev '^(CreateTime|ModifyTime|Online Revision|Offline Latest Revision|Offline Latest DataTable)' <<<"$1"; }
