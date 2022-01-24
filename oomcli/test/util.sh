#!/usr/bin/env bash
set -euo pipefail
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PDIR="$SDIR/../.."

PATH="$SDIR/../build:$PATH"

info() { printf "$(date +'%Y/%m/%d %H:%M:%S') [info] %s\n" "$*" >&2; }
erro() { printf "$(date +'%Y/%m/%d %H:%M:%S') %b[erro]%b %s\n" '\e[0;31m\033[1m' '\e[0m' "$*" >&2; }

BLD=$(tput bold    2>/dev/null || true)
RST=$(tput sgr0    2>/dev/null || true)
GRN=$(tput setaf 2 2>/dev/null || true)
YLW=$(tput setaf 3 2>/dev/null || true)

export OOMCLI_CONFIG="$SDIR/config.yaml"
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

# register features for the sample data
register_features() {
    oomcli register entity device --description "device"
    oomcli register entity user   --description "user"

    oomcli register group phone      --entity device --category "batch"  --description "phone"
    oomcli register group student    --entity user   --category "batch"  --description "student"
    oomcli register group user-click \
      --entity user \
      --category "stream" \
      --snapshot_interval "1s" \
      --description "user click post feature"

    oomcli register feature price  --group phone   --value-type "int64"  --description "price"
    oomcli register feature model  --group phone   --value-type "string" --description "model"

    oomcli register feature name   --group student --value-type "string" --description "name"
    oomcli register feature gender --group student --value-type "string" --description "gender"
    oomcli register feature age    --group student --value-type "int64"  --description "age"

    oomcli register feature last_5_click_posts \
      --group user-click \
      --value-type "string" \
      --description "user last 5 click posts"

    oomcli register feature number_of_user_starred_posts \
      --group user-click \
      --value-type "int64"  \
      --description "number of posts that users starred today"
}

# import sample data
import_device_sample() {
    info "import sample data to offline store..."
    local revision=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
    oomcli import \
    --group phone \
    --revision "$revision" \
    --delimiter "," \
    --input-file device.csv \
    --description 'test data' | grep -o '[0-9]\+'
}

import_student_sample() {
    info "import student sample data to offline store..."
    local revision=${1:-$(perl -MTime::HiRes=time -E 'say int(time * 1000)')}
    oomcli import \
    --group student \
    --revision "$revision" \
    --delimiter "," \
    --input-file student.csv \
    --description 'test data' | grep -o '[0-9]\+'
}

# sync feature values from offline store to online store
sync() {
    info "sync sample data to online store"
    oomcli sync --group-name "$1" --revision-id "$2"
}

init_store() {
    info "initialize feature store"

    # initialize backend database
    oomplay init "$ONLINE_STORE" "$OFFLINE_STORE" "$METADATA_STORE"

    # initialize feature store
    # for db such as mysql, a successful ping in one process
    # doesn't mean the other processes can connect to it
    # so we try 5 times until we give up
    for _i in {1..5}; do
        oomcli init && break
        sleep 2
    done
}
