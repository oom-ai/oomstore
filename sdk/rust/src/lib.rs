//! This crate provides an easy-to-use client for [oomstore](https://github.com/oom-ai/oomstore), a
//! lightweight and fast feature store powered by go.
//!
//! ## Examples
//!
//! ```rust,no_run
//! use oomclient::{Client, OnlineGetFeatures::*};
//!
//! #[tokio::main]
//! async fn main() -> Result<(), Box<dyn std::error::Error>> {
//!     let mut client = Client::connect("http://localhost:50051").await?;
//!
//!     let features = FeatureNames(vec!["account.state".into(), "txn_stats.count_7d".into()]);
//!     let response = client.online_get_raw("48", features.clone()).await?;
//!     println!("RESPONSE={:#?}", response);
//!
//!     let response = client.online_get("48", features).await?;
//!     println!("RESPONSE={:#?}", response);
//!
//!     Ok(())
//! }
//! ```

mod client;
mod error;
mod server;
mod util;

mod oomagent {
    tonic::include_proto!("oomagent");
}

pub use client::{Client, OnlineGetFeatures};
pub use error::OomError;
pub use oomagent::EntityRow;
pub use server::ServerWrapper;

/// Represents a dynamically typed value which can be either
/// an int64, a double, a string, a bool, a unix milliseconds, or a
/// bytes array.
pub type Value = oomagent::value::Value;

/// A result holding an Error.
pub type Result<T> = std::result::Result<T, OomError>;
