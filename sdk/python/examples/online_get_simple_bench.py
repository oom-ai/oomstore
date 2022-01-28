#!/usr/bin/env python3
import asyncio
import time
from oomclient import Client

async def main():
    client = await Client.connect("http://localhost:50051")
    features = ["account.state", "transaction_stats.transaction_count_7d"]
    start = time.time()
    n = 100
    for i in range(n):
        print(await client.online_get(str(i), features))
    spq = (time.time() - start) / n
    print(str(spq) + "s")

asyncio.run(main())
