#!/usr/bin/env python3
import asyncio
from oomclient import Client

async def main():
    client = await Client.connect("http://localhost:50051")
    await client.snapshot("user-click")

asyncio.run(main())
