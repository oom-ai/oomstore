use std::num;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum OomError {
    #[error("a grpc status describing the result of an rpc call")]
    TonicStatus(#[from] tonic::Status),

    #[error("error's that originate from the client or server")]
    TonicTransportError(#[from] tonic::transport::Error),

    #[error("a checked integral type conversion fails.")]
    IntConversionError(#[from] num::TryFromIntError),

    #[error("unknown error")]
    Unknown,
}
