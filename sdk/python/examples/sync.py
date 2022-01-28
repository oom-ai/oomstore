#!/usr/bin/env python3
import asyncio
from oomclient import Client

async def main():
    client = await Client.connect("http://localhost:50051")
    await client.sync("account", 1, 0)

asyncio.run(main())
