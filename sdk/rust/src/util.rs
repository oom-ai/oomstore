use crate::oomagent::{value, FeatureValueMap, Value};
use std::collections::HashMap;

impl From<Value> for value::Value {
    fn from(value: Value) -> Self {
        value.value.expect("`oneof` in protobuffer should not be none")
    }
}

pub fn parse_raw_feature_values(input: FeatureValueMap) -> HashMap<String, value::Value> {
    input.map.into_iter().map(|(k, v)| (k, v.into())).collect()
}

pub fn parse_raw_values(values: Vec<Value>) -> Vec<value::Value> {
    values.into_iter().map(|v| v.into()).collect()
}
