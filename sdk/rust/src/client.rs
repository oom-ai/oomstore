use crate::{
    error::OomError,
    oomagent,
    oomagent::{
        oom_agent_client::OomAgentClient,
        ChannelExportRequest,
        ChannelExportResponse,
        ChannelImportRequest,
        ChannelJoinRequest,
        ChannelJoinResponse,
        ExportRequest,
        FeatureValueMap,
        HealthCheckRequest,
        ImportRequest,
        JoinRequest,
        OnlineGetRequest,
        OnlineMultiGetRequest,
        PushRequest,
        SnapshotRequest,
        SyncRequest,
    },
    server::ServerWrapper,
    util::{parse_raw_feature_values, parse_raw_values},
    EntityRow,
    Result,
    Value,
};
use async_stream::stream;
use futures_core::stream::Stream;
use std::{collections::HashMap, path::Path, sync::Arc};
use tonic::{codegen::StdError, transport, Request};

#[derive(Debug, Clone)]
pub struct Client {
    client: OomAgentClient<transport::Channel>,
    _agent: Option<Arc<ServerWrapper>>,
}

// TODO: Add a Builder to create the client
impl Client {
    pub async fn connect<D>(dst: D) -> Result<Self>
    where
        D: std::convert::TryInto<tonic::transport::Endpoint>,
        D::Error: Into<StdError>,
    {
        Ok(Self { client: OomAgentClient::connect(dst).await?, _agent: None })
    }

    pub async fn with_embedded_oomagent<P1, P2>(bin_path: Option<P1>, cfg_path: Option<P2>) -> Result<Self>
    where
        P1: AsRef<Path>,
        P2: AsRef<Path>,
    {
        let agent = ServerWrapper::new(bin_path, cfg_path, None).await?;
        Ok(Self {
            client: OomAgentClient::connect(format!("http://{}", agent.address())).await?,
            _agent: Some(Arc::new(agent)),
        })
    }

    pub async fn with_default_embedded_oomagent() -> Result<Self> {
        Self::with_embedded_oomagent(None::<String>, None::<String>).await
    }

    pub async fn health_check(&mut self) -> Result<()> {
        Ok(self.client.health_check(HealthCheckRequest {}).await.map(|_| ())?)
    }

    pub async fn online_get_raw(
        &mut self,
        entity_key: impl Into<String>,
        features: Vec<String>,
    ) -> Result<FeatureValueMap> {
        let res = self
            .client
            .online_get(OnlineGetRequest { entity_key: entity_key.into(), features })
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
    ) -> Result<HashMap<String, Option<Value>>> {
        let rs = self.online_get_raw(key, features).await?;
        Ok(parse_raw_feature_values(rs))
    }

    pub async fn online_multi_get_raw(
        &mut self,
        entity_keys: Vec<String>,
        features: Vec<String>,
    ) -> Result<HashMap<String, FeatureValueMap>> {
        let res = self
            .client
            .online_multi_get(OnlineMultiGetRequest { entity_keys, features })
            .await?
            .into_inner();
        Ok(res.result)
    }

    pub async fn online_multi_get(
        &mut self,
        keys: Vec<String>,
        features: Vec<String>,
    ) -> Result<HashMap<String, HashMap<String, Option<Value>>>> {
        let rs = self.online_multi_get_raw(keys, features).await?;
        Ok(rs.into_iter().map(|(k, v)| (k, parse_raw_feature_values(v))).collect())
    }

    pub async fn sync(
        &mut self,
        group: impl Into<String>,
        revision_id: impl Into<Option<u32>>,
        purge_delay: u32,
    ) -> Result<()> {
        let group = group.into();
        let revision_id = revision_id.into().map(i32::try_from).transpose()?;
        let purge_delay = i32::try_from(purge_delay)?;
        self.client
            .sync(SyncRequest { revision_id, group, purge_delay })
            .await?;
        Ok(())
    }

    pub async fn channel_import(
        &mut self,
        group: impl Into<String>,
        revision: impl Into<Option<i64>>,
        description: impl Into<Option<String>>,
        rows: impl Stream<Item = Vec<u8>> + Send + 'static,
    ) -> Result<u32> {
        let mut group = Some(group.into());
        let mut description = description.into();
        let mut revision = revision.into();
        let inbound = stream! {
            for await row in rows {
                yield ChannelImportRequest{group: group.take(), description: description.take(), revision: revision.take(), row};
            }
        };
        let res = self.client.channel_import(Request::new(inbound)).await?.into_inner();
        Ok(res.revision_id as u32)
    }

