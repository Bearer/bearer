struct User {
    username: String,
    email: String,
    password: String,
    age: u32,
}

struct Address {
    street: String,
    city: String,
    zip_code: String,
}

impl User {
    fn new(username: String, email: String, password: String) -> Self {
        User {
            username,
            email,
            password,
            age: 0,
        }
    }

    fn get_email(&self) -> &str {
        &self.email
    }

    fn set_password(&mut self, password: String) {
        self.password = password;
    }
}

fn create_user(name: String, email: String) -> User {
    User::new(name, email, String::from("default"))
}

fn main() {
    let user = User::new(
        String::from("john"),
        String::from("john@example.com"),
        String::from("secret123"),
    );

    println!("User: {}", user.username);
}

