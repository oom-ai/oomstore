use oomstore::Client;

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
    let mut client = Client::connect("http://127.0.0.1:50051").await?;

    let response = client.online_get_raw("48", vec!["account.state".to_string()]).await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
