use std::collections::HashMap;

mod http {
    pub struct HttpClient;

    impl HttpClient {
        pub fn new(url: &str) -> Self {
            HttpClient
        }
    }
}

use http::HttpClient;

fn main() {
    let client = HttpClient::new("https://example.com");
}

