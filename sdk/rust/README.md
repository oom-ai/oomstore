## OomStore Client (Rust)

This crate provides an easy-to-use async rust client for [OomStore](https://github.com/oom-ai/oomstore), a
lightweight and fast feature store powered by go.
It uses [gRPC](https://grpc.io/) protocol to communicate with the oomagent server under the hood.

### Example

```rust
use oomclient::Client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::with_default_embedded_oomagent().await?;

    let features = vec!["account.state".into(), "txn_stats.count_7d".into()];
    let response = client.online_get_raw("48", features.clone()).await?;
    println!("RESPONSE={:#?}", response);

    let response = client.online_get("48", features).await?;
    println!("RESPONSE={:#?}", response);

    Ok(())
}
```

More examples can be found in `examples` directory of the project repo.

**Note**

You need to install the oomagent first following the [guide](https://www.oom.ai/docs/installation).

### Development

Install [cargo-make](https://github.com/sagiegurari/cargo-make) and run `cargo make`.

### License

Apache-2.0
