use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("{0}")]
    Custom(String),
}

#[macro_export]
macro_rules! err {
    ($($tt:tt)*) => { Err(Error::Custom(format!($($tt)*).into())) };
}
