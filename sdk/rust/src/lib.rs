pub mod error;
mod util;

use std::collections::HashMap;

use error::OomError;
use google::protobuf::Empty;
use oomagent::{oom_agent_client::OomAgentClient, value, FeatureValueMap, OnlineGetRequest, OnlineMultiGetRequest};
use tonic::{codegen::StdError, transport};
use util::parse_raw_feature_values;

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

    pub async fn online_get_raw(&mut self, key: impl Into<String>, features: Vec<String>) -> Result<FeatureValueMap> {
        let res = self
            .inner
            .online_get(OnlineGetRequest { entity_key: key.into(), feature_full_names: features })
            .await?
            .into_inner();
        Ok(match res.result {
            Some(res) => res,
            None => FeatureValueMap::default(),
        })
    }

    pub async fn online_get(
        &mut self,
        key: impl Into<String>,
        features: Vec<String>,
    ) -> Result<HashMap<String, value::Kind>> {
        let rs = self.online_get_raw(key, features).await?;
        Ok(parse_raw_feature_values(rs))
    }

    pub async fn online_multi_get_raw(
        &mut self,
        keys: Vec<String>,
        features: Vec<String>,
    ) -> Result<HashMap<String, FeatureValueMap>> {
        let res = self
            .inner
            .online_multi_get(OnlineMultiGetRequest { entity_keys: keys, feature_full_names: features })
            .await?
            .into_inner();
        Ok(res.result)
    }

    pub async fn online_multi_get(
        &mut self,
        keys: Vec<String>,
        features: Vec<String>,
    ) -> Result<HashMap<String, HashMap<String, value::Kind>>> {
        let rs = self.online_multi_get_raw(keys, features).await?;
        Ok(rs.into_iter().map(|(k, v)| (k, parse_raw_feature_values(v))).collect())
    }
}
