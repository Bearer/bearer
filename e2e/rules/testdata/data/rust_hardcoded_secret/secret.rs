const API_KEY: &str = "sk-1234567890abcdef";

fn main() {
    let password = "mysecretpassword123";
    let token = "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx";
    
    let safe_value = "hello world";
    
    connect_api(password, token);
}

fn connect_api(pass: &str, tok: &str) {
    println!("Connecting...");
}

