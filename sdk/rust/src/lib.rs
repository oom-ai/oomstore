pub mod error;

use std::collections::HashMap;

use error::OomError;
use google::protobuf::Empty;
use oomagent::{oom_agent_client::OomAgentClient, value, OnlineGetRequest};
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
        Ok(self.inner.health_check(Empty {}).await.map(|_| ())?)
    }

    pub async fn online_get_raw(
        &mut self,
        key: impl Into<String>,
        features: Vec<String>,
    ) -> Result<HashMap<String, oomagent::Value>> {
        let res = self
            .inner
            .online_get(OnlineGetRequest { entity_key: key.into(), feature_full_names: features })
            .await?
            .into_inner();
        Ok(match res.result {
            Some(res) => res.map,
            None => HashMap::default(),
        })
    }

    pub async fn online_get(
        &mut self,
        key: impl Into<String>,
        features: Vec<String>,
    ) -> Result<HashMap<String, value::Kind>> {
        let rs = self.online_get_raw(key, features).await?;
        Ok(rs
            .into_iter()
            .map(|(k, v)| (k, v.kind.expect("`oneof` should not be none")))
            .collect())
    }
}
