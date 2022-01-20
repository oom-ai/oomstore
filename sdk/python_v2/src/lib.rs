use oomclient::Client as OomClient;
use pyo3::{exceptions::PyException, prelude::*, types::PyType};
use pyo3_asyncio::tokio::future_into_py;

#[pyclass]
pub struct Client {
    inner: OomClient,
}

#[pymethods]
impl Client {
    #[classmethod]
    pub fn connect<'p>(_cls: &PyType, py: Python<'p>, endpoint: String) -> PyResult<&'p PyAny> {
        future_into_py(py, async {
            let inner = OomClient::connect(endpoint).await.map_err(to_py_execption)?;
            let client = Client { inner };
            Python::with_gil(|py| PyCell::new(py, client).map(|py_cell| py_cell.to_object(py)))
        })
    }

    pub fn health_check<'p>(&self, py: Python<'p>) -> PyResult<&'p PyAny> {
        // Don't panic, it's cheap:
        // https://github.com/hyperium/tonic/issues/285#issuecomment-595880400
        let mut inner = OomClient::clone(&self.inner);
        future_into_py(py, async move { inner.health_check().await.map_err(to_py_execption) })
    }
}

pub fn to_py_execption(err: impl std::fmt::Display) -> PyErr {
    PyException::new_err(format!("{}", err))
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
