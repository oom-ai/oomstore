#!/usr/bin/env bash
#SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
#source ./util.sh
#
#init_store
#register_features
#import_sample
#
#case="query single feature"
#expected='
#entity_key,model
#1,xiaomi-mix3
#'
#actual=$(featctl query -g device -n model -k 1)
#assert_eq "$case" "$expected" "$actual"
#
#
#case="query multiple features"
#expected='
#entity_key,model,price
#6,apple-iphone11,4999
#'
#actual=$(featctl query -g device -n model,price -k 6)
#assert_eq "$case" "$expected" "$actual"
