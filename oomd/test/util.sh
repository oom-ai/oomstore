#!/usr/bin/env bash
set -euo pipefail
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

PATH="$SDIR/../build:$PATH"
PATH="$SDIR/../../featctl/build:$PATH"

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

prepare_store() {
    info "initialize feature store"
    execute_sql 'drop database if exists oomstore_test'

    # initialize feature store
    featctl init
    info "create oomstore schema..."
    featctl apply -f ./data/schema.yaml

    info "import sample data to offline store..."
    featctl import \
        --group account \
        --input-file \
        ./data/account.csv \
        --description 'sample account data' >/dev/null
    featctl import \
        --group transaction_stats \
        --input-file \
        ./data/transaction_stats.csv \
        --description 'sample transaction_stats data' >/dev/null

    info "sync sample data to online store"
    featctl sync -r 1
    featctl sync -r 2
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
