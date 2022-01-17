use crate::oomagent::{value, FeatureValueMap, Value};
use std::collections::HashMap;

impl From<Value> for value::Kind {
    fn from(value: Value) -> Self {
        value.kind.expect("`oneof` in protobuffer should not be none")
    }
}

pub fn parse_raw_feature_values(input: FeatureValueMap) -> HashMap<String, value::Kind> {
    input.map.into_iter().map(|(k, v)| (k, v.into())).collect()
}

pub fn parse_raw_values(values: Vec<Value>) -> Vec<value::Kind> {
    values.into_iter().map(|v| v.into()).collect()
}
