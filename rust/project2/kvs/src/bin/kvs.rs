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