/*
 MIT License

 Copyright (c) 2019 Max Kuznetsov <syhpoon@syhpoon.ca>

 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 SOFTWARE.
*/

#[macro_use]
extern crate nix;

mod device;
mod fs;
mod grpc;

use clap::{crate_version, value_t};
use clap::{App, AppSettings, Arg, SubCommand};
use failure::{format_err, Error};

use crate::{
    device::{BlockDevice, FileDevice, DeviceType, device_type}
};
use crate::fs::Filesystem;
use crate::grpc::Server;

fn main() -> Result<(), Error> {
    let app = App::new("learnfs")
        .version(crate_version!())
        .about("LearnFS command line tool")
        .arg(
            Arg::with_name("log")
                .help("Logging configuration file")
                .long("log")
                .takes_value(true)
                .global(true)
                .default_value(""),
        )
        .subcommand(
            SubCommand::with_name("create")
                .about("Create filesystem on a given device")
                .setting(AppSettings::DisableVersion)
                .arg(
                    Arg::with_name("device")
                        .help("File or block device to use")
                        .index(1)
                        .required(true),
                )
                .arg(
                    Arg::with_name("capacity")
                        .help("Explicitly set file capacity in bytes")
                        .short("c")
                        .long("capacity")
                        .takes_value(true)
                        .default_value("0"),
                ),
        )
        .subcommand(
            SubCommand::with_name("info")
                .about("Show existing filesystem info")
                .setting(AppSettings::DisableVersion)
                .arg(
                    Arg::with_name("device")
                        .help("Block device to use")
                        .index(1)
                        .required(true),
                ),
        )
        .subcommand(
            SubCommand::with_name("server")
                .about("Start learnfs server")
                .setting(AppSettings::DisableVersion)
                .arg(
                    Arg::with_name("listen")
                        .help("Grpc listen address")
                        .short("l")
                        .takes_value(true)
                        .default_value("0.0.0.0:4321"),
                ),
        )
        .subcommand(
            SubCommand::with_name("client")
                .about("learnfs client")
                .setting(AppSettings::DisableVersion)
                .subcommand(
                    App::new("rmdir")
                        .about("Remove directory")
                        .setting(AppSettings::DisableVersion)
                        .arg(
                            Arg::with_name("directory")
                                .help("Directory to remove")
                                .index(1)
                                .required(true),
                        )
                        .arg(
                            Arg::with_name("server")
                                .help("Server address")
                                .short("s")
                                .takes_value(true)
                                .default_value("http://127.0.0.1:4321"),
                        ),
                )
        );

    let help = extract_help(&app);
    let matches = app.get_matches();

    init_log(matches.value_of("log").unwrap())?;

    match matches.subcommand() {
        ("create", Some(create_m)) => {
            let path = create_m.value_of("device").unwrap();

            match device_type(path) {
                Ok(DeviceType::File) => {
                    let capacity = value_t!(create_m, "capacity", u64).unwrap();

                    if capacity == 0 {
                        return Err(format_err!(
                            "please explicitly specify file device capacity"
                        ));
                    }

                    let device = FileDevice::new(path, capacity)?;
                    Filesystem::create(device)
                        .expect("failed to create a filesystem");
                }
                Ok(DeviceType::Block) => {
                    let device = BlockDevice::new(path)?;
                    Filesystem::create(device)
                        .expect("failed to create a filesystem");
                }
                Err(e) => {
                    return Err(e);
                }
            }
        }
        ("info", Some(info_m)) => {
            let path = info_m.value_of("device").unwrap();

            match device_type(path) {
                Ok(DeviceType::File) => {
                    let device = FileDevice::open(path)?;
                    let fs = Filesystem::load(device)?;

                    println!("{:#?}", fs.info())
                }
                Ok(DeviceType::Block) => {
                    let device = BlockDevice::open(path)?;
                    let fs = Filesystem::load(device)?;

                    println!("{:#?}", fs.info())
                }
                Err(e) => {
                    return Err(e);
                }
            }
        }
        ("server", Some(server_m)) => {
            let listen = server_m.value_of("listen").unwrap();
            let runtime = tokio::runtime::Runtime::new().unwrap();

            runtime.block_on(Server::run(listen))?
        }
        ("client", Some(client_m)) => {
            match client_m.subcommand() {
                ("rmdir", Some(m)) => {
                    let dir = m.value_of("directory").unwrap().to_string();
                    let server = m.value_of("server").unwrap().to_string();

                    let runtime = tokio::runtime::Runtime::new()?;
                    runtime.block_on(async {
                        let mut client = grpc::Client::connect(server)
                            .await.unwrap();

                        let req = grpc::RmdirRequest {
                            pathname: dir.to_string(),
                        };

                        let resp = client.rmdir(req).await.unwrap();

                        println!("RESPONSE: {:?}", resp);
                    })
                }
                _ => {
                    eprintln!("{}", help);
                }
            }
        }
        _ => {
            eprintln!("{}", help);
        }
    }

    Ok(())
}

fn extract_help(app: &App) -> String {
    let mut vec: Vec<u8> = Vec::new();
    let _ = app.write_help(&mut vec);

    return String::from_utf8(vec).unwrap();
}

fn init_log(config_file: &str) -> Result<(), Error> {
    use log4rs;

    if config_file.len() == 0 {
        // Default logging
        use log::LevelFilter;
        use log4rs::append::console::ConsoleAppender;
        use log4rs::append::file::FileAppender;
        use log4rs::encode::pattern::PatternEncoder;
        use log4rs::config::{Appender, Config, Logger, Root};

        let pattern = "{d(%Y-%m-%dT%H:%M:%S %Z)(utc)} {h([{l}])} at {M}:{L} {m}{n}";

        let stdout = ConsoleAppender::builder()
            .encoder(Box::new(PatternEncoder::new(pattern)))
            .build();

        let lfs_logger = Logger::builder().build("learnfs", LevelFilter::Debug);

        let config = Config::builder()
            .appender(Appender::builder().build("stdout", Box::new(stdout)))
            .logger(lfs_logger)
            .build(Root::builder().appender("stdout").build(LevelFilter::Info))
            .unwrap();

        log4rs::init_config(config)?;
    } else {
        log4rs::init_file(config_file, Default::default())?
    }

    Ok(())
}
