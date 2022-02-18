use oomclient::{Client, OnlineGetFeatures::*};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://localhost:50051").await?;

    let features = FeatureNames(vec!["account.state".into(), "txn_stats.count_7d".into()]);
    let response = client.online_get_raw("48", features.clone()).await?;
    println!("RESPONSE={:#?}", response);

    let response = client.online_get("48", features).await?;
    println!("RESPONSE={:#?}", response);

    Ok(())
}
