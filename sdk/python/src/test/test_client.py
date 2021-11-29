#!/usr/bin/env python3
#
import os
import sys
import tempfile
import shutil
from pathlib import Path

SOURCE_PATH = str(Path(os.getcwd()).parent.absolute())
sys.path.insert(0, SOURCE_PATH)

from oomstore import client


class Data:
    def __init__(self):
        self.temp_dir = tempfile.mkdtemp()

    def __del__(self):
        shutil.rmtree(self.temp_dir)

    # the integration test metadata
    def meta_file(self):
        with open(os.path.join(self.temp_dir, "metadata.yml"), "w") as f:
            f.write(
                """kind: Entity
name: user
length: 8
description: "user ID"
batch-features:
  - group: account
    description: "user account info"
    features:
      - name: state
        db-value-type: varchar(32)
      - name: credit_score
        db-value-type: int
      - name: account_age_days
        db-value-type: int
      - name: has_2fa_installed
        db-value-type: bool
  - group: transaction_stats
    description: "user transaction statistics"
    features:
      - name: transaction_count_7d
        db-value-type: int
      - name: transaction_count_30d
        db-value-type: int"""
            )
            return f.name

    def config_file(self):
        with open(os.path.join(self.temp_dir, "config.yml"), "w") as f:
            f.write(
                """oomstore: &onestore
  backend: postgres
  postgres:
    host: 127.0.0.1
    port: 5432
    user: postgres
    password: postgres
    database: oomstore
online-store: *onestore
offline-store: *onestore
metadata-store: *onestore"""
            )
            return f.name

    def account_file(self):
        with open(os.path.join(self.temp_dir, "account.csv"), "w") as f:
            f.write(
                """user,state,account_age_days,credit_score,has_2fa_installed
1001,Arizona,1547,685,false
1002,Hawaii,861,625,true
1003,Arkansas,958,730,false
1004,Louisiana,1570,610,false
1005,South Dakota,1953,635,false
1006,Louisiana,32,710,false
1007,New Mexico,37,645,true
1008,Nevada,1627,735,false
1009,Kentucky,88,650,true
1010,Delaware,1687,680,false"""
            )
            return f.name

    def transaction_stats_file(self):
        with open(os.path.join(self.temp_dir, "transaction_stats.csv"), "w") as f:
            f.write(
                """user,transaction_count_7d,transaction_count_30d
1001,9,41
1002,11,36
1003,0,16
1004,12,26
1005,7,30
1006,8,22
1007,5,40
1008,12,51
1009,11,23
1010,2,39"""
            )
            return f.name

    def label_file(self):
        with open(os.path.join(self.temp_dir, "label.csv"), "w") as f:
            f.write(
                """1001,1950049136
1002,1950236233
1003,1950411318
1004,1950653614
1005,1950166137
1006,1950403162
1007,1950160030
1008,1950274859
1009,1949958846
1010,1949920686"""
            )
            return f.name


if __name__ == "__main__":
    data = Data()
    c = client.Client(5001, data.config_file())
    revision_id1 = c.import_(
        group_name="account",
        description="sample account data",
        input_file_path=data.account_file(),
        delimiter=",",
    )
    revision_id2 = c.import_(
        group_name="transaction_stats",
        description="sample transaction stat data",
        input_file_path=data.transaction_stats_file(),
        delimiter=",",
    )

    c.sync(revision_id1, 2)
    c.sync(revision_id2, 2)

    print(
        c.online_get(
            entity_key="1006",
            feature_names=[
                "state",
                "credit_score",
                "account_age_days",
                "has_2fa_installed",
                "transaction_count_7d",
                "transaction_count_30d",
            ],
        )
    )
    print(
        c.online_multi_get(
            entity_keys=["1006", "1007"],
            feature_names=[
                "state",
                "credit_score",
                "account_age_days",
                "has_2fa_installed",
                "transaction_count_7d",
                "transaction_count_30d",
            ],
        )
    )
    c.join(
        feature_names=[
            "state",
            "credit_score",
            "account_age_days",
            "has_2fa_installed",
            "transaction_count_7d",
            "transaction_count_30d",
        ],
        input_file_path=data.label_file(),
        output_file_path="/tmp/joined.csv",
    )
