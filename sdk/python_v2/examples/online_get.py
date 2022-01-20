#!/usr/bin/env python3
import asyncio
from oomclient import Client

async def main():
    client = await Client.connect("http://localhost:50051")
    features = ["account.state", "transaction_stats.transaction_count_7d"]
    result = await client.online_get("48", features)
    print(result)

asyncio.run(main())
