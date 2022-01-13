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

oomstore's edges:

- 🍼 Simple. Being serverless and CLI-friendly, users can be productive in hours, not months.
- 🔌 Composable. We support your preferred [databases of choice](https://oom.ai/docs/supported-databases).
- ⚡ Fast. [Benchmark](https://oom.ai/docs/benchmark) shows oomstore performs QPS > 50k and latency < 0.3 ms with Redis.
- 🌊 Streaming. We support streaming features to ensure your predictions are up-to-date.

## Architecture

<p align="center">
  <img src="https://oom.ai/images/architecture/architecture.svg" alt="Architecture">
</p>

You can interact with oomstore with CLI, Go API or Python API. See [Architecture](https://oom.ai/docs/architecture) for more details.

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

3. `oomcli apply -f config.yaml` to register metadata. See config.yaml below.

```yaml
kind: Entity
name: user
description: 'user ID'
groups:
- name: account
  category: batch
  description: 'user account info'
  features:
  - name: state
    value-type: string
  - name: credit_score
    value-type: int64
  - name: account_age_days
    value-type: int64
  - name: has_2fa_installed
    value-type: bool
- name: transaction_stats
  category: batch
  description: 'user transaction statistics'
  features:
  - name: transaction_count_7d
    value-type: int64
  - name: transaction_count_30d
    value-type: int64
```

4. Import CSV data to Offline Store.

```bash
oomcli import \
  --group account \
  --input-file account.csv \
  --description 'sample account data'
oomcli import \
  --group transaction_stats \
  --input-file transaction_stats.csv \
  --description 'sample transaction stat data'
```

5. Sync data from Offline Store to Online Store.

```bash
oomcli sync --revision-id 1
oomcli sync --revision-id 2
```

6. Fetch features by key.

```bash
oomcli get online \
  --entity-keys 1006 \
  --feature account.state,account.credit_score,account.account_age_days,account.has_2fa_installed,transaction_stats.transaction_count_7d,transaction_stats.transaction_count_30d
```

```text
+------+---------------+----------------------+--------------------------+---------------------------+----------------------------------------+-----------------------------------------+
| user | account.state | account.credit_score | account.account_age_days | account.has_2fa_installed | transaction_stats.transaction_count_7d | transaction_stats.transaction_count_30d |
+------+---------------+----------------------+--------------------------+---------------------------+----------------------------------------+-----------------------------------------+
| 1006 | Louisiana     |                  710 |                       32 | false                     |                                      8 |                                      22 |
+------+---------------+----------------------+--------------------------+---------------------------+----------------------------------------+-----------------------------------------+
```

7. Generate training datasets via point-in-time join.

```sh
oomcli join \
	--feature account.state,account.credit_score,account.account_age_days,account.has_2fa_installed,transaction_stats.transaction_count_7d,transaction_stats.transaction_count_30d \
	--input-file label.csv
```

```text
+------------+------------+---------------+----------------------+--------------------------+---------------------------+----------------------------------------+-----------------------------------------+
| entity_key | unix_milli | account.state | account.credit_score | account.account_age_days | account.has_2fa_installed | transaction_stats.transaction_count_7d | transaction_stats.transaction_count_30d |
+------------+------------+---------------+----------------------+--------------------------+---------------------------+----------------------------------------+-----------------------------------------+
|       1002 | 1950236233 | Hawaii        |                  625 |                      861 | true                      |                                     11 |                                      36 |
|       1003 | 1950411318 | Arkansas      |                  730 |                      958 | false                     |                                      0 |                                      16 |
|       1004 | 1950653614 | Louisiana     |                  610 |                     1570 | false                     |                                     12 |                                      26 |
|       1005 | 1950166137 | South Dakota  |                  635 |                     1953 | false                     |                                      7 |                                      30 |
|       1006 | 1950403162 | Louisiana     |                  710 |                       32 | false                     |                                      8 |                                      22 |
|       1007 | 1950160030 | New Mexico    |                  645 |                       37 | true                      |                                      5 |                                      40 |
|       1008 | 1950274859 | Nevada        |                  735 |                     1627 | false                     |                                     12 |                                      51 |
|       1009 | 1949958846 | Kentucky      |                  650 |                       88 | true                      |                                     11 |                                      23 |
|       1010 | 1949920686 | Delaware      |                  680 |                     1687 | false                     |                                      2 |                                      39 |
+------------+------------+---------------+----------------------+--------------------------+---------------------------+----------------------------------------+-----------------------------------------+
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

Feel free to [join the community](https://oom.ai/slack) for questions and requests!
