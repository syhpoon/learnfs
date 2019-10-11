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

mod fs;
mod device;

use clap::{Arg, App, SubCommand, AppSettings};
use clap::{crate_version, value_t};
use failure::{Error, format_err};
use nix::sys::stat::{stat, SFlag};

use crate::device::FileDevice;
use crate::fs::Filesystem;

enum DeviceType {
    File,
    Block,
}

fn main() -> Result<(), Error> {
    let app = App::new("learnfs")
        .version(crate_version!())
        .about("LearnFS command line tool")
        .subcommand(
            SubCommand::with_name("create")
                .about("Create filesystem on a given device")
                .setting(AppSettings::DisableVersion)
                .arg(
                    Arg::with_name("device")
                        .help("Block device to use")
                        .index(1)
                        .required(true)
                )
                .arg(
                    Arg::with_name("capacity")
                        .help("Explicitly set file capacity in bytes")
                        .short("c")
                        .long("capacity")
                        .takes_value(true)
                        .default_value("0")
                )
        ).subcommand(
        SubCommand::with_name("info")
            .about("Show existing filesystem info")
            .setting(AppSettings::DisableVersion)
            .arg(
                Arg::with_name("device")
                    .help("Block device to use")
                    .index(1)
                    .required(true)
            )
    );

    let help = extract_help(&app);

    match app.get_matches().subcommand() {
        ("create", Some(create_match)) => {
            let path = create_match.value_of("device").unwrap();

            match device_type(path) {
                Ok(DeviceType::File) => {
                    let capacity = value_t!(create_match,
                     "capacity", u64).unwrap();

                    if capacity == 0 {
                        return Err(format_err!(
                         "please explicitly specify file device capacity"));
                    }

                    let mut device = FileDevice::new(path, capacity)?;
                    fs::Filesystem::create(device).expect("wut");
                }
                Ok(DeviceType::Block) => {
                    unimplemented!()
                }
                Err(e) => {
                    return Err(e);
                }
            }
        }
        ("info", Some(m)) => {
            let path = m.value_of("device").unwrap();

            match device_type(path) {
                Ok(DeviceType::File) => {
                    let device = FileDevice::open(path)?;
                    let fs = Filesystem::load(device)?;

                    println!("{:#?}", fs.info())
                }
                Ok(DeviceType::Block) => {
                    unimplemented!()
                }
                Err(e) => {
                    return Err(e);
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

fn device_type(path: &str) -> Result<DeviceType, Error> {
    match stat(path) {
        Ok(st) => {
            match SFlag::from_bits(st.st_mode & SFlag::S_IFMT.bits()) {
                // Block device
                Some(SFlag::S_IFBLK) => {
                    Ok(DeviceType::Block)
                }
                // Regular file
                Some(SFlag::S_IFREG) => {
                    Ok(DeviceType::File)
                }
                _ => {
                    Err(format_err!("Only block devices and regular files are currently supported"))
                }
            }
        }
        Err(e) => {
            Err(format_err!("Failed to read device stat: {}", e))
        }
    }
}
