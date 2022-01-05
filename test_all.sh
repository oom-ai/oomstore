#!/usr/bin/env bash

onlineStore=("mysql" "postgres" "sqlite" "redis" "cassandra" "dynamodb" "tidb" "tidbext" "tikv" "tikvext")
offlineStore=("mysql" "postgres" "sqlite" "tidb" "tidbext")
metadataStore=("mysql" "postgres" "sqlite" "tidb" "tidbext")

for online in ${onlineStore[@]}; do
  for offline in ${offlineStore[@]}; do
    for metadata in ${metadataStore[@]}; do
      echo "=== RUN online-$online,offline-$offline,metadata-$metadata ==="
      BACKENDS=$online,$offline,$metadata make integration-test
    done;
  done;
done;
