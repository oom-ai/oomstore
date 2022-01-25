use oomclient::EmbeddedAgent;
use std::time::Duration;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let agent = EmbeddedAgent::new(
        "/home/wenxuan/develop/oom.ai/oomstore/oomagent/build/oomagent",
        "/home/wenxuan/develop/oom.ai/oomstore/oomagent/test/config.yaml",
        None,
    )
    .await?;

    println!("{:?}", agent.address());

    tokio::time::sleep(Duration::from_secs(10)).await;

    Ok(())
}
