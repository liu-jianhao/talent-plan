# Rust Project 1: The Rust toolbox

任务：实现一个内存键值存储程序

由于刚开始学习Rust，感觉这门语言相对于其他语言(比如之前学习的Go语言)来说学习曲线是比较大的，
写一个能编译成功的程序都非常不容易，这也是Rust的特性，把很多错误提到编译期来消除。

其实这个小项目只是用到了clap库(一个简单实现命令行程序的库)和基于collections库中的HashMap简单包了一层接口实现键值存储

先看main函数，就是使用clap库：
```rust
use clap::{Arg, App, SubCommand};
use std::process;

// use kvs::KvStore;

fn main() {
    let matches = App::new("kvs")
        .version("0.1.0")
        .subcommand(
            SubCommand::with_name("set")
                .about("Set the value of a string key to a string")
                .arg(Arg::with_name("KEY").help("A string key").required(true))
                .arg(Arg::with_name("VALUE").help("A string value of the key").required(true)),
        )
        .subcommand(
            SubCommand::with_name("get")
                .about("Get the value of a given key")
                .arg(Arg::with_name("KEY").help("A string key").required(true)),
        )
        .subcommand(
            SubCommand::with_name("rm")
                .about("Remove a given key")
                .arg(Arg::with_name("KEY").help("A string key").required(true)),
        )
        .get_matches();

    match matches.subcommand() {
        ("set", Some(_sub_c)) => {
            // let key = _sub_c.value_of("KEY").unwrap();
            // let value = _sub_c.value_of("VALUE").unwrap();
            // let mut map = KvStore::new();
            // map.set(key.to_string(), value.to_string());
            eprintln!("unimplemented");
            process::exit(1);
        }
        ("get", Some(_sub_c)) => {
            // let key = _sub_c.value_of("KEY").unwrap();
            // let map = KvStore::new();
            // println!("{}", map.get(key.to_string()).unwrap());
            eprintln!("unimplemented");
            process::exit(1);
        }
        ("rm", Some(_sub_c)) => {
            // let key = _sub_c.value_of("KEY").unwrap();
            // let mut map = KvStore::new();
            // map.remove(key.to_string());
            eprintln!("unimplemented");
            process::exit(1);
        }
        _ => unreachable!(),
    }
}
```

再看kvs::KvStore的实现：
```rust
use std::collections::HashMap;

/// The `KvStore` stores string key/value pairs.
///
/// Key/value pairs are stored in a `HashMap` in memory and not persisted to disk.
///
/// Example:
///
/// ```rust
/// # use kvs::KvStore;
/// let mut store = KvStore::new();
/// store.set("key".to_owned(), "value".to_owned());
/// let val = store.get("key".to_owned());
/// assert_eq!(val, Some("value".to_owned()));
/// ```
pub struct KvStore {
    map: HashMap<String, String>,
}

impl KvStore {
    /// Creates a `KvStore`.
    pub fn new() -> KvStore {
        KvStore {
            map: HashMap::new(),
        }
    }

    /// Sets the value of a string key to a string.
    /// 
    /// If the key already exists, the previous value will be overwritten.
    pub fn set(&mut self, key: String, val: String) {
        if let Some(old) = self.map.insert(key, val) {
            eprintln!("Upload old value: {}", old);
        }
    }

    /// Gets the string value of a given string key.
    ///
    /// Returns `None` if the given key does not exist.
    pub fn get(&self, key: String) -> Option<String> {
        self.map.get(&key).cloned()
    }

    /// Remove a given key.
    pub fn remove(&mut self, key: String) {
        if let Some(old) = self.map.remove(&key) {
            eprintln!("Remove old value: {}", old);
        }
    }
}
```