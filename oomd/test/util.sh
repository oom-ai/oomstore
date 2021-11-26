#!/usr/bin/env bash
set -euo pipefail
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

PATH="$SDIR/../build:$PATH"
PATH="$SDIR/../../oomctl/build:$PATH"

PROTO_DIR="$SDIR/../../proto"

info() { printf "$(date +'%Y/%m/%d %H:%M:%S') [info] %s\n" "$*" >&2; }
erro() { printf "$(date +'%Y/%m/%d %H:%M:%S') %b[erro]%b %s\n" '\e[0;31m\033[1m' '\e[0m' "$*" >&2; }

BLD=$(tput bold    2>/dev/null || true)
RST=$(tput sgr0    2>/dev/null || true)
GRN=$(tput setaf 2 2>/dev/null || true)
YLW=$(tput setaf 3 2>/dev/null || true)

export FEATCTL_CONFIG="$SDIR/config.yaml"
export OOMD_CONFIG="$SDIR/config.yaml"

trim() {
    local var="$*"
    # remove leading whitespace characters
    var="${var#"${var%%[![:space:]]*}"}"
    # remove trailing whitespace characters
    var="${var%"${var##*[![:space:]]}"}"
    printf '%s' "$var"
}

assert_eq() {
  local case expected actual
  case="case - $1"
  expected="$(trim "$2")"
  actual="$(trim "$3")"

  if [ "$expected" == "$actual" ]; then
      info "${BLD}${GRN}Passed $case${RST}"
      return 0
  else
      erro "${BLD}${GRN}Failed $case${RST}"
      echo "${BLD}${YLW}=> expected:${RST}"
      echo "$expected"
      echo "${BLD}${YLW}=> actual:${RST}"
      echo "$actual"
      echo "${BLD}${YLW}=> diff:${RST}"
      diff --color=auto <(echo "$expected" ) <(echo "$actual")
      return 1
  fi
}

assert_json_eq() {
  local case expected actual
  case="case - $1"
  expected="$(jq <<< "$2")"
  actual="$(jq <<< "$3")"

  if [ "$expected" == "$actual" ]; then
      info "${BLD}${GRN}Passed $case${RST}"
      return 0
  else
      erro "${BLD}${GRN}Failed $case${RST}"
      echo "${BLD}${YLW}=> expected:${RST}"
      echo "$expected"
      echo "${BLD}${YLW}=> actual:${RST}"
      echo "$actual"
      echo "${BLD}${YLW}=> diff:${RST}"
      diff --color=auto <(echo "$expected" ) <(echo "$actual")
      return 1
  fi
}

import_sample() {
    local group=$1
    local file=$2
    info "import sample data '$file' into group '$group'..."
    oomctl import \
        --group "$group" \
        --revision "$(perl -MTime::HiRes=time -E 'say int(time * 1000)')" \
        --input-file \
        "$file" \
        --description 'sample account data' >/dev/null
}

prepare_store() {
    info "initialize feature store"
    execute_sql 'drop database if exists oomstore_test'

    # initialize feature store
    oomctl init
    info "create oomstore schema..."
    oomctl apply -f ./data/schema.yaml

    info "import sample data to offline store..."
    import_sample account ./data/account_100.csv
    import_sample transaction_stats ./data/transaction_stats_100.csv

    info "sync sample data to online store"
    oomctl sync -r 1
    oomctl sync -r 2
}

testgrpc() {
    grpcurl --import-path "$PROTO_DIR" --proto "$PROTO_DIR/oomd.proto" -plaintext -d @ localhost:50051 "oomd.OomD/$1" "${@:1:}"
}

prepare_oomd() {
    info "start oomd server..."
    trap 'kill $(jobs -p)' EXIT INT TERM HUP
    oomd &
    sleep 1
}

execute_sql() {
    PGPASSWORD=postgres psql -h localhost -U postgres -c "$1" >/dev/null
}
