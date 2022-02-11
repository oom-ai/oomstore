## OomStore Client (Python)

This library provides an easy-to-use async python client for [OomStore](https://github.com/oom-ai/oomstore), a
lightweight and fast feature store powered by go.

It's built on top of [oomstore client in rust](https://github.com/oom-ai/oomstore/tree/main/sdk/rust) via [pyo3](https://github.com/PyO3/pyo3), the rust bindings for python.

### Install

This package requires Python 3.7+. MacOS and Linux are supported.

```sh
pip3 install oomclient
```

### Example

```python
import asyncio
from oomclient import Client

async def main():
    client = await Client.with_embedded_oomagent()
    features = ["account.state", "transaction_stats.transaction_count_7d"]
    result = await client.online_get("48", features)
    print(result)

asyncio.run(main())
```

More examples can be found in the [examples](https://github.com/oom-ai/oomstore/tree/main/sdk/python/examples) directory.

### Development

**Install maturin**

```sh
pip3 install maturin
```

**Init venv**

```sh
python -m venv .env
source .env/bin/activate
```

**Build and install**

```sh
maturin develop
```

Then you can enter into the python interpreter or run the example scripts to test the library.


There are also some [cargo-make](https://github.com/sagiegurari/cargo-make) tasks defined in `Makefile.toml`.
They can be executed by `cargo make <job name>`.

### License

Apache-2.0
