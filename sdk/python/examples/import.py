#!/usr/bin/env python3
import asyncio
from oomclient import Client

async def main():
    client = await Client.connect("http://localhost:50051")
    group = "account"
    path = "/tmp/demo.csv"
    with open(path, "w") as f:
        rows = """
user,state,credit_score,account_age_days,has_2fa_installed
1,Nevada,530,242,true
2,South Carolina,520,268,false
3,New Jersey,655,84,false
4,Ohio,677,119,true
5,California,566,289,false
6,North Carolina,533,155,true
7,North Dakota,605,334,true
8,West Virginia,664,282,false
9,Alabama,577,150,true
10,Idaho,693,212,true
""".lstrip()
        f.write(rows)
    await client.import_(group, None, None, path, None)

asyncio.run(main())
