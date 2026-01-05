pub mod models {
    pub struct Customer {
        pub id: u64,
        pub name: String,
        pub email: String,
    }

    impl Customer {
        pub fn new(id: u64, name: String, email: String) -> Self {
            Customer { id, name, email }
        }
    }
}

pub mod services {
    use super::models::Customer;

    pub fn get_customer(id: u64) -> Option<Customer> {
        Some(Customer::new(id, String::from("Test"), String::from("test@example.com")))
    }
}

