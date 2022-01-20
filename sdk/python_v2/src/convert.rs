use oomclient::Value;
use pyo3::prelude::*;

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
