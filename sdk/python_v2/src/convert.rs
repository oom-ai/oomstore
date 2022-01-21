use std::collections::HashMap;

use oomclient::Value;
use pyo3::{exceptions::PyException, prelude::*};

pub fn value_to_py(value: Option<&Value>, py: Python) -> PyObject {
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
        .map(|(k, v)| (k, value_to_py(v.as_ref(), py)))
        .collect::<HashMap<_, _>>()
        .into_py(py)
}

pub fn err_to_py(err: impl std::error::Error) -> PyErr {
    PyException::new_err(format!("{:?}", err))
}

pub fn py_to_value(obj: &PyAny) -> PyResult<Value> {
    obj.extract::<ValueWrapper>().map(|wrapper| match wrapper {
        ValueWrapper::Int64(v) => Value::Int64(v),
        ValueWrapper::Double(v) => Value::Double(v),
        ValueWrapper::String(v) => Value::String(v),
        ValueWrapper::Bool(v) => Value::Bool(v),
        ValueWrapper::UnixMilli(v) => Value::UnixMilli(v),
        ValueWrapper::Bytes(v) => Value::Bytes(v),
    })
}

#[derive(FromPyObject, Debug)]
enum ValueWrapper {
    Int64(i64),
    Double(f64),
    String(String),
    Bool(bool),
    UnixMilli(i64),
    Bytes(Vec<u8>),
}
