use oomagent::{oom_agent_client::OomAgentClient, HealthCheckRequest};

pub mod google {
    pub mod protobuf {
        tonic::include_proto!("google.protobuf");
    }
    pub mod rpc {
        tonic::include_proto!("google.rpc");
    }
}

pub mod oomagent {
    tonic::include_proto!("oomagent");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = OomAgentClient::connect("http://127.0.0.1:50051").await?;

    let request = tonic::Request::new(HealthCheckRequest {});

    let response = client.health_check(request).await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
