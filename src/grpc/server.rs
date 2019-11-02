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

use failure::Error;
use tonic::transport as tr;
use tonic::{Request, Response, Status};

mod proto {
    tonic::include_proto!("learnfs");
}

use proto::{
    server::{LearnFs, LearnFsServer},
    RmdirRequest, CallResult,
};

type GrpcResult = Result<Response<CallResult>, Status>;

pub struct Server {}

#[tonic::async_trait]
impl LearnFs for Server {
    async fn rmdir(
        &self, request: Request<RmdirRequest>) -> GrpcResult {
        println!("Got a request: {:?}", request);

        let reply = CallResult {
            code: 1,
        };

        Ok(Response::new(reply))
    }
}

impl Server {
    pub async fn run(addr: &str) -> Result<(), Error> {
        let paddr = addr.parse()?;
        let server = Server{};

        tr::Server::builder()
            .add_service(LearnFsServer::new(server))
            .serve(paddr)
            .await?;

        Ok(())
    }
}

