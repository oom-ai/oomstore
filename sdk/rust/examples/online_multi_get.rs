use oomstore::Client;

macro_rules! svec { ($($part:expr),* $(,)?) => {{ vec![$($part.to_string(),)*] }} }

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://127.0.0.1:50051").await?;

    let response = client
        .online_multi_get_raw(svec!["19", "48", "38"], svec![
            "account.state",
            "transaction_stats.transaction_count_7d",
        ])
        .await?;
    println!("RESPONSE={:#?}", response);

    let response = client
        .online_multi_get(svec!["19", "48", "38"], svec![
            "account.state",
            "transaction_stats.transaction_count_7d",
        ])
        .await?;
    println!("RESPONSE={:#?}", response);

    Ok(())
}
