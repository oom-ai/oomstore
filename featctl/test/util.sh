#!/usr/bin/env bash
set -euo pipefail

PATH="$SDIR/../build:$HOME/bin:$PATH"

info() { printf "$(date +'%Y/%m/%d %H:%M:%S') [info] %s\n" "$*" >&2; }
erro() { printf "$(date +'%Y/%m/%d %H:%M:%S') [erro] %s\n" "$*" >&2; }

BLD=$(tput bold)
RST=$(tput sgr0)
GRN=$(tput setaf 2)

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
      erro "expected:"
      echo "$expected"
      erro "actual:"
      echo "$actual"
      erro "diff:"
      diff  <(echo "$expected" ) <(echo "$actual")
      return 1
  fi
}

# import sample data with specified revision
import_sample() {
    local revision=$1
    info "import sample data $revision"
    trap 'command rm -rf -- lightning.* *.tmp' EXIT INT TERM HUP

    featctl import \
    --group device \
    --revision "$revision" \
    --schema-template schema.sql \
    --input-file device.csv \
    --has-header \
    --description 'test data' &>/dev/null
}

# register features in sample data with specified revision
register_features() {
    local revision=$1
    info "register features with revision: $revision"
    # register feature price
    featctl register feature --group device --category batch --revision "$revision" --name price --description 'device average price'
    # register feature model
    featctl register feature --group device --category batch --revision "$revision" --name model --description 'device model name'
}

execute_sql() {
    mysql -h 127.0.0.1 -u root -P 4000 -e "$1"
}

init_store() {
    info "initilaize feature store"

    # create test user
    execute_sql "CREATE USER IF NOT EXISTS 'test'@'%' IDENTIFIED BY 'test'"
    execute_sql "GRANT ALL PRIVILEGES ON *.* TO 'test'@'%' WITH GRANT OPTION"

    # destroy database
    execute_sql 'drop database if exists onestore'

    # initilaize feature store
    featctl init
}
