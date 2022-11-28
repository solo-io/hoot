// Copyright 2022 Daniel Hawton
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

use log::{info, debug, trace};
use proxy_wasm as wasm;
use wasm::{types::Action, types::ContextType};

wasm::main! {{
    wasm::set_log_level(wasm::types::LogLevel::Trace);
    wasm::set_root_context(|_| -> Box<dyn wasm::traits::RootContext> {
        Box::new(RustTest)
    });
}}

struct RustTest;

struct HttpHeaders {
    context_id: u32,
}

impl wasm::traits::Context for RustTest {}

impl wasm::traits::RootContext for RustTest {
    fn on_vm_start(&mut self, _vm_configuration_size: usize) -> bool {
        info!("on_vm_start");
        true
    }

    fn get_type(&self) -> Option<ContextType> {
        Some(ContextType::HttpContext)
    }

    fn create_http_context(&self, context_id: u32) -> Option<Box<dyn wasm::traits::HttpContext>> {
        Some(Box::new(HttpHeaders { context_id }))
    }
}

const TEAPOT_ASCII = b"I'm a teapot

                       (
            _           ) )
         _,(_)._        ((
    ___,(_______).        )
  ,'__.   /       \\    /\\_
 /,' /  |\"\"|       \\  /  /
| | |   |__|       |,'  /
 \\`.|                  /
  `. :           :    /
    `.            :.,'
      `-.________,-'
";

impl wasm::traits::Context for HttpHeaders {}

impl wasm::traits::HttpContext for HttpHeaders {
    fn on_http_request_headers(&mut self, _: usize, _: bool) -> wasm::types::Action {
        info!("on_http_request_headers: {}", self.context_id);
        for (name, value) in &self.get_http_request_headers() {
            trace!("#{} - {} = {}", self.context_id, name, value);
        }

        const path = self.get_http_request_header(":path");
        const method = self.get_http_request_header(":method");

        match self.get_http_request_header(":path") {
            Some(path) if path == "/get" => {
                info!("on_http_request_headers: {} - /get intercepted", self.context_id);
                self.send_http_response(
                    418,
                    vec![("x-powered-by", "rust"), ("content-type", "text/plain")],
                    Some(TEAPOT_ASCII),
                );
                Action::Pause
            }
            _ => Action::Continue,
        }
    }

    fn on_log(&mut self) {
        info!("#{} completed.", self.context_id);
    }
}