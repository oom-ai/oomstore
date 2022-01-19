use oomclient::Client;
use std::time::{SystemTime, UNIX_EPOCH};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://localhost:50051").await?;

    let features = vec![
        "account.state".into(),
        "account.has_2fa_installed".into(),
        "transaction_stats.transaction_count_7d".into(),
    ];

    let now = SystemTime::now().duration_since(UNIX_EPOCH)?.as_millis().try_into()?;
    client.export(features, now, "/tmp/output.csv", 20).await?;

    std::fs::copy("/tmp/output.csv", "/dev/stdout")?;

    Ok(())
}
