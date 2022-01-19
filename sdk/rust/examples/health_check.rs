use oomrpc::Client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://localhost:50051").await?;

    let response = client.health_check().await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
