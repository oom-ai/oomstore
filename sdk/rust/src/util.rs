use crate::oomagent::{value, FeatureValueMap, Value};
use std::collections::HashMap;

pub fn parse_raw_feature_values(input: FeatureValueMap) -> HashMap<String, value::Kind> {
    input
        .map
        .into_iter()
        .map(|(k, v)| (k, v.kind.expect("`oneof` should not be none")))
        .collect()
}

pub fn parse_raw_values(values: Vec<Value>) -> Vec<value::Kind> {
    values
        .into_iter()
        .map(|v| v.kind.expect("`oneof` should not be none"))
        .collect()
}
