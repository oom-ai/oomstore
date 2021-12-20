#!/usr/bin/env bash
set -euo pipefail
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

PATH="$SDIR/../build:$PATH"

info() { printf "$(date +'%Y/%m/%d %H:%M:%S') [info] %s\n" "$*" >&2; }
erro() { printf "$(date +'%Y/%m/%d %H:%M:%S') %b[erro]%b %s\n" '\e[0;31m\033[1m' '\e[0m' "$*" >&2; }

BLD=$(tput bold    2>/dev/null || true)
RST=$(tput sgr0    2>/dev/null || true)
GRN=$(tput setaf 2 2>/dev/null || true)
YLW=$(tput setaf 3 2>/dev/null || true)

export OOMCLI_CONFIG="$SDIR/config.yaml"

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

# register features for the sample data
register_features() {
    oomcli register entity device --length 32     --description "device"
    oomcli register entity user   --length 64     --description "user"

    oomcli register group phone    --entity device  --description "phone"
    oomcli register group student  --entity user    --description "student"

    oomcli register batch-feature price --group phone    --db-value-type "int"          --description "price"
    oomcli register batch-feature model --group phone    --db-value-type "varchar(32)"  --description "model"
    oomcli register batch-feature age   --group student  --db-value-type "int"          --description "age"
}

# import sample data
import_sample() {
    info "import sample data to offline store..."
    local revision=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
    oomcli import \
    --group phone \
    --revision "$revision" \
    --delimiter "," \
    --input-file device.csv \
    --description 'test data' | grep -o '[0-9]\+'
}

# sync feature values from offline store to online store
sync() {
    info "sync sample data to online store"
    oomcli sync -r "$1"
}

init_store() {
    info "initialize feature store"

    # initialize backend database
    oomplay init -c "$OOMCLI_CONFIG"

    # initialize feature store
    oomcli init
}
