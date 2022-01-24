#!/usr/bin/env bash
set -euo pipefail
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PDIR="$SDIR/../.."

PATH="$SDIR/../build:$PATH"
PATH="$SDIR/../../oomcli/build:$PATH"

PROTO_DIR="$SDIR/../../proto"

info() { printf "$(date +'%Y/%m/%d %H:%M:%S') [info] %s\n" "$*" >&2; }
erro() { printf "$(date +'%Y/%m/%d %H:%M:%S') %b[erro]%b %s\n" '\e[0;31m\033[1m' '\e[0m' "$*" >&2; }

BLD=$(tput bold    2>/dev/null || true)
RST=$(tput sgr0    2>/dev/null || true)
GRN=$(tput setaf 2 2>/dev/null || true)
YLW=$(tput setaf 3 2>/dev/null || true)

export OOMCLI_CONFIG="$SDIR/config.yaml"
export OOMAGENT_CONFIG="$OOMCLI_CONFIG"

BACKENDS=${BACKENDS:-"postgres,postgres,postgres"}
ONLINE_STORE=$(cut -d, -f1 <<<"$BACKENDS")
OFFLINE_STORE=$(cut -d, -f2 <<<"$BACKENDS")
METADATA_STORE=$(cut -d, -f3 <<<"$BACKENDS")
"$PDIR/scripts/config_gen.sh" "$ONLINE_STORE" "$OFFLINE_STORE" "$METADATA_STORE" > "$OOMCLI_CONFIG"

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
    local revision=${3:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
    info "import sample data '$file' into group '$group'..."
    oomcli import \
        --group "$group" \
        --revision "$revision" \
        --input-file \
        "$file" \
        --description "sample $group data from $(basename "$file")" >/dev/null
}

prepare_store() {
    info "initialize feature store"
    oomplay init "$ONLINE_STORE" "$OFFLINE_STORE" "$METADATA_STORE"

    # initialize feature store
    # for db such as mysql, a successful ping in one process
    # doesn't mean the other processes can connect to it
    # so we try 5 times until we give up
    for _i in {1..5}; do
        oomcli init && break
        sleep 2
    done

    info "create oomstore schema..."
    oomcli apply -f ./data/fraud_detection.yaml

    info "import sample data to offline store..."
    import_sample account ./data/account_100.csv
    import_sample transaction_stats ./data/transaction_stats_100.csv

    info "sync sample data to online store"
    oomcli sync --group-name account --revision-id 1
    oomcli sync --group-name transaction_stats --revision-id 2
}

testgrpc() {
    grpcurl --import-path "$PROTO_DIR" --proto "$PROTO_DIR/oomagent.proto" -plaintext -d @ localhost:50051 "oomagent.OomAgent/$1" "${@:1:}"
}

wait_for_ready() {
    local interval=$1
    local retries=$2
    local message
    for ((i = 0; i < retries; i++ )); do
        if message=$(testgrpc HealthCheck <<<"" 2>&1); then
            return
        fi
        sleep "$interval"
    done
    erro "oomagent is still not ready: $message"
}

prepare_oomagent() {
    info "start oomagent server..."
    trap 'kill $(jobs -p)' EXIT INT TERM HUP
    oomagent &
    wait_for_ready 0.1 10
}
