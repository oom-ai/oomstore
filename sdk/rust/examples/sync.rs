use oomclient::Client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://localhost:50051").await?;

    client.sync("account", 1, 0).await?;

    Ok(())
}
