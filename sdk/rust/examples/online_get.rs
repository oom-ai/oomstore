use oomrpc::Client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://127.0.0.1:50051").await?;

    let features = vec!["account.state".into(), "transaction_stats.transaction_count_7d".into()];
    let response = client.online_get_raw("48", features.clone()).await?;
    println!("RESPONSE={:#?}", response);

    let response = client.online_get("48", features).await?;
    println!("RESPONSE={:#?}", response);

    Ok(())
}
