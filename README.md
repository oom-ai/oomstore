<h1 align="center">OomStore</h1>
<p align="center">
    <em>Lightweight and Fast Feature Store Powered by Go</em>
</p>

<p align="center">
    <a href="https://github.com/oom-ai/oomstore/actions/workflows/ci.yml">
        <img src="https://github.com/oom-ai/oomstore/actions/workflows/ci.yml/badge.svg" alt="CICD"/>
    </a>
    <a href="https://goreportcard.com/report/oom-ai/oomstore">
        <img src="https://goreportcard.com/badge/oom-ai/oomstore" alt="Version">
    </a>
    <a href="http://godoc.org/github.com/oom-ai/oomstore">
        <img src="https://godoc.org/github.com/oom-ai/oomstore?status.png" alt="Version">
    </a>
    <a href="https://codecov.io/gh/oom-ai/oomstore">
        <img src="https://codecov.io/gh/oom-ai/oomstore/branch/main/graph/badge.svg?token=C59L7LTRM4" alt="Platform"/>
    </a>
</p>

<p align="center">
  <a href="https://oom.ai/docs/quickstart">Quickstart</a>
  <span> · </span>
  <a href="https://oom.ai/docs/architecture">Architecture</a>
  <span> · </span>
  <a href="https://oom.ai/docs/benchmark">Benchmark</a>
</p>

## Overview

oomstore allows you to:

- Define features with YAML.
- Store features in databases of choice.
- Retrieve features for both online serving and offline training, **fast**.

Please see our [docs](https://oom.ai/docs) for more details.

## Features

- 🍼 Simple. Being serverless and CLI-friendly, users can be productive in hours, not months.
- 🔌 Composable. We support your preferred [databases of choice](https://oom.ai/docs/supported-databases).
- ⚡ Fast. [Benchmark](https://oom.ai/docs/benchmark) shows oomstore performs QPS > 50k and latency < 0.3 ms with Redis.
- 🌊 Streaming. We support streaming features to ensure your predictions are up-to-date.

## Architecture

<p align="center">
  <img src="https://oom.ai/images/architecture/architecture.svg" alt="Architecture">
</p>

See [Architecture](https://oom.ai/docs/architecture) for more details.

## Quickstart

1. Install `oomcli` following the [guide](https://oom.ai/docs/installation#cli).

2. `oomcli init` to initialize oomstore. Make sure there is a `~/.config/oomstore/config.yaml` as below.

```yaml
online-store:
  sqlite:
    db-file: /tmp/oomstore.db

offline-store:
  sqlite:
    db-file: /tmp/oomstore.db

metadata-store:
  sqlite:
    db-file: /tmp/oomstore.db
```

3. `oomcli apply -f metadata.yaml` to register metadata. See metadata.yaml below.

```yaml
kind: Entity
name: user
description: 'user ID'
groups:
- name: account
  category: batch
  description: 'user account'
  features:
  - name: state
    value-type: string
  - name: credit_score
    value-type: int64
  - name: account_age_days
    value-type: int64
  - name: 2fa_installed
    value-type: bool
- name: txn_stats
  category: batch
  description: 'user txn stats'
  features:
  - name: count_7d
    value-type: int64
  - name: count_30d
    value-type: int64
- name: recent_txn_stats
  category: stream
  snapshot-interval: 24h
  description: 'user recent txn stats'
  features:
  - name: count_10min
    value-type: int64
```

4. Import CSV data to Offline Store, then sync from Offline to Online Store

```bash
oomcli import \
  --group account \
  --input-file account.csv \
  --description 'sample account data'
oomcli import \
  --group txn_stats \
  --input-file txn_stats.csv \
  --description 'sample txn stats data'
oomcli sync --group-name account --revision-id 2
oomcli sync --group-name txn_stats --revision-id 4
```

5. Push stream data to both Online and Offline Store.

```bash
oomcli push --group recent_txn_stats --entity-key 1006 --feature count_10min=1
```

6. Fetch features by key.

```bash
oomcli get online \
  --entity-key 1006 \
  --feature account.state,account.credit_score,account.account_age_days,account.2fa_installed,txn_stats.count_7d,txn_stats.count_30d,recent_txn_stats.count_10min
```

```text
+------+---------------+----------------------+--------------------------+-----------------------+--------------------+---------------------+------------------------------+
| user | account.state | account.credit_score | account.account_age_days | account.2fa_installed | txn_stats.count_7d | txn_stats.count_30d | recent_txn_stats.count_10min |
+------+---------------+----------------------+--------------------------+-----------------------+--------------------+---------------------+------------------------------+
| 1006 | Louisiana     |                  710 |                       32 | false                 |                  8 |                  22 |                            1 |
+------+---------------+----------------------+--------------------------+-----------------------+--------------------+---------------------+------------------------------+
```

7. Generate training datasets via Point-in-Time Join.

```sh
oomcli join \
	--feature account.state,account.credit_score,account.account_age_days,account.2fa_installed,txn_stats.count_7d,txn_stats.count_30d,recent_txn_stats.count_10min \
	--input-file label.csv
```

```text
+------------+---------------+---------------+----------------------+--------------------------+-----------------------+--------------------+---------------------+------------------------------+
| entity_key |  unix_milli   | account.state | account.credit_score | account.account_age_days | account.2fa_installed | txn_stats.count_7d | txn_stats.count_30d | recent_txn_stats.count_10min |
+------------+---------------+---------------+----------------------+--------------------------+-----------------------+--------------------+---------------------+------------------------------+
|       1002 | 1950236233000 | Hawaii        |                  625 |                      861 | true                  |                 11 |                  36 |                              |
|       1003 | 1950411318000 | Arkansas      |                  730 |                      958 | false                 |                  0 |                  16 |                              |
|       1004 | 1950653614000 | Louisiana     |                  610 |                     1570 | false                 |                 12 |                  26 |                              |
|       1005 | 1950166137000 | South Dakota  |                  635 |                     1953 | false                 |                  7 |                  30 |                              |
|       1006 | 1950403162000 | Louisiana     |                  710 |                       32 | false                 |                  8 |                  22 |                            1 |
|       1007 | 1950160030000 | New Mexico    |                  645 |                       37 | true                  |                  5 |                  40 |                              |
|       1008 | 1950274859000 | Nevada        |                  735 |                     1627 | false                 |                 12 |                  51 |                              |
|       1009 | 1949958846000 | Kentucky      |                  650 |                       88 | true                  |                 11 |                  23 |                              |
|       1010 | 1949920686000 | Delaware      |                  680 |                     1687 | false                 |                  2 |                  39 |                              |
+------------+---------------+---------------+----------------------+--------------------------+-----------------------+--------------------+---------------------+------------------------------+
```

See [Quickstart](https://oom.ai/docs/quickstart) for more complete details.

## Supported Databases

### Online Store

- Amazon DynamoDB
- Redis
- TiKV
- Cassandra
- PostgreSQL
- MySQL
- TiDB
- SQLite

### Offline Store

- Snowflake
- Amazon Redshift
- Google BigQuery
- PostgreSQL
- MySQL
- TiDB
- SQLite

### Metadata Store

- PostgreSQL
- MySQL
- TiDB
- SQLite

## Community

Feel free to [join the community](https://oom.ai/docs/community) for questions and feature requests!

## Credits

oomstore is highly inspired by [Feast](https://github.com/feast-dev/feast).
