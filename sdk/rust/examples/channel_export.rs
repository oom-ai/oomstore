use futures_util::{pin_mut, StreamExt};
use oomrpc::Client;
use std::time::{SystemTime, UNIX_EPOCH};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://127.0.0.1:50051").await?;

    let features = vec![
        "account.state".into(),
        "account.has_2fa_installed".into(),
        "transaction_stats.transaction_count_7d".into(),
    ];

    let now = SystemTime::now().duration_since(UNIX_EPOCH)?.as_millis().try_into()?;
    let (header, rows) = client.channel_export(features, now, 20).await?;

    println!("RESPONSE HEADER={:?}", header);

    pin_mut!(rows); // needed for iteration

    while let Some(row) = rows.next().await {
        println!("RESPONSE ROWS={:?}", row?);
    }

    Ok(())
}
