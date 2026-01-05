struct User {
    name: String,
    email: String,
    age: u32,
}

impl User {
    fn new(name: String, email: String, age: u32) -> Self {
        User { name, email, age }
    }

    fn get_name(&self) -> &str {
        &self.name
    }
}

fn main() {
    let user = User {
        name: String::from("John"),
        email: String::from("john@example.com"),
        age: 30,
    };

    println!("{}", user.name);
    println!("{}", user.get_name());
}

