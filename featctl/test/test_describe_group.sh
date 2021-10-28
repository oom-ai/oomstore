#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample

case='featctl describe group works'
expected='
Name:                     phone
Entity:                   device
Description:
Online Revision:          <NULL>
Offline Latest Revision:  <NULL>
Offline Latest DataTable: <NULL>
CreateTime:
ModifyTime:
'
actual=$(featctl describe group phone)
ignore_time() { grep -Ev '^(CreateTime|ModifyTime|Online Revision|Offline Latest Revision|Offline Latest DataTable)' <<<"$1"; }
