#!/usr/bin/env bash
SDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) && cd "$SDIR" || exit 1
source ./util.sh

echo '=== VERSION INFO ==='
oomcli --version
echo

export DEBUG='enable'
for test_file in test_*.sh; do
    echo "=== RUN $test_file ==="
    "./$test_file"
    echo
done
