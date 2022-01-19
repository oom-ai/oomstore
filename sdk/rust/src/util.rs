use crate::oomagent::{value, FeatureValueMap, Value};
use std::collections::HashMap;

pub fn parse_raw_feature_values(input: FeatureValueMap) -> HashMap<String, Option<value::Value>> {
    input.map.into_iter().map(|(k, v)| (k, v.value)).collect()
}

pub fn parse_raw_values(values: Vec<Value>) -> Vec<Option<value::Value>> {
    values.into_iter().map(|v| v.value).collect()
}
