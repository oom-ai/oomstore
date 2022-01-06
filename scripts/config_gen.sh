#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
set -euo pipefail
IFS=$'\n\t'

usage() { echo "Usage: $(basename "$0") <online_store> <offline_store> <metadata_store>" >&2; }

info() { printf "%b[info]%b %s\n" '\e[0;32m\033[1m' '\e[0m' "$*" >&2; }
warn() { printf "%b[warn]%b %s\n" '\e[0;33m\033[1m' '\e[0m' "$*" >&2; }
erro() { printf "%b[erro]%b %s\n" '\e[0;31m\033[1m' '\e[0m' "$*" >&2; }

[ $# -ne 3 ] && usage && exit 1

online=$1
offline=$2
metadata=$3

filter_config() {
    local name=$1
    for source in ./config_source*.yaml; do
        cfg=$(sed -n "/$name:/, /^\$/p" "$source" | grep -v '^$')
        [[ -n $cfg ]] && echo "$cfg" && return
    done

    erro "config not found: '$name'"
    exit 1
}

indent() { sed 's/^/  /'; }

online_cfg=$(filter_config   "$online"   | sed 's/ext:$/:/')
offline_cfg=$(filter_config  "$offline"  | sed 's/ext:$/:/')
metadata_cfg=$(filter_config "$metadata" | sed 's/ext:$/:/')

cat <<-EOF
online-store:
$(indent <<< "$online_cfg")

offline-store:
$(indent <<< "$offline_cfg")

metadata-store:
$(indent <<< "$metadata_cfg")
EOF
