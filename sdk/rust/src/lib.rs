pub mod error;
mod util;
mod oomagent {
    tonic::include_proto!("oomagent");
}

use error::OomError;
use futures_core::stream::Stream;
use oomagent::{
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
};
use std::{collections::HashMap, path::Path};
use tonic::{codegen::StdError, transport, Request};
use util::{parse_raw_feature_values, parse_raw_values};

pub use oomagent::{value::Value, EntityRow};

type Result<T> = std::result::Result<T, OomError>;

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
        Ok(self.inner.health_check(HealthCheckRequest {}).await.map(|_| ())?)
    }

    pub async fn online_get_raw(
        &mut self,
        entity_key: impl Into<String>,
        features: Vec<String>,
    ) -> Result<FeatureValueMap> {
        let res = self
            .inner
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
            .inner
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

    pub async fn sync(&mut self, revision_id: u32, purge_delay: u32) -> Result<()> {
        self.inner
            .sync(SyncRequest {
                revision_id: i32::try_from(revision_id)?,
                purge_delay: i32::try_from(purge_delay)?,
            })
            .await?;
        Ok(())
    }

    pub async fn channel_import(
        &mut self,
        group: impl Into<Option<String>>,
        description: impl Into<Option<String>>,
        revision: impl Into<Option<i64>>,
        rows: impl Iterator<Item = Vec<u8>> + Send + 'static,
    ) -> Result<u32> {
        let mut group = group.into();
        let mut description = description.into();
        let mut revision = revision.into();
        let inbound = async_stream::stream! {
            for row in rows {
                yield ChannelImportRequest{group: group.take(), description: description.take(), revision: revision.take(), row};
            }
        };
        let res = self.inner.channel_import(Request::new(inbound)).await?.into_inner();
        Ok(res.revision_id as u32)
    }

    pub async fn import(
        &mut self,
        group: impl Into<String>,
        description: impl Into<Option<String>>,
        revision: impl Into<Option<i64>>,
        input_file: impl AsRef<Path>,
        delimiter: impl Into<Option<char>>,
    ) -> Result<u32> {
        let res = self
            .inner
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
        kv_pairs: Vec<(impl Into<String>, impl Into<Value>)>,
    ) -> Result<()> {
        let mut keys = Vec::with_capacity(kv_pairs.len());
        let mut vals = Vec::with_capacity(kv_pairs.len());
        kv_pairs.into_iter().for_each(|(k, v)| {
            keys.push(k.into());
            vals.push(oomagent::Value { value: Some(v.into()) });
        });
        self.inner
            .push(PushRequest {
                entity_key:     entity_key.into(),
                group:          group.into(),
                features:       keys,
                feature_values: vals,
            })
            .await?
            .into_inner();
        Ok(())
    }

    pub async fn channel_join(
        &mut self,
        join_features: Vec<String>,
        existed_features: Vec<String>,
        entity_rows: impl Iterator<Item = EntityRow> + Send + 'static,
    ) -> Result<(Vec<String>, impl Stream<Item = Result<Vec<Option<Value>>>>)> {
        let inbound = async_stream::stream! {
            for (i, row) in entity_rows.enumerate() {
                let (join_features, existed_features) = match i {
                    0 => (join_features.clone(), existed_features.clone()),
                    _ => (Vec::new(), Vec::new()),
                };
                yield ChannelJoinRequest {
                    join_features,
                    existed_features,
                    entity_row: Some(row),
                };
            }
        };

        let mut outbound = self.inner.channel_join(Request::new(inbound)).await?.into_inner();

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
        self.inner
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
            .inner
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
        self.inner
            .export(ExportRequest { features, unix_milli, output_file, limit })
            .await?;
        Ok(())
    }

    pub async fn snapshot(&mut self, group: impl Into<String>) -> Result<()> {
        self.inner.snapshot(SnapshotRequest { group: group.into() }).await?;
        Ok(())
    }
}
