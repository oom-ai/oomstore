use oomclient::error::OomError;
use std::collections::HashMap;

use oomclient::Value;
use pyo3::{exceptions::PyException, prelude::*};

pub fn value_to_py(value: Option<Value>, py: Python) -> PyObject {
    value
        .map(|value| match value {
            Value::Int64(v) => v.to_object(py),
            Value::Double(v) => v.to_object(py),
            Value::String(v) => v.to_object(py),
            Value::Bool(v) => v.to_object(py),
            Value::UnixMilli(v) => v.to_object(py),
            Value::Bytes(v) => v.to_object(py),
        })
        .to_object(py)
}

pub fn value_map_to_py(m: HashMap<String, Option<Value>>, py: Python) -> PyObject {
    m.into_iter()
        .map(|(k, v)| (k, value_to_py(v, py)))
        .collect::<HashMap<_, _>>()
        .into_py(py)
}

pub fn oom_err_to_py(err: OomError) -> PyErr {
    PyException::new_err(format!("{:?}", err))
}
