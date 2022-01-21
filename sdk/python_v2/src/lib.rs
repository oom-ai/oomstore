mod convert;

use convert::{oom_err_to_py, value_map_to_py};
use oomclient::Client as OomClient;
use pyo3::{prelude::*, types::PyType};
use pyo3_asyncio::tokio::future_into_py;
use std::collections::HashMap;

#[pyclass]
pub struct Client {
    inner: OomClient,
}

#[pymethods]
impl Client {
    #[classmethod]
    pub fn connect<'p>(_cls: &PyType, py: Python<'p>, endpoint: String) -> PyResult<&'p PyAny> {
        future_into_py(py, async {
            let inner = OomClient::connect(endpoint).await.map_err(oom_err_to_py)?;
            let client = Client { inner };
            Python::with_gil(|py| PyCell::new(py, client).map(|py_cell| py_cell.to_object(py)))
        })
    }

    pub fn health_check<'p>(&self, py: Python<'p>) -> PyResult<&'p PyAny> {
        // Don't panic, it's cheap:
        // https://github.com/hyperium/tonic/issues/285#issuecomment-595880400
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move { inner.health_check().await.map_err(oom_err_to_py) })
    }

    pub fn online_get<'p>(&self, py: Python<'p>, entity_key: String, features: Vec<String>) -> PyResult<&'p PyAny> {
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move {
            inner
                .online_get(entity_key, features)
                .await
                .map_err(oom_err_to_py)
                .map(|m| Python::with_gil(|py| value_map_to_py(m, py)))
        })
    }

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
                .map_err(oom_err_to_py)
                .map(|r| {
                    Python::with_gil(|py| {
                        r.into_iter()
                            .map(|(k, m)| (k, value_map_to_py(m, py)))
                            .collect::<HashMap<_, _>>()
                    })
                })
        })
    }

    pub fn sync<'p>(&mut self, py: Python<'p>, revision_id: u32, purge_delay: u32) -> PyResult<&'p PyAny> {
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move {
            inner.sync(revision_id, purge_delay).await.map_err(oom_err_to_py)
        })
    }

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
