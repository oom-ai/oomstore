//! This crate provides an easy-to-use python client for [oomstore](https://github.com/oom-ai/oomstore), a
//! lightweight and fast feature store powered by go.

mod convert;
mod error;

use convert::{err_to_py, py_to_value, value_map_to_py};
use error::Error;
use oomclient::Client as OomClient;
use pyo3::{
    prelude::*,
    types::{PyDict, PyType},
};
use pyo3_asyncio::tokio::future_into_py;
use std::collections::HashMap;

/// A python client for [oomstore](https://github.com/oom-ai/oomstore),
/// using the grpc protocol to communicate with oomagent under the hood.
#[pyclass]
pub struct Client {
    inner: OomClient,
}

#[pymethods]
impl Client {
    /// Connect to an oomagent instance running on the given endpoint.
    ///
    /// Args:
    ///     endpoint (str): The endpoint to connect to.
    ///
    /// Returns:
    ///     The client instance.
    #[classmethod]
    #[pyo3(text_signature = "(endpoint)")]
    pub fn connect<'p>(_cls: &PyType, py: Python<'p>, endpoint: String) -> PyResult<&'p PyAny> {
        future_into_py(py, async {
            let inner = OomClient::connect(endpoint).await.map_err(err_to_py)?;
            let client = Client { inner };
            Python::with_gil(|py| PyCell::new(py, client).map(|py_cell| py_cell.to_object(py)))
        })
    }

    /// Connect to an oomagent instance embedded with the client.
    ///
    /// Args:
    ///     bin_path (str): The path to the oomagent binary.
    ///     cfg_path (str): The path to the oomstore configuration file.
    ///
    /// Returns:
    ///    The client instance.
    #[classmethod]
    #[pyo3(text_signature = "(bin_path, cfg_path)")]
    pub fn with_embedded_oomagent<'p>(
        _cls: &PyType,
        py: Python<'p>,
        bin_path: Option<String>,
        cfg_path: Option<String>,
    ) -> PyResult<&'p PyAny> {
        future_into_py(py, async {
            let inner = OomClient::with_embedded_oomagent(bin_path, cfg_path)
                .await
                .map_err(err_to_py)?;
            let client = Client { inner };
            Python::with_gil(|py| PyCell::new(py, client).map(|py_cell| py_cell.to_object(py)))
        })
    }

    /// Check if oomagent is ready to serve requests.
    #[pyo3(text_signature = "()")]
    pub fn health_check<'p>(&self, py: Python<'p>) -> PyResult<&'p PyAny> {
        // Don't panic, it's cheap:
        // https://github.com/hyperium/tonic/issues/285#issuecomment-595880400
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move { inner.health_check().await.map_err(err_to_py) })
    }

    /// Get online features for an entity.
    ///
    /// Args:
    ///     entity_key (str): An entity identifier, could be device ID, user ID, etc.
    ///     features: A list of feature full names.
    ///       A feature full name has the following format: &lt;group_name&gt;.&lt;feature_name&gt;,
    ///       for example, txn_stats.count_7d.
    ///
    /// Returns:
    ///     A dict mapping feature full names to feature values.
    #[pyo3(text_signature = "(entity_key, features)")]
    pub fn online_get<'p>(&self, py: Python<'p>, entity_key: String, features: Vec<String>) -> PyResult<&'p PyAny> {
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move {
            inner
                .online_get(entity_key, features)
                .await
                .map_err(err_to_py)
                .map(|m| Python::with_gil(|py| value_map_to_py(m, py)))
        })
    }

    /// Get online features for multiple entities.
    ///
    /// Args:
    ///     entity_keys: A list of entity identifiers, could be device IDs, user IDs, etc.
    ///     features: A list of feature full names.
    ///       A feature full name has the following format: &lt;group_name&gt;.&lt;feature_name&gt;,
    ///       for example, txn_stats.count_7d.
    ///
    /// Returns:
    ///     A dict mapping entity keys to a dict, that maps feature full names to feature values.
    #[pyo3(text_signature = "(entity_keys, features)")]
    pub fn online_multi_get<'p>(
        &self,
        py: Python<'p>,
        entity_keys: Vec<String>,
        features: Vec<String>,
    ) -> PyResult<&'p PyAny> {
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move {
            inner
                .online_multi_get(entity_keys, features)
                .await
                .map_err(err_to_py)
                .map(|r| {
                    Python::with_gil(|py| {
                        r.into_iter()
                            .map(|(k, m)| (k, value_map_to_py(m, py)))
                            .collect::<HashMap<_, _>>()
                    })
                })
        })
    }

    /// Sync a certain revision of batch features from offline to online store.
    ///
    /// Args:
    ///     group: The group to sync from offline store to online store.
    ///     revision_id: The revision to sync, it only applies to batch feature.
    ///       For batch feature: if null, will sync the latest revision;
    ///       otherwise, sync the designated revision.
    ///       For streaming feature, revision ID is not required, will always
    ///       sync the latest values.
    ///     purge_delay: PurgeDelay represents the seconds to sleep before purging
    ///       the previous revision in online store.
    ///       It only applies to batch feature group.
    ///
    /// Returns:
    ///     None
    #[pyo3(text_signature = "(group, revision_id, purge_delay)")]
    pub fn sync<'p>(
        &mut self,
        py: Python<'p>,
        group: String,
        revision_id: Option<u32>,
        purge_delay: u32,
    ) -> PyResult<&'p PyAny> {
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move {
            inner.sync(group, revision_id, purge_delay).await.map_err(err_to_py)
        })
    }

    /// Import features from external (batch and stream) data sources to offline store through files.
    ///
    /// Args:
    ///     group: The group to be imported from data source to offline store.
    ///     revision: The revision of the imported data, it only applies to
    ///       batch feature (not required).
    ///       For batch features, if revision is null, will use the
    ///       timestamp when it starts serving in online store; otherwise,
    ///       use the designated revision.
    ///     description: Description of this import.
    ///     input_file: The path of data source.
    ///     delimiter: Delimiter of data source.
    ///
    /// Returns:
    ///     revision_id: The revision ID of this import, it only applies to batch feature.
    #[pyo3(name = "import_")]
    #[pyo3(text_signature = "(group, revision, description, input_file, delimiter)")]
    pub fn import<'p>(
        &mut self,
        py: Python<'p>,
        group: String,
        revision: Option<i64>,
        description: Option<String>,
        input_file: String,
        delimiter: Option<String>,
    ) -> PyResult<&'p PyAny> {
        let delimiter = delimiter
            .map(|s| {
                let mut chars = s.chars();
                match (chars.next(), chars.next()) {
                    (Some(c), None) => Ok(c),
                    _ => err!("delimiter must be exactly one character"),
                }
            })
            .transpose()
            .map_err(err_to_py)?;
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move {
            inner
                .import(group, revision, description, input_file, delimiter)
                .await
                .map_err(err_to_py)
        })
    }

    /// Push stream features from stream data source to both offline and online stores.
    ///
    /// Args:
    ///     group: The group to be pushed from data source to offline store.
    ///     entity_key: An entity identifier.
    ///     kv_pairs: Feature values maps feature name to feature value.
    ///
    /// Returns:
    ///     None
    #[pyo3(text_signature = "(group, entity_key, kv_pairs)")]
    pub fn push<'p>(
        &mut self,
        py: Python<'p>,
        group: String,
        entity_key: String,
        kv_pairs: &PyDict,
    ) -> PyResult<&'p PyAny> {
        let kvs = kv_pairs
            .into_iter()
            .map(|(k, v)| Ok((k.extract::<String>()?, py_to_value(v)?)))
            .collect::<PyResult<_>>()?;
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move {
            inner.push(entity_key, group, kvs).await.map_err(err_to_py)
        })
    }

    /// Point-in-Time Join features against labeled entity rows through files.
    ///
    /// Args:
    ///     features: A list of feature full names, their feature values will be
    ///       joined and fetched from offline store.
    ///     input_file: File path of entity rows.
    ///     output_file: File path of joined result.
    ///
    /// Returns:
    ///     None
    #[pyo3(text_signature = "(features, input_file, output_file)")]
    pub fn join<'p>(
        &mut self,
        py: Python<'p>,
        features: Vec<String>,
        input_file: String,
        output_file: String,
    ) -> PyResult<&'p PyAny> {
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move {
            inner.join(features, input_file, output_file).await.map_err(err_to_py)
        })
    }

    /// Export certain features to a file.
    ///
    /// Args:
    ///     features: A list of feature full names.
    ///     unix_milli: A unix milliseconds, export the feature value before this timestamp.
    ///     output_file: File path of export result.
    ///     limit: Limit the size of export data.
    ///
    /// Returns:
    ///     None
    #[pyo3(text_signature = "(features, unix_milli, output_file, limit)")]
    pub fn export<'p>(
        &mut self,
        py: Python<'p>,
        features: Vec<String>,
        unix_milli: u64,
        output_file: String,
        limit: Option<usize>,
    ) -> PyResult<&'p PyAny> {
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move {
            inner
                .export(features, unix_milli, output_file, limit)
                .await
                .map_err(err_to_py)
        })
    }

    /// Take snapshot for a stream feature group in offline store.
    ///
    /// Args:
    ///     group: A streaming feature group.
    ///
    /// Returns:
    ///     None
    #[pyo3(text_signature = "(group)")]
    pub fn snapshot<'p>(&mut self, py: Python<'p>, group: String) -> PyResult<&'p PyAny> {
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move { inner.snapshot(group).await.map_err(err_to_py) })
    }
}

/// OomClient python module implemented in Rust.
#[pymodule]
fn oomclient(_py: Python, m: &PyModule) -> PyResult<()> {
    unsafe {
        pyo3::ffi::PyEval_InitThreads();
    }
    m.add_class::<Client>()?;
    Ok(())
}
