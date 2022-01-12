use google::protobuf::Empty;
use oomagent::oom_agent_client::OomAgentClient;

pub mod google {
    pub mod protobuf {
        tonic::include_proto!("google.protobuf");
    }
}

pub mod oomagent {
    tonic::include_proto!("oomagent");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = OomAgentClient::connect("http://127.0.0.1:50051").await?;

    let response = client.health_check(Empty {}).await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
