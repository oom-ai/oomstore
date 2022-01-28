#!/usr/bin/env python3
import asyncio
import time
from oomclient import Client


async def main():
    client = await Client.connect("http://localhost:50051")
    features = [
        "account.state",
        "account.has_2fa_installed",
        "transaction_stats.transaction_count_7d",
    ]
    output_path = "/tmp/output.csv"
    now = round(time.time() * 1000)
    await client.export(features, now, output_path, 20)
    print(open(output_path).read())


asyncio.run(main())
