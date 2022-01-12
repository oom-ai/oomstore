use oomstore::Client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://127.0.0.1:50051").await?;

    let response = client.health_check().await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
