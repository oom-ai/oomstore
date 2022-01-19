use oomrpc::Client;
use std::time::{SystemTime, UNIX_EPOCH};

macro_rules! svec { ($($part:expr),* $(,)?) => {{ vec![$($part.to_string(),)*] }} }

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://127.0.0.1:50051").await?;

    let features = svec![
        "account.state",
        "account.has_2fa_installed",
        "transaction_stats.transaction_count_7d"
    ];

    let now = SystemTime::now().duration_since(UNIX_EPOCH)?.as_millis().try_into()?;
    client.export(features, now, "/tmp/output.csv", 20).await?;

    std::fs::copy("/tmp/output.csv", "/dev/stdout")?;

    Ok(())
}
