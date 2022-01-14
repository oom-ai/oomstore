#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

init_store
register_features
import_sample >> /dev/null

case='get all'
expected='device,phone.price
1,3999
2,5299
3,3999
4,1999
5,999
6,4999
7,5999
8,6500
9,4500
'
timestamp=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
actual=$(oomcli export --feature phone.price --unix-milli $timestamp -o csv)
assert_eq "$case" "$(sort <<< "$expected")" "$(sort <<< "$actual")"
