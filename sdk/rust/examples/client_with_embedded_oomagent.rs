use oomclient::Client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::with_default_embedded_oomagent().await?;

    let response = client.health_check().await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
