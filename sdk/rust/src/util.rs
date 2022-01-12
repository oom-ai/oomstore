use crate::oomagent::{value, FeatureValueMap};
use std::collections::HashMap;

pub fn parse_raw_feature_values(input: FeatureValueMap) -> HashMap<String, value::Kind> {
    input
        .map
        .into_iter()
        .map(|(k, v)| (k, v.kind.expect("`oneof` should not be none")))
        .collect()
}
