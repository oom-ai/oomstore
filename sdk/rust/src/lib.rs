pub mod error;

use error::OomError;
use google::protobuf::Empty;
use oomagent::oom_agent_client::OomAgentClient;
use tonic::{codegen::StdError, transport};

type Result<T> = std::result::Result<T, OomError>;

pub mod google {
    pub mod protobuf {
        tonic::include_proto!("google.protobuf");
    }
}

pub mod oomagent {
    tonic::include_proto!("oomagent");
}

pub struct Client {
    inner: OomAgentClient<transport::Channel>,
}

impl Client {
    pub async fn connect<D>(dst: D) -> Result<Self>
    where
        D: std::convert::TryInto<tonic::transport::Endpoint>,
        D::Error: Into<StdError>,
    {
        Ok(Self { inner: OomAgentClient::connect(dst).await? })
    }

    pub async fn health_check(&mut self) -> Result<()> {
        match self.inner.health_check(Empty {}).await {
            Ok(_) => Ok(()),
            Err(e) => Err(e.into()),
        }
    }
}
