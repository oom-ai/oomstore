mod client;
mod error;
mod server;
mod util;

mod oomagent {
    tonic::include_proto!("oomagent");
}

pub use client::Client;
pub use error::OomError;
pub use oomagent::{value::Value, EntityRow};
pub use server::ServerWrapper;

pub type Result<T> = std::result::Result<T, OomError>;