    pub async fn import(
        &mut self,
        group: impl Into<String>,
        revision: impl Into<Option<i64>>,
        description: impl Into<Option<String>>,
        input_file: impl AsRef<Path>,
        delimiter: impl Into<Option<char>>,
    ) -> Result<u32> {
        let res = self
            .client
            .import(ImportRequest {
                group:       group.into(),
                description: description.into(),
                revision:    revision.into(),
                input_file:  input_file.as_ref().display().to_string(),
                delimiter:   delimiter.into().map(String::from),
            })
            .await?
            .into_inner();
        Ok(res.revision_id as u32)
    }

    pub async fn push(
        &mut self,
        entity_key: impl Into<String>,
        group: impl Into<String>,
        kv_pairs: HashMap<String, Value>,
    ) -> Result<()> {
        let kv_pairs = kv_pairs
            .into_iter()
            .map(|(k, v)| (k, oomagent::Value { value: Some(v) }))
            .collect();
        self.client
            .push(PushRequest {
                entity_key:     entity_key.into(),
                group:          group.into(),
                feature_values: kv_pairs,
            })
            .await?
            .into_inner();
        Ok(())
    }

    pub async fn channel_join(
        &mut self,
        join_features: Vec<String>,
        existed_features: Vec<String>,
        entity_rows: impl Stream<Item = EntityRow> + Send + 'static,
    ) -> Result<(Vec<String>, impl Stream<Item = Result<Vec<Option<Value>>>>)> {
        let mut join_features = Some(join_features);
        let mut existed_features = Some(existed_features);
        let inbound = stream! {
            for await row in entity_rows {
                let (join_features, existed_features) = match (join_features.take(), existed_features.take()) {
                    (Some(join_features), Some(existed_features)) => (join_features, existed_features),
                    _ => (Vec::new(), Vec::new()),
                };
                yield ChannelJoinRequest {
                    join_features,
                    existed_features,
                    entity_row: Some(row),
                };
            }
        };

        let mut outbound = self.client.channel_join(Request::new(inbound)).await?.into_inner();

        let ChannelJoinResponse { header, joined_row } = outbound
            .message()
            .await?
            .ok_or_else(|| OomError::Unknown(String::from("stream finished with no response")))?;

        let row = parse_raw_values(joined_row);

        let outbound = async_stream::try_stream! {
            yield row;
            while let Some(ChannelJoinResponse { joined_row, .. }) = outbound.message().await? {
                yield parse_raw_values(joined_row)
            }
        };
        Ok((header, outbound))
    }

    pub async fn join(
        &mut self,
        features: Vec<String>,
        input_file: impl AsRef<Path>,
        output_file: impl AsRef<Path>,
    ) -> Result<()> {
        self.client
            .join(JoinRequest {
                features,
                input_file: input_file.as_ref().display().to_string(),
                output_file: output_file.as_ref().display().to_string(),
            })
            .await?;
        Ok(())
    }

    pub async fn channel_export(
        &mut self,
        features: Vec<String>,
        unix_milli: u64,
        limit: impl Into<Option<usize>>,
    ) -> Result<(Vec<String>, impl Stream<Item = Result<Vec<Option<Value>>>>)> {
        let unix_milli = unix_milli.try_into()?;
        let limit = limit.into().map(|n| n.try_into()).transpose()?;
        let mut outbound = self
            .client
            .channel_export(ChannelExportRequest { features, unix_milli, limit })
            .await?
            .into_inner();

        let ChannelExportResponse { header, row } = outbound
            .message()
            .await?
            .ok_or_else(|| OomError::Unknown(String::from("stream finished with no response")))?;

        let row = parse_raw_values(row);
        let outbound = async_stream::try_stream! {
            yield row;
            while let Some(ChannelExportResponse{row, ..}) = outbound.message().await? {
                yield parse_raw_values(row)
            }
        };
        Ok((header, outbound))
    }

    pub async fn export(
        &mut self,
        features: Vec<String>,
        unix_milli: u64,
        output_file: impl AsRef<Path>,
        limit: impl Into<Option<usize>>,
    ) -> Result<()> {
        let unix_milli = unix_milli.try_into()?;
        let limit = limit.into().map(|n| n.try_into()).transpose()?;
        let output_file = output_file.as_ref().display().to_string();
        self.client
            .export(ExportRequest { features, unix_milli, output_file, limit })
            .await?;
        Ok(())
    }

    pub async fn snapshot(&mut self, group: impl Into<String>) -> Result<()> {
        self.client.snapshot(SnapshotRequest { group: group.into() }).await?;
        Ok(())
    }
}
