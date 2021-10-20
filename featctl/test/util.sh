#!/usr/bin/env bash
set -euo pipefail

PATH="$SDIR/../build:$PATH"

info() { printf "$(date +'%Y/%m/%d %H:%M:%S') [info] %s\n" "$*" >&2; }
erro() { printf "$(date +'%Y/%m/%d %H:%M:%S') %b[erro]%b %s\n" '\e[0;31m\033[1m' '\e[0m' "$*" >&2; }

BLD=$(tput bold    2>/dev/null || true)
RST=$(tput sgr0    2>/dev/null || true)
GRN=$(tput setaf 2 2>/dev/null || true)
YLW=$(tput setaf 3 2>/dev/null || true)

export FEATCTL_HOST=127.0.0.1
export FEATCTL_PORT=4000
export FEATCTL_USER=test
export FEATCTL_PASS=test

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
    featctl register entity device --length 32
    featctl register entity user --length 64 --description "all users"
    featctl register group phone --entity device
    featctl register batch-feature price --group phone --value-type "int"
    featctl register batch-feature model --group phone --value-type "varchar(32)"
}

# import sample data
import_sample() {
    info "import sample data..."
    featctl import \
    --group phone \
    --separator "," \
    --input-file device.csv \
    --description 'test data'
}

execute_sql() {
    mysql -h 127.0.0.1 -u root -P 4000 -e "$1"
}

init_store() {
    info "initialize feature store"

    # create test user
    execute_sql "CREATE USER IF NOT EXISTS 'test'@'%' IDENTIFIED BY 'test'"
    execute_sql "GRANT ALL PRIVILEGES ON *.* TO 'test'@'%' WITH GRANT OPTION"

    # destroy database
    execute_sql 'drop database if exists onestore'

    # initialize feature store
    featctl init
}
