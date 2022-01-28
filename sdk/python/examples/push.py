#!/usr/bin/env python3
import asyncio
from oomclient import Client


async def main():
    client = await Client.connect("http://localhost:50051")
    kv_pairs = {
        "last_5_click_posts": "1,3,5,7,9",
        "number_of_user_starred_posts": 30,
    }
    await client.push("user-click", "929", kv_pairs)


asyncio.run(main())
