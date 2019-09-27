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

use std::fs::OpenOptions;

use clap::{Arg, App, SubCommand, AppSettings};
use clap::crate_version;

mod cmd;
mod fs;

fn main() {
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
        );

    let help = extract_help(&app);

    match app.get_matches().subcommand() {
        ("create", Some(create_match)) => {
            let device = create_match.value_of("device").unwrap();

            let mut file = OpenOptions::new()
                .write(true)
                .truncate(true)
                .create(true)
                .open(device).unwrap();

            cmd::create::create_fs(&mut file).unwrap();
        }
        _ => {
            eprintln!("{}", help);
        }
    }

    fn extract_help(app: &App) -> String {
        let mut vec: Vec<u8> = Vec::new();
        let _ = app.write_help(&mut vec);

        return String::from_utf8(vec).unwrap();
    }
}
