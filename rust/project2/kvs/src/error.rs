use failure::Fail;
use std::io;

/// error type for kvs
#[derive(Debug, Fail)]
pub enum KvsError {
    /// IO error
    #[fail(display = "{}", _0)]
    Io(#[cause] io::Error),
    /// Serialization or Deserialization error
    #[fail(display = "{}", _0)]
    Serde(#[cause] serde_json::Error),
    /// Removing non-existent ey error
    #[fail(display = "Key not found")]
    KeyNotFound,
    /// Unexpected command type error
    #[fail(display = "Unexpected command type")]
    UnexpectedCommandType,
}

impl From<io::Error> for KvsError {
    fn from(err: io::Error) -> KvsError {
        KvsError::Io(err)
    }
}

impl From<serde_json::Error> for KvsError {
    fn from(err: serde_json::Error) -> KvsError {
        KvsError::Serde(err)
    }
}

/// usage
pub type Result<T> = std::result::Result<T, KvsError>;
