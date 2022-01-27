use oomclient::EmbeddedAgent;
use std::time::Duration;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let agent = EmbeddedAgent::default().await?;

    println!("{:?}", agent.address());

    tokio::time::sleep(Duration::from_secs(10)).await;

    Ok(())
}
