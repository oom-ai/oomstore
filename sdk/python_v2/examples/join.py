#!/usr/bin/env python3
import asyncio
from oomclient import Client


async def main():
    client = await Client.connect("http://localhost:50051")
    features = [
        "driver_stats.conv_rate",
        "driver_stats.acc_rate",
        "driver_stats.avg_daily_trips",
    ]
    input_path = "/tmp/input.csv"
    output_path = "/tmp/output.csv"
    with open(input_path, "w") as f:
        rows = """
entity_key,unix_milli,age
1,0,20
1,3,44
1,4,37
2,3,28
3,3,48
4,4,22
5,2,41
5,2,43
5,3,47
5,4,59
5,5,55
6,1,27
6,3,46
7,0,42
7,4,40
7,4,49
8,4,39
9,3,35
10,4,54
10,4,33
""".lstrip()
        f.write(rows)
    await client.join(features, input_path, output_path)
    print(open(output_path).read())

asyncio.run(main())
