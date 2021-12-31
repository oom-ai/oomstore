<h1 align="center">OomStore</h1>
<p align="center">
    <em>A Fast Feature Store Powered by Go.</em>
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

The name "OOM" is derived from **O**nline Store, **O**ffline Store, and **M**etadata Store that put together oomstore.
It allows you to define features with YAML,
store features in databases of choice,
and retrieve features in both online and offline use cases, fast.


Please see our [docs](https://oom.ai/docs) for more details.

## Architecture

<p align="center">
  <img src="https://oom.ai/images/architecture/architecture.svg" alt="Architecture">
</p>

You can interact with oomstore with CLI, Go API or Python API. See [Architecture](https://oom.ai/docs/architecture) for more details.

## Features

Compared to other feature store implementations, oomstore has its edges:

- Fast. Benchmark shows oomstore performs QPS > 50k and latency < 0.3 ms using Redis as the Online Store. For more details, see [benchmark](https://oom.ai/docs/benchmark).
- Pluggable. We support a wide range of databases already (see below), and there are more to come.
  - Online Store: DynamoDB, Redis, TiKV, Cassandra, TiDB, PostgreSQL, MySQL, SQLite.
  - Offline Store: Snowflake, Redshift, BigQuery, TiDB, PostgreSQL, MySQL, SQLite.
  - Metadata Store: TiDB, PostgreSQL, MySQL, SQLite.
- Simple. In the minimal, oomstore can run aganist a single MySQL/PostgreSQL database. This helps get started quickly, and you can always switch to a different database later without having to rewrite your code.

## Quickstart

1. Install `oomcli` following the [guide](https://oom.ai/docs/installation#cli).

2. `oomcli init` to initialize oomstore. Make sure there is a `~/.config/oomstore/config.yaml` as below.

```yaml
store: &pg
  backend: postgres
  postgres:
    host: 127.0.0.1
    port: 5432
    user: postgres
    password: postgres
    database: oomstore

online-store: *pg
offline-store: *pg
metadata-store: *pg
```

3. `oomcli apply -f config.yaml` to register metadata. See config.yaml below.

```yaml
kind: Entity
name: user
length: 8
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
  --entity-key 1006 \
  --feature state,credit_score,account_age_days,has_2fa_installed,transaction_count_7d,transaction_count_30d
```

```text
+------+-----------+--------------+------------------+-------------------+----------------------+-----------------------+
| user |   state   | credit_score | account_age_days | has_2fa_installed | transaction_count_7d | transaction_count_30d |
+------+-----------+--------------+------------------+-------------------+----------------------+-----------------------+
| 1006 | Louisiana |          710 |               32 | false             |                    8 |                    22 |
+------+-----------+--------------+------------------+-------------------+----------------------+-----------------------+
```

7. Generate training datasets via point-in-time join.

```sh
oomcli join \
  --feature state,credit_score,account_age_days,has_2fa_installed,transaction_count_7d,transaction_count_30d \
  --input-file label.csv
```

```text
+------------+------------+--------------+--------------+------------------+-------------------+----------------------+-----------------------+
| entity_key | unix_time  |    state     | credit_score | account_age_days | has_2fa_installed | transaction_count_7d | transaction_count_30d |
+------------+------------+--------------+--------------+------------------+-------------------+----------------------+-----------------------+
|       1001 | 1950049136 | Arizona      |          685 |             1547 | false             |                    9 |                    41 |
|       1002 | 1950236233 | Hawaii       |          625 |              861 | true              |                   11 |                    36 |
|       1003 | 1950411318 | Arkansas     |          730 |              958 | false             |                    0 |                    16 |
|       1004 | 1950653614 | Louisiana    |          610 |             1570 | false             |                   12 |                    26 |
|       1005 | 1950166137 | South Dakota |          635 |             1953 | false             |                    7 |                    30 |
|       1006 | 1950403162 | Louisiana    |          710 |               32 | false             |                    8 |                    22 |
|       1007 | 1950160030 | New Mexico   |          645 |               37 | true              |                    5 |                    40 |
|       1008 | 1950274859 | Nevada       |          735 |             1627 | false             |                   12 |                    51 |
|       1009 | 1949958846 | Kentucky     |          650 |               88 | true              |                   11 |                    23 |
|       1010 | 1949920686 | Delaware     |          680 |             1687 | false             |                    2 |                    39 |
+------------+------------+--------------+--------------+------------------+-------------------+----------------------+-----------------------+
```

See [Quickstart](https://oom.ai/docs/quickstart) for more complete details.

## Roadmap

We are looking to support stream features. See [Roadmap](https://oom.ai/docs/roadmap) for more details.

## Community

Feel free to join our [Slack Community](https://oom.ai/slack) for questions and suggestions!
