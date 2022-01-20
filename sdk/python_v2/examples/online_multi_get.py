#!/usr/bin/env python3
import asyncio
from oomclient import Client

async def main():
    client = await Client.connect("http://localhost:50051")
    features = ["account.state", "transaction_stats.transaction_count_7d"]
    entity_keys = ["48", "1", "5", "9", "27", "6", "41", "72"]
    result = await client.online_multi_get(entity_keys, features)
    print(result)

asyncio.run(main())
