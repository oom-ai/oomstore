#!/usr/bin/env python3
import asyncio
from oomclient import Client

def gen_cfg(cfg_path):
    with open(cfg_path, "w") as f:
        f.write("""
online-store:
  sqlite:
    db-file: /tmp/oomplay.db

offline-store:
  sqlite:
    db-file: /tmp/oomplay.db

metadata-store:
  sqlite:
    db-file: /tmp/oomplay.db
""")

async def main():
    cfg_path = "/tmp/oomclient-py-oomagent.yaml"
    gen_cfg(cfg_path)
    client = await Client.with_embedded_oomagent(cfg_path=cfg_path)
    await client.health_check()


asyncio.run(main())
