use std::{io, net::AddrParseError, num};
use thiserror::Error;

/// An error originating from the oomclient or it's dependencies.
#[derive(Error, Debug)]
pub enum OomError {
    #[error("a grpc status describing the result of an rpc call")]
    TonicStatus(#[from] tonic::Status),

    #[error("error's that originate from the client or server")]
    TonicTransportError(#[from] tonic::transport::Error),

    #[error("a checked integral type conversion fails.")]
    IntConversionError(#[from] num::TryFromIntError),

    #[error(transparent)]
    IoError(#[from] io::Error),

    #[error(transparent)]
    AddrParseError(#[from] AddrParseError),

    #[error("unknown error")]
    Unknown(String),
}
